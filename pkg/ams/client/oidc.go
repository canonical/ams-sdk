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
	restapi "github.com/anbox-cloud/ams-sdk/pkg/ams/shared/rest/api"
	"github.com/anbox-cloud/ams-sdk/pkg/ams/shared/rest/client"
)

// GetOIDCConfig returns the OIDC configuration of the AMS service
func (c *clientImpl) GetOIDCConfig(grantType string) (*restapi.OIDCResponse, string, error) {
	hasOIDCSupport, err := c.HasExtension("oidc_support")
	if err != nil {
		return nil, "", err
	}
	if !hasOIDCSupport {
		return nil, "", errs.NewErrNotSupported("OIDC Authentication")
	}
	config := &restapi.OIDCResponse{}
	params := client.QueryParams{}
	if grantType == "device_code" {
		params["grant_type"] = grantType
	}
	etag, err := c.QueryStruct("GET", client.APIPath("auth/oidc"), params, nil, nil, "", config)
	return config, etag, err
}

// CreateOIDCIdentity creates a new OIDC based identity in AMS
func (c *clientImpl) CreateOIDCIdentity(details *api.OIDCIdentityPost) (client.Operation, error) {
	hasOIDCSupport, err := c.HasExtension("oidc_support")
	if err != nil {
		return nil, err
	}
	if !hasOIDCSupport {
		return nil, errs.NewErrNotSupported("OIDC Authentication")
	}
	b, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}
	header := http.Header{"Content-Type": []string{"application/json"}}
	op, _, err := c.QueryOperation("POST", client.APIPath("auth/identities"), nil, header, bytes.NewReader(b), "")
	return op, err
}
