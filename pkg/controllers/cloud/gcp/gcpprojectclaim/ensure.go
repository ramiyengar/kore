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
	"errors"
	"time"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	gcp "github.com/appvia/kore/pkg/apis/gcp/v1alpha1"
	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	cloudbilling "google.golang.org/api/cloudbilling/v1"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
	iam "google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
	servicemanagement "google.golang.org/api/servicemanagement/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// EnsurePermitted is responsible for checking the project has access to the credentials
func (t ctrl) EnsurePermitted(ctx context.Context, project *gcp.GCPProjectClaim) error {
	// @step: we check if the gcp organization has been allocated to us
	permitted, err := t.Teams().Team(project.Namespace).Allocations().IsPermitted(ctx, project.Spec.Organization)
	if err != nil {
		return err
	}
	if !permitted {
		project.Status.Conditions.SetCondition(corev1.Component{
			Name:    "provision",
			Message: "GCP Organization has not been allocated to you",
			Status:  corev1.FailureStatus,
		})

		return errors.New("gcp organization has not been allocated to team")
	}

	return nil
}

// EnsureOrganization is responsible for checking and retrieving the gcp org
func (t ctrl) EnsureOrganization(ctx context.Context, project *gcp.GCPProjectClaim) (*gcp.GCPAdminProject, error) {
	org := &gcp.GCPAdminProject{}

	key := types.NamespacedName{
		Namespace: project.Spec.Organization.Namespace,
		Name:      project.Spec.Organization.Name,
	}

	if err := t.mgr.GetClient().Get(ctx, key, org); err != nil {
		project.Status.Conditions.SetCondition(corev1.Component{
			Name:    "provision",
			Detail:  err.Error(),
			Message: "Attempting to retrieve the GCP Organization resources from API",
			Status:  corev1.FailureStatus,
		})

		return nil, err
	}

	// @step: check if the admin project exists and if successful
	if org.Status.Status != corev1.SuccessStatus {
		project.Status.Conditions.SetCondition(corev1.Component{
			Name:    "provision",
			Detail:  "resource is in failing state",
			Message: "GCP Admin Project is in a failing state, cannot provision projects",
			Status:  corev1.FailureStatus,
		})

		return nil, errors.New("admin project still provisioning or failed")
	}

	return org, nil

}

// EnsureUnclaimed is responsible for making sure the project is unclaimed
func (t ctrl) EnsureUnclaimed(ctx context.Context, project *gcp.GCPProjectClaim) error {
	logger := log.WithFields(log.Fields{
		"project": project.Name,
		"team":    project.Namespace,
	})

	// @step: check if the project claim has already been claimed else where
	claimed, err := t.IsProjectClaimed(ctx, project)
	if err != nil {
		logger.WithError(err).Error("trying to check if the project is already claimed")

		project.Status.Status = corev1.FailureStatus
		project.Status.Conditions.SetCondition(corev1.Component{
			Name:    "provision",
			Message: "Unable to fulfil request, project name has already been claimed in the organization",
			Status:  corev1.FailureStatus,
		})

		return errors.New("failed to check if project is already claimed")
	}
	if claimed {
		logger.Warn("attempting to claim gcp project which has already been provisioned")

		project.Status.Status = corev1.FailureStatus
		project.Status.Conditions.SetCondition(corev1.Component{
			Name:    "provision",
			Message: "Project has already been claimed by another team in kore",
			Status:  corev1.FailureStatus,
		})

		return errors.New("gcp project name already provisioned")
	}

	return nil
}

// EnsureOrganizationCredentials is responsible for retrieving the credentials
func (t ctrl) EnsureOrganizationCredentials(ctx context.Context, org *gcp.GCPAdminProject, project *gcp.GCPProjectClaim) (*v1.Secret, error) {
	// @TODO we probably shouldn't rely on the parent name here
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      org.Name,
			Namespace: org.Namespace,
		},
	}

	found, err := kubernetes.GetIfExists(ctx, t.mgr.GetClient(), secret)
	if err != nil {
		project.Status.Conditions.SetCondition(corev1.Component{
			Name:    "provision",
			Detail:  err.Error(),
			Message: "Attempting to retrieve the GCP Organization credentials",
			Status:  corev1.FailureStatus,
		})

		return nil, err
	}
	if !found {
		project.Status.Conditions.SetCondition(corev1.Component{
			Name:    "provision",
			Detail:  "credentials not found",
			Message: "GCP Organization credentials either not provisioned or failed",
			Status:  corev1.FailureStatus,
		})

		return nil, errors.New("credentials not found")
	}

	return secret, nil
}

// EnsureProject is responsible for ensuring the project is there
func (t ctrl) EnsureProject(ctx context.Context,
	credentials *v1.Secret,
	org *gcp.GCPAdminProject,
	project *gcp.GCPProjectClaim) error {

	logger := log.WithFields(log.Fields{
		"project": project.Name,
		"team":    project.Namespace,
	})
	stage := "provision"

	// @step: create the client
	client, err := cloudresourcemanager.NewService(ctx, option.WithCredentialsJSON(credentials.Data["key"]))
	if err != nil {
		logger.WithError(err).Error("trying to create the cloud resource client")

		project.Status.Conditions.SetCondition(corev1.Component{
			Name:    stage,
			Detail:  err.Error(),
			Message: "Failed to create a projects client, please check credentials",
			Status:  corev1.FailureStatus,
		})

		return err
	}

	// @step: we check if the project exists and if not create it
	_, found, err := IsProject(ctx, client, project.Name)
	if found {
		logger.Debug("gcp project already exists, checking if it was created by us")

		// @TODO we need something to check in the project to see if we create this project

		return nil
	}

	logger.Info("gcp project does not exist, creating it now")

	// @step: create the project in gcp
	resp, err := client.Projects.Create(&cloudresourcemanager.Project{
		Name: project.Name,
		Parent: &cloudresourcemanager.ResourceId{
			Id:   org.Spec.ParentID,
			Type: org.Spec.ParentType,
		},
	}).Context(ctx).Do()

	if err != nil {
		logger.WithError(err).Error("trying to create the gcp project")

		project.Status.Conditions.SetCondition(corev1.Component{
			Name:    stage,
			Detail:  err.Error(),
			Message: "Unable to request the project in GCP",
			Status:  corev1.FailureStatus,
		})

		return err
	}

	// @step: wait for the operation to complete or fail
	if err := utils.WaitUntilComplete(ctx, 5*time.Minute, 10*time.Second, func() (bool, error) {
		status, err := client.Operations.Get(resp.Name).Context(ctx).Do()
		if err != nil {
			logger.WithError(err).Error("checking the status of the project operation")

			return false, nil
		}
		if !status.Done {
			return false, nil
		}
		if status.Error != nil {
			return false, errors.New(status.Error.Message)
		}

		return true, nil
	}); err != nil {
		project.Status.Conditions.SetCondition(corev1.Component{
			Name:    stage,
			Detail:  err.Error(),
			Message: "Unable to provision project in GCP",
			Status:  corev1.FailureStatus,
		})

		return err
	}

	return nil
}

// EnsureBilling is responsible for ensuring the billing account
func (t ctrl) EnsureBilling(
	ctx context.Context,
	credentials *v1.Secret,
	organization *gcp.GCPAdminProject,
	project *gcp.GCPProjectClaim) error {

	logger := log.WithFields(log.Fields{
		"project": project.Name,
		"team":    project.Namespace,
	})
	stage := "billing"

	err := func() error {
		client, err := cloudbilling.NewService(ctx, option.WithCredentialsJSON(credentials.Data["key"]))
		if err != nil {
			logger.WithError(err).Error("trying to create cloud resource client")

			return err
		}

		resp, err := client.Projects.GetBillingInfo(project.Name).Context(ctx).Do()
		if err != nil {
			logger.WithError(err).Error("trying to retrieve the billing details for account")

			return err
		}

		// @if they are the same we can return
		if resp.BillingAccountName == organization.Spec.BillingAccountName {
			return nil
		}

		if resp.BillingAccountName == "" {
			logger.Info("billing account not set, attempting to set now")
		}
		if resp.BillingAccountName != "" && resp.BillingAccountName != organization.Spec.BillingAccountName {
			logger.Warn("project billing account differs, trying to reconcile now")
		}

		if _, err := client.Projects.UpdateBillingInfo(project.Name, &cloudbilling.ProjectBillingInfo{
			BillingAccountName: "billingAccounts/" + organization.Spec.BillingAccountName,
			BillingEnabled:     true,
		}).Context(ctx).Do(); err != nil {
			logger.WithError(err).Error("trying to update the project billing details")

			return err
		}

		return err
	}()
	if err != nil {
		project.Status.Conditions.SetCondition(corev1.Component{
			Name:    stage,
			Detail:  err.Error(),
			Message: "Failed to link the billing account to project",
			Status:  corev1.FailureStatus,
		})

		return err
	}

	project.Status.Conditions.SetCondition(corev1.Component{
		Name:    stage,
		Message: "GCP Project has been linked billing account",
		Status:  corev1.SuccessStatus,
	})

	return nil
}

// EnsureAPIs is responsible for ensuing the apis are enabled in the account
func (t ctrl) EnsureAPIs(ctx context.Context, credentials *v1.Secret, project *gcp.GCPProjectClaim) error {
	stage := "iam"

	logger := log.WithFields(log.Fields{
		"project": project.Name,
		"team":    project.Namespace,
	})

	client, err := servicemanagement.NewService(ctx, option.WithCredentialsJSON(credentials.Data["key"]))
	if err != nil {
		logger.WithError(err).Error("trying to create the service management client")

		project.Status.Conditions.SetCondition(corev1.Component{
			Name:    stage,
			Detail:  err.Error(),
			Message: "Failed to create the service management client, please check credentials",
			Status:  corev1.FailureStatus,
		})

		return err
	}

	for _, name := range t.GetRequiredAPI() {
		logger.WithField(
			"api", name,
		).Debug("attempting to enable the api in the project")

		request := &servicemanagement.EnableServiceRequest{
			ConsumerId: "project:" + project.Name,
		}

		resp, err := client.Services.Enable(name, request).Context(ctx).Do()
		if err != nil {
			logger.WithError(err).Error("trying to enable the api")

			project.Status.Conditions.SetCondition(corev1.Component{
				Name:    stage,
				Detail:  err.Error(),
				Message: "Failed to enable " + name + " api in the project",
				Status:  corev1.FailureStatus,
			})

			return err
		}
		logger.Debug("successfully enabled the api in the project")

		if err := utils.WaitUntilComplete(ctx, 3*time.Minute, 5*time.Second, func() (bool, error) {
			status, err := client.Operations.Get(resp.Name).Context(ctx).Do()
			if err != nil {
				logger.WithError(err).Error("trying to retrieve status of operation")

				return false, nil
			}
			if !status.Done {
				return false, nil
			}
			if status.Error != nil {
				return false, errors.New(status.Error.Message)
			}

			return true, nil
		}); err != nil {
			logger.WithError(err).Error("waiting on the api enabling operation")

			project.Status.Conditions.SetCondition(corev1.Component{
				Name:    stage,
				Detail:  err.Error(),
				Message: "Failed to enable " + name + " api in the project",
				Status:  corev1.FailureStatus,
			})

			return err
		}
	}

	project.Status.Conditions.SetCondition(corev1.Component{
		Name:    stage,
		Message: "Successfully enabled all the APIs in project",
		Status:  corev1.SuccessStatus,
	})

	return nil
}

// EnsureServiceAccount is responsible for creating the service account in the project
func (t ctrl) EnsureServiceAccount(ctx context.Context, credentials *v1.Secret, project *gcp.GCPProjectClaim) error {
	stage := "iam"
	account := project.Spec.ServiceAccountName
	if account == "" {
		account = t.DefaultServiceAccountName()
	}

	logger := log.WithFields(log.Fields{
		"account": account,
		"project": project.Name,
		"team":    project.Namespace,
	})

	err := func() error {
		// @step: create the iam client
		client, err := iam.NewService(ctx, option.WithCredentialsJSON(credentials.Data["key"]))
		if err != nil {
			logger.WithError(err).Error("trying to create the client")

			return err
		}

		// @step: ensure the service account exists in the project
		list, err := client.Projects.ServiceAccounts.List(account).Context(ctx).Do()
		if err != nil {
			logger.WithError(err).Error("trying to retrieve the service account list")

			return err
		}
		if len(list.Accounts) <= 0 {
			logger.Debug("service account does not exist, creating now")

			if _, err := client.Projects.ServiceAccounts.Create("projects/"+project.Name, &iam.CreateServiceAccountRequest{
				AccountId: account,
				ServiceAccount: &iam.ServiceAccount{
					DisplayName: "Kore Service Account",
				},
			}).Context(ctx).Do(); err != nil {

				logger.WithError(err).Error("trying to create the service account in project")

				return err
			}
		}

		return nil
	}()
	if err != nil {
		logger.WithError(err).Error("attempting to provision the service account")

		project.Status.Conditions.SetCondition(corev1.Component{
			Name:    stage,
			Detail:  err.Error(),
			Message: "Failed to provision the IAM credentials in the project",
			Status:  corev1.FailureStatus,
		})

		return err
	}

	project.Status.Conditions.SetCondition(corev1.Component{
		Name:    stage,
		Message: "Successfully provision the IAM in project",
		Status:  corev1.SuccessStatus,
	})

	return nil
}

// EnsureServiceAccountKey is responsible for ensuring the account key exists
func (t ctrl) EnsureServiceAccountKey(
	ctx context.Context,
	credentials *v1.Secret,
	organization *gcp.GCPAdminProject,
	project *gcp.GCPProjectClaim) error {

	stage := "permissions"

	account := project.Spec.ServiceAccountName
	if account == "" {
		account = t.DefaultServiceAccountName()
	}

	logger := log.WithFields(log.Fields{
		"account": account,
		"project": project.Name,
		"team":    project.Namespace,
	})
	var key *iam.ServiceAccountKey

	err := func() error {
		// @step: create the iam client
		client, err := iam.NewService(ctx, option.WithCredentialsJSON(credentials.Data["key"]))
		if err != nil {
			logger.WithError(err).Error("trying to create iam client for project")

			return err
		}

		// @step: check if the service account key exists already
		list, err := client.Projects.ServiceAccounts.Keys.List(account).Context(ctx).Do()
		if err != nil {
			logger.WithError(err).Error("trying to check if service account key exists already")

			return err
		}
		if len(list.Keys) > 0 {
			logger.Debug("service account key exists already, skipping creation")
			key = list.Keys[0]

			return nil
		} else {
			key, err = client.Projects.ServiceAccounts.Keys.Create(account, &iam.CreateServiceAccountKeyRequest{
				KeyAlgorithm:   "KEY_ALG_RSA_2048",
				PrivateKeyType: "TYPE_GOOGLE_CREDENTIALS_FILE",
			}).Context(ctx).Do()

			if err != nil {
				logger.WithError(err).Error("trying to provision the service account key")

				return err
			}
		}

		encoded, err := key.MarshalJSON()
		if err != nil {
			logger.WithError(err).Error("trying to marshal the service account key")

			return err
		}

		// @step: create the gke credentials in the account
		gkecreds := &gke.GKECredentials{
			ObjectMeta: metav1.ObjectMeta{
				Name:      project.Name,
				Namespace: project.Namespace,
			},
			Spec: gke.GKECredentialsSpec{
				Account: key.M,
			},
		}

		return nil
	}()
	if err != nil {
		project.Status.Conditions.SetCondition(corev1.Component{
			Name:    stage,
			Detail:  err.Error(),
			Message: "Failed to provision the service account key in the project",
			Status:  corev1.FailureStatus,
		})

		return nil, err
	}

	project.Status.Conditions.SetCondition(corev1.Component{
		Name:    stage,
		Message: "Successfully provision the service account in project",
		Status:  corev1.SuccessStatus,
	})

	return key, nil
}

// DefaultServiceAccountName is the default name of the service account
func (t ctrl) DefaultServiceAccountName() string {
	return "kore"
}

// GetRequiredAPI returns a list of required apis
func (t ctrl) GetRequiredAPI() []string {
	return []string{
		"cloudbilling.googleapis.com",
		"cloudresourcemanager.googleapis.com",
		"compute.googleapis.com",
		"iam.googleapis.com",
		"serviceusage.googleapis.com",
	}
}
