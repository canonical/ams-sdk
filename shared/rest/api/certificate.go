// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package api

// CertificatesPost represents the fields of a new auth provided certificate
type CertificatesPost struct {
	Certificate   string `json:"certificate" yaml:"certificate"`
	TrustPassword string `json:"trust-password,omitempty" yaml:"trust-password,omitempty"`
}

// Certificate represents an available client certificate
type Certificate struct {
	Certificate string `json:"certificate" yaml:"certificate"`
	Fingerprint string `json:"fingerprint" yaml:"fingerprint"`
}
