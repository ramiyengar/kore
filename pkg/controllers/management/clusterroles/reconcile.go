/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of hub-apiserver.
 *
 * hub-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * hub-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with hub-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package clusterroles

import (
	"context"
	"errors"
	"fmt"
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/hub"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	rbacv1 "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	rolesFinalizer = "managedclusterroles.hub.appvia.io"
)

// Reconcile ensures the clusters roles across all the managed clusters
func (a crCtrl) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	logger.Debug("attempting to reconcile managed cluster roles")

	ctx := context.Background()

	// @step: retrieve the resource from the api
	role := &clustersv1.ManagedClusterRole{}
	err := a.mgr.GetClient().Get(ctx, request.NamespacedName, role)
	if err != nil {
		if !kerrors.IsNotFound(err) {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}
	original := role.DeepCopy()

	failed := &clustersv1.KubernetesList{}
	var list *clustersv1.KubernetesList

	finalizer := kubernetes.NewFinalizer(a.mgr.GetClient(), rolesFinalizer)
	if finalizer.IsDeletionCandidate(role) {
		return a.Delete(request)
	}

	logger.Debug("attempting to retrieve a list of cluster applicable")

	// @step: retrieve a list of cluster which this role applies
	list, err = a.FilterClustersBySource(ctx,
		role.Spec.Clusters,
		role.Spec.Teams,
		role.Namespace)
	if err != nil {
		logger.WithError(err).Error("trying to retrieve a list of clusters")

		role.Status.Status = corev1.FailureStatus
		role.Status.Conditions = []corev1.Condition{{
			Detail:  err.Error(),
			Message: "failed trying to retrieve list of clusters to apply",
		}}

		if err := a.mgr.GetClient().Status().Patch(ctx, role, client.MergeFrom(original)); err != nil {
			logger.WithError(err).Error("trying to update the resource")

			return reconcile.Result{}, err
		}

		return reconcile.Result{}, err
	}

	logger.WithField("clusters", len(list.Items)).Debug("applying the change to x clusters")

	err = func() error {
		// @step: we iterate the clusters and apply the roles
		for _, cluster := range list.Items {
			logger := logger.WithFields(log.Fields{
				"cluster": cluster.Name,
				"team":    cluster.Namespace,
			})
			logger.Debug("attempting to reconcile the managed role in cluster")

			err := func() error {
				client, err := controllers.CreateClientFromSecret(ctx, a.mgr.GetClient(),
					cluster.Name, cluster.Namespace)
				if err != nil {
					logger.WithError(err).Error("trying to create kubernetes client")

					return err
				}
				logger.Debug("creating the managed cluster role in cluster")

				rules := role.Spec.Rules
				if !role.Spec.Enabled {
					rules = []rbacv1.PolicyRule{}
				}

				// @step: update or create the role
				if _, err := kubernetes.CreateOrUpdateManagedClusterRole(ctx, client, &rbacv1.ClusterRole{
					ObjectMeta: metav1.ObjectMeta{
						Name: role.Name,
						Labels: map[string]string{
							hub.Label("owned"): "true",
						},
					},
					Rules: rules,
				}); err != nil {
					logger.WithError(err).Error("trying to update or crreate the managed role")

					return err
				}

				return nil
			}()
			if err != nil {
				failed.Items = append(failed.Items, cluster)
			}
		}

		if len(failed.Items) > 0 {
			logger.WithFields(log.Fields{
				"failed": len(failed.Items),
				"total":  len(list.Items),
			}).Warn("we failed to provision on all the clusters")

			role.Status.Status = corev1.WarningStatus
			role.Status.Conditions = []corev1.Condition{{
				Message: fmt.Sprintf("failed to provision on all clusters, %d out of %d failed", len(failed.Items), len(list.Items)),
			}}

			return errors.New("provisioning managed cluster roles")
		}

		logger.Debug("successfully updated the managed role in the clusters")

		role.Status.Status = corev1.SuccessStatus
		role.Status.Conditions = []corev1.Condition{}

		return nil
	}()
	if err != nil {
		logger.WithError(err).Error("trying to reconcile the managed cluster role")
	} else {
		if finalizer.NeedToAdd(role) {
			if err := finalizer.Add(role); err != nil {
				logger.WithError(err).Error("trying to add the finalizer")

				return reconcile.Result{}, err
			}

			return reconcile.Result{Requeue: true}, nil
		}
	}

	if err := a.mgr.GetClient().Status().Patch(ctx, role, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("trying to update the managed cluster role status")

		return reconcile.Result{}, err
	}

	return reconcile.Result{RequeueAfter: 15 * time.Minute}, err
}