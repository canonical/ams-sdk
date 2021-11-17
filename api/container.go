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
	Port int `json:"port,omitempty" yaml:"port,omitempty"`
	// PortEnd is the end of the port range set for a service. If empty, only a
	// single port is opened
	PortEnd int `json:"port_end,omitempty" yaml:"port_end,omitempty"`
	// List of network protocols (tcp, udp) the port should be exposed for
	Protocols []NetworkProtocol `json:"protocols" yaml:"protocols"`
	// Expose defines wether the service is exposed on the public endpoint of the node
	// or if it is only available on the private endpoint. To expose the service set to
	// true and to false otherwise.
	Expose bool `json:"expose" yaml:"expose"`
	// Name gives the container a hint what the exposed port is being used for. This
	// allows further tweaks inside the container to expose the service correctly.
	Name string `json:"name" yaml:"name"`
}

// ContainerService describes a single service the container exposes to the outside world.
// While NetworkServiceSpec defines what the user requests, ContainerService is what is actually
// opened on the container.
type ContainerService struct {
	// Port is the port the container provides a service on
	Port int `json:"port,omitempty" yaml:"port,omitempty"`
	// PortEnd, if specified, denotes the end of the port range starting at Port
	PortEnd int `json:"port_end,omitempty" yaml:"port_end,omitempty"`
	// NodePort is the port used on the LXD node to map to the service port
	// If left empty the node port is automatically selected.
	NodePort *int `json:"node_port,omitempty" yaml:"node_port,omitempty"`
	// NodePortEnd, if specified, denotes the end of the port range on the node starting
	// at NodePort
	NodePortEnd *int `json:"node_port_end,omitempty" yaml:"node_port_end,omitempty"`
	// List of network protocols (tcp, udp) the port should be exposed for
	Protocols []NetworkProtocol `json:"protocols" yaml:"protocols"`
	// Expose defines wether the service is exposed on the public endpoint of the node
	// or if it is only available on the private endpoint. To expose the service set to
	// true and to false otherwise.
	Expose bool `json:"expose" yaml:"expose"`
	// Name gives the container a hint what the exposed port is being used for. This
	// allows further tweaks inside the container to expose the service correctly.
	Name string `json:"name" yaml:"name"`
}

// ContainerType describes the type of a container
type ContainerType string

const (
	// ContainerTypeUnknown describes a container whichs type is not known
	ContainerTypeUnknown ContainerType = "unknown"
	// ContainerTypeRegular describes a regular container
	ContainerTypeRegular ContainerType = "regular"
	// ContainerTypeBase describes a base container
	ContainerTypeBase ContainerType = "base"
)

// Container represents a AMS container
type Container struct {
	ID            string             `json:"id" yaml:"id"`
	Name          string             `json:"name" yaml:"name"`
	Type          ContainerType      `json:"type" yaml:"type"`
	StatusCode    ContainerStatus    `json:"status_code" yaml:"status_code"`
	Status        string             `json:"status" yaml:"status"`
	Node          string             `json:"node" yaml:"node"`
	AppID         string             `json:"app_id" yaml:"app_id"`
	AppVersion    int                `json:"app_version" yaml:"app_version"`
	ImageID       string             `json:"image_id" yaml:"image_id"`
	ImageVersion  int                `json:"image_version" yaml:"image_version"`
	CreatedAt     int64              `json:"created_at" yaml:"created_at"`
	Address       string             `json:"address" yaml:"address"`
	PublicAddress string             `json:"public_address" yaml:"public_address"`
	Services      []ContainerService `json:"services" yaml:"services"`
	StoredLogs    []string           `json:"stored_logs" yaml:"stored_logs"`

	// Only filled with content when Status == ContainerStatusError
	ErrorMessage string `json:"error_message" yaml:"error_message"`

	Config struct {
		Platform        string `json:"platform,omitempty" yaml:"platform,omitempty"`
		BootPackage     string `json:"boot_package,omitempty" yaml:"boot_package,omitempty"`
		BootActivity    string `json:"boot_activity,omitempty" yaml:"boot_activity,omitempty"`
		MetricsServer   string `json:"metrics_server,omitempty" yaml:"metrics_server,omitempty"`
		DisableWatchdog bool   `json:"disable_watchdog,omitempty" yaml:"disable_watchdog,omitempty"`
	} `json:"config,omitempty"`

	Resources struct {
		CPUs     int    `json:"cpus,omitempty" yaml:"cpus,omitempty"`
		Memory   string `json:"memory,omitempty" yaml:"memory,omitempty"`
		DiskSize string `json:"disk-size,omitempty" yaml:"disk-size,omitempty"`
		GPUSlots int    `json:"gpu-slots,omitempty" yaml:"gpu-slots,omitempty"`
	} `json:"resources,omitempty"`

	Architecture string `json:"architecture,omitempty" yaml:"architecture,omitempty"`

	Tags []string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// ContainersPost represents the fields required to launch a new container for
// a specific application
type ContainersPost struct {
	ApplicationID      string               `json:"app_id" yaml:"app_id"`
	ApplicationVersion *int                 `json:"app_version" yaml:"app_version"`
	ImageID            string               `json:"image_id" yaml:"image_id"`
	ImageVersion       *int                 `json:"image_version" yaml:"image_version"`
	InstanceType       string               `json:"instance_type" yaml:"instance_type"`
	Node               string               `json:"node" yaml:"node"`
	Userdata           *string              `json:"user_data,omitempty" yaml:"user_data,omitempty"`
	Addons             []string             `json:"addons,omitempty" yaml:"addons,omitempty"`
	Services           []NetworkServiceSpec `json:"services" yaml:"services"`
	DiskSize           *int64               `json:"disk_size" yaml:"disk_size"`
	CPUs               *int                 `json:"cpus,omitempty" yaml:"cpus,omitempty"`
	Memory             *int64               `json:"memory,omitempty" yaml:"memory,omitempty"`
	GPUSlots           *int                 `json:"gpu-slots,omitempty" yaml:"gpu-slots,omitempty"`
	Tags               []string             `json:"tags,omitempty" yaml:"tags,omitempty"`
	Config             struct {
		Platform        string `json:"platform,omitempty" yaml:"platform,omitempty"`
		BootPackage     string `json:"boot_package,omitempty" yaml:"boot_package,omitempty"`
		BootActivity    string `json:"boot_activity,omitempty" yaml:"boot_activity,omitempty"`
		MetricsServer   string `json:"metrics_server,omitempty" yaml:"metrics_server,omitempty"`
		DisableWatchdog bool   `json:"disable_watchdog,omitempty" yaml:"disable_watchdog,omitempty"`
		Features        string `json:"features,omitempty" yaml:"features,omitempty"`
	} `json:"config,omitempty"`
}

// ContainersGet represents a list of containers
type ContainersGet struct {
	Containers []Container `json:"containers" yaml:"containers"`
}

// ContainerPatch describes the fields which can be changed for an existing container
type ContainerPatch struct {
	Status     *string `json:"status"`
	ErrCode    *int    `json:"error_code,omitempty"`
	ErrMessage *string `json:"error_message,omitempty"`
}

// ContainerExecPost represents a AMS container shell request
//
// API extension: container_exec
type ContainerExecPost struct {
	Command     []string          `json:"command" yaml:"command"`
	Environment map[string]string `json:"environment" yaml:"environment"`
	Interactive bool              `json:"interactive" yaml:"interactive"`
	Width       int               `json:"width" yaml:"width"`
	Height      int               `json:"height" yaml:"height"`
}

// ContainerExecControl represents a message on the container shell "control" socket
//
// API extension: container_exec
type ContainerExecControl struct {
	Command string            `json:"command" yaml:"command"`
	Args    map[string]string `json:"args" yaml:"args"`
	Signal  int               `json:"signal" yaml:"signal"`
}

// ContainerDelete describes a request used to delete a container
type ContainerDelete struct {
	Force bool `json:"force"`
}
