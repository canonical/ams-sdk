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
	"strconv"

	"github.com/anbox-cloud/ams-sdk/api"
	errs "github.com/anbox-cloud/ams-sdk/shared/errors"
	"github.com/anbox-cloud/ams-sdk/shared/rest/client"
)

// AddAddon adds a new addon and uploads the given addon package to AMS
func (c *clientImpl) AddAddon(name string, packagePath string, sentBytes chan float64) (client.Operation, error) {
	details := api.AddonsPost{Name: name}
	return c.upload("POST", client.APIPath("addons"), nil, packagePath, details, sentBytes)
}

// UpdateAddon updates an existing addon
func (c *clientImpl) UpdateAddon(name, packagePath string, sentBytes chan float64) (client.Operation, error) {
	if len(name) == 0 {
		return nil, errs.NewInvalidArgument("name")
	}
	details := api.AddonPatch{}
	return c.upload("PATCH", client.APIPath("addons", name), nil, packagePath, details, sentBytes)
}

// RetrieveAddon loads an addon from the connected AMS service
func (c *clientImpl) RetrieveAddon(name string) (*api.Addon, string, error) {
	if len(name) == 0 {
		return nil, "", errs.NewInvalidArgument("name")
	}
	var details api.Addon
	etag, err := c.QueryStruct("GET", client.APIPath("addons", name), nil, nil, nil, "", &details)
	return &details, etag, err
}

// DeleteAddon deletes an existing addon
func (c *clientImpl) DeleteAddon(name string) (client.Operation, error) {
	if len(name) == 0 {
		return nil, errs.NewInvalidArgument("name")
	}
	op, _, err := c.QueryOperation("DELETE", client.APIPath("addons", name), nil, nil, nil, "")
	return op, err
}

// DeleteAddonVersion deletes a specific version of the given addon
func (c *clientImpl) DeleteAddonVersion(name string, version int) (client.Operation, error) {
	if len(name) == 0 {
		return nil, errs.NewInvalidArgument("name")
	}
	op, _, err := c.QueryOperation("DELETE", client.APIPath("addons", name, strconv.Itoa(version)), nil, nil, nil, "")
	return op, err
}

// ListAddons lists all currently available addons of the connected AMS service
func (c *clientImpl) ListAddons() ([]api.Addon, error) {
	addons := []api.Addon{}
	params := client.QueryParams{
		"recursion": "1",
	}
	_, err := c.QueryStruct("GET", client.APIPath("addons"), params, nil, nil, "", &addons)
	return addons, err
}
