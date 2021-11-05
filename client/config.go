// -*- Mode: Go; indent-tabs-mode: t -*-
/*
 * This file is part of AMS SDK
 * Copyright 2021 Canonical Ltd.
 *
 * This program is free software: you can redistribute it and/or modify it under
 * the terms of the GNU Lesser General Public License version 3, as published
 * by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranties of MERCHANTABILITY, SATISFACTORY
 * QUALITY, or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public
 * License for more details.
 *
 * You should have received a copy of the Lesser GNU General Public License along
 * with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package client

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/anbox-cloud/ams-sdk/api"
	"github.com/anbox-cloud/ams-sdk/shared/rest/client"
)

// SetConfigItem sets the specified config item to the given value
func (c *clientImpl) SetConfigItem(name, value string) error {
	req := api.ConfigPost{
		Name:  name,
		Value: value,
	}

	b, err := json.Marshal(req)
	if err != nil {
		return err
	}

	op, _, err := c.QueryOperation("PATCH", client.APIPath("config"), nil, nil, bytes.NewReader(b), "")
	if err != nil {
		return err
	}
	return op.Wait(context.Background())
}

// RetrieveConfigItems returns a list of configuration items available on the AMS service
func (c *clientImpl) RetrieveConfigItems() (map[string]interface{}, error) {
	resp := api.ConfigGet{}
	_, err := c.QueryStruct("GET", client.APIPath("config"), nil, nil, nil, "", &resp)
	return resp.Config, err
}
