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
	"github.com/urfave/cli"
)

// GetCommands returns all the commands
func GetCommands(config *Config) []cli.Command {
	return []cli.Command{
		GetLoginCommand(config),
		GetLogoutCommand(config),
		GetProfilesCommand(config),
		GetLocalCommand(config),
		GetAutoCompleteCommand(config),
		GetCreateCommand(config),
		GetEditCommand(config),
		GetApplyCommand(config),
		GetDeleteCommand(config),
		GetClustersCommand(config),
		GetGetCommand(config),
		GetWhoamiCommand(config),
	}
}
