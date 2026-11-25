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
	errs "github.com/anbox-cloud/ams-sdk/pkg/ams/shared/errors"
	"github.com/anbox-cloud/ams-sdk/pkg/ams/shared/rest/client"
)

const projectsBasePath = "projects"

func checkProjectSupport(c *clientImpl) error {
	hasProjectSupport, err := c.HasExtension("project_support")
	if err != nil {
		return err
	}
	if !hasProjectSupport {
		return errs.NewErrNotSupported("Projects")
	}
	return nil
}

// CreateProject creates a new project in AMS
func (c *clientImpl) CreateProject(details *api.ProjectsPost) (client.Operation, error) {
	if err := checkProjectSupport(c); err != nil {
		return nil, err
	}
	b, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}
	op, _, err := c.QueryOperation("POST", client.APIPath(projectsBasePath), nil, nil, bytes.NewReader(b), "")
	return op, err
}

// ListProjectsWithFilters lists all projects AMS currently manages, optionally
// filtered by the given filter expressions.
func (c *clientImpl) ListProjectsWithFilters(filters []string) ([]api.Project, error) {
	if err := checkProjectSupport(c); err != nil {
		return nil, err
	}
	projects := []api.Project{}
	params, err := convertFiltersToParams(filters)
	if err != nil {
		return nil, err
	}
	params["recursion"] = "1"
	_, err = c.QueryStruct("GET", client.APIPath(projectsBasePath), params, nil, nil, "", &projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

// RetrieveProjectByName retrieves a single project specified by name from AMS
func (c *clientImpl) RetrieveProjectByName(name string) (*api.Project, string, error) {
	if len(name) == 0 {
		return nil, "", errs.NewInvalidArgument("name")
	}
	if err := checkProjectSupport(c); err != nil {
		return nil, "", err
	}
	project := &api.Project{}
	etag, err := c.QueryStruct("GET", client.APIPath(projectsBasePath, name), nil, nil, nil, "", project)
	return project, etag, err
}

// UpdateProject updates an existing project. If etag is non-empty it is sent as
// an If-Match header so the update fails if the project was modified
// concurrently since it was last read.
func (c *clientImpl) UpdateProject(name string, details *api.ProjectsPatch, etag string) (client.Operation, error) {
	if len(name) == 0 {
		return nil, errs.NewInvalidArgument("name")
	}
	if details == nil {
		return nil, errs.NewInvalidArgument("details")
	}
	if err := checkProjectSupport(c); err != nil {
		return nil, err
	}
	b, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}
	op, _, err := c.QueryOperation("PATCH", client.APIPath(projectsBasePath, name), nil, nil, bytes.NewReader(b), etag)
	return op, err
}

// DeleteProject deletes an existing project. If force is set, any instances
// contained in the project are deleted first.
func (c *clientImpl) DeleteProject(name string, force bool) (client.Operation, error) {
	if len(name) == 0 {
		return nil, errs.NewInvalidArgument("name")
	}
	if err := checkProjectSupport(c); err != nil {
		return nil, err
	}
	details := api.ProjectDelete{Force: force}
	b, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}
	op, _, err := c.QueryOperation("DELETE", client.APIPath(projectsBasePath, name), nil, nil, bytes.NewReader(b), "")
	return op, err
}
