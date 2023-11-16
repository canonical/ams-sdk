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
	"fmt"
	"net/http"
	"strconv"

	"github.com/anbox-cloud/ams-sdk/api"
	errs "github.com/anbox-cloud/ams-sdk/shared/errors"
	"github.com/anbox-cloud/ams-sdk/shared/rest/client"
)

// ListImages lists all available images the AMS service currently has
func (c *clientImpl) ListImages() ([]api.Image, error) {
	images := []api.Image{}
	params := client.QueryParams{
		"recursion": "1",
	}
	_, err := c.QueryStruct("GET", client.APIPath("images"), params, nil, nil, "", &images)
	return images, err
}

// AddImage adds a new image with the given payload
func (c *clientImpl) AddImage(name, packagePath string, isDefault bool, sentBytes chan float64) (client.Operation, error) {
	details := api.ImagesPost{
		Name:    name,
		Default: isDefault,
	}
	return c.upload("POST", client.APIPath("images"), nil, packagePath, details, sentBytes)
}

// ImportImage imports a new image from the image server
func (c *clientImpl) ImportImage(name, path string, isDefault bool) (client.Operation, error) {
	return c.ImportImageByType(name, path, api.ImageTypeAny, isDefault)
}

// ImportImageByType imports a new image of the given type from the image server
func (c *clientImpl) ImportImageByType(name, path string, imgType api.ImageType, isDefault bool) (client.Operation, error) {
	details := api.ImagesPost{
		Name:    name,
		Path:    path,
		Default: isDefault,
		Type:    imgType,
	}

	b, err := json.Marshal(&details)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request body: %v", err)
	}

	header := http.Header{"Content-Type": []string{"application/json"}}
	op, _, err := c.QueryOperation("POST", client.APIPath("images"), nil, header, bytes.NewReader(b), "")
	return op, err

}

// UpdateImage updates an existing image with the given payload
func (c *clientImpl) UpdateImage(id, packagePath string, sentBytes chan float64) (client.Operation, error) {
	details := api.ImagePatch{}
	return c.upload("PATCH", client.APIPath("images", id), nil, packagePath, details, sentBytes)
}

func (c *clientImpl) SetDefaultImage(id string) error {
	d := new(bool)
	*d = true

	details := api.ImagePatch{
		Default: d,
	}

	b, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("could not marshal request body: %v", err)
	}

	header := http.Header{"Content-Type": []string{"application/json"}}
	op, _, err := c.QueryOperation("PATCH", client.APIPath("images", id), nil, header, bytes.NewReader(b), "")
	if err != nil {
		return err
	}
	return op.Wait(context.Background())
}

// DeleteImageByIDOrName deletes an image identified by the given id or name
func (c *clientImpl) DeleteImageByIDOrName(id string, force bool) (client.Operation, error) {
	if len(id) == 0 {
		return nil, errs.NewInvalidArgument("id")
	}

	details := api.ImageDelete{Force: force}
	b, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}

	op, _, err := c.QueryOperation("DELETE", client.APIPath("images", id), nil, nil, bytes.NewReader(b), "")
	return op, err
}

// DeleteImageVersion deletes a single image version
func (c *clientImpl) DeleteImageVersion(id string, version int) (client.Operation, error) {
	if len(id) == 0 {
		return nil, errs.NewInvalidArgument("id")
	}
	op, _, err := c.QueryOperation("DELETE", client.APIPath("images", id, strconv.Itoa(version)), nil, nil, nil, "")
	return op, err
}

func (c *clientImpl) RetrieveDefaultImage() (*api.Image, string, error) {
	i := []*api.Image{}
	params := client.QueryParams{
		"default": "true",
	}
	etag, err := c.QueryStruct("GET", client.APIPath("images"), params, nil, nil, "", &i)
	if len(i) != 1 {
		return nil, "", fmt.Errorf("Failed to retrieve default image")
	}
	return i[0], etag, err
}

// RetrieveImageByIDOrName retrieves a single image by its ID or Name
func (c *clientImpl) RetrieveImageByIDOrName(id string) (*api.Image, string, error) {
	if len(id) == 0 {
		return nil, "", errs.NewInvalidArgument("id")
	}
	i := &api.Image{}
	etag, err := c.QueryStruct("GET", client.APIPath("images", id), nil, nil, nil, "", i)
	return i, etag, err
}
