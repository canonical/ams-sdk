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
	Resources     map[string][]string    `json:"resources,omitempty" yaml:"resources,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	MayCancel     bool                   `json:"may_cancel" yaml:"may_cancel"`
	Err           string                 `json:"err,omitempty" yaml:"err,omitempty"`
	ServerAddress string                 `json:"server_address,omitempty" yaml:"server_address,omitempty"`
}
