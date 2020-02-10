package v1alpha1

import (
	core "github.com/appvia/hub-apis/pkg/apis/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GCPAdminProjectSpec defines the desired state of GCPAdminProject
// +k8s:openapi-gen=true
type GCPAdminProjectSpec struct {
	// Project is the GCP project ID
	// +kubebuilder:validation:Minimum=3
	// +kubebuilder:validation:Required
	Project string `json:"project"`
	// ProjectName is the GCP project name
	// +kubebuilder:validation:Minimum=3
	// +kubebuilder:validation:Required
	ProjectName string `json:"projectName"`
	// ParentType is the type of parent this project has
	// Valid types are: "organization", "folder", and "project"
	// +kubebuilder:validation:Enum=organization;folder;project
	// +kubebuilder:validation:Required
	ParentType string `json:"parentType"`
	// ParentId is the type specific ID of the parent this project has
	// +kubebuilder:validation:Required
	ParentId string `json:"parentId"`
	// BillingAccountName is the resource name of the billing account associated with the project
	// e.g. '012345-567890-ABCDEF'
	// +kubebuilder:validation:Required
	BillingAccountName string `json:"billingAccountName"`
	// ServiceAccountName is the name used when creating the service account
	// e.g. 'hub-admin'
	// +kubebuilder:validation:Minimum=3
	// +kubebuilder:validation:Required
	ServiceAccountName string `json:"serviceAccountName"`
	// Credentials is a reference to the gcp token object to use
	// +kubebuilder:validation:Required
	Credentials core.Ownership `json:"credentials"`
}

// GCPAdminProjectStatus defines the observed state of GCPAdminProject
// +k8s:openapi-gen=true
type GCPAdminProjectStatus struct {
	// Conditions is a collection of conditions of errors
	// +kubebuilder:validation:Optional
	// +listType
	Conditions []core.Condition `json:"conditions,omitempty"`
	// Status provides a overall status
	Status core.Status `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GCPAdminProject is the Schema for the gcpadminprojects API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=gcpadminprojects,scope=Namespaced
type GCPAdminProject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GCPAdminProjectSpec   `json:"spec,omitempty"`
	Status GCPAdminProjectStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GCPAdminProjectList contains a list of GCPAdminProject
type GCPAdminProjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GCPAdminProject `json:"items"`
}
