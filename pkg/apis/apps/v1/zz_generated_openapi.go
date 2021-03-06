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

package v1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/appvia/kore/pkg/apis/apps/v1.AppDeployment":       schema_pkg_apis_apps_v1_AppDeployment(ref),
		"github.com/appvia/kore/pkg/apis/apps/v1.AppDeploymentSpec":   schema_pkg_apis_apps_v1_AppDeploymentSpec(ref),
		"github.com/appvia/kore/pkg/apis/apps/v1.AppDeploymentStatus": schema_pkg_apis_apps_v1_AppDeploymentStatus(ref),
		"github.com/appvia/kore/pkg/apis/apps/v1.InstallPlan":         schema_pkg_apis_apps_v1_InstallPlan(ref),
		"github.com/appvia/kore/pkg/apis/apps/v1.InstallPlanSpec":     schema_pkg_apis_apps_v1_InstallPlanSpec(ref),
		"github.com/appvia/kore/pkg/apis/apps/v1.InstallPlanStatus":   schema_pkg_apis_apps_v1_InstallPlanStatus(ref),
	}
}

func schema_pkg_apis_apps_v1_AppDeployment(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "AppDeployment is the Schema for the allocations API",
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
							Ref: ref("github.com/appvia/kore/pkg/apis/apps/v1.AppDeploymentSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/appvia/kore/pkg/apis/apps/v1.AppDeploymentStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/apps/v1.AppDeploymentSpec", "github.com/appvia/kore/pkg/apis/apps/v1.AppDeploymentStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_apps_v1_AppDeploymentSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "AppDeploymentSpec defines the desired state of Allocation",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"cluster": {
						SchemaProps: spec.SchemaProps{
							Description: "Cluster is the cluster the application should be deployed on",
							Ref:         ref("github.com/appvia/kore/pkg/apis/core/v1.Ownership"),
						},
					},
					"summary": {
						SchemaProps: spec.SchemaProps{
							Description: "Summary is a summary of what the application is",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"description": {
						SchemaProps: spec.SchemaProps{
							Description: "Decription is a longer description of what the application provides",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"package": {
						SchemaProps: spec.SchemaProps{
							Description: "Package is the name of the resource being shared",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"version": {
						SchemaProps: spec.SchemaProps{
							Description: "Version is the version of the package to install",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"source": {
						SchemaProps: spec.SchemaProps{
							Description: "Source is the source of the package",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"capabilities": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-type": "set",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "Capabilities defines the features supported by the package",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"keywords": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-type": "set",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "Keywords keywords whuch describe the application",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"vendor": {
						SchemaProps: spec.SchemaProps{
							Description: "Vendor is the entity whom published the package",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"official": {
						SchemaProps: spec.SchemaProps{
							Description: "Official indicates if the applcation is officially published by Appvia",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"replaces": {
						SchemaProps: spec.SchemaProps{
							Description: "Replaces indicates the version this replaces",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"subscription": {
						SchemaProps: spec.SchemaProps{
							Description: "Subscription is the nature of upgrades i.e manual or automatic",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"values": {
						SchemaProps: spec.SchemaProps{
							Description: "Values are optional values suppilied to the application deployment",
							Ref:         ref("k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1.JSON"),
						},
					},
				},
				Required: []string{"cluster", "summary", "description", "package", "version", "source", "keywords", "vendor", "official", "replaces", "subscription"},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/core/v1.Ownership", "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1.JSON"},
	}
}

func schema_pkg_apis_apps_v1_AppDeploymentStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "AppDeploymentStatus defines the observed state of Allocation",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"status": {
						SchemaProps: spec.SchemaProps{
							Description: "Status is the general status of the resource",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"conditions": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-type": "set",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "Conditions is a collection of potential issues",
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
					"installPlan": {
						SchemaProps: spec.SchemaProps{
							Description: "InstallPlan in the name of the installplan which this deployment has deployed from",
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

func schema_pkg_apis_apps_v1_InstallPlan(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "InstallPlan is the Schema for the allocations API",
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
							Ref: ref("github.com/appvia/kore/pkg/apis/apps/v1.InstallPlanSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/appvia/kore/pkg/apis/apps/v1.InstallPlanStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/apps/v1.InstallPlanSpec", "github.com/appvia/kore/pkg/apis/apps/v1.InstallPlanStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_apps_v1_InstallPlanSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "InstallPlanSpec defines the desired state of Allocation",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"approved": {
						SchemaProps: spec.SchemaProps{
							Description: "Approved indicates if the update has been approved",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
				},
			},
		},
	}
}

func schema_pkg_apis_apps_v1_InstallPlanStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "InstallPlanStatus defines the observed state of Allocation",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"conditions": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-type": "set",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "Conditions is a collection of potential issues",
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
					"deployed": {
						SchemaProps: spec.SchemaProps{
							Description: "Deployed is the applciation deployment parameters",
							Ref:         ref("github.com/appvia/kore/pkg/apis/apps/v1.AppDeployment"),
						},
					},
					"update": {
						SchemaProps: spec.SchemaProps{
							Description: "Update is the incoming deployment is requiring approval",
							Ref:         ref("github.com/appvia/kore/pkg/apis/apps/v1.AppDeployment"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Description: "Status is the general status of the resource",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"deployed"},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/apps/v1.AppDeployment", "github.com/appvia/kore/pkg/apis/core/v1.Condition"},
	}
}
