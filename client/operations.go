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
	"github.com/anbox-cloud/ams-sdk/shared/rest/api"
	"github.com/anbox-cloud/ams-sdk/shared/rest/client"
)

// ListOperations lists all operations arranged by their status
func (c *clientImpl) ListOperations() (map[string][]*api.Operation, error) {
	params := client.QueryParams{
		"recursion": "1",
	}
	var operations map[string][]*api.Operation
	_, err := c.QueryStruct("GET", client.APIPath("operations"), params, nil, nil, "", &operations)
	return operations, err
}

// ShowOperation shows details about a single operation
func (c *clientImpl) ShowOperation(id string) (*api.Operation, error) {
	var operation *api.Operation
	_, err := c.QueryStruct("GET", client.APIPath("operations", id), nil, nil, nil, "", &operation)
	return operation, err
}

// CancelOperation cancels an operation if it supports it
func (c *clientImpl) CancelOperation(id string) error {
	_, _, err := c.CallAPI("DELETE", client.APIPath("operations", id), nil, nil, nil, "")
	return err
}
