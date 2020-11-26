// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2020 Canonical Ltd.  All rights reserved.

package client

import (
	"github.com/anbox-cloud/ams-sdk/api"
	"github.com/anbox-cloud/ams-sdk/shared/rest/client"
)

// ListApplicationsFromRegistry returns a list of all availables applications through the
// registered application registry
func (c *clientImpl) ListApplicationsFromRegistry() ([]api.RegistryApplication, error) {
	apps := []api.RegistryApplication{}
	_, err := c.QueryStruct("GET", client.APIPath("registry", "applications"), nil, nil, nil, "", &apps)
	return apps, err
}

// PushApplicationToRegistry pushes an application to the configured application registry
func (c *clientImpl) PushApplicationToRegistry(id string) (client.Operation, error) {
	op, _, err := c.QueryOperation("POST", client.APIPath("registry", "applications", id, "push"), nil, nil, nil, "")
	return op, err
}

// PullApplicationFromRegistry pulls an application from the configured application registry
func (c *clientImpl) PullApplicationFromRegistry(id string) (client.Operation, error) {
	op, _, err := c.QueryOperation("POST", client.APIPath("registry", "applications", id, "pull"), nil, nil, nil, "")
	return op, err
}

// DeleteApplicationFromRegistry deletes an application from the configured application registry
func (c *clientImpl) DeleteApplicationFromRegistry(id string) (client.Operation, error) {
	op, _, err := c.QueryOperation("DELETE", client.APIPath("registry", "applications", id), nil, nil, nil, "")
	return op, err
}
