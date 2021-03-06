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
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/urfave/cli"
)

var deleteLongDescription = `
The object type accepts both singular and plural nouns (e.g. "user" and "users").

Example to delete a user:
  $ korectl delete user joe@example.com

Example to delete multiple resources from a file:
  $ korectl delete --file resources.yaml
`

// GetDeleteCommand creates and returns the delete command
func GetDeleteCommand(config *Config) cli.Command {
	return cli.Command{
		Name:        "delete",
		Aliases:     []string{"rm", "del"},
		Usage:       "Deletes one or more resources",
		Description: formatLongDescription(deleteLongDescription),
		ArgsUsage:   "-f <file> | [TYPE] [NAME]",
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:  "file,f",
				Usage: "The path to the file containing the resources definitions `PATH`",
			},
			cli.StringFlag{
				Name:  "team,t",
				Usage: "Used to filter the results by team `TEAM`",
			},
		},
		Subcommands: []cli.Command{
			// @note once we figure out the global flag issue we will place this back in
			//GetDeleteTeamMemberCommand(config),
		},
		Before: func(ctx *cli.Context) error {
			if !ctx.IsSet("file") && !ctx.Args().Present() {
				return cli.ShowCommandHelp(ctx.Parent(), "delete")
			}
			return nil
		},
		Action: func(ctx *cli.Context) error {
			team := GlobalStringFlag(ctx, "team")

			for _, file := range ctx.StringSlice("file") {
				// @step: read in the content of the file
				content, err := ioutil.ReadFile(file)
				if err != nil {
					return err
				}
				documents, err := ParseDocument(bytes.NewReader(content), team)
				if err != nil {
					return err
				}
				for _, x := range documents {
					gvk := x.Object.GetObjectKind().GroupVersionKind()
					err := NewRequest().
						WithConfig(config).
						WithContext(ctx).
						WithEndpoint(x.Endpoint).
						WithRuntimeObject(x.Object).
						Delete()
					if err != nil {
						fmt.Printf("%s/%s failed with error: %s\n", gvk.Group, x.Endpoint, err)

						return err
					}

					fmt.Printf("%s/%s deleted\n", gvk.Group, x.Endpoint)
				}
			}
			if len(ctx.StringSlice("file")) <= 0 {
				if ctx.NArg() < 2 {
					return errors.New("you need to specify a resource type and name")
				}

				req, _, err := NewRequestForResource(config, ctx)
				if err != nil {
					return err
				}

				exists, err := req.Exists()
				if err != nil {
					return err
				}

				if !exists {
					return fmt.Errorf("%q does not exist", ctx.Args().Get(1))
				}

				if err := req.Delete(); err != nil {
					return err
				}

				fmt.Printf("%q was successfully deleted\n", ctx.Args().Get(1))
			}

			return nil
		},
	}
}
