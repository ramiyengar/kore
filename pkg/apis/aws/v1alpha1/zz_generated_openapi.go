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
		"github.com/appvia/kore/pkg/apis/aws/v1alpha1.AWSCredentials":       schema_pkg_apis_aws_v1alpha1_AWSCredentials(ref),
		"github.com/appvia/kore/pkg/apis/aws/v1alpha1.AWSCredentialsSpec":   schema_pkg_apis_aws_v1alpha1_AWSCredentialsSpec(ref),
		"github.com/appvia/kore/pkg/apis/aws/v1alpha1.AWSCredentialsStatus": schema_pkg_apis_aws_v1alpha1_AWSCredentialsStatus(ref),
		"github.com/appvia/kore/pkg/apis/aws/v1alpha1.EKSCluster":           schema_pkg_apis_aws_v1alpha1_EKSCluster(ref),
		"github.com/appvia/kore/pkg/apis/aws/v1alpha1.EKSClusterSpec":       schema_pkg_apis_aws_v1alpha1_EKSClusterSpec(ref),
		"github.com/appvia/kore/pkg/apis/aws/v1alpha1.EKSClusterStatus":     schema_pkg_apis_aws_v1alpha1_EKSClusterStatus(ref),
	}
}

func schema_pkg_apis_aws_v1alpha1_AWSCredentials(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "AWSCredential is the Schema for the awscredentials API",
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
							Ref: ref("github.com/appvia/kore/pkg/apis/aws/v1alpha1.AWSCredentialsSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/appvia/kore/pkg/apis/aws/v1alpha1.AWSCredentialsStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/aws/v1alpha1.AWSCredentialsSpec", "github.com/appvia/kore/pkg/apis/aws/v1alpha1.AWSCredentialsStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_aws_v1alpha1_AWSCredentialsSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "AWSCredentialsSpec defines the desired state of AWSCredential",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"accessKeyID": {
						SchemaProps: spec.SchemaProps{
							Description: "AccessKeyID is the AWS access key credentials",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"accountID": {
						SchemaProps: spec.SchemaProps{
							Description: "AccountID is the AWS account these credentials reside",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"secretAccessKey": {
						SchemaProps: spec.SchemaProps{
							Description: "SecretAccessKey is the AWS secret key credentials containing the permissions to provision EKS",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
	}
}

func schema_pkg_apis_aws_v1alpha1_AWSCredentialsStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "AWSCredentialsStatus defines the observed state of AWSCredential",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"verified": {
						SchemaProps: spec.SchemaProps{
							Description: "Verified checks that the credentials are ok and valid",
							Type:        []string{"boolean"},
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
				},
				Required: []string{"verified", "status"},
			},
		},
	}
}

func schema_pkg_apis_aws_v1alpha1_EKSCluster(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "EKSCluster is the Schema for the eksclusters API",
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
							Ref: ref("github.com/appvia/kore/pkg/apis/aws/v1alpha1.EKSClusterSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/appvia/kore/pkg/apis/aws/v1alpha1.EKSClusterStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/aws/v1alpha1.EKSClusterSpec", "github.com/appvia/kore/pkg/apis/aws/v1alpha1.EKSClusterStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_aws_v1alpha1_EKSClusterSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "EKSClusterSpec defines the desired state of EKSCluster",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"credentials": {
						SchemaProps: spec.SchemaProps{
							Description: "Credentials is a reference to an AWSCredentials object to use for authentication",
							Ref:         ref("github.com/appvia/kore/pkg/apis/core/v1.Ownership"),
						},
					},
					"roleARN": {
						SchemaProps: spec.SchemaProps{
							Description: "RoleARN is the role arn which provides permissions to EKS.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"version": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"region": {
						SchemaProps: spec.SchemaProps{
							Description: "Region is the AWS region which the EKS cluster should be provisioned.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"subnetID": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-type": "set",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "SubnetID is a collection of subnet id's which the EKS cluster should be attached to - if not defined we will provision on behalf of",
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
					"securityGroupID": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-type": "set",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "SecurityGroupID is a list of security group IDs which the EKS cluster should be attached to - If not defined we will provision on behalf of",
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
					"vpc": {
						SchemaProps: spec.SchemaProps{
							Description: "VPC is the AWS VPC Id which the EKS cluster should reside. If not defined we will provision on your behalf.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"region"},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/core/v1.Ownership"},
	}
}

func schema_pkg_apis_aws_v1alpha1_EKSClusterStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "EKSClusterStatus defines the observed state of EKSCluster",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"conditions": {
						SchemaProps: spec.SchemaProps{
							Description: "Conditions is the status of the components",
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
					"caCertificate": {
						SchemaProps: spec.SchemaProps{
							Description: "CACertificate is the certificate for this cluster",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"endpoint": {
						SchemaProps: spec.SchemaProps{
							Description: "Endpoint is the endpoint of the cluster",
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
				},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/core/v1.Component"},
	}
}
