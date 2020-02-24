/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
 */

package gcpprojectclaim

import (
	"context"
	"time"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	gcp "github.com/appvia/kore/pkg/apis/gcp/v1alpha1"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	finalizerName = "gcp-project-claims.kore.appvia.io"
)

// Reconcile is the entrypoint for the reconciliation logic
func (t ctrl) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	logger := log.WithFields(log.Fields{
		"project": request.NamespacedName.Name,
		"team":    request.NamespacedName.Namespace,
	})
	logger.Info("attempting to reconcile gcp project")

	project := &gcp.GCPProjectClaim{}
	if err := t.mgr.GetClient().Get(ctx, request.NamespacedName, project); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}
	original := project.DeepCopy()

	// @step: ensure we have components in the status
	if project.Status.Conditions == nil {
		project.Status.Conditions = &corev1.Components{}
	}

	finalizer := kubernetes.NewFinalizer(t.mgr.GetClient(), finalizerName)
	if finalizer.IsDeletionCandidate(project) {
		return t.Delete(request)
	}

	result, err := func() (reconcile.Result, error) {
		// @step: ensure the project has access to the org
		if err := t.EnsurePermitted(ctx, project); err != nil {
			logger.WithError(err).Error("checking if project has permission to gcp organization")

			return reconcile.Result{}, err
		}

		// @step: ensure thr project has not been claimed already
		if err := t.EnsureUnclaimed(ctx, project); err != nil {
			logger.WithError(err).Error("checking if project is claimed")

			return reconcile.Result{}, err
		}

		// @step: ensure the gcp organization
		org, err := t.EnsureOrganization(ctx, project)
		if err != nil {
			logger.WithError(err).Error("trying to retrieve the gcp organization")

			return reconcile.Result{RequeueAfter: 2 * time.Minute}, err
		}

		// @step: we need to grab the credentials from the organization and create clients
		secret, err := t.EnsureOrganizationCredentials(ctx, org, project)
		if err != nil {
			logger.WithError(err).Error("trying to retrieve the gcp organization")

			return reconcile.Result{}, err
		}

		// @step: ensure the project is created
		if err := t.EnsureProject(ctx, secret, org, project); err != nil {
			logger.WithError(err).Error("trying to ensure the project")

			return reconcile.Result{}, err
		}

		// @step: ensure the project is linked to the billing account
		if err := t.EnsureBilling(ctx, secret, org, project); err != nil {
			logger.WithError(err).Error("trying to ensure the billing account it linked")

			return reconcile.Result{}, err
		}
		// @step: ensure the project apis are enabled
		if err := t.EnsureAPIs(ctx, secret, project); err != nil {
			logger.WithError(err).Error("trying to toggle the apis in the project")

			return reconcile.Result{}, err
		}

		// @step: ensure the service account in the project
		if err := t.EnsureServiceAccount(ctx, secret, project); err != nil {
			logger.WithError(err).Error("trying to enable the service account in project")

			return reconcile.Result{}, err
		}

		// @step: ensure the service account key in the project
		key, err := t.EnsureServiceAccountKey(ctx, secret, org, project)
		if err != nil {
			logger.WithError(err).Error("trying to ensure the service account key")

			return reconcile.Result{}, err
		}

		// @step: we need to mint the GKE credentials for them

		project.Status.Status = corev1.SuccessStatus

		return reconcile.Result{}, nil
	}()
	if err != nil {
		logger.WithError(err).Error("trying to reconcile the gcp project")

		project.Status.Status = corev1.FailureStatus
	}
	if err == nil {
		if finalizer.NeedToAdd(project) {
			if err := finalizer.Add(project); err != nil {
				logger.WithError(err).Error("trying to add the finalizer to resource")

				return reconcile.Result{}, err
			}

			return reconcile.Result{Requeue: true}, nil
		}
	}

	if err := t.mgr.GetClient().Status().Patch(ctx, project, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("updating the project status")

		return reconcile.Result{}, err
	}

	return result, nil
}

// IsProjectClaimed checks if the project name has already been claimed by another team
func (t ctrl) IsProjectClaimed(ctx context.Context, project *gcp.GCPProjectClaim) (bool, error) {
	list := &gcp.GCPProjectClaimList{}

	if err := t.mgr.GetClient().List(ctx, list, client.InNamespace("")); err != nil {
		return false, err
	}

	for _, x := range list.Items {
		if x.Name == project.Name && x.Namespace != project.Namespace {
			return true, nil
		}
	}

	return false, nil
}
