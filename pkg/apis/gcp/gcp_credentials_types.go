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

package v1alpha1

import (
	core "github.com/appvia/kore/pkg/apis/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GCPCredentialsSpec defines the desired state of GCPCredentials
// +k8s:openapi-gen=true
type GCPCredentialsSpec struct {
	// ServiceAccount is the credential used to create GCP projects
	// You must create a service account with resourcemanager.projectCreator
	// and billing.user roles at the organization level and use the JSON payload here
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Required
	ServiceAccount string `json:"serviceAccount,omitempty"`
	// Project is the GCP project ID these credentials belong to
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Required
	Project string `json:"project"`
	// Organization is the GCP org you wish the projects to reside within
	// +kubebuilder:validation:Minimum=2
	// +kubebuilder:validation:Required
	Organization string `json:"organization"`
}

// GCPCredentialsStatus defines the observed state of GCPCredentials
// +k8s:openapi-gen=true
type GCPCredentialsStatus struct {
	// Conditions is a collection of conditions of errors
	// +kubebuilder:validation:Optional
	// +listType
	Conditions []core.Condition `json:"conditions,omitempty"`
	// Status provides a overall status
	// +kubebuilder:validation:Optional
	Status string `json:"status,omitempty"`
	// Verified checks that the credentials are ok and valid
	// +kubebuilder:validation:Optional
	Verified bool `json:"verified,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GCPCredentials is the Schema for the gcpcredentials API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=gcpcredentials,scope=Namespaced
type GCPCredentials struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GCPCredentialsSpec   `json:"spec,omitempty"`
	Status GCPCredentialsStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GCPCredentialsList contains a list of GCPCredentials
type GCPCredentialsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GCPCredentials `json:"items"`
}
