package v1alpha1

import (
	core "github.com/appvia/kore/pkg/apis/core/v1"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GCPProjectClaimSpec defines the desired state of GCPProjectClaim
// +k8s:openapi-gen=true
type GCPProjectClaimSpec struct {
	// Organization is a reference to the gcp admin project to use
	// +kubebuilder:validation:Required
	Organization core.Ownership `json:"organization"`
}

// GCPProjectClaimStatus defines the observed state of GCPProject
// +k8s:openapi-gen=true
type GCPProjectClaimStatus struct {
	// CredentialRef is the reference to the credentials secret
	CredentialRef *v1.SecretReference `json:"credentialRef,omitempty"`
	// ProjectID is the GCP project id
	ProjectID string `json:"projectID,omitempty"`
	// Status provides a overall status
	Status core.Status `json:"status,omitempty"`
	// Conditions is a set of components conditions
	Conditions *core.Components `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GCPProjectClaim is the Schema for the gcpprojects API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=gcpprojectclaim,scope=Namespaced
type GCPProjectClaim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GCPProjectClaimSpec   `json:"spec,omitempty"`
	Status GCPProjectClaimStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GCPProjectClaimList contains a list of GCPProjectClaim
type GCPProjectClaimList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GCPProjectClaim `json:"items"`
}
