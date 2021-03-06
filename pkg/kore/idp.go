/*
 * Copyright (C) 2019 Appvia Ltd. <gambol99@gmail.com>
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
	"errors"
	"fmt"
	"time"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// IDPClientMaxRetries specifies the maximum times to try connecting to an IDP provider
	IDPClientMaxRetries int = 50
	// IDPClientBackOff is the time to wait before retrying
	IDPClientBackOff time.Duration = 5 * time.Second
	// DefaultIDP is the name to use for the IDP that is used as "default"
	DefaultIDP = "default"
)

// IDP is the kore api idp interface
type IDP interface {
	// Get returns the idp from the kore
	Get(context.Context, string) (*corev1.IDP, error)
	// Default gets the default identity provider
	Default(context.Context) (*corev1.IDP, error)
	// Exists checks if the idp provider exists
	Exists(context.Context, string) (bool, error)
	// List returns a list of configured idps
	List(context.Context) (*corev1.IDPList, error)
	// ConfigTypes returns the idp config types available
	ConfigTypes(context.Context) []*corev1.IDPConfig
	// Update is responsible for updating / creating a confgured IDP
	Update(context.Context, *corev1.IDP) error
	// UpdateClient is responsible for updating / creating an idp client
	UpdateClient(context.Context, *corev1.IDPClient) error
}

// authImpl provides the implementation for Identity providers
type idpImpl struct {
	Interface
}

// Default returns the default identity configuration
func (a idpImpl) Default(ctx context.Context) (idp *corev1.IDP, err error) {
	return a.Get(ctx, DefaultIDP)
}

// Exists checks if an identity provider exists
func (a idpImpl) Exists(ctx context.Context, name string) (bool, error) {
	if _, err := a.Get(ctx, name); err != nil {
		if err == ErrNotFound {
			return false, nil
		}
		return false, nil
	}
	return true, nil
}

// Get returns the identity configurations in the kore
func (a idpImpl) Get(ctx context.Context, name string) (*corev1.IDP, error) {
	// Is DEX enabled?
	if a.Config().DEX.EnabledDex {
		// Get from DEX...
		idp, err := getDEXConnector(a.Config().DEX, name)
		if err != nil {
			return nil, err
		}
		return idp, nil
	}
	// DEX not enabled - there is only a default configured
	if name == DefaultIDP {
		return a.getDirectIDP()
	}
	return nil, ErrNotFound
}

// List returns a list of identity providers
func (a idpImpl) List(ctx context.Context) (*corev1.IDPList, error) {
	var err error
	var items []corev1.IDP
	if a.Config().DEX.EnabledDex {
		items, err = getDEXConnectors(a.Config().DEX)
	} else {
		// No DEX
		var d *corev1.IDP
		d, err = a.getDirectIDP()
		items = []corev1.IDP{
			*d,
		}
	}
	return &corev1.IDPList{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "IDPList",
		},
		Items: items,
	}, err
}

// ConfigTypes provides sample identity configurations
// should swagger be used by UX to know this?
func (a idpImpl) ConfigTypes(ctx context.Context) []*corev1.IDPConfig {
	if a.Config().DEX.EnabledDex {
		return []*corev1.IDPConfig{
			{
				Google: &corev1.GoogleIDP{},
			},
			{
				SAML: &corev1.SAMLIDP{},
			},
			{
				OIDC: &corev1.OIDCIDP{},
			},
			{
				Github: &corev1.GithubIDP{},
			},
		}
	}
	// DEX disabled:
	return []*corev1.IDPConfig{
		{
			OIDCDirect: &corev1.StaticOIDCIDP{},
		},
	}
}

// Update is responsible for updating / creating a identity providers
func (a idpImpl) Update(ctx context.Context, idp *corev1.IDP) error {
	if !a.Config().DEX.EnabledDex {
		return NewErrNotAllowed("a static IDP is configured so no change possible")
	}
	// Update DEX
	err := updateDEXConector(a.Config().DEX, idp)
	if err != nil {
		return err
	}
	return nil
}

// UpdateClient is responsible for updating / creating a idp clients
func (a idpImpl) UpdateClient(ctx context.Context, c *corev1.IDPClient) error {
	if !a.Config().DEX.EnabledDex {
		return NewErrNotAllowed("a static IDP is configured so no change possible")
	}
	// Update DEX
	if err := updateDEXClient(a.Config().DEX, c); err != nil {
		return err
	}
	return nil
}

// UpdateUser will create / update a static user in IDP
func (a idpImpl) UpdateUser(ctx context.Context, username string, password string) error {
	if !a.Config().DEX.EnabledDex {
		return NewErrNotAllowed("a static IDP is configured so no change possible")
	}
	if len(password) <= 0 {
		return fmt.Errorf("must set a non 0 length password for user %s", username)
	}
	if err := updateDexUser(a.Config().DEX, username, password); err != nil {
		return err
	}
	return nil
}

func (a idpImpl) getDirectIDP() (*corev1.IDP, error) {
	if !a.Config().HasOpenID() {
		return nil, errors.New("Only OIDC providers supported without IDP broker")
	}
	d := corev1.StaticOIDCIDP{}
	d.ClientID = a.Config().ClientID
	d.ClientSecret = a.Config().ClientSecret
	d.DiscoveryURL = a.Config().DiscoveryURL
	d.Issuer = a.Config().DiscoveryURL
	d.ClientScopes = a.Config().ClientScopes
	d.UserClaims = a.Config().UserClaims

	return &corev1.IDP{
		ObjectMeta: metav1.ObjectMeta{
			Name:      DefaultIDP,
			Namespace: HubNamespace,
		},
		Spec: corev1.IDPSpec{
			DisplayName: "Kore configured IDP",
			Config: corev1.IDPConfig{
				OIDCDirect: &d,
			},
		},
	}, nil
}
