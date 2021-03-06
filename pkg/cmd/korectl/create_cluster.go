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

package korectl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	"github.com/appvia/kore/pkg/utils"
	"gopkg.in/yaml.v2"

	"github.com/urfave/cli"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	createClusterLongDescription = `
Provides the ability to provision a kubernetes cluster in the team. The cluster 
itself is provisioned from a predefined plan (a template). You can view the plans
available to you via $ korectl get plans. Once the cluster has been built the 
members of your team can gain access via running $ korectl login. 

Note: you retrieve a list of all the plans available to you via:
$ korectl get plans 
$ korectl get plans <name> -o yaml

Examples:
$ korectl -t <myteam> create cluster dev --plan gke-development --allocation <name>

# Create a cluster and provision some namespaces on there as well
$ korectl -t <myteam> create cluster dev --plan gke-development -a <name> --namespace=app1,app2

# Check the status of the cluster
$ korectl -t <myteam> get cluster dev -o yaml

Once you have created the cluster you can login via 
$ korectl clusters auth -t <myteam>

This will generate your ${HOME}/.kube/config for you with the clusters from team.
`
)

// GetCreateClusterCommand returns the command to create clusters
// @Note: we probably need to move this cluster provisioning off a plan into the API itself
// and offload it from the CLI - but needs discussion first.
func GetCreateClusterCommand(config *Config) cli.Command {
	return cli.Command{
		Name:        "clusters",
		Aliases:     []string{"cluster"},
		Description: createClusterLongDescription,
		Usage:       "Create a kubernetes cluster within the team",
		ArgsUsage:   "<name> [options]",

		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "team,t",
				Usage: "Used to select the team context you are operating in `NAME`",
			},
			cli.StringFlag{
				Name:  "plan,p",
				Usage: "the plan which this cluster will be templated from `NAME`",
			},
			cli.StringFlag{
				Name:  "description",
				Usage: "provides a short description for the cluster `DESCRIPTION`",
			},
			cli.StringFlag{
				Name:  "team-role",
				Usage: "the default role inherited by all members in the team on the cluster `NAME`",
				Value: "viewer",
			},
			cli.StringSliceFlag{
				Name:  "namespace",
				Usage: "you can preprovision a collection namespaces on this cluster as well `NAMES`",
			},
			cli.StringFlag{
				Name:  "allocation,a",
				Usage: "the name of the allocated credentials to use for this cluster `NAME`",
			},
			cli.BoolFlag{
				Name:  "show-time",
				Usage: "shows the time it took to successfully provision a new cluster `BOOL`",
			},
			cli.BoolTFlag{
				Name:  "wait",
				Usage: "indicates we should wait for the cluster to be build (defaults: true) `BOOL`",
			},
			cli.BoolFlag{
				Name:  "dry-run",
				Usage: "generate the cluster specification but does not apply `BOOL`",
			},
		},

		Before: func(ctx *cli.Context) error {
			if !ctx.Args().Present() {
				return fmt.Errorf("the cluster should have a name")
			}

			return nil
		},

		Action: func(ctx *cli.Context) error {
			name := ctx.Args().First()
			plan := ctx.String("plan")
			allocation := ctx.String("allocation")
			namespaces := ctx.StringSlice("namespace")
			team := GlobalStringFlag(ctx, "team")
			role := ctx.String("team-role")
			waitfor := ctx.Bool("wait")
			dry := ctx.Bool("dry-run")

			if team == "" {
				return errTeamParameterMissing
			}

			// @step: check for an allocation
			if allocation == "" {
				return fmt.Errorf("no allocation defined, please use $ korectl get allocations -t %s", team)
			}
			if plan == "" {
				return fmt.Errorf("no plan defined, please use: $ korectl get plans")
			}

			// @step: create the cloud provider
			provider, err := CreateClusterProviderFromPlan(config, team, name, plan, allocation, dry)
			if err != nil {
				return err
			}

			cluster, err := CreateKubernetesClusterFromProvider(config, provider, team, name, role, dry)
			if err != nil {
				return err
			}

			// @step: we need to construct the provider type
			if waitfor {
				now := time.Now()

				err := func() error {
					// @step: lets try and short cut the wait
					cluster, err := GetCluster(config, team, name)
					if err == nil {
						if cluster.Status.Status == corev1.SuccessStatus {
							return nil
						}
					}

					fmt.Printf("Waiting for %q to provision (usually takes around 5 minutes, ctrl-c to background)\n", name)

					// @step: allow for cancellation of the block - and probably wrap this up into a common framework
					sig := make(chan os.Signal, 1)
					signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

					c, cancel := context.WithCancel(context.Background())
					defer cancel()

					go func() {
						<-sig
						cancel()
					}()

					for {
						cluster, err = GetCluster(config, team, name)
						if err == nil {
							switch cluster.Status.Status {
							case corev1.SuccessStatus:
								fmt.Println("Cluster", cluster.Name, "has been successfully provisioned")
								return nil
							case corev1.FailureStatus:
								return fmt.Errorf("failed to provision cluster: %q, please check via $ korectl get clusters -o yaml", name)
							}
						}
						if utils.Sleep(c, 5*time.Second) {
							fmt.Printf("\nProvisioning has been backgrounded, you can check the status via: $ korectl get clusters -t %s\n", team)
							return nil
						}
					}
				}()
				if err != nil {
					return fmt.Errorf("has failed to provision, use: $ korectl get clusters %s -t %s -o yaml to view status", name, team)
				}
				if ctx.Bool("show-time") {
					fmt.Printf("Provisioning took: %s\n", time.Since(now))
				}

			} else {
				fmt.Printf("Cluster provisioning in background: you can check the status via: $ korectl get clusters %s -t %s\n", name, team)
			}

			// @step: create the cluster ownership
			ownership := corev1.Ownership{
				Group:     clustersv1.GroupVersion.Group,
				Version:   clustersv1.GroupVersion.Version,
				Kind:      "Kubernetes",
				Namespace: cluster.Namespace,
				Name:      cluster.Name,
			}

			// @step: do we need to provision any namespaces? - note the split and joining
			// allows for --namespace a,b,c
			var list []string
			for _, x := range namespaces {
				list = append(list, strings.Split(x, ",")...)
			}

			for _, x := range list {
				if err := CreateClusterNamespace(config, ownership, team, x, dry); err != nil {
					return fmt.Errorf("trying to provision namespace claim: %s on cluster: %s", x, err)
				}
			}

			// @step: print a the message
			fmt.Printf("\nYou can retrieve your kubeconfig via: $ korectl clusters auth -t %s\n", team)

			return nil
		},
	}
}

// CreateKubernetesClusterFromProvider is used to provision a k8s cluster from a provider
func CreateKubernetesClusterFromProvider(config *Config, provider *unstructured.Unstructured, team, name, role string, dry bool) (*clustersv1.Kubernetes, error) {
	whoami, err := GetWhoAmI(config)
	if err != nil {
		return nil, err
	}
	kind := "Kubernetes"

	// @step: create the cluster on top of
	object := &clustersv1.Kubernetes{
		TypeMeta: metav1.TypeMeta{
			APIVersion: clustersv1.GroupVersion.String(),
			Kind:       kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: team,
		},
		Spec: clustersv1.KubernetesSpec{
			InheritTeamMembers: true,
			DefaultTeamRole:    role,
			Provider:           utils.GetOwnership(provider, "credentials"),
			ClusterUsers: []clustersv1.ClusterUser{
				{
					Username: whoami.Username,
					Roles:    []string{"cluster-admin"},
				},
			},
		},
	}
	if dry {
		return object, yaml.NewEncoder(os.Stdout).Encode(object)
	}

	found, err := TeamResourceExists(config, team, "clusters", name)
	if err != nil {
		return nil, fmt.Errorf("trying to check if cluster exists: %s", err)
	}
	if found {
		return object, nil
	}

	return object, CreateTeamResource(config, team, "clusters", name, object)
}

// CreateClusterNamespace is called to provision a namespace on the cluster
func CreateClusterNamespace(config *Config, cluster corev1.Ownership, team, name string, dry bool) error {
	resourceName := fmt.Sprintf("%s-%s", cluster.Name, name)
	kind := "namespaceclaims"

	object := &clustersv1.NamespaceClaim{
		TypeMeta: metav1.TypeMeta{
			APIVersion: clustersv1.GroupVersion.String(),
			Kind:       "NamespaceClaim",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: resourceName,
		},
		Spec: clustersv1.NamespaceClaimSpec{
			Name:    name,
			Cluster: cluster,
		},
	}
	if dry {
		return yaml.NewEncoder(os.Stdout).Encode(object)
	}

	found, err := TeamResourceExists(config, team, kind, resourceName)
	if err != nil {
		return err
	}
	if found {
		fmt.Printf("--> Namespace: %s already exists, skipping creation\n", name)

		return nil
	}
	fmt.Printf("--> Attempting to create namespace: %s\n", name)

	return CreateTeamResource(config, team, kind, name, object)
}

// CreateClusterProviderFromPlan is used to provision a cluster in kore
// @TODO need to be revisited once we have autogeneration of resources
func CreateClusterProviderFromPlan(config *Config, team, name, plan, allocation string, dry bool) (*unstructured.Unstructured, error) {
	// @step: we need to check if the plan exists in the kore
	if found, err := ResourceExists(config, "plan", plan); err != nil {
		return nil, fmt.Errorf("trying to retrieve plan from api: %s", err)
	} else if !found {
		return nil, fmt.Errorf("plan %q does not exist, you can view plans via $ korectl get plans", plan)
	}
	template := &configv1.Plan{}
	if err := GetResource(config, "plan", plan, template); err != nil {
		return nil, fmt.Errorf("trying to retrieve plan from api: %s", err)
	}

	// @step: decode the plan values into a map
	kv := make(map[string]interface{})
	if err := json.NewDecoder(bytes.NewReader(template.Spec.Values.Raw)).Decode(&kv); err != nil {
		return nil, fmt.Errorf("trying to decode plan values: %s", err)
	}
	kv["description"] = fmt.Sprintf("%s cluster", plan)

	kind := strings.ToLower(utils.ToPlural(template.Spec.Kind))

	object := &unstructured.Unstructured{}
	object.SetGroupVersionKind(schema.GroupVersionKind{
		Kind: template.Spec.Kind,
		// needs to be change by added by expanding to the plans to apply to a specific resource
		// @TODO in another pull_request
		Group:   gke.GroupVersion.Group,
		Version: gke.GroupVersion.Version,
	})
	object.SetName(name)
	object.SetNamespace(team)
	// @TODO: we need to fix this up later, much like above
	object.SetAPIVersion(gke.GroupVersion.String())

	utils.InjectValuesIntoUnstructured(kv, object)

	// @step: ensure the allocation exists and retrieve it
	if found, err := TeamResourceExists(config, team, "allocation", allocation); err != nil {
		return nil, fmt.Errorf("retrieving the allocation from api: %s", err)
	} else if !found {
		return nil, fmt.Errorf("allocation: %s has not been assigned to team", allocation)
	}
	permit := &configv1.Allocation{}
	if err := GetTeamResource(config, team, "allocation", allocation, permit); err != nil {
		return nil, fmt.Errorf("retrieving the allocation from api: %s", err)
	}

	utils.InjectOwnershipIntoUnstructured("credentials", permit.Spec.Resource, object)

	if dry {
		return object, yaml.NewEncoder(os.Stdout).Encode(object)
	}

	// @step: check the cluster already exists
	if found, err := TeamResourceExists(config, team, kind, name); err != nil {
		return nil, fmt.Errorf("trying to check if cluster exists: %s", err)
	} else if found {
		fmt.Printf("Cluster: %q already exists, skipping the creation\n", name)

		return object, nil
	}

	fmt.Printf("Attempting to create cluster: %q, plan: %s\n", name, plan)

	return object, CreateTeamResource(config, team, kind, name, object)
}
