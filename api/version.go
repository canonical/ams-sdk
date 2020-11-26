// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package api

// VersionGet is the JSON response from the API version request method
type VersionGet struct {
	Version string `json:"version"`
}
