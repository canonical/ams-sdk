// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

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
