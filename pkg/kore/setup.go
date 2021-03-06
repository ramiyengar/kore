/*
 * Copyright (C) 2019 Appvia Ltd <info@appvia.io>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package kore

import (
	"context"
	"fmt"
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	core "github.com/appvia/kore/pkg/apis/core/v1"
	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/kore/assets"
	"github.com/appvia/kore/pkg/services/users/model"
	"github.com/appvia/kore/pkg/store"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Setup is called one on initialization and used to provision and empty kore
func (h hubImpl) Setup(ctx context.Context) error {
	log.Info("initializing the kore")

	// @step: ensure the kore namespaces are there
	for _, x := range []string{HubNamespace, HubAdminTeam, HubDefaultTeam} {
		if err := h.ensureNamespace(ctx, x); err != nil {
			return err
		}
	}

	// @step: ensure the default user is there
	for _, x := range []string{HubAdminUser} {
		if err := h.ensureHubAdminUser(ctx, x, "admin@local"); err != nil {
			return err
		}
	}

	// @step: ensure the kore admin team exists
	for _, x := range []string{HubAdminTeam, HubDefaultTeam} {
		if err := h.ensureHubTeam(ctx, x, "Team for "+x); err != nil {
			return err
		}
	}

	// @step: ensure the kore admin user
	for _, x := range []string{HubAdminUser} {
		if err := h.ensureHubAdminMembership(ctx, x, HubAdminTeam); err != nil {
			return err
		}
	}

	// @step: ensure an OIDC client is created in IDP broker
	if h.Config().DEX.EnabledDex {
		if err := h.ensureHubIDPClientExists(ctx); err != nil {
			return err
		}
	}

	// @step: ensure some default plans
	for _, x := range assets.GetDefaultPlans() {
		if err := h.Plans().Update(getAdminContext(ctx), x); err != nil {
			return err
		}
	}
	for _, x := range assets.GetDefaultClusterRoles() {
		x.Namespace = HubAdminTeam

		found, err := h.Store().Client().Has(ctx,
			store.HasOptions.From(&clustersv1.ManagedClusterRole{}),
			store.HasOptions.InNamespace(HubAdminTeam),
			store.HasOptions.WithName(x.Name),
		)
		if err != nil {
			return err
		}
		if !found {
			if err := h.Store().Client().Create(ctx, store.CreateOptions.From(&x)); err != nil {
				return err
			}
		}
	}

	return nil
}

// ensureHubAdminMembership ensures the user is there
func (h hubImpl) ensureHubAdminMembership(ctx context.Context, name, team string) error {
	return h.usermgr.Teams().AddUser(ctx, name, team, []string{"member", "admin"})
}

// ensureHubAdminUser ensures the user is there
func (h hubImpl) ensureHubAdminUser(ctx context.Context, name, email string) error {
	logger := log.WithFields(log.Fields{
		"username": name,
	})

	found, err := h.Users().Exists(ctx, name)
	if err != nil {
		return err
	}

	if !found {
		logger.Info("provisioning the default kore team in api")

		err := h.usermgr.Users().Update(ctx, &model.User{Email: email, Username: name})
		if err != nil {
			logger.WithError(err).Error("trying to create admin user")

			return err
		}
	}
	// Add or update user to IDP broker:
	if h.Config().DEX.EnabledDex {
		if err = h.idp.UpdateUser(ctx, name, h.Config().AdminPass); err != nil {
			logger.WithError(err).Error("trying to update idp password")

			return err
		}
	}

	if h.Config().AdminPass != "" {
		user, err := h.usermgr.Users().Get(ctx, name)
		if err != nil {
			return err
		}

		return h.usermgr.Identities().Update(ctx, &model.Identity{
			Provider:      "basicauth",
			ProviderEmail: email,
			ProviderToken: h.Config().AdminPass,
			UserID:        user.ID,
		})
	}

	return nil
}

// ensureHubTeam ensure a kore team exists in the kore
func (h hubImpl) ensureHubTeam(ctx context.Context, name, description string) error {
	nc := getAdminContext(ctx)

	log.WithFields(log.Fields{
		"team": name,
	}).Info("provisioning the default kore team in api")

	_, err := h.Teams().Update(nc, &orgv1.Team{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: HubNamespace,
		},
		Spec: orgv1.TeamSpec{
			Description: description,
			Summary:     description,
		},
	})

	return err
}

// ensureNamespace ensures the namespace exists in the kore
func (h hubImpl) ensureNamespace(ctx context.Context, namespace string) error {
	found, err := h.Store().Client().Has(ctx,
		store.HasOptions.From(&corev1.Namespace{}),
		store.HasOptions.InNamespace(HubNamespace),
		store.HasOptions.WithName(namespace),
	)
	if err != nil || found {
		return err
	}

	log.WithFields(log.Fields{
		"namespace": namespace,
	}).Info("provisioning the namespace in api")

	// @step: we need to create it
	return h.Store().Client().Create(ctx,
		store.CreateOptions.From(&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}),
	)
}

func (h hubImpl) ensureHubIDPClientExists(ctx context.Context) error {
	for i := 0; i < IDPClientMaxRetries; i++ {
		// Ensure there is a client created in DEX for the API server
		err := h.idp.UpdateClient(ctx, &core.IDPClient{
			Spec: core.IDPClientSpec{
				DisplayName: "The API server OIDC client",
				ID:          h.Config().ClientID,
				Secret:      h.Config().ClientSecret,
				RedirectURIs: []string{
					h.Config().PublicAPIURL + "/oauth/callback",
				},
			},
		})
		if err != nil {
			if err == ErrServerNotAvailable {
				// loop for now
				time.Sleep(IDPClientBackOff)
				log.Warn("IDP broker not available so waiting")
				continue
			} else {
				return fmt.Errorf("error configuring IDP client for IDP broker")
			}
		}
	}
	log.Info("API server OIDC client configured in IDP broker")

	return nil
}
