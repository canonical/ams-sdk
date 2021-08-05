// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

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
)

func (s ImageStatus) String() string {
	switch s {
	case ImageStatusError:
		return "error"
	case ImageStatusCreated:
		return "created"
	case ImageStatusActive:
		return "active"
	case ImageStatusInitializing:
		return "initializing"
	}
	return "unknown"
}

// ImageVersion describes a single version of an image
type ImageVersion struct {
	Number      int         `json:"version" yaml:"version"`
	Fingerprint string      `json:"fingerprint" yaml:"fingerprint"`
	Size        int64       `json:"size" yaml:"size"`
	CreatedAt   int64       `json:"created_at" yaml:"upload_time"`
	StatusCode  ImageStatus `json:"status_code" yaml:"status_code"`
	Status      string      `json:"status" yaml:"status"`
	RemoteID    string      `json:"remote_id" yaml:"remote_id"`
}

// Image represents an image available in AMS
type Image struct {
	ID           string         `json:"id" yaml:"id"`
	Name         string         `json:"name" yaml:"name"`
	Versions     []ImageVersion `json:"versions" yaml:"versions"`
	StatusCode   ImageStatus    `json:"status_code" yaml:"status_code"`
	Status       string         `json:"status" yaml:"status"`
	UsedBy       []string       `json:"used_by" yaml:"used_by"`
	Immutable    bool           `json:"immutable" yaml:"immutable"`
	Default      bool           `json:"default" yaml:"default"`
	Architecture string         `json:"architecture,omitempty" yaml:"architecture,omitempty"`
}

// ImagesPost represents the fields to upload a new image
type ImagesPost struct {
	Name    string `json:"name" yaml:"name"`
	Path    string `json:"path" yaml:"path"`
	Default bool   `json:"default" yaml:"default"`
}

// ImagePatch represents the fields to update an existing image
type ImagePatch struct {
	Default *bool `json:"default" yaml:"default"`
}

// ImagesGet represents a list of images
type ImagesGet struct {
	Images []Image `json:"images" yaml:"images"`
}

// ImageDelete describes a request used to delete an image
type ImageDelete struct {
	Force bool `json:"force"`
}
