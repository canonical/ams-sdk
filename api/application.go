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

// ApplicationStatus represents the status an application can be in
type ApplicationStatus int

const (
	// ApplicationStatusError represents the state when an application encountered an error.
	ApplicationStatusError ApplicationStatus = -1
	// ApplicationStatusUnknown represents the state when the real state of an application
	// cannot be determined.
	ApplicationStatusUnknown ApplicationStatus = 0
	// ApplicationStatusInitializing represents the state when an application was created
	// and is correctly being created.
	ApplicationStatusInitializing ApplicationStatus = 1
	// ApplicationStatusReady represents the state when an application was successfully
	// created and is ready to be used.
	ApplicationStatusReady ApplicationStatus = 2
	// ApplicationStatusDeleted represents the status an application has when it is currently
	// being deleted
	ApplicationStatusDeleted ApplicationStatus = 3
)

func (s *ApplicationStatus) String() string {
	switch *s {
	case ApplicationStatusError:
		return "error"
	case ApplicationStatusInitializing:
		return "initializing"
	case ApplicationStatusReady:
		return "ready"
	case ApplicationStatusDeleted:
		return "deleted"
	}
	return "unknown"
}

// VideoEncoderType describes the type of a video encoder is used by an application
type VideoEncoderType string

const (
	// VideoEncoderTypeGPU describes a GPU-based video encoder
	VideoEncoderTypeGPU VideoEncoderType = "gpu"
	// VideoEncoderTypeGPUPreferred describes a gpu preferred video encoder if the video
	// encoder slots is not filled up, otherwise it fallbacks to software video encoder.
	VideoEncoderTypeGPUPreferred VideoEncoderType = "gpu-preferred"
	// VideoEncoderTypeSoftware describes a software-based video encoder
	VideoEncoderTypeSoftware VideoEncoderType = "software"
	// VideoEncoderTypeUnknown describes a unknown video encoder
	VideoEncoderTypeUnknown VideoEncoderType = "unknown"
)

// VideoEncoderFromString return a video encode type from the given value
func VideoEncoderFromString(value string) VideoEncoderType {
	switch value {
	case "gpu":
		return VideoEncoderTypeGPU
	case "gpu-preferred":
		return VideoEncoderTypeGPUPreferred
	case "software":
		return VideoEncoderTypeSoftware
	}
	return VideoEncoderTypeUnknown
}

// ApplicationAddon describes a specific version of an addon an application version uses
type ApplicationAddon struct {
	Name    string `json:"name" yaml:"name"`
	Version int    `json:"version" yaml:"version"`
}

// ApplicationVersion describes a single version of an application
type ApplicationVersion struct {
	Number              int                             `json:"number" yaml:"number"`
	ManifestVersion     string                          `json:"manifest_version" yaml:"manifest-version"`
	ParentImageID       string                          `json:"parent_image_id" yaml:"parent_image_id"`
	ParentImageVersion  int                             `json:"parent_image_version" yaml:"parent_image_version"`
	StatusCode          ImageStatus                     `json:"status_code" yaml:"status_code"`
	Status              string                          `json:"status" yaml:"status"`
	Published           bool                            `json:"published" yaml:"published"`
	CreatedAt           int64                           `json:"created_at" yaml:"created_at"`
	BootActivity        string                          `json:"boot_activity" yaml:"boot-activity"`
	RequiredPermissions []string                        `json:"required_permissions" yaml:"required_permissions"`
	Addons              []ApplicationAddon              `json:"addons" yaml:"addons"`
	ExtraData           map[string]ApplicationExtraData `json:"extra_data" yaml:"extra_data"`
	ErrorMessage        string                          `json:"error_message" yaml:"error_message"`
	VideoEncoder        VideoEncoderType                `json:"video_encoder,omitempty" yaml:"video-encoder,omitempty"`
	Watchdog            ApplicationWatchdog             `json:"watchdog" yaml:"watchdog"`
	Services            []NetworkServiceSpec            `json:"services,omitempty" yaml:"services,omitempty"`
	Features            []string                        `json:"features" yaml:"features"`
	Hooks               ApplicationHooks                `json:"hooks,omitempty" yaml:"hooks,omitempty"`
	Bootstrap           ApplicationBootstrap            `json:"bootstrap,omitempty" yaml:"bootstrap,omitempty"`
}

// ApplicationExtraData represents an extra application data
type ApplicationExtraData struct {
	Target      string `json:"target" yaml:"target"`
	Owner       string `json:"owner" yaml:"owner"`
	Permissions string `json:"permissions" yaml:"permissions"`
}

// ApplicationResources describes resources allocated for an application
type ApplicationResources struct {
	CPUs     int    `json:"cpus,omitempty" yaml:"cpus,omitempty"`
	Memory   string `json:"memory,omitempty" yaml:"memory,omitempty"`
	DiskSize string `json:"disk-size,omitempty" yaml:"disk-size,omitempty"`
	GPUSlots int    `json:"gpu-slots,omitempty" yaml:"gpu-slots,omitempty"`
}

// ToApplicationResources returns a valid application resource from an application resource patch
func (a *ApplicationResourcesPost) ToApplicationResources() ApplicationResources {
	// NOTE: GPUSlots = 0 is a valid resource option, which means no gpu
	// slot will be plugged for the container launching from the application
	resources := ApplicationResources{GPUSlots: -1}
	if a.CPUs != nil && *a.CPUs > 0 {
		resources.CPUs = *a.CPUs
	}
	if a.Memory != nil && len(*a.Memory) > 0 {
		resources.Memory = *a.Memory
	}
	if a.DiskSize != nil && len(*a.DiskSize) > 0 {
		resources.DiskSize = *a.DiskSize
	}
	if a.GPUSlots != nil && *a.GPUSlots >= 0 {
		resources.GPUSlots = *a.GPUSlots
	}
	return resources
}

// ApplicationResourcesPost represents the fields used to update an application resource
type ApplicationResourcesPost struct {
	CPUs     *int    `json:"cpus,omitempty" yaml:"cpus,omitempty"`
	Memory   *string `json:"memory,omitempty" yaml:"memory,omitempty"`
	DiskSize *string `json:"disk-size,omitempty" yaml:"disk-size,omitempty"`
	GPUSlots *int    `json:"gpu-slots,omitempty" yaml:"gpu-slots,omitempty"`
}

// ApplicationWatchdog describes the fields used to update an application watchdog
type ApplicationWatchdog struct {
	Disabled        bool     `json:"disabled" yaml:"disabled"`
	AllowedPackages []string `json:"allowed-packages" yaml:"allowed-packages"`
}

// ApplicationHooks describes the fields used to configure the hooks of an application
type ApplicationHooks struct {
	Timeout string `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

// ApplicationBootstrap describes the fields used to configure the application bootstrap
type ApplicationBootstrap struct {
	Keep []string `json:"keep,omitempty" yaml:"keep,omitempty"`
}

// Application represents an AMS application
type Application struct {
	ID                 string               `json:"id" yaml:"id"`
	Name               string               `json:"name" yaml:"name"`
	StatusCode         ApplicationStatus    `json:"status_code" yaml:"status_code"`
	Status             string               `json:"status" yaml:"status"`
	InstanceType       string               `json:"instance_type" yaml:"instance_type"`
	BootPackage        string               `json:"boot_package" yaml:"boot_package"`
	ParentImageID      string               `json:"parent_image_id" yaml:"parent_image_id"`
	Published          bool                 `json:"published" yaml:"published"`
	Versions           []ApplicationVersion `json:"versions" yaml:"versions"`
	Addons             []string             `json:"addons" yaml:"addons"`
	CreatedAt          int64                `json:"created_at" yaml:"created_at"`
	Immutable          bool                 `json:"immutable" yaml:"immutable"`
	Tags               []string             `json:"tags" yaml:"tags"`
	Resources          ApplicationResources `json:"resources,omitempty" yaml:"resources,omitempty"`
	ABI                string               `json:"abi,omitempty" yaml:"abi,omitempty"`
	InhibitAutoUpdates bool                 `json:"inhibit_auto_updates" yaml:"inhibit_auto_updates"`
	NodeSelector       []string             `json:"node_selector" yaml:"node_selector"`
}

// ApplicationVersionPatch represents the fields used to update an application version
type ApplicationVersionPatch struct {
	Published *bool `json:"published" yaml:"published"`
}

// ApplicationPatch represents the fields a user can modify
type ApplicationPatch struct {
	Image              *string                   `json:"image" yaml:"image"`
	InstanceType       *string                   `json:"instance-type" yaml:"instance-type"`
	Tags               *[]string                 `json:"tags" yaml:"tags"`
	Addons             *[]string                 `json:"addons" yaml:"addons"`
	Resources          *ApplicationResourcesPost `json:"resources,omitempty" yaml:"resources,omitempty"`
	InhibitAutoUpdates *bool                     `json:"inhibit_auto_updates" yaml:"inhibit_auto_updates"`
	// For application version update, changing those values would trigger a new application version creation
	Services            *[]NetworkServiceSpec `json:"services,omitempty" yaml:"services,omitempty"`
	Watchdog            *ApplicationWatchdog  `json:"watchdog" yaml:"watchdog"`
	BootActivity        *string               `json:"boot_activity" yaml:"boot-activity"`
	RequiredPermissions *[]string             `json:"required_permissions" yaml:"required_permissions"`
	VideoEncoder        *VideoEncoderType     `json:"video_encoder,omitempty" yaml:"video-encoder,omitempty"`
	ManifestVersion     *string               `json:"manifest_version" yaml:"manifest-version"`
	Features            *[]string             `json:"features" yaml:"features"`
	Hooks               *ApplicationHooks     `json:"hooks,omitempty" yaml:"hooks,omitempty"`
	Bootstrap           *ApplicationBootstrap `json:"bootstrap,omitempty" yaml:"bootstrap,omitempty"`
	NodeSelector        *[]string             `json:"node_selector,omitempty" yaml:"node-selector,omitempty"`
}

// ApplicationDelete represents the fields used to delete an application
type ApplicationDelete struct {
	Force bool `json:"force" yaml:"force"`
}

// ApplicationVersionDelete represents the fields used to delete an application version
type ApplicationVersionDelete struct {
	Force bool `json:"force" yaml:"force"`
}
