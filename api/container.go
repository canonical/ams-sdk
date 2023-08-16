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

// ContainerStatus represents the status a container can be in
type ContainerStatus int

const (
	// ContainerStatusError represents the state a container is in when an error occurred.
	ContainerStatusError ContainerStatus = -1
	// ContainerStatusUnknown represents the state a container is in when its real state
	// cannot be determined.
	ContainerStatusUnknown ContainerStatus = 0
	// ContainerStatusCreated represents the status a container is in when its object
	// was created on the store but no further operation was performed on it yet.
	ContainerStatusCreated ContainerStatus = 1
	// ContainerStatusPrepared represents the status a container is in when it got all
	// its resources assigned and is ready to be constructed and started on the LXD cluster.
	ContainerStatusPrepared ContainerStatus = 2
	// ContainerStatusStarted represents the state a container is in when it is currently
	// starting.
	ContainerStatusStarted ContainerStatus = 3
	// ContainerStatusRunning represents the state a container is in when it is running
	ContainerStatusRunning ContainerStatus = 4
	// ContainerStatusStopped represents the state a container is in when it is stopped
	ContainerStatusStopped ContainerStatus = 5
	// ContainerStatusDeleted represents the state a container is currently being deleted
	ContainerStatusDeleted ContainerStatus = 6
)

func (s ContainerStatus) String() string {
	switch s {
	case ContainerStatusCreated:
		return "created"
	case ContainerStatusPrepared:
		return "prepared"
	case ContainerStatusStarted:
		return "started"
	case ContainerStatusStopped:
		return "stopped"
	case ContainerStatusRunning:
		return "running"
	case ContainerStatusError:
		return "error"
	case ContainerStatusDeleted:
		return "deleted"
	}
	return "unknown"
}

// NetworkProtocol describes a specific network protocol like TCP or UDP a ContainerService can use
type NetworkProtocol string

const (
	// NetworkProtocolUnknown describes a unknown network protocol
	NetworkProtocolUnknown NetworkProtocol = "unknown"
	// NetworkProtocolTCP describes the TCP network protocol
	NetworkProtocolTCP NetworkProtocol = "tcp"
	// NetworkProtocolUDP describes the UDP network protocol
	NetworkProtocolUDP NetworkProtocol = "udp"
)

// NetworkProtocolFromString parses a network protocol from the given value. Returns NetworkProtocolUnknown
// if not known value can be parsed.
func NetworkProtocolFromString(value string) NetworkProtocol {
	switch value {
	case "tcp":
		return NetworkProtocolTCP
	case "udp":
		return NetworkProtocolUDP
	}
	return "unknown"
}

// NetworkServiceSpec is used to define the user defined network services that should be opened on a container
type NetworkServiceSpec struct {
	// Port is the port the container provides a service on
	// Example: 3000
	Port int `json:"port,omitempty" yaml:"port,omitempty"`
	// PortEnd is the end of the port range set for a service. If empty, only a
	// single port is opened
	// Example: 3010
	PortEnd int `json:"port_end,omitempty" yaml:"port_end,omitempty"`
	// List of network protocols (tcp, udp) the port should be exposed for
	// Example: ["tcp", "udp"]
	Protocols []NetworkProtocol `json:"protocols" yaml:"protocols"`
	// Expose defines wether the service is exposed on the public endpoint of the node
	// or if it is only available on the private endpoint. To expose the service set to
	// true and to false otherwise.
	Expose bool `json:"expose" yaml:"expose"`
	// Name gives the container a hint what the exposed port is being used for. This
	// allows further tweaks inside the container to expose the service correctly.
	// Exampe: ssh
	Name string `json:"name" yaml:"name"`
}

// ContainerService describes a single service the container exposes to the outside world.
//
// While NetworkServiceSpec defines what the user requests, ContainerService is what is actually
// opened on the container.
//
// swagger:model
type ContainerService struct {
	// Port is the port the container provides a service on
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
	// Name gives the container a hint what the exposed port is being used for. This
	// allows further tweaks inside the container to expose the service correctly.
	// Example: myservice
	Name string `json:"name" yaml:"name"`
}

// ContainerType describes the type of a container.
// Possible values are: regular, base, unknown
type ContainerType string

const (
	// ContainerTypeUnknown describes a container whichs type is not known
	ContainerTypeUnknown ContainerType = "unknown"
	// ContainerTypeRegular describes a regular container
	ContainerTypeRegular ContainerType = "regular"
	// ContainerTypeBase describes a base container
	ContainerTypeBase ContainerType = "base"
)

// Container represents a single container
//
// swagger:model
type Container struct {
	// ID of the container
	// Example: cilsreunfpfec9b1ktg0
	ID string `json:"id" yaml:"id"`
	// Name of the container. Typically in the format "ams-<ID>".
	// Example: ams-cilsreunfpfec9b1ktg0
	Name string `json:"name" yaml:"name"`
	// Type of the container. Possible values are: regular, base, unknown
	// Example: regular
	Type ContainerType `json:"type" yaml:"type"`
	// StatusCode of the container. Matches the Status field.
	// Example: 4
	StatusCode ContainerStatus `json:"status_code" yaml:"status_code"`
	// Status of the container
	// Example: running
	Status string `json:"status" yaml:"status"`
	// Node the container is running on
	// Example: lxd0
	Node string `json:"node" yaml:"node"`
	// AppID is the ID of the application the container is created from. Empty if the
	// container has not been created from an application.
	// Example: cilsiomnfpfec9b1kteg
	AppID string `json:"app_id" yaml:"app_id"`
	// AppName is the name of the application the container is created from. Empty if the
	// container has not been created from an application.
	// Example: myapp
	AppName string `json:"app_name" yaml:"app_name"`
	// AppVersion is the version of the application the container is created from. Empty if the
	// container has not been created from an application.
	// Example: 0
	AppVersion int `json:"app_version" yaml:"app_version"`
	// ImageID is the ID of the image the container is created from. Empty if the
	// container has not been created from an image.
	// Example: cilshrmnfpfec9b1kte0
	ImageID string `json:"image_id" yaml:"image_id"`
	// ImageVersion is the version of the image the container is created from. Empty if the
	// container has not been created from an image.
	// Example: 0
	ImageVersion int `json:"image_version" yaml:"image_version"`
	// CreatedAt specifies the time at which the container was created
	// Example: 1689604498
	CreatedAt int64 `json:"created_at" yaml:"created_at"`
	// Address is the IP address of the container
	// Example: 192.168.1.74
	Address string `json:"address" yaml:"address"`
	// PublicAddress is the external IP address the container is accessible on (in most
	// cases the IP of the node it is running on)
	// Example: 1.2.3.4
	PublicAddress string `json:"public_address" yaml:"public_address"`
	// Services the container exposes
	Services []ContainerService `json:"services" yaml:"services"`
	// StoredLogs lists log files AMS stores for the container.
	// Example: ["android.log", "system.log"]
	StoredLogs []string `json:"stored_logs" yaml:"stored_logs"`
	// ErrorMessage provides an error message when the container status is set to error.
	// Example: container failed to boot
	ErrorMessage string `json:"error_message" yaml:"error_message"`
	// Config summarizes the configuration the container uses
	Config struct {
		// Platform specifies the Anbox platform the container is running with
		// Example: webrtc
		Platform string `json:"platform,omitempty" yaml:"platform,omitempty"`
		// BootPackage specifies the Android application package name which is started by default
		// Example: com.android.settings
		BootPackage string `json:"boot_package,omitempty" yaml:"boot_package,omitempty"`
		// BootActivity specifies the Android activity which is started by default
		// Example: com.android.settings/.DevSettings
		BootActivity string `json:"boot_activity,omitempty" yaml:"boot_activity,omitempty"`
		// MetricsServer specifies a metrics server the container will use
		// Example: 10.0.0.45:8086
		MetricsServer string `json:"metrics_server,omitempty" yaml:"metrics_server,omitempty"`
		// DisableWatchdog defines whether the watchdog is disabled
		DisableWatchdog bool `json:"disable_watchdog,omitempty" yaml:"disable_watchdog,omitempty"`
		// DevMode specifies if development mode has been turned on for the container
		DevMode bool `json:"devmode,omitempty" yaml:"devmode,omitempty"`
	} `json:"config,omitempty"`
	// Resources specifies the resources allocated for the container
	Resources struct {
		// CPUs cores assigned to the container
		// Example: 2
		CPUs int `json:"cpus,omitempty" yaml:"cpus,omitempty"`
		// Memory assigned to the container
		// Example: 3GB
		Memory string `json:"memory,omitempty" yaml:"memory,omitempty"`
		// DiskSize specifies the amount of storage assigned to the container
		// Example: 3GB
		DiskSize string `json:"disk-size,omitempty" yaml:"disk-size,omitempty"`
		// GPUSlots specifies the number of GPU slots the container got allocated
		// Example: 1
		GPUSlots int `json:"gpu-slots,omitempty" yaml:"gpu-slots,omitempty"`
		// VPUSlots specifies the number of VPU slots the container
		VPUSlots int `json:"vpu-slots,omitempty" yaml:"vpu-slots,omitempty"`
	} `json:"resources,omitempty"`
	// Architecture describes the CPU archtitecture the container is using
	// Example: aarch64
	Architecture string `json:"architecture,omitempty" yaml:"architecture,omitempty"`
	// Tags specifies the tags the container has assigned
	// Example: ["foo", "bar"]
	Tags []string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// GetContainerFilters returns an array of attributes available on the api to
// filter containers
func GetContainerFilters() []string {
	return []string{
		"id",
		"name",
		"status",
		"type",
		"node",
		"app_id",
		"app_version",
		"image_id",
		"image_version",
		"app_name",
		"tags",
	}
}

// ContainersPost represents the fields required to launch a new container for
// a specific application
//
// swagger:model
type ContainersPost struct {
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
	// Instance type to use for the container.
	// Example a2.3
	InstanceType string `json:"instance_type" yaml:"instance_type"`
	// Node to start the container on. If empty node will be automatically selected.
	// Example: lxd0
	Node string `json:"node" yaml:"node"`
	// User data to pass to the container.
	// Example: {\"key\":\"value\"}
	Userdata *string `json:"user_data,omitempty" yaml:"user_data,omitempty"`
	// Addons to enable for the container
	// Example: ["addon0", "addon1"]
	Addons []string `json:"addons,omitempty" yaml:"addons,omitempty"`
	// Services to enable for the container
	Services []NetworkServiceSpec `json:"services" yaml:"services"`
	// Disk size the container should get allocated in bytes
	// Example: 3221225472
	DiskSize *int64 `json:"disk_size" yaml:"disk_size"`
	// Number of CPU cores the container should get assigned.
	// Example: 4
	CPUs *int `json:"cpus,omitempty" yaml:"cpus,omitempty"`
	// Memory the container should get assigned in bytes.
	// Example: 3221225472
	Memory *int64 `json:"memory,omitempty" yaml:"memory,omitempty"`
	// Number of GPU slots the container should get assigned.
	// Example: 1
	GPUSlots *int `json:"gpu-slots,omitempty" yaml:"gpu-slots,omitempty"`
	// Number of VPU slots the container should get assigned
	// Example: 1
	VPUSlots *int `json:"vpu-slots,omitempty" yaml:"vpu-slots,omitempty"`
	// Tags which will be assigned to the container
	// Example: ["tag0", "tag1"]
	Tags   []string `json:"tags,omitempty" yaml:"tags,omitempty"`
	Config struct {
		// Platform specifies the Anbox platform the container is running with
		// Example: webrtc
		Platform string `json:"platform,omitempty" yaml:"platform,omitempty"`
		// BootPackage specifies the Android application package name which is started by default
		// Example: com.android.settings
		BootPackage string `json:"boot_package,omitempty" yaml:"boot_package,omitempty"`
		// BootActivity specifies the Android activity which is started by default
		// Example: com.android.settings/.DevSettings
		BootActivity string `json:"boot_activity,omitempty" yaml:"boot_activity,omitempty"`
		// MetricsServer specifies a metrics server the container will use
		// Example: 10.0.0.45:8086
		MetricsServer string `json:"metrics_server,omitempty" yaml:"metrics_server,omitempty"`
		// DisableWatchdog defines whether the watchdog is disabled
		DisableWatchdog bool `json:"disable_watchdog,omitempty" yaml:"disable_watchdog,omitempty"`
		// Feature flags to enable for the container.
		// Example: feature0, feature1
		Features string `json:"features,omitempty" yaml:"features,omitempty"`
		// DevMode specifies if development mode has been turned on for the container
		DevMode bool `json:"devmode,omitempty" yaml:"devmode,omitempty"`
	} `json:"config,omitempty"`
	// Do not start the container after creation.
	NoStart bool `json:"no_start,omitempty" yaml:"no_start,omitempty"`
}

// ContainersDelete represents a list of containers to delete together
//
// swagger:model
//
// API extensions: bulk_delete.containers
type ContainersDelete struct {
	// IDs of the containers to delete
	// Example: ["cilsreunfpfec9b1ktg0", "cilsreunfpfec9b1ktg1"]
	IDs []string `json:"ids" yaml:"ids"`
	// Whether deletion of the containers should be forced
	Force bool `json:"force" yaml:"force"`
}

// ContainerExecPost represents an container execution request
//
// swagger:model
//
// API extension: container_exec
type ContainerExecPost struct {
	// Command inside the container to execute
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

// ContainerExecControl represents a message on the container shell "control" socket
//
// swagger:model
//
// API extension: container_exec
type ContainerExecControl struct {
	// Command to execute. Possible values are: window-resize, signal
	// Example: window-resize
	Command string `json:"command" yaml:"command"`
	// Arguments to pass to the command
	// Example: {"width": 1280, "height": 720}
	Args map[string]string `json:"args" yaml:"args"`
	// Signal to send
	Signal int `json:"signal" yaml:"signal"`
}

// ContainerDelete describes a request used to delete a container
//
// swagger:model
type ContainerDelete struct {
	// Whether deletion of the container should be forced
	Force bool `json:"force"`
}

// ContainerPatch describes the fields which can be changed for an existing container
//
// swagger:model
type ContainerPatch struct {
	// Desired status of the container
	DesiredStatus *string `json:"desired_status"`
}
