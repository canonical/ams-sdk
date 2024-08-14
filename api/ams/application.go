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

import (
	"fmt"
	"regexp"

	"github.com/anbox-cloud/ams-sdk/pkg/ams/constants"
	"github.com/anbox-cloud/ams-sdk/pkg/ams/shared"
	"github.com/anbox-cloud/ams-sdk/pkg/ams/shared/errors"
)

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
	// VideoEncoderTypeVPU describes a VPU based video encoder
	VideoEncoderTypeVPU VideoEncoderType = "vpu"
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
	case "vpu":
		return VideoEncoderTypeVPU
	}
	return VideoEncoderTypeUnknown
}

// ApplicationAddon describes a specific version of an addon an application version uses
//
// swagger:model
type ApplicationAddon struct {
	// Name of the application addon
	// Example: ssh
	Name string `json:"name" yaml:"name"`
	// Version of the application addon
	// Example: 0
	Version int `json:"version" yaml:"version"`
}

// ApplicationVersion describes a single version of an application
//
// swagger:model
type ApplicationVersion struct {
	// Version of the application
	// Example: 0
	Number int `json:"number" yaml:"number"`
	// Version of the manifest used to create the application
	// Example: 0.1
	ManifestVersion string `json:"manifest_version" yaml:"manifest-version"`
	// Parent image used for the application version
	// Example: btavtegj1qm58qg7ru50
	ParentImageID string `json:"parent_image_id" yaml:"parent_image_id"`
	// Parent image version to use for the application version
	// Example: 0
	ParentImageVersion int `json:"parent_image_version" yaml:"parent_image_version"`
	// Status of the application as an integer value
	// Example: 3
	StatusCode ImageStatus `json:"status_code" yaml:"status_code"`
	// Current status of the version
	// Enum: initializing,unknown,ready,deleted,error
	// Example: deleted
	Status string `json:"status" yaml:"status"`
	// Whether or not the version is published
	// Example: false
	Published bool `json:"published" yaml:"published"`
	// Creation UTC timestamp of the version
	// Example: 1532150640
	CreatedAt int64 `json:"created_at" yaml:"created_at"`
	// Name of the boot package for the version
	// Example: com.foo.bar.MainActivity
	BootActivity string `json:"boot_activity" yaml:"boot-activity"`
	// Required android application permissions
	// Example: ["android.permission.WRITE_EXTERNAL_STORAGE","android.permission.READ_EXTERNAL_STORAGE"]
	RequiredPermissions []string `json:"required_permissions" yaml:"required_permissions"`
	// List of addons enabled for the version
	Addons []ApplicationAddon `json:"addons" yaml:"addons"`
	// Extra data to be setup on the instance on application creation
	// Example: {}
	ExtraData map[string]ApplicationExtraData `json:"extra_data" yaml:"extra_data"`
	// Error message in case the application ran into an error
	// Example: {}
	ErrorMessage string `json:"error_message" yaml:"error_message"`
	// Encoder type to use for the application
	// Enum: [ gpu, gpu-preferred, software, unknown ]
	// Example: gpu
	VideoEncoder VideoEncoderType `json:"video_encoder,omitempty" yaml:"video-encoder,omitempty"`
	// Watchdog settings for the application
	Watchdog ApplicationWatchdog `json:"watchdog" yaml:"watchdog"`
	// List of services exposed by the application that should be exposed on the instance
	Services []NetworkServiceSpec `json:"services,omitempty" yaml:"services,omitempty"`
	// List of features supported by the application
	// Example: ["feature1", "feature2"]
	Features []string `json:"features" yaml:"features"`
	// Hook settings for the application
	Hooks ApplicationHooks `json:"hooks,omitempty" yaml:"hooks,omitempty"`
	// Boostrap settings for the application
	Bootstrap ApplicationBootstrap `json:"bootstrap,omitempty" yaml:"bootstrap,omitempty"`
}

// ApplicationExtraData represents an extra application data
//
// swagger:model
type ApplicationExtraData struct {
	// Path to the target file on the android filesystem
	// Example: /sdcard/Android/obb/com.foo.bar/
	Target string `json:"target" yaml:"target"`
	// Owner and group for the files
	// Example: root:root
	Owner string `json:"owner" yaml:"owner"`
	// Unix permissions for the files
	// Example: 0644
	Permissions string `json:"permissions" yaml:"permissions"`
}

// ApplicationResources describes resources allocated for an application
//
// swagger:model
type ApplicationResources struct {
	// Number of CPUs required by the application
	// Example: 2
	CPUs int `json:"cpus,omitempty" yaml:"cpus,omitempty"`
	// Memory to be assigned to the application
	// Example: 3GB
	Memory string `json:"memory,omitempty" yaml:"memory,omitempty"`
	// Storage size required by the application
	// Example: 3GB
	DiskSize string `json:"disk-size,omitempty" yaml:"disk-size,omitempty"`
	// Number of GPU slots required by the application
	// Example: 2
	GPUSlots int `json:"gpu-slots,omitempty" yaml:"gpu-slots,omitempty"`
	// Number of VPU slots required by the application
	// Example: 1
	VPUSlots int `json:"vpu-slots,omitempty" yaml:"vpu-slots,omitempty"`
}

// ToApplicationResources returns a valid application resource from an application resource patch
func (a *ApplicationResourcesPost) ToApplicationResources() ApplicationResources {
	// NOTE: GPUSlots or VPUSlots = 0 is a valid resource option, which means no gpu
	// or vpu slot will be plugged for the container launching from the application
	resources := ApplicationResources{GPUSlots: -1, VPUSlots: -1}
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
	if a.VPUSlots != nil && *a.VPUSlots >= 0 {
		resources.VPUSlots = *a.VPUSlots
	}
	return resources
}

// ApplicationResourcesPost represents the fields used to update an application resource
//
// swagger:model
type ApplicationResourcesPost struct {
	// Number of CPUs required by the application
	// Example: 2
	CPUs *int `json:"cpus,omitempty" yaml:"cpus,omitempty"`
	// Memory to be assigned to the application
	// Example: 3GB
	Memory *string `json:"memory,omitempty" yaml:"memory,omitempty"`
	// Storage size required by the application
	// Example: 3GB
	DiskSize *string `json:"disk-size,omitempty" yaml:"disk-size,omitempty"`
	// Number of GPU slots required by the application
	// Example: 2
	GPUSlots *int `json:"gpu-slots,omitempty" yaml:"gpu-slots,omitempty"`
	// Number of VPU slots required by the application
	// Example: 1
	VPUSlots *int `json:"vpu-slots,omitempty" yaml:"vpu-slots,omitempty"`
}

// ApplicationWatchdog describes the fields used to update an application watchdog
//
// swagger:model
type ApplicationWatchdog struct {
	// Whether or not to enable the watchdog for the application
	// Example: true
	Disabled bool `json:"disabled" yaml:"disabled"`
	// List of android packages to enable the watchdog for
	// Example: ["com.android.settings"]
	AllowedPackages []string `json:"allowed-packages" yaml:"allowed-packages"`
}

// ValidateAllowedPackages checks the value for given allowed packages
func (a ApplicationWatchdog) ValidateAllowedPackages() error {
	if len(a.AllowedPackages) > 0 {
		if shared.StringInSlice("*", a.AllowedPackages) {
			if len(a.AllowedPackages) > 1 {
				return fmt.Errorf("specifying an asterisk ('*') and package names is not possible as they are mutually exclusive")
			}
		} else {
			for _, allowPackage := range a.AllowedPackages {
				if match, _ := regexp.MatchString(constants.AndroidPackageNamePattern, allowPackage); !match {
					return errors.NewInvalidArgument(fmt.Sprintf("'allowed-packages': %s", allowPackage))
				}
			}
		}
	}
	return nil
}

// ApplicationHooks describes the fields used to configure the hooks of an application
type ApplicationHooks struct {
	// Time limit to wait for the hook to complete before timing out
	// Example: 10m
	Timeout string `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

// ApplicationBootstrap describes the fields used to configure the application bootstrap
type ApplicationBootstrap struct {
	// List of files to keep after the bootstrap
	// Example: ["app.apk"]
	Keep []string `json:"keep,omitempty" yaml:"keep,omitempty"`
}

// Application represents an AMS application
//
// swagger:model
type Application struct {
	// ID of the application
	// Example: btavtegj1qm58qg7ru50
	ID string `json:"id" yaml:"id"`
	// Name of the application
	// Example: my-app
	Name string `json:"name" yaml:"name"`
	// Status of the application as an integer value
	// Example: 3
	StatusCode ApplicationStatus `json:"status_code" yaml:"status_code"`
	// Current status of the application
	// Enum: initializing,unknown,ready,deleted,error
	// Example: deleted
	Status string `json:"status" yaml:"status"`
	// Instance type required by the application
	// Example: a2.3
	InstanceType string `json:"instance_type" yaml:"instance_type"`
	// Name of the boot package for the application
	// Example: com.foo.bar
	BootPackage string `json:"boot_package" yaml:"boot_package"`
	// Parent image used for the application
	// Example: btavtegj1qm58qg7ru50
	ParentImageID string `json:"parent_image_id" yaml:"parent_image_id"`
	// Whether or not the application is published
	// Example: false
	Published bool `json:"published" yaml:"published"`
	// List of versions for the application
	Versions []ApplicationVersion `json:"versions" yaml:"versions"`
	// List of addons enabled for the application
	// Example: ["ssh", "gms"]
	Addons []string `json:"addons" yaml:"addons"`
	// Creation UTC timestamp of the application
	// Example: 1532150640
	CreatedAt int64 `json:"created_at" yaml:"created_at"`
	// Flag to show whether the application can be edited
	// Example: false
	Immutable bool `json:"immutable" yaml:"immutable"`
	// Tags to attach to the application
	// Example: ["created_by=anbox"]
	Tags []string `json:"tags" yaml:"tags"`
	// Resources required by the application. Overrides the instance_type requirements.
	Resources ApplicationResources `json:"resources,omitempty" yaml:"resources,omitempty"`
	// ABI supported by the application
	// Example: x86_64
	ABI string `json:"abi,omitempty" yaml:"abi,omitempty"`
	// Whether or not to auto update the application's base image
	// Example: false
	InhibitAutoUpdates bool `json:"inhibit_auto_updates" yaml:"inhibit_auto_updates"`
	// List of tags for filtering the nodes to run the application on
	// Example: ["gpu=nvidia", "cpu=intel"]
	NodeSelector []string `json:"node_selector" yaml:"node_selector"`
	// Whether the application is based on virtual machines or containers
	VM bool `json:"vm" yaml:"vm"`
}

// GetApplicationFilters returns an array of attributes available on the api to
// filter applications
func GetApplicationFilters() []string {
	return []string{
		"id",
		"name",
		"status",
		"instance_type",
		"boot_package",
		"published",
		"immutable",
		"abi",
		"addons",
		"inhibit_auto_updates",
		"tags",
	}
}

// ApplicationVersionPatch represents the fields used to update an application version
//
// swagger:model
type ApplicationVersionPatch struct {
	// Used to publish a specific version of the application
	// Example: true
	Published *bool `json:"published" yaml:"published"`
}

// ApplicationPatch represents the fields a user can modify
//
// swagger:model
type ApplicationPatch struct {
	// Base image id or name to use for the applicaiton
	// Example: btavtegj1qm58qg7ru50
	Image *string `json:"image" yaml:"image"`
	// Instance type to use for the application
	// Example: a3.4
	InstanceType *string `json:"instance-type" yaml:"instance-type"`
	// Tags to attach to the application
	// Example: ["created_by=anbox"]
	Tags *[]string `json:"tags" yaml:"tags"`
	// List of addons enabled for the application
	// Example: ["ssh", "gms"]
	Addons *[]string `json:"addons" yaml:"addons"`
	// Resources required by the application. Overrides the instance_type requirements.
	Resources *ApplicationResourcesPost `json:"resources,omitempty" yaml:"resources,omitempty"`
	// Whether or not to auto update the application's base image
	// Example: false
	InhibitAutoUpdates *bool `json:"inhibit_auto_updates" yaml:"inhibit_auto_updates"`
	// List of services exposed by the application that should be expose on the instance
	// For application version update, changing those values would trigger a new application version creation
	Services *[]NetworkServiceSpec `json:"services,omitempty" yaml:"services,omitempty"`
	// Watchdog settings for the application
	Watchdog *ApplicationWatchdog `json:"watchdog" yaml:"watchdog"`
	// Name of the boot package for the version
	// Example: com.foo.bar.MainActivity
	BootActivity *string `json:"boot_activity" yaml:"boot-activity"`
	// Required android application permissions
	// Example: ["android.permission.WRITE_EXTERNAL_STORAGE","android.permission.READ_EXTERNAL_STORAGE"]
	RequiredPermissions *[]string `json:"required_permissions" yaml:"required_permissions"`
	// Encoder type to use for the application
	// Enum: [ gpu, gpu-preferred, software, unknown ]
	// Example: gpu
	VideoEncoder *VideoEncoderType `json:"video_encoder,omitempty" yaml:"video-encoder,omitempty"`
	// Version of the manifest used to create the application
	// Example: 0.1
	ManifestVersion *string `json:"manifest_version" yaml:"manifest-version"`
	// List of features supported by the application
	// Example: ["feature1", "feature2"]
	Features *[]string `json:"features" yaml:"features"`
	// Hook settings for the application
	Hooks *ApplicationHooks `json:"hooks,omitempty" yaml:"hooks,omitempty"`
	// Bootstrap settings for the application
	Bootstrap *ApplicationBootstrap `json:"bootstrap,omitempty" yaml:"bootstrap,omitempty"`
	// List of tags for filtering the nodes to run the application on
	// Example: ["gpu=nvidia", "cpu=intel"]
	NodeSelector *[]string `json:"node_selector,omitempty" yaml:"node-selector,omitempty"`
}

// ApplicationDelete represents the fields used to delete an application
//
// swagger:model
type ApplicationDelete struct {
	// Whether deletion of the application should be forced
	// Example: false
	Force bool `json:"force" yaml:"force"`
}

// ApplicationsDelete represents a list of application to delete together
//
// swagger:model
//
// API extensions: bulk_delete.applications
type ApplicationsDelete struct {
	// IDs or names of the applications to delete
	// Example: ["cilsreunfpfec9b1ktg0", "cilsreunfpfec9b1ktg1", "myapp"]
	IDs []string `json:"ids" yaml:"ids"`
	// Whether deletion of the applications should be forced
	// Example: false
	Force bool `json:"force" yaml:"force"`
}

// ApplicationVersionDelete represents the fields used to delete an application version
//
// swagger:model
type ApplicationVersionDelete struct {
	// Whether deletion of the application version should be forced
	// Example: false
	Force bool `json:"force" yaml:"force"`
}
