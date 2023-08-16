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

// EventType is used to describe the type of an event
//
// swagger:enum EventType
type EventType string

const (
	// EventTypeOperation is the event type sent when a operation event is reported
	EventTypeOperation EventType = "operation"
	// EventTypeLifecycle is the event type sent when a lifecycle event is reported
	EventTypeLifecycle EventType = "lifecycle"
)

// Event defines the structure of an event sent on the events API endpoint
//
// swagger:model
type Event struct {
	// Type defines the type of event. Listeners can watch specific
	// event types
	Type EventType `json:"type"`
	// Timestamp (in ISO8601 format) is filled when sending the event if empty
	// Example: 2017-07-28T05:02:22.92201407Z
	Timestamp time.Time `json:"timestamp"`
	// Metadata represents the actual event data
	// Example: { "class": "task", "created_at": "2017-07-28T05:02:22.92201407Z", "description": "Deleting container", "err": "", "id": "bc85137b-b20d-470a-a6ea-daa9a2b8506a", "may_cancel": false, "metadata": null, "resources": { "containers": [ "/1.0/containers/c0946voj1qm6t2783db0" ] }, "server_address": "", "status": "Success", "status_code": 200, "updated_at": "2017-07-28T05:02:22.92201407Z" }
	Metadata interface{} `json:"metadata"`
}

// LifecycleEventAction describes a single lifecycle action
type LifecycleEventAction string

const (
	// LifecycleEventActionContainerCreated is sent when a container was created
	LifecycleEventActionContainerCreated LifecycleEventAction = "container-created"
	// LifecycleEventActionContainerScheduled is sent when a container was scheduled
	LifecycleEventActionContainerScheduled LifecycleEventAction = "container-scheduled"
	// LifecycleEventActionContainerStarted is sent when a container was started
	LifecycleEventActionContainerStarted LifecycleEventAction = "container-started"
	// LifecycleEventActionContainerStopped is sent when a container was stopped
	LifecycleEventActionContainerStopped LifecycleEventAction = "container-stopped"
	// LifecycleEventActionContainerRemoved is sent when a container was removed
	LifecycleEventActionContainerRemoved LifecycleEventAction = "container-removed"
	// LifecycleEventActionContainerFailed is sent when a container failed
	LifecycleEventActionContainerFailed LifecycleEventAction = "container-failed"
)

// LifecycleEvent contains information about a lifecycle event
type LifecycleEvent struct {
	Action LifecycleEventAction `json:"action"`
	Source string               `json:"source"`
}
