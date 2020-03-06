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
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"

	"github.com/urfave/cli"
)

var (
	createSecretLongDescription = `
Create secret is used to provision a secret in Kore. The command 
 $ korectl create secret <name> -t <team> [options]

Examples:
 # Create a secret from a file
 $ korectl create secret gke --from-file=<key>=<filename>
`
)

// GetCreateSecretCommand creates and returns the create secret command
func GetCreateSecretCommand(config *Config) cli.Command {
	return cli.Command{
		Name:        "secret",
		Aliases:     []string{"secrets"},
		Description: createSecretLongDescription,
		Usage:       "Creates a secret in kore",
		ArgsUsage:   "<name> [options]",

		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "team,t",
				Usage: "Used to select the team context you are operating in `NAME`",
			},
			cli.StringFlag{
				Name:  "type",
				Usage: "indicates the type of secret you are generating `NAME`",
				Value: "generic",
			},
			cli.StringFlag{
				Name:  "from-literal",
				Usage: "adding a literal to the secret `KEY=NAME`",
			},
			cli.StringFlag{
				Name:  "from-file",
				Usage: "builds a secret from the key reference `KEY=NAME`",
			},
			cli.StringFlag{
				Name:  "from-env-file",
				Usage: "builds a secret from the environment file, format NAME=VALUE `PATH`",
			},
			cli.BoolFlag{
				Name:  "dry-run",
				Usage: "generate the cluster specification but does not apply `BOOL`",
			},
		},

		Before: func(ctx *cli.Context) error {
			if !ctx.Args().Present() {
				return errors.New("the secret should have a name")
			}
			if GlobalStringFlag(ctx, "team") == "" {
				return errors.New("you need to specify a team")
			}

			return nil
		},

		Action: func(ctx *cli.Context) error {
			name := ctx.Args().First()
			team := GlobalStringFlag(ctx, "team")

			var secret *configv1.Secret
			var err error

			switch {
			case ctx.String("from-env-file") != "":
				secret, err = createSecretFromEnvironmentFile(ctx.String("from-env-file"))
			case ctx.String("from-file") != "":

			}
			if err != nil {
				return fmt.Errorf("failed to create secret: %s", err)
			}

			if err := CreateTeamResource(config, team, "secret", name, secret); err != nil {
				return err
			}

			return nil
		},
	}
}

// createSecretFromEnvironmentFile generates a secret from the environment file
func createSecretFromEnvironmentFile(path string) (*configv1.Secret, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	filter := regexp.MustCompile("^.*=.*$")

	secret := configv1.Secret{
		Spec: configv1.SecretSpec{
			Data: make(map[string][]byte),
		},
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		}
		if !filter.MatchString(scanner.Text()) {
			return nil, fmt.Errorf("invalid format: %s, must be name=value", scanner.Text())
		}
		e := strings.Split(scanner.Text(), "=")
		secret.Spec.Data[e[0]] = []byte(e[1])
	}

	return nil, secret
}
