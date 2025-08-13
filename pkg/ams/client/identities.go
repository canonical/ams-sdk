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
	errs "github.com/anbox-cloud/ams-sdk/pkg/ams/shared/errors"
	"github.com/anbox-cloud/ams-sdk/pkg/ams/shared/rest/client"
)

const identitiesBasePath = "auth/identities"

// ListIdentitiesWithFilters lists all available identities the AMS service currently manages with filters
func (c *clientImpl) ListIdentitiesWithFilters(filters []string) ([]api.Identity, error) {
	hasOIDCSupport, err := c.HasExtension("oidc_support")
	if err != nil {
		return nil, err
	}
	if !hasOIDCSupport {
		return nil, errs.NewErrNotSupported("OIDC Authentication")
	}
	identities := []api.Identity{}
	params, err := convertFiltersToParams(filters)
	if err != nil {
		return nil, err
	}
	params["recursion"] = "1"
	_, err = c.QueryStruct("GET", client.APIPath(identitiesBasePath), params, nil, nil, "", &identities)
	if err != nil {
		return nil, err
	}
	return identities, nil
}

// RetrieveIdentityByID retrieves a single identity by its ID
func (c *clientImpl) RetrieveIdentityByID(id string) (*api.Identity, string, error) {
	if len(id) == 0 {
		return nil, "", errs.NewInvalidArgument("id")
	}

	hasOIDCSupport, err := c.HasExtension("oidc_support")
	if err != nil {
		return nil, "", err
	}
	if !hasOIDCSupport {
		return nil, "", errs.NewErrNotSupported("OIDC Authentication")
	}
	instance := &api.Identity{}
	etag, err := c.QueryStruct("GET", client.APIPath(identitiesBasePath, id), nil, nil, nil, "", instance)
	return instance, etag, err
}

// DeleteIdentity deletes an existing identity
func (c *clientImpl) DeleteIdentity(id string) (client.Operation, error) {
	if len(id) == 0 {
		return nil, errs.NewInvalidArgument("id")
	}
	hasOIDCSupport, err := c.HasExtension("oidc_support")
	if err != nil {
		return nil, err
	}
	if !hasOIDCSupport {
		return nil, errs.NewErrNotSupported("OIDC Authentication")
	}
	op, _, err := c.QueryOperation("DELETE", client.APIPath(identitiesBasePath, id), nil, nil, nil, "")
	return op, err
}
