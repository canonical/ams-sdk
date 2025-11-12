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

// AuthGroupPost is used for creating a group.
//
// swagger:model
type AuthGroupPost struct {
	// Name is the name of the group.
	// Example: default-c1-viewers
	Name string `json:"name" yaml:"name"`

	// Description is a short description of the group.
	// Example: Developers of application foo.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

// AuthGroupDelete represent the body in a delete method
//
// swagger:model
type AuthGroupDelete struct {
	// Force allows the group to be deleted even if it is a part of groups
	// Example: true
	Force bool `json:"force" yaml:"force"`
}

// AuthGroup is the type for an AMS authorization group.
//
// swagger:model
type AuthGroup struct {
	AuthGroupPost `json:",inline" yaml:",inline"`

	// CreatedAt specifies the time at which the group was created. The value is a UTC timestamp.
	// Example: 1689604498
	CreatedAt int64 `json:"created_at" yaml:"created_at"`

	// UpdatedAt specifies the time at which the group was last updated. The value is a UTC timestamp.
	// Example: 1689604498
	UpdatedAt int64 `json:"updated_at" yaml:"updated_at"`

	// Identities is a list of authentication method to slice of group identifiers.
	Identities []string `json:"identities" yaml:"identities"`

	// Permissions represent the list of permissions assigned to an auth group
	// Example: [{"entitlement": "can_view", "resource": "application:foo"}, {"entitlement": "can_edit", "resource": "instance:bar"}]
	Permissions []Permission `json:"permissions" yaml:"permissions"`

	// Immutable defines whether a group is immutable and cannot be deleted or renamed
	// Example: false
	Immutable bool `json:"immutable" yaml:"immutable"`
}

// AuthGroupPut replaces the editable fields for an auth group with the required values.
//
// swagger:model
type AuthGroupPut struct {
	// Description is a short description of the group.
	// Example: Developers of application foo.
	Description *string `json:"description,omitempty" yaml:"description,omitempty"`

	// Permissions represent the list of permissions assigned to an auth group
	// Example: [{"entitlement": "can_view", "resource": "application:foo"}, {"entitlement": "can_edit", "resource": "instance:bar"}]
	Permissions *[]Permission `json:"permissions,omitempty" yaml:"permissions,omitempty"`
}

// GetGroupFilters returns an array of attributes available on the api to
// filter groups
func GetGroupFilters() []string {
	return []string{
		"name",
	}
}
