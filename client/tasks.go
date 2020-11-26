// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package client

import (
	"github.com/anbox-cloud/ams-sdk/api"
	"github.com/anbox-cloud/ams-sdk/shared/rest/client"
)

// ListTasks returns a list of all availables tasks
func (c *clientImpl) ListTasks() ([]api.Task, error) {
	tasks := []api.Task{}
	params := client.QueryParams{
		"recursion": "1",
	}
	_, err := c.QueryStruct("GET", client.APIPath("tasks"), params, nil, nil, "", &tasks)
	return tasks, err
}
