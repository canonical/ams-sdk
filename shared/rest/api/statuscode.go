// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package api

const swaggerModelStatusCode = `
rest-api-status-code:
  type: integer
  description: |
    * 100: Operation Created
    * 101: Started
    * 102: Stopped
    * 103: Running
    * 104: Cancelling
    * 105: Pending
    * 106: Startings
    * 107: Stopping
    * 108: Aborting
    * 109: Freezing
    * 110: Frozen
    * 111: Thawed
    * 112: Error
    * 200: Success
    * 400: Failure
    * 401: Cancelled
  enum:
    - 100
    - 101
    - 102
    - 103
    - 104
    - 105
    - 106
    - 107
    - 108
    - 109
    - 110
    - 111
    - 112
    - 200
    - 400
    - 401
`

// StatusCode represents a valid REST operation
type StatusCode int

// status codes
const (
	OperationCreated StatusCode = 100
	Started          StatusCode = 101
	Stopped          StatusCode = 102
	Running          StatusCode = 103
	Cancelling       StatusCode = 104
	Pending          StatusCode = 105
	Starting         StatusCode = 106
	Stopping         StatusCode = 107
	Aborting         StatusCode = 108
	Freezing         StatusCode = 109
	Frozen           StatusCode = 110
	Thawed           StatusCode = 111
	Error            StatusCode = 112

	Success StatusCode = 200

	Failure   StatusCode = 400
	Cancelled StatusCode = 401
)

// String returns a suitable string representation for the status code
func (o StatusCode) String() string {
	return map[StatusCode]string{
		OperationCreated: "Operation created",
		Started:          "Started",
		Stopped:          "Stopped",
		Running:          "Running",
		Cancelling:       "Cancelling",
		Pending:          "Pending",
		Success:          "Success",
		Failure:          "Failure",
		Cancelled:        "Cancelled",
		Starting:         "Starting",
		Stopping:         "Stopping",
		Aborting:         "Aborting",
		Freezing:         "Freezing",
		Frozen:           "Frozen",
		Thawed:           "Thawed",
		Error:            "Error",
	}[o]
}

// IsFinal will return true if the status code indicates an end state
func (o StatusCode) IsFinal() bool {
	return int(o) >= 200
}
