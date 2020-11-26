// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

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
