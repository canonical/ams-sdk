// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package api

import "time"

// We don't put metadata here as it should be in
const swaggerModelOperation = `
operation:
  type: object
  properties:
    id:
      type: string
      description: UUID of the operation
      example: "c6832c58-0867-467e-b245-2962d6527876"
    class:
      type: string
      description: Class of the operation
      enum: [task, websocket, token]
      example: task
    description:
      type: string
      description: Human readable description of the operation
      example: updating addon 3apqo5te
    created_at:
      type: string
      format: date-time
      description: When the operation was created
      example: 2018-04-02T16:49:36.341463206+02:00
    updated_at:
      type: string
      format: date-time
      description: Last time the operation was updated
      example: 2018-04-02T16:50:36.341463206+02:00
    status:
      type: string
      description: String version of the operation status
      example: Running
    status_code:
      $ref: '#/definitions/rest-api-status-code'
    resources:
      type: object
      description: |
        Dictionnary of resource types (containers, snapshots, images)
        and affected resources
      additionalProperties:
        type: array
        items:
          type: string
    may_cancel:
      type: boolean
      description: Whether this operation can be canceled (DELETE over REST)
    err:
      type: string
      description: The error string if the operation failed
`

// Operation represents a LXD background operation
type Operation struct {
	ID            string                 `json:"id" yaml:"id"`
	Class         string                 `json:"class" yaml:"class"`
	Description   string                 `json:"description" yaml:"description"`
	CreatedAt     time.Time              `json:"created_at" yaml:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at" yaml:"updated_at"`
	Status        string                 `json:"status" yaml:"status"`
	StatusCode    StatusCode             `json:"status_code" yaml:"status_code"`
	Resources     map[string][]string    `json:"resources" yaml:"resources"`
	Metadata      map[string]interface{} `json:"metadata" yaml:"metadata"`
	MayCancel     bool                   `json:"may_cancel" yaml:"may_cancel"`
	Err           string                 `json:"err" yaml:"err"`
	ServerAddress string                 `json:"server_address" yaml:"server_address"`
}
