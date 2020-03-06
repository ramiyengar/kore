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

	gcp "github.com/appvia/kore/pkg/apis/gcp/v1alpha1"

	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateCredentialsSecret returns a project credentials secret
func CreateCredentialsSecret(project *gcp.GCPProjectClaim, name string, key []byte) *v1.Secret {
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: project.Namespace,
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion:         gcp.GroupVersion.String(),
				BlockOwnerDeletion: &isTrue,
				Controller:         &isTrue,
				Kind:               "GCPProjectClaim",
				Name:               project.GetName(),
				UID:                project.GetUID(),
			}},
		},
		Data: map[string][]byte{"key": key},
	}
}

// IsProject checks if the project exists
func IsProject(ctx context.Context, client *cloudresourcemanager.Service, name string) (*cloudresourcemanager.Project, bool, error) {
	list, err := client.Projects.List().Context(ctx).Do()
	if err != nil {
		return nil, false, err
	}
	for _, x := range list.Projects {
		if x.Name == name {
			return x, true, nil
		}
	}

	return nil, false, nil
}

func newClient() (*cloudresourcemanager.Service, error) {
	return nil, nil
}
