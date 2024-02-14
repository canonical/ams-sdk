// -*- Mode: Go; indent-tabs-mode: t -*-
/*
 * This file is part of AMS SDK
 * Copyright 2021 Canonical Ltd.
 *
 * This program is free software: you can redistribute it and/or modify it under
 * the terms of the Lesser GNU General Public License version 3, as published
 * by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranties of MERCHANTABILITY, SATISFACTORY
 * QUALITY, or FITNESS FOR A PARTICULAR PURPOSE.  See the Lesser GNU General Public
 * License for more details.
 *
 * You should have received a copy of the Lesser GNU General Public License along
 * with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package client

import (
	api "github.com/anbox-cloud/ams-sdk/api/ams"
	"github.com/anbox-cloud/ams-sdk/internal/ams/shared/rest/client"
)

// ListApplicationsFromRegistry returns a list of all availables applications through the
// registered application registry
func (c *clientImpl) ListApplicationsFromRegistry() ([]api.RegistryApplication, error) {
	apps := []api.RegistryApplication{}
	_, err := c.QueryStruct("GET", client.APIPath("registry", "applications"), nil, nil, nil, "", &apps)
	return apps, err
}

// PushApplicationToRegistry pushes an application to the configured application registry
func (c *clientImpl) PushApplicationToRegistry(id string) (client.Operation, error) {
	op, _, err := c.QueryOperation("POST", client.APIPath("registry", "applications", id, "push"), nil, nil, nil, "")
	return op, err
}

// PullApplicationFromRegistry pulls an application from the configured application registry
func (c *clientImpl) PullApplicationFromRegistry(id string) (client.Operation, error) {
	op, _, err := c.QueryOperation("POST", client.APIPath("registry", "applications", id, "pull"), nil, nil, nil, "")
	return op, err
}

// DeleteApplicationFromRegistry deletes an application from the configured application registry
func (c *clientImpl) DeleteApplicationFromRegistry(id string) (client.Operation, error) {
	op, _, err := c.QueryOperation("DELETE", client.APIPath("registry", "applications", id), nil, nil, nil, "")
	return op, err
}
