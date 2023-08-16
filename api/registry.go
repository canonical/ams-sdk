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

// RegistryApplicationHooks describes the fields used to configure the hooks of an application
type RegistryApplicationHooks struct {
	Timeout string `json:"timeout" yaml:"timeout"`
}

// RegistryApplicationBootstrap describes the fields used to configure the application bootstrap
type RegistryApplicationBootstrap struct {
	Keep []string `json:"keep" yaml:"keep"`
}

// RegistryApplicationVersion describes a version of an application available
// in the registry
type RegistryApplicationVersion struct {
	BootActivity string                       `json:"boot_activity" yaml:"boot_activity"`
	Features     []string                     `json:"features" yaml:"features"`
	Hooks        RegistryApplicationHooks     `json:"hooks" yaml:"hooks"`
	Bootstrap    RegistryApplicationBootstrap `json:"bootstrap" yaml:"bootstrap"`
	VideoEncoder string                       `json:"video_encoder" yaml:"video_encoder"`
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
