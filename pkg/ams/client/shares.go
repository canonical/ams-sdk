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
	"fmt"

	api "github.com/anbox-cloud/ams-sdk/api/ams"
	errs "github.com/anbox-cloud/ams-sdk/pkg/ams/shared/errors"
	restclient "github.com/anbox-cloud/ams-sdk/pkg/ams/shared/rest/client"
)

const (
	shareSupportExtension = "instance_share_support"
)

var (
	sharesSupported       *bool = nil
	errSharesNotSupported       = fmt.Errorf("instance shares are not supported by the connected AMS service")
)

func (c *clientImpl) hasSharesSupport() (bool, error) {
	if sharesSupported != nil {
		return *sharesSupported, nil
	}

	supported, err := c.HasExtension(shareSupportExtension)
	if err != nil {
		return false, err
	}

	sharesSupported = &supported
	return supported, nil
}

// CreateShare creates a share for an instance
func (c *clientImpl) CreateInstanceShare(id string, details *api.InstanceSharesPost) (*api.InstanceSharesPostResponse, error) {
	b, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}

	supported, err := c.hasSharesSupport()
	if err != nil {
		return nil, err
	}
	if !supported {
		return nil, errSharesNotSupported
	}

	var resp api.InstanceSharesPostResponse
	_, err = c.QueryStruct("POST", restclient.APIPath("instances", id, "shares"), nil, nil, bytes.NewReader(b), "", &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// UpdateInstanceShareByID updates a share
func (c *clientImpl) UpdateInstanceShareByID(instanceID, shareID string, details *api.InstanceSharePatch) error {
	if details == nil || len(instanceID) == 0 || len(shareID) == 0 {
		return errs.NewInvalidArgument("details")
	}

	supported, err := c.hasSharesSupport()
	if err != nil {
		return err
	}
	if !supported {
		return errSharesNotSupported
	}

	b, err := json.Marshal(details)
	if err != nil {
		return err
	}

	_, err = c.QueryStruct("PATCH", restclient.APIPath("instances", instanceID, "shares", shareID),
		nil, nil, bytes.NewReader(b), "", nil)
	return err
}

// DeleteInstanceShareByID deletes a share by its ID
func (c *clientImpl) DeleteInstanceShareByID(instanceID, shareID string) (restclient.Operation, error) {
	if len(instanceID) == 0 || len(shareID) == 0 {
		return nil, errs.NewInvalidArgument("id")
	}

	supported, err := c.hasSharesSupport()
	if err != nil {
		return nil, err
	}
	if !supported {
		return nil, errSharesNotSupported
	}

	op, _, err := c.QueryOperation("DELETE", restclient.APIPath("instances", instanceID, "shares", shareID),
		nil, nil, nil, "")
	return op, err
}
