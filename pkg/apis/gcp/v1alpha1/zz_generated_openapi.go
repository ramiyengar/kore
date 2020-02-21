// +build !ignore_autogenerated

/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */
// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/appvia/kore/pkg/apis/gcp/v1alpha1.GCPAdminProject":       schema_pkg_apis_gcp_v1alpha1_GCPAdminProject(ref),
		"github.com/appvia/kore/pkg/apis/gcp/v1alpha1.GCPAdminProjectSpec":   schema_pkg_apis_gcp_v1alpha1_GCPAdminProjectSpec(ref),
		"github.com/appvia/kore/pkg/apis/gcp/v1alpha1.GCPAdminProjectStatus": schema_pkg_apis_gcp_v1alpha1_GCPAdminProjectStatus(ref),
		"github.com/appvia/kore/pkg/apis/gcp/v1alpha1.GCPProjectClaim":       schema_pkg_apis_gcp_v1alpha1_GCPProjectClaim(ref),
		"github.com/appvia/kore/pkg/apis/gcp/v1alpha1.GCPProjectClaimSpec":   schema_pkg_apis_gcp_v1alpha1_GCPProjectClaimSpec(ref),
		"github.com/appvia/kore/pkg/apis/gcp/v1alpha1.GCPProjectClaimStatus": schema_pkg_apis_gcp_v1alpha1_GCPProjectClaimStatus(ref),
	}
}

func schema_pkg_apis_gcp_v1alpha1_GCPAdminProject(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "GCPAdminProject is the Schema for the gcpadminprojects API",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/appvia/kore/pkg/apis/gcp/v1alpha1.GCPAdminProjectSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/appvia/kore/pkg/apis/gcp/v1alpha1.GCPAdminProjectStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/gcp/v1alpha1.GCPAdminProjectSpec", "github.com/appvia/kore/pkg/apis/gcp/v1alpha1.GCPAdminProjectStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_gcp_v1alpha1_GCPAdminProjectSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "GCPAdminProjectSpec defines the desired state of GCPAdminProject",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"project": {
						SchemaProps: spec.SchemaProps{
							Description: "Project is the GCP project ID",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"projectName": {
						SchemaProps: spec.SchemaProps{
							Description: "ProjectName is the GCP project name",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"parentType": {
						SchemaProps: spec.SchemaProps{
							Description: "ParentType is the type of parent this project has Valid types are: \"organization\", \"folder\", and \"project\"",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"parentID": {
						SchemaProps: spec.SchemaProps{
							Description: "ParentID is the type specific ID of the parent this project has",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"billingAccountName": {
						SchemaProps: spec.SchemaProps{
							Description: "BillingAccountName is the resource name of the billing account associated with the project e.g. '012345-567890-ABCDEF'",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"serviceAccountName": {
						SchemaProps: spec.SchemaProps{
							Description: "ServiceAccountName is the name used when creating the service account e.g. 'hub-admin'",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"credentials": {
						SchemaProps: spec.SchemaProps{
							Description: "Credentials is a reference to the gcp token object to use",
							Ref:         ref("github.com/appvia/kore/pkg/apis/core/v1.Ownership"),
						},
					},
				},
				Required: []string{"project", "projectName", "parentType", "parentID", "billingAccountName", "serviceAccountName", "credentials"},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/core/v1.Ownership"},
	}
}

func schema_pkg_apis_gcp_v1alpha1_GCPAdminProjectStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "GCPAdminProjectStatus defines the observed state of GCPAdminProject",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"conditions": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-type": "",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "Conditions is a collection of conditions of errors",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("github.com/appvia/kore/pkg/apis/core/v1.Condition"),
									},
								},
							},
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Description: "Status provides a overall status",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/core/v1.Condition"},
	}
}

func schema_pkg_apis_gcp_v1alpha1_GCPProjectClaim(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "GCPProjectClaim is the Schema for the gcpprojects API",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/appvia/kore/pkg/apis/gcp/v1alpha1.GCPProjectClaimSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/appvia/kore/pkg/apis/gcp/v1alpha1.GCPProjectClaimStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/gcp/v1alpha1.GCPProjectClaimSpec", "github.com/appvia/kore/pkg/apis/gcp/v1alpha1.GCPProjectClaimStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_gcp_v1alpha1_GCPProjectClaimSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "GCPProjectClaimSpec defines the desired state of GCPProjectClaim",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"serviceAccountName": {
						SchemaProps: spec.SchemaProps{
							Description: "ServiceAccountName is an optional name of the service account provision - else we default to the name kore",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"organization": {
						SchemaProps: spec.SchemaProps{
							Description: "Organization is a reference to the gcp admin project to use",
							Ref:         ref("github.com/appvia/kore/pkg/apis/core/v1.Ownership"),
						},
					},
				},
				Required: []string{"organization"},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/core/v1.Ownership"},
	}
}

func schema_pkg_apis_gcp_v1alpha1_GCPProjectClaimStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "GCPProjectClaimStatus defines the observed state of GCPProject",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"projectID": {
						SchemaProps: spec.SchemaProps{
							Description: "ProjectID is the GCP project ID",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Description: "Status provides a overall status",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"conditions": {
						SchemaProps: spec.SchemaProps{
							Description: "Conditions is a set of components conditions",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("github.com/appvia/kore/pkg/apis/core/v1.Component"),
									},
								},
							},
						},
					},
				},
				Required: []string{"status"},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/core/v1.Component"},
	}
}
