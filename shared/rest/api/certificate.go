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

// CertificatesPost represents the fields of a new auth provided certificate
//
// swagger:model
type CertificatesPost struct {
	// Base64 encoded certificate content without the header or the footer
	// Example: MIIFUTCCAzmgAw...xjKoUEEQOzJ9
	Certificate string `json:"certificate" yaml:"certificate"`
	// TrustPassword is used to register a new client with the service
	// Example: sUp3rs3cr3t
	TrustPassword string `json:"trust-password,omitempty" yaml:"trust-password,omitempty"`
}

// Certificate represents an available client certificate
//
// swagger:model
type Certificate struct {
	// Base64 encoded certificate content without the header or the footer
	// Example: MIIFUTCCAzmgAw...xjKoUEEQOzJ9
	Certificate string `json:"certificate" yaml:"certificate"`
	// SHA-256 fingerprint of the certificate
	// Example: b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9
	Fingerprint string `json:"fingerprint" yaml:"fingerprint"`
}
