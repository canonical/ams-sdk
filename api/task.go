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

// TaskStatus holds the status of a task
type TaskStatus int

const (
	// TaskStatusUnknown is the status of a task when its real status is unknown
	TaskStatusUnknown TaskStatus = -1
	// TaskStatusCreated is the status of a task when it was created by the orchestrator
	// but not processed any further.
	TaskStatusCreated TaskStatus = 0
	// TaskStatusAssigned is the status of a task when it got assigned to a worker which
	// will process it.
	TaskStatusAssigned TaskStatus = 1
	// TaskStatusPrepared is the status of a task when all necessary resources for the
	// associated container are allocated
	TaskStatusPrepared TaskStatus = 2
	// TaskStatusStarted is the status of a task when the associated container was started
	// and is starting up.
	TaskStatusStarted TaskStatus = 3
	// TaskStatusRunning is the status of a task when the associated container is up and running
	TaskStatusRunning TaskStatus = 4
	// TaskStatusShutdown is the status of a task when its in the process of shutting down.
	TaskStatusShutdown TaskStatus = 5
	// TaskStatusStopped is the status of a task when the associated container got stopped.
	TaskStatusStopped TaskStatus = 6
	// TaskStatusCompleted is the status of a task when all operations are completed.
	TaskStatusCompleted TaskStatus = 7
	// TaskStatusError is the status of a task when an error occurred.
	TaskStatusError TaskStatus = 8
	// TaskStatusDeleted is the status of a task when it and its associated container should
	// be deleted.
	TaskStatusDeleted TaskStatus = 9
)

// String returns the textual representation of the task status
func (s TaskStatus) String() string {
	switch s {
	case TaskStatusCreated:
		return "created"
	case TaskStatusPrepared:
		return "prepared"
	case TaskStatusStarted:
		return "started"
	case TaskStatusRunning:
		return "running"
	case TaskStatusStopped:
		return "stopped"
	case TaskStatusShutdown:
		return "shutdown"
	case TaskStatusCompleted:
		return "completed"
	case TaskStatusError:
		return "error"
	case TaskStatusDeleted:
		return "deleted"
	}
	return "unknown"
}

// Task is the scheduling unit AMS uses for container launches
//
// swagger:model
type Task struct {
	// ID of the task
	// Example: c055dl0j1qm027422feg
	ID string `json:"id"`
	// Status of the task
	// Enum: created,prepared,started,running,stopped,shutdown,completed,error,deleted,unknown
	// Example: running
	Status string `json:"status"`
	// ID of the object that the task is operating on
	// Example: c055dl0j1qm027422fe0
	ObjectID string `json:"object_id"`
	// Type of the object that the task is operating on
	// Example: container
	ObjectType string `json:"object_type"`
}
