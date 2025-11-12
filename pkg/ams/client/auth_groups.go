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
	"bytes"
	"encoding/json"
	"net/http"

	api "github.com/anbox-cloud/ams-sdk/api/ams"
	errs "github.com/anbox-cloud/ams-sdk/pkg/ams/shared/errors"
	"github.com/anbox-cloud/ams-sdk/pkg/ams/shared/rest/client"
)

const groupsBasePath = "auth/groups"

func checkGroupSupport(c *clientImpl) error {
	hasGroupSupport, err := c.HasExtension("auth.groups")
	if err != nil {
		return err
	}
	if !hasGroupSupport {
		return errs.NewErrNotSupported("Auth Groups")
	}
	return nil
}

// CreateAuthGroup creates a new authorization group in AMS
func (c *clientImpl) CreateAuthGroup(details *api.AuthGroupPost) (client.Operation, error) {
	if err := checkGroupSupport(c); err != nil {
		return nil, err
	}
	b, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}
	header := http.Header{"Content-Type": []string{"application/json"}}
	op, _, err := c.QueryOperation("POST", client.APIPath(groupsBasePath), nil, header, bytes.NewReader(b), "")
	return op, err
}

// RetrieveAuthGroup retrieves a single group by its name
func (c *clientImpl) RetrieveAuthGroup(name string) (*api.AuthGroup, string, error) {
	if len(name) == 0 {
		return nil, "", errs.NewInvalidArgument("name")
	}

	if err := checkGroupSupport(c); err != nil {
		return nil, "", err
	}
	group := &api.AuthGroup{}
	etag, err := c.QueryStruct("GET", client.APIPath(groupsBasePath, name), nil, nil, nil, "", group)
	return group, etag, err
}

// ListAuthGroupsWithFilters lists all available groups the AMS service currently manages with filters
func (c *clientImpl) ListAuthGroupsWithFilters(filters []string) ([]api.AuthGroup, error) {
	if err := checkGroupSupport(c); err != nil {
		return nil, err
	}
	groups := []api.AuthGroup{}
	params, err := convertFiltersToParams(filters)
	if err != nil {
		return nil, err
	}
	params["recursion"] = "1"
	_, err = c.QueryStruct("GET", client.APIPath(groupsBasePath), params, nil, nil, "", &groups)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

// DeleteAuthGroup deletes an existing group
func (c *clientImpl) DeleteAuthGroup(name string, force bool) (client.Operation, error) {
	if len(name) == 0 {
		return nil, errs.NewInvalidArgument("name")
	}
	if err := checkGroupSupport(c); err != nil {
		return nil, err
	}
	details := api.AuthGroupDelete{Force: force}
	b, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}
	op, _, err := c.QueryOperation("DELETE", client.APIPath(groupsBasePath, name), nil, nil, bytes.NewReader(b), "")
	return op, err
}

// UpdateAuthGroupDescription updates an existing auth group's description
func (c *clientImpl) UpdateAuthGroupDescription(name, description string) (client.Operation, error) {
	details := &api.AuthGroupPut{
		Description: &description,
	}
	return c.updateAuthGroup(name, details)
}

func (c *clientImpl) updateAuthGroup(name string, details *api.AuthGroupPut) (client.Operation, error) {
	if len(name) == 0 {
		return nil, errs.NewInvalidArgument("name")
	}
	if err := checkGroupSupport(c); err != nil {
		return nil, err
	}
	b, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}
	op, _, err := c.QueryOperation("PUT", client.APIPath(groupsBasePath, name), nil, nil, bytes.NewReader(b), "")
	return op, err
}

// SetPermissionsForGroup sets the permissions for an auth group
func (c *clientImpl) SetPermissionsForGroup(name string, permissions []api.Permission) (client.Operation, error) {
	details := &api.AuthGroupPut{
		Permissions: &permissions,
	}
	return c.updateAuthGroup(name, details)
}
