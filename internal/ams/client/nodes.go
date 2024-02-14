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

	api "github.com/anbox-cloud/ams-sdk/api/ams"
	errs "github.com/anbox-cloud/ams-sdk/internal/ams/shared/errors"
	"github.com/anbox-cloud/ams-sdk/internal/ams/shared/rest/client"
)

// ListNodes returns a list of all availables LXD nodes AMS knows about
func (c *clientImpl) ListNodes() ([]api.Node, error) {
	nodes := []api.Node{}
	params := client.QueryParams{
		"recursion": "1",
	}
	_, err := c.QueryStruct("GET", client.APIPath("nodes"), params, nil, nil, "", &nodes)
	return nodes, err
}

// AddNode adds a new node to AMS
func (c *clientImpl) AddNode(node *api.NodesPost) (client.Operation, error) {
	b, err := json.Marshal(node)
	if err != nil {
		return nil, err
	}
	op, _, err := c.QueryOperation("POST", client.APIPath("nodes"), nil, nil, bytes.NewReader(b), "")
	return op, err
}

// RemoveNode removes a single node
func (c *clientImpl) RemoveNode(name string, force, keepInCluster bool) (client.Operation, error) {
	if len(name) == 0 {
		return nil, errs.NewInvalidArgument("name")
	}

	details := api.NodeDelete{
		Force:         force,
		KeepInCluster: keepInCluster,
	}
	b, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}

	op, _, err := c.QueryOperation("DELETE", client.APIPath("nodes", name), nil, nil, bytes.NewReader(b), "")
	return op, err
}

// RetrieveNodeByName retrieves a node specified by name from AMS
func (c *clientImpl) RetrieveNodeByName(name string) (*api.Node, string, error) {
	if len(name) == 0 {
		return nil, "", errs.NewInvalidArgument("name")
	}
	node := &api.Node{}
	etag, err := c.QueryStruct("GET", client.APIPath("nodes", name), nil, nil, nil, "", node)
	return node, etag, err
}

// UpdateNode updates an existing node
func (c *clientImpl) UpdateNode(name string, details *api.NodePatch) (client.Operation, error) {
	if len(name) == 0 {
		return nil, errs.NewInvalidArgument("name")
	}
	if details == nil {
		return nil, errs.NewInvalidArgument("details")
	}
	b, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}
	op, _, err := c.QueryOperation("PATCH", client.APIPath("nodes", name), nil, nil, bytes.NewReader(b), "")
	return op, err
}
