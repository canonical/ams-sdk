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

	"github.com/anbox-cloud/ams-sdk/shared"
)

// InstanceStatus represents the status a instance can be in
type InstanceStatus int

const (
	// InstanceStatusError represents the state a instance is in when an error occurred.
	InstanceStatusError InstanceStatus = -1
	// InstanceStatusUnknown represents the state a instance is in when its real state
	// cannot be determined.
	InstanceStatusUnknown InstanceStatus = 0
	// InstanceStatusCreated represents the status a instance is in when its object
	// was created on the store but no further operation was performed on it yet.
	InstanceStatusCreated InstanceStatus = 1
	// InstanceStatusPrepared represents the status a instance is in when it got all
	// its resources assigned and is ready to be constructed and started on the LXD cluster.
	InstanceStatusPrepared InstanceStatus = 2
	// InstanceStatusStarted represents the state a instance is in when it is currently
	// starting.
	InstanceStatusStarted InstanceStatus = 3
	// InstanceStatusRunning represents the state a instance is in when it is running
	InstanceStatusRunning InstanceStatus = 4
	// InstanceStatusStopped represents the state a instance is in when it is stopped
	InstanceStatusStopped InstanceStatus = 5
	// InstanceStatusDeleted represents the state a instance is currently being deleted
	InstanceStatusDeleted InstanceStatus = 6
)

func (s InstanceStatus) String() string {
	switch s {
	case InstanceStatusCreated:
		return "created"
	case InstanceStatusPrepared:
		return "prepared"
	case InstanceStatusStarted:
		return "started"
	case InstanceStatusStopped:
		return "stopped"
	case InstanceStatusRunning:
		return "running"
	case InstanceStatusError:
		return "error"
	case InstanceStatusDeleted:
		return "deleted"
	}
	return "unknown"
}

// InstanceService describes a single service the instance exposes to the outside world.
//
// While NetworkServiceSpec defines what the user requests, InstanceService is what is actually
// opened on the instance.
//
// swagger:model
type InstanceService struct {
	// Port is the port the instance provides a service on
	// Example: 3000
	Port int `json:"port,omitempty" yaml:"port,omitempty"`
	// PortEnd, if specified, denotes the end of the port range starting at Port
	// Example: 3010
	PortEnd int `json:"port_end,omitempty" yaml:"port_end,omitempty"`
	// NodePort is the port used on the LXD node to map to the service port
	// If left empty the node port is automatically selected.
	// Example: 4000
	NodePort *int `json:"node_port,omitempty" yaml:"node_port,omitempty"`
	// NodePortEnd, if specified, denotes the end of the port range on the node starting
	// at NodePort
	// Example: 4010
	NodePortEnd *int `json:"node_port_end,omitempty" yaml:"node_port_end,omitempty"`
	// List of network protocols (tcp, udp) the port should be exposed for
	// Example: ["tcp", "udp"]
	Protocols []NetworkProtocol `json:"protocols" yaml:"protocols"`
	// Expose defines wether the service is exposed on the public endpoint of the node
	// or if it is only available on the private endpoint. To expose the service set to
	// true and to false otherwise.
	Expose bool `json:"expose" yaml:"expose"`
	// Name gives the instance a hint what the exposed port is being used for. This
	// allows further tweaks inside the instance to expose the service correctly.
	// Example: myservice
	Name string `json:"name" yaml:"name"`
}

// InstanceType describes the type of an instance (container, vm)
type InstanceType string

const (
	// InstanceTypeAny specifies that the type of the instance does not matter
	InstanceTypeAny InstanceType = ""
	// InstanceTypeContainer specifies that the instance is a container
	InstanceTypeContainer InstanceType = "container"
	// InstanceTypeVM specifies that the instance is a VM
	InstanceTypeVM InstanceType = "vm"
)

// InstanceResources represents resources assigned to an instance
//
// swagger:model
type InstanceResources struct {
	// CPUs cores assigned to the instance
	// Example: 2
	CPUs int `json:"cpus,omitempty" yaml:"cpus,omitempty"`
	// Memory assigned to the instance in bytes
	// Example: 3221225472
	Memory int64 `json:"memory,omitempty" yaml:"memory,omitempty"`
	// DiskSize specifies the amount of storage assigned to the instance in bytes
	// Example: 3221225472
	DiskSize int64 `json:"disk-size,omitempty" yaml:"disk-size,omitempty"`
	// GPUSlots specifies the number of GPU slots the instance got allocated
	// Example: 1
	GPUSlots int `json:"gpu-slots,omitempty" yaml:"gpu-slots,omitempty"`
	// VPUSlots specifies the number of VPU slots the instance
	VPUSlots int `json:"vpu-slots,omitempty" yaml:"vpu-slots,omitempty"`
}

// Instance represents a single instance
//
// swagger:model
type Instance struct {
	// ID of the instance
	// Example: cilsreunfpfec9b1ktg0
	ID string `json:"id" yaml:"id"`
	// Name of the instance. Typically in the format "ams-<ID>".
	// Example: ams-cilsreunfpfec9b1ktg0
	Name string `json:"name" yaml:"name"`
	// If this is a base instance or not
	IsBase bool `json:"base" yaml:"base"`
	// Type specifies the type of the instance (container, vm)
	Type InstanceType `json:"type" yaml:"type"`
	// StatusCode of the instance. Matches the Status field.
	// Example: 4
	StatusCode InstanceStatus `json:"status_code" yaml:"status_code"`
	// Status of the instance
	// Example: running
	Status string `json:"status" yaml:"status"`
	// Node the instance is running on
	// Example: lxd0
	Node string `json:"node" yaml:"node"`
	// AppID is the ID of the application the instance is created from. Empty if the
	// instance has not been created from an application.
	// Example: cilsiomnfpfec9b1kteg
	AppID string `json:"app_id" yaml:"app_id"`
	// AppName is the name of the application the instance is created from. Empty if the
	// instance has not been created from an application.
	// Example: myapp
	AppName string `json:"app_name" yaml:"app_name"`
	// AppVersion is the version of the application the instance is created from. Empty if the
	// instance has not been created from an application.
	// Example: 0
	AppVersion int `json:"app_version" yaml:"app_version"`
	// ImageID is the ID of the image the instance is created from. Empty if the
	// instance has not been created from an image.
	// Example: cilshrmnfpfec9b1kte0
	ImageID string `json:"image_id" yaml:"image_id"`
	// ImageVersion is the version of the image the instance is created from. Empty if the
	// instance has not been created from an image.
	// Example: 0
	ImageVersion int `json:"image_version" yaml:"image_version"`
	// CreatedAt specifies the time at which the instance was created
	// Example: 1689604498
	CreatedAt int64 `json:"created_at" yaml:"created_at"`
	// Address is the IP address of the instance
	// Example: 192.168.1.74
	Address string `json:"address" yaml:"address"`
	// PublicAddress is the external IP address the instance is accessible on (in most
	// cases the IP of the node it is running on)
	// Example: 1.2.3.4
	PublicAddress string `json:"public_address" yaml:"public_address"`
	// Services the instance exposes
	Services []InstanceService `json:"services" yaml:"services"`
	// StoredLogs lists log files AMS stores for the instance.
	// Example: ["android.log", "system.log"]
	StoredLogs []string `json:"stored_logs" yaml:"stored_logs"`
	// ErrorMessage provides an error message when the instance status is set to error.
	// Example: instance failed to boot
	ErrorMessage string `json:"error_message" yaml:"error_message"`
	// StatusMessage describes the current status of the instance
	// Example: "Waiting for image download"
	StatusMessage string `json:"status_message" yaml:"status_message"`
	// Config summarizes the configuration the instance uses
	Config struct {
		// Platform specifies the Anbox platform the instance is running with
		// Example: webrtc
		Platform string `json:"platform,omitempty" yaml:"platform,omitempty"`
		// BootPackage specifies the Android application package name which is started by default
		// Example: com.android.settings
		BootPackage string `json:"boot_package,omitempty" yaml:"boot_package,omitempty"`
		// BootActivity specifies the Android activity which is started by default
		// Example: com.android.settings/.DevSettings
		BootActivity string `json:"boot_activity,omitempty" yaml:"boot_activity,omitempty"`
		// MetricsServer specifies a metrics server the instance will use
		// Example: 10.0.0.45:8086
		MetricsServer string `json:"metrics_server,omitempty" yaml:"metrics_server,omitempty"`
		// DisableWatchdog defines whether the watchdog is disabled
		DisableWatchdog bool `json:"disable_watchdog,omitempty" yaml:"disable_watchdog,omitempty"`
		// DevMode specifies if development mode has been turned on for the instance
		DevMode bool `json:"devmode,omitempty" yaml:"devmode,omitempty"`
	} `json:"config,omitempty"`
	// Resources specifies the resources allocated for the instance
	Resources InstanceResources `json:"resources,omitempty"`
	// Architecture describes the CPU archtitecture the instance is using
	// Example: aarch64
	Architecture string `json:"architecture,omitempty" yaml:"architecture,omitempty"`
	// Tags specifies the tags the instance has assigned
	// Example: ["foo", "bar"]
	Tags []string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// GetInstanceFilters returns an array of attributes available on the api to
// filter instances
func GetInstanceFilters() []string {
	return []string{
		"id",
		"name",
		"status",
		"type",
		"base",
		"node",
		"app_id",
		"app_version",
		"image_id",
		"image_version",
		"app_name",
		"tags",
	}
}

// InstancesPost represents the fields required to launch a new instance for
// a specific application
//
// swagger:model
type InstancesPost struct {
	// Type of the instance (container, vm)
	Type InstanceType `json:"type" yaml:"type"`
	// ID of the application to use. Can be empty if an image ID is specified instead
	// Example: cilsiomnfpfec9b1kteg
	ApplicationID string `json:"app_id" yaml:"app_id"`
	// Version of the application to use. If not specified, the latest version is used.
	// Example: 0
	ApplicationVersion *int `json:"app_version" yaml:"app_version"`
	// ID of the image to use. Can be empty if an application ID is specified instead.
	// Example: cilshrmnfpfec9b1kte0
	ImageID string `json:"image_id" yaml:"image_id"`
	// Version of the image to use. If not specified, the latest version is used.
	// Example: 0
	ImageVersion *int `json:"image_version" yaml:"image_version"`
	// Node to start the instance on. If empty node will be automatically selected.
	// Example: lxd0
	Node string `json:"node" yaml:"node"`
	// User data to pass to the instance.
	// Example: {\"key\":\"value\"}
	Userdata *string `json:"user_data,omitempty" yaml:"user_data,omitempty"`
	// Addons to enable for the instance
	// Example: ["addon0", "addon1"]
	Addons []string `json:"addons,omitempty" yaml:"addons,omitempty"`
	// Services to enable for the instance
	Services []NetworkServiceSpec `json:"services" yaml:"services"`
	// Resources specifies the resources allocated for the instance
	Resources struct {
		// Number of CPU cores the instance should get assigned.
		// Example: 4
		CPUs *int `json:"cpus,omitempty" yaml:"cpus,omitempty"`
		// Disk size the instance should get allocated in bytes
		// Example: 3221225472
		DiskSize *int64 `json:"disk_size" yaml:"disk_size"`
		// Memory the instance should get assigned in bytes.
		// Example: 3221225472
		Memory *int64 `json:"memory,omitempty" yaml:"memory,omitempty"`
		// Number of GPU slots the instance should get assigned.
		// Example: 1
		GPUSlots *int `json:"gpu-slots,omitempty" yaml:"gpu-slots,omitempty"`
		// Number of VPU slots the instance should get assigned
		// Example: 1
		VPUSlots *int `json:"vpu-slots,omitempty" yaml:"vpu-slots,omitempty"`
	} `json:"resources" yaml:"resources"`
	// Tags which will be assigned to the instance
	// Example: ["tag0", "tag1"]
	Tags   []string `json:"tags,omitempty" yaml:"tags,omitempty"`
	Config struct {
		// Platform specifies the Anbox platform the instance is running with
		// Example: webrtc
		Platform string `json:"platform,omitempty" yaml:"platform,omitempty"`
		// BootPackage specifies the Android application package name which is started by default
		// Example: com.android.settings
		BootPackage string `json:"boot_package,omitempty" yaml:"boot_package,omitempty"`
		// BootActivity specifies the Android activity which is started by default
		// Example: com.android.settings/.DevSettings
		BootActivity string `json:"boot_activity,omitempty" yaml:"boot_activity,omitempty"`
		// MetricsServer specifies a metrics server the instance will use
		// Example: 10.0.0.45:8086
		MetricsServer string `json:"metrics_server,omitempty" yaml:"metrics_server,omitempty"`
		// DisableWatchdog defines whether the watchdog is disabled
		DisableWatchdog bool `json:"disable_watchdog,omitempty" yaml:"disable_watchdog,omitempty"`
		// Feature flags to enable for the instance.
		// Example: feature0, feature1
		Features string `json:"features,omitempty" yaml:"features,omitempty"`
		// DevMode specifies if development mode has been turned on for the instance
		DevMode bool `json:"devmode,omitempty" yaml:"devmode,omitempty"`
	} `json:"config,omitempty"`
	// Do not start the instance after creation.
	NoStart bool `json:"no_start,omitempty" yaml:"no_start,omitempty"`
}

// InstancesDelete represents a list of instances to delete together
//
// swagger:model
//
// API extensions: bulk_delete.instances
type InstancesDelete struct {
	// IDs of the instances to delete
	// Example: ["cilsreunfpfec9b1ktg0", "cilsreunfpfec9b1ktg1"]
	IDs []string `json:"ids" yaml:"ids"`
	// Whether deletion of the instances should be forced
	Force bool `json:"force" yaml:"force"`
}

// InstanceExecPost represents an instance execution request
//
// swagger:model
//
// API extension: instance_exec
type InstanceExecPost struct {
	// Command inside the instance to execute
	// Example: /bin/ls
	Command []string `json:"command" yaml:"command"`
	// Environment to setup when the command is executed.
	// Example: {"FOO": "bar"}
	Environment map[string]string `json:"environment" yaml:"environment"`
	// Whether the command is executed interactively or not
	Interactive bool `json:"interactive" yaml:"interactive"`
	// Width of the terminal. Only required when `interactive` is set to `true`.
	Width int `json:"width" yaml:"width"`
	// Height of the terminal. Only required when `interactive` is set to `true`.
	Height int `json:"height" yaml:"height"`
}

// InstanceExecControl represents a message on the instance shell "control" socket
//
// swagger:model
//
// API extension: instance_exec
type InstanceExecControl struct {
	// Command to execute. Possible values are: window-resize, signal
	// Example: window-resize
	Command string `json:"command" yaml:"command"`
	// Arguments to pass to the command
	// Example: {"width": 1280, "height": 720}
	Args map[string]string `json:"args" yaml:"args"`
	// Signal to send
	Signal int `json:"signal" yaml:"signal"`
}

// InstanceDelete describes a request used to delete a instance
//
// swagger:model
type InstanceDelete struct {
	// Whether deletion of the instance should be forced
	Force bool `json:"force"`
}

// InstancePatch describes the fields which can be changed for an existing instance
//
// swagger:model
type InstancePatch struct {
	// Desired status of the instance
	DesiredStatus *string `json:"desired_status"`
}

// MapInstanceToContainer converts an instance API object to a container one
func MapInstanceToContainer(inst *Instance) Container {
	c := Container{
		ID:            inst.ID,
		Name:          inst.Name,
		Type:          ContainerTypeRegular,
		StatusCode:    ContainerStatus(inst.StatusCode),
		Status:        inst.Status,
		Node:          inst.Node,
		AppID:         inst.AppID,
		AppName:       inst.AppName,
		AppVersion:    inst.AppVersion,
		ImageID:       inst.ImageID,
		ImageVersion:  inst.ImageVersion,
		CreatedAt:     inst.CreatedAt,
		Address:       inst.Address,
		PublicAddress: inst.PublicAddress,
		StoredLogs:    inst.StoredLogs,
		ErrorMessage:  inst.ErrorMessage,
		Architecture:  inst.Architecture,
		Tags:          inst.Tags,
	}
	for _, service := range inst.Services {
		c.Services = append(c.Services, ContainerService{
			Port:        service.Port,
			PortEnd:     service.PortEnd,
			NodePort:    service.NodePort,
			NodePortEnd: service.NodePortEnd,
			Protocols:   service.Protocols,
			Expose:      service.Expose,
			Name:        service.Name,
		})
	}
	c.Config.Platform = inst.Config.Platform
	c.Config.BootPackage = inst.Config.BootPackage
	c.Config.BootActivity = inst.Config.BootActivity
	c.Config.MetricsServer = inst.Config.MetricsServer
	c.Config.DisableWatchdog = inst.Config.DisableWatchdog
	c.Config.DevMode = inst.Config.DevMode
	c.Resources.CPUs = inst.Resources.CPUs

	c.Resources.Memory = shared.GetByteSizeString(inst.Resources.Memory, 0)
	c.Resources.DiskSize = shared.GetByteSizeString(inst.Resources.DiskSize, 0)
	c.Resources.GPUSlots = inst.Resources.GPUSlots
	c.Resources.VPUSlots = inst.Resources.VPUSlots
	if inst.IsBase {
		c.Type = ContainerTypeBase
	}
	return c
}

// MapContainerToInstance maps a container to an instance object
func MapContainerToInstance(c *Container) (Instance, error) {
	inst := Instance{
		ID:            c.ID,
		Name:          c.Name,
		StatusCode:    InstanceStatus(c.StatusCode),
		Status:        c.Status,
		Node:          c.Node,
		AppID:         c.AppID,
		AppName:       c.AppName,
		AppVersion:    c.AppVersion,
		ImageID:       c.ImageID,
		ImageVersion:  c.ImageVersion,
		CreatedAt:     c.CreatedAt,
		Address:       c.Address,
		PublicAddress: c.PublicAddress,
		StoredLogs:    c.StoredLogs,
		ErrorMessage:  c.ErrorMessage,
		Architecture:  c.Architecture,
		Tags:          c.Tags,
	}
	for _, service := range c.Services {
		inst.Services = append(inst.Services, InstanceService{
			Port:        service.Port,
			PortEnd:     service.PortEnd,
			NodePort:    service.NodePort,
			NodePortEnd: service.NodePortEnd,
			Protocols:   service.Protocols,
			Expose:      service.Expose,
			Name:        service.Name,
		})
	}
	inst.Config.Platform = c.Config.Platform
	inst.Config.BootPackage = c.Config.BootPackage
	inst.Config.BootActivity = c.Config.BootActivity
	inst.Config.MetricsServer = c.Config.MetricsServer
	inst.Config.DisableWatchdog = c.Config.DisableWatchdog
	inst.Config.DevMode = c.Config.DevMode
	inst.Resources.CPUs = c.Resources.CPUs
	inst.Resources.GPUSlots = c.Resources.GPUSlots
	inst.Resources.VPUSlots = c.Resources.VPUSlots
	inst.IsBase = c.Type == ContainerTypeBase

	memory, err := shared.ParseByteSizeString(c.Resources.Memory)
	if err != nil {
		return Instance{}, fmt.Errorf("failed to parse memory resource value: %w", err)
	}
	inst.Resources.Memory = memory

	diskSize, err := shared.ParseByteSizeString(c.Resources.DiskSize)
	if err != nil {
		return Instance{}, fmt.Errorf("failed to parse disk size resource value: %w", err)
	}
	inst.Resources.DiskSize = diskSize

	return inst, nil
}
