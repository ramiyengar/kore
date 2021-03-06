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
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/utils"

	jwt "github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

// Invitations is the contract for the invitation
type Invitations interface {
	// HandleGenerateLink is responsible for handling the link
	HandleGenerateLink(context.Context, string) error
	// VerifyGenerateLink checks the token is valid
	VerifyGenerateLink(context.Context, string) (bool, error)
}

// ivImpl implements the invitation contract
type ivImpl struct {
	Interface
}

// HandleGenerateLink is responsible for handling the link
func (t ivImpl) HandleGenerateLink(ctx context.Context, encoded string) error {
	// @step: extract the user from the context
	u, found := authentication.GetIdentity(ctx)
	if !found {
		log.Warn("no user found in request")
	}

	// @step: we parse the token and grab the claims
	c, err := t.ParseToken(ctx, encoded)
	if err != nil {
		return err
	}
	claims := utils.NewClaims(c)

	// @step: extract claims from the token
	team, found := claims.GetString("team")
	if !found {
		return ErrNotAllowed{message: "no team found in the invitation claim"}
	}

	// @step: check if the claims content a user
	user, found := claims.GetString("user")
	if found && u != nil {
		// @check if user found in context and user found in the token - we need
		// to ensure they are the same
		if user != u.Username() {
			return ErrNotAllowed{message: "invitition link is for another user"}
		}
	} else if !found && u == nil {
		// @check no user in the token and no user in the context
		return ErrNotAllowed{message: "no user found in the request context"}
	} else if !found {
		// @check no user in the token but we have a context
		user = u.Username()
	}

	logger := log.WithFields(log.Fields{
		"user": user,
		"team": team,
	})
	logger.Info("handling invitation link for user")

	// @step: create the user if the dont exists
	/*
		if _, err := t.Users().CreateIfNotExists(ctx, user); err != nil {
			logger.WithError(err).Error("failed to create user if they didn't exist")

			return err
		}
	*/

	return t.Teams().Team(team).Members().Add(ctx, user)
}

// VerifyGenerateLink is responsible for checking the link is valid
func (t ivImpl) VerifyGenerateLink(ctx context.Context, token string) (bool, error) {
	_, err := t.ParseToken(ctx, token)

	return err != nil, err
}

// ParseToken is responsible for extract the claims from the token
func (t ivImpl) ParseToken(ctx context.Context, encoded string) (jwt.MapClaims, error) {
	var claims jwt.MapClaims

	// @step: ensure we have a hmac
	err := func() error {
		if !t.Config().HasHMAC() {
			return errors.New("kore has no hmac token configured")
		}

		// @step: we need to base54 decode the token
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return err
		}

		// @step: parse and extract the token from the payload
		token, err := jwt.Parse(string(decoded), func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(t.Config().HMAC), nil
		})
		if err != nil {
			return err
		}

		// @step: cast extract and check it's valid
		c, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return fmt.Errorf("expected: jwt.MapClaims but got: %T", token.Claims)
		}
		if err := c.Valid(); err != nil {
			return ErrNotAllowed{message: "generated link is invalid: " + err.Error()}
		}
		claims = c

		return nil
	}()

	return claims, err
}
