// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2019 Canonical Ltd.  All rights reserved.

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
