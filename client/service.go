// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package client

import (
	"github.com/anbox-cloud/ams-sdk/api"
	"github.com/anbox-cloud/ams-sdk/shared/rest/client"
)

// RetrieveServiceStatus returns the status of the AMS service
func (c *clientImpl) RetrieveServiceStatus() (*api.ServiceStatus, string, error) {
	status := &api.ServiceStatus{}
	etag, err := c.QueryStruct("GET", client.APIPath(""), nil, nil, nil, "", status)
	return status, etag, err
}

// HasExtension checks if the AMS service the client is connected to supports
// the given API extension. Returns true if the API extension is supported and
// false otherwise.
func (c *clientImpl) HasExtension(name string) bool {
	if c.serviceStatus == nil {
		status, _, err := c.RetrieveServiceStatus()
		if err != nil {
			return false
		}
		c.serviceStatus = status
	}

	for _, ext := range c.serviceStatus.APIExtensions {
		if ext == name {
			return true
		}
	}

	return false
}
