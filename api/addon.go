// -*- Mode: Go; indent-tabs-mode: t -*-
/*
 * This file is part of AMS SDK
 * Copyright 2021 Canonical Ltd.
 *
 * This program is free software: you can redistribute it and/or modify it under
 * the terms of the GNU Lesser General Public License version 3, as published
 * by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranties of MERCHANTABILITY, SATISFACTORY
 * QUALITY, or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public
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
type AddonVersion struct {
	Number      int    `json:"version" yaml:"version"`
	Fingerprint string `json:"fingerprint" yaml:"fingerprint"`
	Size        int64  `json:"size" yaml:"size"`
	CreatedAt   int64  `json:"created_at" yaml:"created_at"`
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
type Addon struct {
	Name     string         `json:"name" yaml:"name"`
	Versions []AddonVersion `json:"versions" yaml:"versions"`
	UsedBy   []string       `json:"used_by" yaml:"used_by"`
}

// AddonsPost is used to create a new addon
type AddonsPost struct {
	Name string `json:"name" yaml:"name"`
}

// AddonPatch allows updating an existing addon with a new version
type AddonPatch struct{}
