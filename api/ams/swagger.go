//go:build docs
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
// +build docs

// -*- Mode: Go; indent-tabs-mode: t -*-

// AMS external REST API
//
// all AMS clients. Note that internal endpoints are not included in this
// documentation.
//
// The AMS API is available over both a local unix+http and a remote https API.
// Authentication for local users relies on group membership and access to the
// unix socket. For remote users, the default authentication method is TLS client
// certificates.
//
//	Version: 1.0
//
// swagger:meta
package api

import (
	"encoding/json"

	"github.com/anbox-cloud/ams-sdk/internal/ams/shared/rest/api"
)

func SwaggerModels() []string {
	return []string{
		swaggerModelAddon,
		swaggerModelAddonVersion,
	}
}

// define swagger schemas for API models here

// Operation
//
// swagger:model OperationResponse
type swaggerOperation struct {
	// Example: async
	Type string `json:"type"`

	// Example: Operation created
	Status string `json:"status"`

	// Example: 100
	StatusCode int `json:"status_code"`

	// Example: /1.0/operations/66e83638-9dd7-4a26-aef2-5462814869a1
	Operation string `json:"operation"`

	// Operation Metadata
	Metadata api.Operation `json:"metadata"`
}

// Swagger Synchronous response without metadata field
//
// swagger:model NoMetaSyncResponse
type swaggerNoMetaResponse struct {
	// Type of operation response
	// Example: sync
	Type string `json:"type"`

	// Status of requested operation
	// Example: Success
	Status string `json:"status"`

	// Status code of the request
	// Example: 200
	StatusCode int `json:"status_code"`

	// Error code for the operation
	// Example: 0
	Code int `json:"error_code" yaml:"error_code"`
}

// Collection Response
//
// swagger:model CollectionResponse
type swaggerCollectionResponse struct {
	// swagger:allOf
	swaggerNoMetaResponse
	// Total Count of the collection
	// Example: 99
	TotalSize int `json:"total_size"`
}

// Empty sync response
//
// swagger:response EmptySyncResponse
type swaggerEmptySyncResponse struct {
	// Empty sync response
	// in: body
	Body struct {
		// Example: sync
		Type string `json:"type"`

		// Example: Success
		Status string `json:"status"`

		// Example: 200
		StatusCode int `json:"status_code"`

		// Example: {}
		Metadata json.RawMessage `json:"metadata"`
	}
}

// Unauthorized
//
// swagger:response ErrorUnauthorized
type swaggerErrorUnauthorized struct {
	// Bad Request
	// in: body
	Body struct {
		// Example: error
		Type string `json:"type"`

		// Example: missing secret
		Error string `json:"error"`

		// Example: 401
		ErrorCode int `json:"error_code"`
	}
}

// Forbidden
//
// swagger:response ErrorForbidden
type swaggerErrorForbidden struct {
	// Bad Request
	// in: body
	Body struct {
		// Example: error
		Type string `json:"type"`

		// Example: Not Authorized
		Error string `json:"error"`

		// Example: 403
		ErrorCode int `json:"error_code"`
	}
}

// Bad Request
//
// swagger:response ErrorBadRequest
type swaggerErrorBadRequest struct {
	// Bad Request
	// in: body
	Body struct {
		// Example: error
		Type string `json:"type"`

		// Example: bad request
		Error string `json:"error"`

		// Example: 400
		ErrorCode int `json:"error_code"`

		// Example: {}
		Metadata json.RawMessage `json:"metadata"`
	}
}

// Internal Server Error
//
// swagger:response InternalServerError
type swaggerInternalServerError struct {
	// Internal server Error
	// in: body
	Body struct {
		// Example: error
		Type string `json:"type"`

		// Example: internal server error
		Error string `json:"error"`

		// Example: 500
		ErrorCode int `json:"error_code"`

		// Example: {}
		Metadata json.RawMessage `json:"metadata"`
	}
}

// Not found
//
// swagger:response ErrorNotFound
type swaggerErrorNotFound struct {
	// Not found
	// in: body
	Body struct {
		// Example: error
		Type string `json:"type"`

		// Example: not found
		Error string `json:"error"`

		// Example: 404
		ErrorCode int `json:"error_code"`
	}
}

// Service Unavailable
//
// swagger:response ErrorServiceUnavailable
type swaggerErrorServiceUnavailable struct {
	// Service Unavailable
	// in: body
	Body struct {
		// Example: error
		Type string `json:"type"`

		// Example: service unavailable
		Error string `json:"error"`

		// Example: 503
		ErrorCode int `json:"error_code"`
	}
}

// Already Exists
//
// swagger:response ErrorAlreadyExists
type swaggerErrorConflict struct {
	// in: body
	Body struct {
		// Example: error
		Type string `json:"type"`

		// Example: already exists
		Error string `json:"error"`

		// Example: 409
		ErrorCode int `json:"error_code"`
	}
}
