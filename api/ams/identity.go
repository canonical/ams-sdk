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

// IdentityPost is the input body adding for an identity
//
// swagger:model
type IdentityPost struct {
	// AuthenticationMethod is the method used to authenticate the identity
	// Example: oidc
	AuthenticationMethod string `json:"authentication_method" yaml:"authentication_method"`
	// Name of the new identity
	// Example: john-doe
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// Email of the new user (for OIDC based identities)
	// Example: john.doe@example.com
	Email string `json:"email,omitempty" yaml:"email,omitempty"`

	// Base64 encoded certificate content without the header or the footer (for TLS based identities)
	// Example: MIIFUTCCAzmgAw...xjKoUEEQOzJ9
	Certificate string `json:"certificate,omitempty" yaml:"certificate,omitempty"`
}

// Identity represents an authenticated party that can make requests to the HTTPS API.
//
// swagger:model
type Identity struct {
	// A unique identifier for the identity
	// Example: btavtegj1qm58qg7ru50
	ID string `json:"id" yaml:"id"`
	// Name of the new user
	// Example: john-doe
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// The list of groups for which the identity is a member.
	// Example: ["foo", "bar"]
	Groups []string `json:"groups" yaml:"groups"`
	// The authentication method that the identity
	// authenticates to AMS with.
	// Example: tls
	AuthenticationMethod string `json:"authentication_method" yaml:"authentication_method"`
	// CreatedAt specifies the time at which the identity was created
	// Example: 1689604498
	CreatedAt int64 `json:"created_at" yaml:"created_at"`
	// UpdatedAt specifies the time at which the identity was last updated
	// Example: 1689604498
	UpdatedAt int64 `json:"updated_at" yaml:"updated_at"`
	// The certificate fingerprint for TLS based identities
	Fingerprint string `json:"fingerprint,omitempty" yaml:"fingerprint,omitempty"`
	// Email of the new user for OIDC based identities
	// Example: john.doe@example.com
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
	// Base64 encoded certificate content without the header or the footer (for TLS based identities)
	// Example: MIIFUTCCAzmgAw...xjKoUEEQOzJ9
	Certificate string `json:"certificate,omitempty" yaml:"certificate,omitempty"`
}

// IdentityPut contains the editable fields of an Identity.
//
// swagger:model
type IdentityPut struct {
	// Groups is the list of groups for which the identity is a member.
	// Example: ["foo", "bar"]
	Groups []string `json:"groups" yaml:"groups"`
}

// IdentityDelete represent the body in a delete method
//
// swagger:model
type IdentityDelete struct {
	// Force allows the identity to be deleted even if it is a part of groups
	// Example: true
	Force bool `json:"force" yaml:"force"`
}

// GetIdentityFilters returns an array of attributes available on the api to
// filter identites
func GetIdentityFilters() []string {
	return []string{
		"id",
		"auth_type",
	}
}
