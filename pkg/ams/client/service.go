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
	"github.com/anbox-cloud/ams-sdk/pkg/ams/shared/rest/client"
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
func (c *clientImpl) HasExtension(name string) (bool, error) {
	if c.serviceStatus == nil {
		status, _, err := c.RetrieveServiceStatus()
		if err != nil {
			return false, err
		}
		c.serviceStatus = status
	}

	for _, ext := range c.serviceStatus.APIExtensions {
		if ext == name {
			return true, nil
		}
	}

	return false, nil
}
