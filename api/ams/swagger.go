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
// ## Authentication
//
// Not all REST API endpoints require authentication, for example, the following API calls are allowed for everyone:
//
// * `GET /`
// * `GET /1.0`
// * `GET /1.0/version`
//
// Some endpoints require an additional authentication token to ensure that the requester is authorised to access the resource, for example:
//
// * `GET /1.0/artifacts`
// * `PATCH /1.0/instances/<name>`
//
// ## API versioning
//
// The details of a version of the API can be retrieved using `GET /<version>`. For example, `GET /1.0`.
//
// If an API version is bumped to a major version, it indicates that backward compatibility is affected.
//
// Feature additions done without breaking backward compatibility only result in additions to `api_extensions` which can be used by the client to check if a given feature is supported by the server.
//
// ## Return values
//
// There are three standard return types:
//
// * Standard return value
// * Background operation
// * Error
//
// ### Standard return value
//
// For a standard synchronous operation, the following dict is returned:
//
// ```json
// {
// "type": "sync",
// "status": "Success",
// "status_code": 200,
// "metadata": {}
// }
// ```
//
// HTTP response status code is 200.
//
// ### Background operation
//
// When a request results in a background operation, the HTTP code is set to 202 (Accepted) and the Location HTTP header is set to the operation URL.
//
// The response body is a dict with the following structure:
//
// ```json
// {
// "type": "async",
// "status": "OK",
// "status_code": 100,
// "operation": "/1.0/containers/<id>",
// "metadata": {}
// }
// ```
//
// The operation metadata structure looks like:
//
// ```json
// {
// "id": "c6832c58-0867-467e-b245-2962d6527876",
// "class": "task",
// "created_at": "2018-04-02T16:49:36.341463206+02:00",
// "updated_at": "2018-04-02T16:49:36.341463206+02:00",
// "status": "Running",
// "status_code": 103,
// "resources": {
// "containers": [
// "/1.0/containers/3apqo5te"
// ]
// },
// "metadata": null,
// "may_cancel": false,
// "err": ""
// }
// ```
//
// The body is mostly provided as a user friendly way of seeing what's going on without having to pull the target operation, all information in the body can also be retrieved from the background operation URL.
//
// ### Error
//
// There are various situations in which something may immediately go wrong, in those cases, the following return value is used:
//
// ```json
// {
// "type": "error",
// "error": "Failure",
// "error_code": 400,
// "metadata": {}
// }
// ```
//
// HTTP response status code is one of 400, 401, 403, 404, 409, 412 or 500.
//
// ## Status codes
//
// The REST API often has to return status information, which could be the reason for an error, the current state of an operation or the state of the various resources it exports.
//
// To make it simple to debug, there are two ways in which such information is represented - a numeric representation of the state which is guaranteed never to change and can be relied on by API clients and a text version so that it is easier for people manually using the API to understand better. In most cases, those will be called `status` and `status_code`, the former being the user friendly string representation and the latter being the fixed numeric value.
//
// The codes are always 3 digits, with the following ranges:
//
// * 100 to 199: resource state (started, stopped, ready, ...)
// * 200 to 399: positive action result
// * 400 to 599: negative action result
// * 600 to 999: future use
//
// ### List of current status codes
//
// | Code  | Meaning |
// |------|------ |
// | 100   | Operation created |
// | 101   | Started |
// | 102   | Stopped |
// | 103   | Running |
// | 104   | Cancelling |
// | 105   | Pending |
// | 106   | Starting |
// | 107   | Stopping |
// | 108   | Aborting |
// | 109   | Freezing |
// | 110   | Frozen |
// | 111   | Thawed |
// | 200   | Success |
// | 400   | Failure |
// | 401   | Cancelled |
//
// ## Recursion
//
// To optimise queries of large lists, recursion is implemented for collections. A `recursion` argument can be passed to a GET query against a collection.
//
// The default value is 0 which means that collection member URLs are returned. Setting it to 1 will have those URLs be replaced by the object they point to (typically a dict).
//
// Recursion is implemented by simply replacing any pointer to a job (URL) by the object itself.
//
// ## Async operations
//
// Any operation which take more than a second must be done in the background, returning a background operation ID to the client. With this ID, the client is able to either poll for a status update or wait for a notification using the long-poll API.
//
// ## Notifications
//
// A web-socket based API is available for notifications. Different notification types exist to limit the traffic going to the client. It is recommended that the client always subscribes to the *operations* notification type before triggering remote operations so that it doesn't have to continually poll for their status.
//
// ## PUT vs PATCH
//
// PUT and PATCH APIs are supported to modify existing objects.
//
// PUT replaces the entire object with a new definition, it's typically called after the current object state was retrieved through GET.
//
// To avoid race conditions, the ETag header should be read from the GET response and sent as If-Match for the PUT request. Doing so makes the request fail if the object was modified between GET and PUT.
//
// PATCH can be used to modify a single field inside an object by only specifying the property that you want to change. To unset a key, setting it to empty will usually do the trick, but there are cases where PATCH won't work and PUT needs to be used instead.
//
// ## Authorisation
//
// Some operation may require a token to be included in the HTTP Authorisation header even if the request is already authenticated using a trusted certificate. If the token is not valid, the request is rejected by the server. This ensures that only authorised clients can access those endpoints.
//
// Authorization: bearer <token>
//
// ## File upload
//
// Some operations require uploading a payload. To prevent the difficulties of handling multipart requests, a unique file is uploaded and its bytes are included in the body of the request. The following metadata associated with the file is included in extra HTTP headers:
//
// * X-AMS-Fingerprint: Fingerprint of the payload being added
// * X-AMS-Request: Metadata for the payload. This is a JSON, specific for the operation.
//
// ## Instances and Containers
//
// The documentation shows paths such as `/1.0/instances/...`, which were introduced with Anbox Cloud version 1.20.0. Older releases that supported only containers and not virtual machines supply the exact same API at `/1.0/containers/...`.
//
// Although deprecated, the `1.0/containers/...` API is still available for backward compatibility.
//
// Version: 1.0
//
// swagger:meta
package api

import (
	"encoding/json"

	"github.com/anbox-cloud/ams-sdk/pkg/ams/shared/rest/api"
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
