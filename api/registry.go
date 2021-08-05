// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2020 Canonical Ltd.  All rights reserved.

package api

// RegistryApplicationVersion describes a version of an application available
// in the registry
type RegistryApplicationVersion struct {
	BootActivity string   `json:"boot_activity" yaml:"boot_activity"`
	Features     []string `json:"features"`
}

// RegistryApplication describes a single application available in the registry
type RegistryApplication struct {
	Name          string                              `json:"name" yaml:"name"`
	Architectures []string                            `json:"architectures" yaml:"architectures"`
	BootPackage   string                              `json:"boot_package" yaml:"boot_package"`
	InstanceType  string                              `json:"instance_type" yaml:"instance_type"`
	Tags          []string                            `json:"tags" yaml:"tags"`
	Versions      map[int]*RegistryApplicationVersion `json:"versions" yaml:"versions"`
}
