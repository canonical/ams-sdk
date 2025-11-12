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

// InstanceSharesPost contains information required to create a share for an instance
type InstanceSharesPost struct {
	// Type is the type of the session to be shared
	Type string `json:"type" example:"adb"`
	// Description is the user provided description for the share to be created
	Description string `json:"description" example:"instance shared with john"`
	// ExpiryAt specifies the exact expiration time for the share.
	// The value is a Unix timestamp (UTC) indicating when the share will expire.
	ExpiryAt *int64 `json:"expiry_at" example:"1610645117"`
}

// InstanceSharesPostResponse is returned when sharing an instance
//
// swagger:model
type InstanceSharesPostResponse struct {
	Response
	Metadata InstanceSharesPostResponseData `json:"metadata"`
}

// InstanceSharesPostResponseData contains the share URL
type InstanceSharesPostResponseData struct {
	// ID of the share
	ID string `json:"id" example:"ctigqirc209urni862a0"`
	// Type of the share
	Type string `json:"type" example:"adb"`
	// URL is the endpoint to reach to start the WebRTC signaling process
	URL string `json:"url" example:"https://api.example.com/1.0/sessions/foo/connect?token=bar"`
	// ExpiryAt specifies the exact expiration time for the share
	ExpiryAt int64 `json:"expiry_at" example:"1610645117"`
}

// InstanceSharePatch contains information required to update a share
type InstanceSharePatch struct {
	// Description is the user provided description for the share to be created
	Description *string `json:"description,omitempty" example:"instance shared with john"`
	// ExpiryAt specifies the exact expiration time for the share.
	// The value is a Unix timestamp (UTC) indicating when the share will expire.
	ExpiryAt *int64 `json:"expiry_at,omitempty" example:"1610645117"`
}
