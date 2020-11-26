// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

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
type Task struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	ObjectID   string `json:"object_id"`
	ObjectType string `json:"object_type"`
}
