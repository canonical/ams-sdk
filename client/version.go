// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package client

import (
	"github.com/anbox-cloud/ams-sdk/api"
	"github.com/anbox-cloud/ams-sdk/shared/rest/client"
)

// GetVersion returns the version of the AMS server (not the API version)
func (c *clientImpl) GetVersion() (string, error) {
	v := api.VersionGet{}
	_, err := c.QueryStruct("GET", client.APIPath("version"), nil, nil, nil, "", &v)
	return v.Version, err
}
