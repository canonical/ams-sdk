// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package api

// ConfigPost contains the field necessary to set or update a config item
type ConfigPost struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ConfigGet describes a list of config items
type ConfigGet struct {
	Config map[string]interface{} `json:"config"`
}
