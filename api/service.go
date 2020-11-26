// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package api

// ServiceStatus represents the status of the AMS service
type ServiceStatus struct {
	APIExtensions []string `json:"api_extensions" yaml:"api_extensions"`
	APIStatus     string   `json:"api_status" yaml:"api_status"`
	APIVersion    string   `json:"api_version" yaml:"api_version"`
	Auth          string   `json:"auth" yaml:"auth"`
	AuthMethods   []string `json:"auth_methods" yaml:"auth_methods"`
}
