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

// AuthenticationMethod describes the type of an authentication to use with a remote (tls, oidc)
type AuthenticationMethod = string

const (
	// AuthenticationMethodTLS specifies that the authentication type to be
	// used with a remote should be mTLS
	AuthenticationMethodTLS AuthenticationMethod = "tls"
	// AuthenticationMethodOIDC specifies that the authentication type to be
	// used with a remote should be OIDC.
	AuthenticationMethodOIDC AuthenticationMethod = "oidc"
)

// OIDCResponse represents the OIDC configuration accepted by AMS
//
// swagger:model OIDCResponse
type OIDCResponse struct {
	// Issuer URL as configured in AMS
	// Example: "https://myauth.auth0.com/"
	IssuerURL string `json:"issuer" yaml:"issuer"`
	// Required Scope for AMS to work. These scopes should be request by clients
	// when performing OIDC based authentication with an Identity Provider.
	// Example: ["openid", "email", "profile"]
	RequiredScopes []string `json:"required_scopes" yaml:"required_scopes"`
	// ClientID to use when using flows that do not require a client_secret on
	// identity provider e.g Device flow grant.
	// Example: web
	ClientID string `json:"client_id,omitempty" yaml:"client_id,omitempty"`
	// Audience to use when using flows that do not require a client_secret on
	// identity provider e.g Device flow grant.
	// Example: http://ams.example.com
	Audience string `json:"audience,omitempty" yaml:"audience,omitempty"`
}
