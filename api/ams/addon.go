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

const swaggerModelAddonVersion = `
addon-version:
  type: object
  properties:
    number:
      type: integer
      example: 3
    fingerprint:
      type: string
      example: e25f3e3bcead096461d89d8ab7043f14bdb1ecd39
    size:
      type: integer
      format: int64
      example: 28171923
    created_at:
      type: integer
      format: int64
      example: 1583946268
`

// AddonVersion describes a single version of an addon
//
// swagger:model
type AddonVersion struct {
	// Version for the addon
	// Example: 0
	Number int `json:"version" yaml:"version"`
	// SHA-256 fingerprint of the addon version
	// Example: 0791cfc011f67c60b7bd0f852ddb686b79fa46083d9d43ef9845c9235c67b261
	Fingerprint string `json:"fingerprint" yaml:"fingerprint"`
	// Size (in bytes) of the addon payload
	// Example: 529887868
	Size int64 `json:"size" yaml:"size"`
	// Creation timestamp of the addon
	// Example: 1610641117
	CreatedAt int64 `json:"created_at" yaml:"created_at"`
}

const swaggerModelAddon = `
addon:
  type: object
  properties:
    name:
      type: string
    versions:
      type: array
      items:
        $ref: "#/definitions/addon-version"
`

// Addon describes a package with additional functionality to be added to containers
//
// swagger:model
type Addon struct {
	// Name of the addon
	// Example: my-addon
	Name string `json:"name" yaml:"name"`
	// List of versions of the addon
	Versions []AddonVersion `json:"versions" yaml:"versions"`
	// List of applications using this addon
	// Example: ["app1", "app2"]
	UsedBy []string `json:"used_by" yaml:"used_by"`
}

// AddonsPost is used to create a new addon
//
// swagger:model
type AddonsPost struct {
	// Name of the addon
	// Example: my-addon
	Name string `json:"name" yaml:"name"`
}

// AddonPatch allows updating an existing addon with a new version
//
// swagger:model
type AddonPatch struct{}
