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

// ImageStatus represents the status of an image
type ImageStatus int

const (
	// ImageStatusError represents the state an image is in when an error occurred
	// during a operation
	ImageStatusError ImageStatus = -1
	// ImageStatusUnknown represents the state an image is in when its real state
	// cannot be determined.
	ImageStatusUnknown ImageStatus = 0
	// ImageStatusInitializing represents the state an image is in when its currently being created
	ImageStatusInitializing ImageStatus = 1
	// ImageStatusCreated represents the state an image is in when it was uploaded to
	// AMS but not yet available on all LXD nodes.
	ImageStatusCreated ImageStatus = 2
	// ImageStatusActive represents the state an image is in when it is available on
	// all LXD nodes.
	ImageStatusActive ImageStatus = 3
	// ImageStatusDeleted represents the state an image is in when it is currently being deleted
	ImageStatusDeleted ImageStatus = 4
	// ImageStatusAvailable represents the state when an image is present on the remote but not in the LXD cluster
	ImageStatusAvailable ImageStatus = 5
)

func (s ImageStatus) String() string {
	switch s {
	case ImageStatusError:
		return "error"
	case ImageStatusCreated:
		return "created"
	case ImageStatusAvailable:
		return "available"
	case ImageStatusActive:
		return "active"
	case ImageStatusInitializing:
		return "initializing"
	case ImageStatusDeleted:
		return "deleted"
	}
	return "unknown"
}

// ImageVersion describes a single version of an image
//
// swagger:model
type ImageVersion struct {
	// Version for the image
	// Example: 0
	Number int `json:"version" yaml:"version"`
	// Fingerprint of the image version
	// Example: 0791cfc011f67c60b7bd0f852ddb686b79fa46083d9d43ef9845c9235c67b261
	Fingerprint string `json:"fingerprint" yaml:"fingerprint"`
	// Size (in bytes) of the image version
	// Example: 529887868
	Size int64 `json:"size" yaml:"size"`
	// Creation UTC timestamp of the image
	// Example: 1610641117
	CreatedAt int64 `json:"created_at" yaml:"upload_time"`
	// Status of the image as an integer value
	// Example: 3
	StatusCode ImageStatus `json:"status_code" yaml:"status_code"`
	// Current status of the image
	// Enum: error,created,active,initializing,unknown
	// Example: active
	Status string `json:"status" yaml:"status"`
	// Version of the image in the remote server
	// Example: 1.2.3
	RemoteID string `json:"remote_id" yaml:"remote_id"`
	// ErrorMessage describes a potential error which has occured while
	// downloading the image.
	ErrorMessage string `json:"error_message,omitempty" yaml:"error_message,omitempty"`
}

// ImageType specifies the type of an image
type ImageType string

const (
	// ImageTypeAny is used when the type of the image is not relevant
	ImageTypeAny ImageType = ""
	// ImageTypeUnknown is returned when the type of the image is not known
	ImageTypeUnknown ImageType = "unknown"
	// ImageTypeContainer specifies that the image is used for containers
	ImageTypeContainer ImageType = "container"
	// ImageTypeVM specifies that the image is used for virtual machines
	ImageTypeVM ImageType = "vm"
)

// Image represents an image available in AMS
//
// swagger:model
type Image struct {
	// ID of the image
	// Example: btavtegj1qm58qg7ru50
	ID string `json:"id" yaml:"id"`
	// Name of the image
	// Example: my-image
	Name string `json:"name" yaml:"name"`
	// List of versions for the image
	Versions []ImageVersion `json:"versions" yaml:"versions"`
	// Status of the image as an integer value
	// Example: 3
	StatusCode ImageStatus `json:"status_code" yaml:"status_code"`
	// Current status of the image
	// Enum: error,created,active,initializing,unknown
	// Example: active
	Status string `json:"status" yaml:"status"`
	// List of application ids using the image as a base
	// Example: ["btavtegj1qm58asf123"]
	UsedBy []string `json:"used_by" yaml:"used_by"`
	// Flag to show whether the image can be edited by an AMS instance or not
	// Example: false
	Immutable bool `json:"immutable" yaml:"immutable"`
	// Flag to show whether the image is used by default if no image name is provided
	// Example: false
	Default bool `json:"default" yaml:"default"`
	// CPU architecture supported by the image
	// Example: x86_64
	Architecture string `json:"architecture,omitempty" yaml:"architecture,omitempty"`
	// Type of the image. Possible values are: container, vm
	Type ImageType `json:"type" yaml:"type"`
	// Variant of the image. Possible values are: android, aaos, generic, unknown
	Variant string `json:"variant" yaml:"variant"`
}

// ImagesPost represents the fields to upload a new image
//
// swagger:model
type ImagesPost struct {
	// Name of the image
	// Example: my-image
	Name string `json:"name" yaml:"name"`
	// Path to store the image
	// Example: /save/image
	Path string `json:"path" yaml:"path"`
	// Make the image as default
	// Example: false
	Default bool `json:"default" yaml:"default"`
	// Type specifies if the type of the to be imported image. Only valid
	// when no image payload is provided with the request and the image
	// is meant to be imported from a remote image server. If not specified
	// all available image types will be imported.
	Type ImageType `json:"type" yaml:"type"`
}

// ImagePatch represents the fields to update an existing image
//
// swagger:model
type ImagePatch struct {
	// Make the image as default
	// Example: true
	Default *bool `json:"default" yaml:"default"`

	// ForceSync forces synchronization of the image from the remote image server
	// Examle: true
	ForceSync bool `json:"force_sync" yaml:"force_sync"`
}

// ImagesGet represents a list of images
//
// swagger:model
type ImagesGet struct {
	Images []Image `json:"images" yaml:"images"`
}

// ImageDelete describes a request used to delete an image
//
// swagger:model
type ImageDelete struct {
	Force bool `json:"force"`
}
