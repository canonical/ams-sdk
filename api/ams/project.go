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

// ProjectStatus represents the lifecycle status of a Project.
type ProjectStatus int

const (
	// ProjectStatusError represents the state when a project encountered an error
	// while its hypervisor level resources were being provisioned.
	ProjectStatusError ProjectStatus = -1
	// ProjectStatusUnknown represents the state when the real state of a project
	// cannot be determined.
	ProjectStatusUnknown ProjectStatus = 0
	// ProjectStatusCreated represents the state when a project record was persisted
	// and hypervisor provisioning is in progress.
	ProjectStatusCreated ProjectStatus = 1
	// ProjectStatusReady represents the state when a project is fully provisioned
	// and operational.
	ProjectStatusReady ProjectStatus = 2
	// ProjectStatusDeleted represents the state when a project has been removed.
	ProjectStatusDeleted ProjectStatus = 3
)

func (s *ProjectStatus) String() string {
	switch *s {
	case ProjectStatusError:
		return "error"
	case ProjectStatusCreated:
		return "created"
	case ProjectStatusReady:
		return "ready"
	case ProjectStatusDeleted:
		return "deleted"
	}
	return "unknown"
}

// Project is the public API representation of an Anbox Project.
//
// swagger:model
type Project struct {
	// Name of the project
	// Example: team-a
	Name string `json:"name" yaml:"name"`
	// Description of the project
	// Example: Project for team A
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Subnet the project is configured with in CIDR notation
	// Example: 100.64.0.0/28
	Subnet string `json:"subnet,omitempty" yaml:"subnet,omitempty"`
	// Status of the project as an integer value
	// Example: 2
	StatusCode ProjectStatus `json:"status_code" yaml:"status_code"`
	// Current status of the project
	// Enum: unknown,created,ready,deleted,error
	// Example: ready
	Status string `json:"status" yaml:"status"`
	// Creation UTC timestamp of the project
	// Example: 1532150640
	CreatedAt int64 `json:"created_at" yaml:"created_at"`
	// Entitlements available to the user for this project
	// Example: ["can_view", "can_edit", "can_delete"]
	Entitlements []string `json:"entitlements,omitempty" yaml:"entitlements,omitempty"`
	// Immutable indicates that the project cannot be updated or deleted,
	// regardless of any force flag. Only true for the built-in "default" project.
	Immutable bool `json:"immutable" yaml:"immutable"`
	// InstanceCount is the number of instances currently attached to this project.
	// Example: 3
	InstanceCount int `json:"instance_count" yaml:"instance_count"`
}

// ProjectsPost is the request body for creating a Project.
//
// swagger:model
type ProjectsPost struct {
	// Name of the project
	// Example: team-a
	Name string `json:"name" yaml:"name"`
	// Description of the project
	// Example: Project for team A
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Subnet to assign to the project in CIDR notation. Only validated for
	// well-formedness; no allocation from a pool happens yet.
	// Example: 100.64.0.0/28
	Subnet string `json:"subnet,omitempty" yaml:"subnet,omitempty"`
}

// ProjectsPatch is the request body for updating a Project.
//
// swagger:model
type ProjectsPatch struct {
	// Description of the project
	// Example: Project for team A
	Description *string `json:"description,omitempty" yaml:"description,omitempty"`
}

// ProjectDelete is the request body for deleting a Project.
//
// swagger:model
type ProjectDelete struct {
	// Whether deletion of the project should be forced, deleting any
	// contained instances first. Never overrides immutability.
	// Example: false
	Force bool `json:"force" yaml:"force"`
}

// GetProjectFilters returns an array of attributes available on the api to
// filter projects.
func GetProjectFilters() []string {
	return []string{
		"name",
		"status",
	}
}
