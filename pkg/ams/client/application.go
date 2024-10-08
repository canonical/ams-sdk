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
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	api "github.com/anbox-cloud/ams-sdk/api/ams"
	errs "github.com/anbox-cloud/ams-sdk/pkg/ams/shared/errors"
	"github.com/anbox-cloud/ams-sdk/pkg/ams/shared/rest/client"
)

// ApplicationCreateArgs provides details on how to create a new application
type ApplicationCreateArgs struct {
	PackagePath   string
	VM            bool
	SentBytesChan chan float64
}

// CreateApplication creates a new application
func (c *clientImpl) CreateApplication(packagePath string, sentBytes chan float64) (client.Operation, error) {
	return c.CreateApplicationWithArgs(&ApplicationCreateArgs{
		PackagePath:   packagePath,
		SentBytesChan: sentBytes,
	})
}

// CreateApplicationWithArgs creates a new application based on the provided arguments
func (c *clientImpl) CreateApplicationWithArgs(args *ApplicationCreateArgs) (client.Operation, error) {
	hasVMSupport, err := c.HasExtension("vm_support")
	if err != nil {
		return nil, err
	}
	if args.VM && !hasVMSupport {
		return nil, errs.NewErrNotSupported("VM")
	}
	params := client.QueryParams{
		"vm": strconv.FormatBool(args.VM),
	}
	return c.upload("POST", client.APIPath("applications"), params,
		args.PackagePath, nil, args.SentBytesChan)
}

// UpdateApplicationWithDetails updates specific fields of an existing application
func (c *clientImpl) UpdateApplicationWithDetails(id string, details api.ApplicationPatch) error {
	if len(id) == 0 {
		return errs.NewInvalidArgument("id")
	}

	b, err := json.Marshal(details)
	if err != nil {
		return err
	}

	header := http.Header{"Content-Type": []string{"application/json"}}
	op, _, err := c.QueryOperation("PATCH", client.APIPath("applications", id), nil, header, bytes.NewReader(b), "")
	if err != nil {
		return err
	}

	return op.Wait(context.Background())
}

// UpdateApplicationWithPackage updates an existing application
func (c *clientImpl) UpdateApplicationWithPackage(id, packagePath string, sentBytes chan float64) (client.Operation, error) {
	if len(id) == 0 {
		return nil, errs.NewInvalidArgument("id")
	}
	return c.upload("PATCH", client.APIPath("applications", id), nil, packagePath, nil, sentBytes)
}

// UpdateApplication updates an existing application
func (c *clientImpl) UpdateApplication(id string) (client.Operation, error) {
	if len(id) == 0 {
		return nil, errs.NewInvalidArgument("id")
	}

	header := http.Header{"Content-Type": []string{"application/json"}}
	op, _, err := c.QueryOperation("PATCH", client.APIPath("applications", id), nil, header, nil, "")
	if err != nil {
		return nil, err
	}
	return op, err
}

// ListApplications lists all available applications the AMS service knows about
func (c *clientImpl) ListApplications() ([]api.Application, error) {
	params := client.QueryParams{
		"recursion": "1",
	}
	return c.queryApplications(params)
}

// ListApplicationsWithFilters lists all available applications using AMS's API level filters
func (c *clientImpl) ListApplicationsWithFilters(filters []string) ([]api.Application, error) {
	params, err := convertFiltersToParams(filters)
	if err != nil {
		return nil, err
	}
	params["recursion"] = "1"
	return c.queryApplications(params)
}

// FindApplicationsByName list all applications matching provided pattern
func (c *clientImpl) FindApplicationsByName(pattern string) ([]api.Application, error) {
	params := client.QueryParams{
		"recursion": "1",
		"name":      pattern,
	}
	return c.queryApplications(params)
}

func (c *clientImpl) queryApplications(params client.QueryParams) ([]api.Application, error) {
	var apps []api.Application
	_, err := c.QueryStruct("GET", client.APIPath("applications"), params, nil, nil, "", &apps)
	return apps, err
}

// RetrieveApplicationByID retrieves a single application by its ID
func (c *clientImpl) RetrieveApplicationByID(id string) (*api.Application, string, error) {
	if len(id) == 0 {
		return nil, "", errs.NewInvalidArgument("id")
	}
	details := api.Application{}
	etag, err := c.QueryStruct("GET", client.APIPath("applications", id), nil, nil, nil, "", &details)
	return &details, etag, err
}

// DeleteApplicationByID deletes an existing application identified by its ID
func (c *clientImpl) DeleteApplicationByID(id string, force bool) (client.Operation, error) {
	if len(id) == 0 {
		return nil, errs.NewInvalidArgument("id")
	}
	details := api.ApplicationDelete{Force: force}
	b, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}
	op, _, err := c.QueryOperation("DELETE", client.APIPath("applications", id), nil, nil, bytes.NewReader(b), "")
	return op, err
}

// DeleteApplications deletes multiple applications identified by their ID
func (c *clientImpl) DeleteApplications(ids []string, force bool) (client.Operation, error) {
	if len(ids) == 0 {
		return nil, errs.NewInvalidArgument("ids")
	}
	details := api.ApplicationsDelete{IDs: ids, Force: force}
	b, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}
	op, _, err := c.QueryOperation("DELETE", client.APIPath("applications"), nil, nil, bytes.NewReader(b), "")
	return op, err
}

// ExportApplicationByVersion exports an existing application identified by its version
func (c *clientImpl) ExportApplicationByVersion(id string, version int, downloader func(header *http.Header, body io.ReadCloser) error) error {
	if len(id) == 0 {
		return errs.NewInvalidArgument("id")
	}
	if version < 0 {
		return errs.NewInvalidArgument("version")
	}
	hasAppImageExportSupport, err := c.HasExtension("application_image_export")
	if err != nil {
		return err
	}
	if !hasAppImageExportSupport {
		return errs.NewErrNotSupported("api extension \"application_image_export\"")
	}

	return c.download(client.APIPath("applications", id, strconv.Itoa(version)), nil, nil, downloader)
}

func (c *clientImpl) changeApplicationVersion(id string, version int, details *api.ApplicationVersionPatch) (client.Operation, error) {
	if len(id) == 0 {
		return nil, errs.NewInvalidArgument("id")
	}
	b, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}
	op, _, err := c.QueryOperation("PATCH", client.APIPath("applications", id, strconv.Itoa(version)), nil, nil, bytes.NewReader(b), "")
	return op, err
}

// PublishApplicationVersion publishes an existing application version
func (c *clientImpl) PublishApplicationVersion(id string, version int) (client.Operation, error) {
	published := true
	details := api.ApplicationVersionPatch{Published: &published}
	return c.changeApplicationVersion(id, version, &details)
}

// RevokeApplicationVersion revokes an existing and previously published application version
func (c *clientImpl) RevokeApplicationVersion(id string, version int) (client.Operation, error) {
	published := false
	details := api.ApplicationVersionPatch{Published: &published}
	return c.changeApplicationVersion(id, version, &details)
}

// DeleteApplicationVersion deletes a specific version of the given application
func (c *clientImpl) DeleteApplicationVersion(id string, version int, force bool) (client.Operation, error) {
	if len(id) == 0 {
		return nil, errs.NewInvalidArgument("id")
	}
	if version < 0 {
		return nil, errs.NewInvalidArgument("version")
	}
	details := api.ApplicationVersionDelete{Force: force}
	b, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}
	op, _, err := c.QueryOperation("DELETE", client.APIPath("applications", id, strconv.Itoa(version)), nil, nil, bytes.NewReader(b), "")
	return op, err
}
