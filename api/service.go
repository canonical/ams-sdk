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

// ServiceStatus represents the status of the AMS service
//
// swagger:model ServiceStatus
type ServiceStatus struct {
	// List of supported features for the service version running
	// Example: ["addon_restore_hook", "container_logs", "registry"]
	APIExtensions []string `json:"api_extensions" yaml:"api_extensions"`
	// Shows the API stability status
	// Example: stable
	APIStatus string `json:"api_status" yaml:"api_status"`
	// API version for the service
	// Example: 1.0
	APIVersion string `json:"api_version" yaml:"api_version"`
	// Used to see if the client is trusted. Can be `trusted` or `untrusted`.
	// Example: untrusted
	Auth string `json:"auth" yaml:"auth"`
	// Authentication method used for the requests.
	// Example: 2waySSL
	AuthMethods []string `json:"auth_methods" yaml:"auth_methods"`
}
