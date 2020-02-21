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

package gcpprojectclaim

import (
	"context"

	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

// IsProject checks if the project exists
func IsProject(ctx context.Context, client *cloudresourcemanager.Service, name string) (*cloudresourcemanager.Project, bool, error) {
	list, err := client.Projects.List().Context(ctx).Do()
	if err != nil {
		return nil, false, err
	}
	for _, x := range list.Projects {
		if x.Name == name {
			return x, true, nil
		}
	}

	return nil, false, nil
}

func newClient() (*cloudresourcemanager.Service, error) {
	return nil, nil
}
