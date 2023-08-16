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

package api

import (
	"encoding/json"
)

// "there are three standard return types: Standard return value,
// Background operation, Error", each returning a JSON object with the
// following "type" field:
const (
	ResponseTypeSync  ResponseType = "sync"
	ResponseTypeAsync ResponseType = "async"
	ResponseTypeError ResponseType = "error"
)

// ResponseType represents a valid LXD response type
type ResponseType string

// ResponseRaw represents a REST operation in its original form
type ResponseRaw struct {
	Response `yaml:",inline"`
	Metadata interface{} `json:"metadata" yaml:"metadata"`
}

// Response represents a LXD operation
type Response struct {
	Type ResponseType `json:"type" yaml:"type"`

	// Total Count of the collection
	TotalSize *int `json:"total_size,omitempty" yaml:"total_size,omitempty"`

	// Valid only for Sync responses
	Status     string `json:"status" yaml:"status"`
	StatusCode int    `json:"status_code" yaml:"status_code"`

	// Valid only for Async responses
	Operation string `json:"operation,omitempty" yaml:"operation,omitempty"`

	// Valid only for Error responses
	Code  int    `json:"error_code" yaml:"error_code"`
	Error string `json:"error,omitempty" yaml:"error,omitempty"`

	// Valid for Sync and Error responses
	Metadata json.RawMessage `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// MetadataAsMap parses the Response metadata into a map
func (r *Response) MetadataAsMap() (map[string]interface{}, error) {
	ret := map[string]interface{}{}
	err := r.MetadataAsStruct(&ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// MetadataAsOperation turns the Response metadata into an Operation
func (r *Response) MetadataAsOperation() (*Operation, error) {
	op := Operation{}
	err := r.MetadataAsStruct(&op)
	if err != nil {
		return nil, err
	}

	return &op, nil
}

// MetadataAsStringSlice parses the Response metadata into a slice of string
func (r *Response) MetadataAsStringSlice() ([]string, error) {
	sl := []string{}
	err := r.MetadataAsStruct(&sl)
	if err != nil {
		return nil, err
	}

	return sl, nil
}

// MetadataAsStruct parses the Response metadata into a provided struct
func (r *Response) MetadataAsStruct(target interface{}) error {
	return json.Unmarshal(r.Metadata, &target)
}
