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

// NodeStatus describes the current status of a node
//
// swagger:model
type NodeStatus int

const (
	// NodeStatusError represents the status a node is in when it can not be used
	// because of an error.
	NodeStatusError NodeStatus = -1
	// NodeStatusUnknown describes the status of a node when the status is not known
	NodeStatusUnknown NodeStatus = 0
	// NodeStatusCreated represents the status a node is in after it was created
	NodeStatusCreated NodeStatus = 1
	// NodeStatusInitializing represents the status a node is in when its corresponding
	// LXD instance is currently being configured.
	NodeStatusInitializing NodeStatus = 2
	// NodeStatusInitialized represents the status a node is in after its corresponding
	// LXD instance was initialized
	NodeStatusInitialized NodeStatus = 3
	// NodeStatusOnline represents the status a node is in when it available and can
	// be used by AMS.
	NodeStatusOnline NodeStatus = 4
	// NodeStatusOffline represents the status a node is in when it is not available
	// and can't be used by AMS.
	NodeStatusOffline NodeStatus = 5
	// NodeStatusDeleted represents the status a node is in when it got deleted by the user
	NodeStatusDeleted NodeStatus = 6
)

// String returns a textual representation of a node status
func (s *NodeStatus) String() string {
	switch *s {
	case NodeStatusCreated:
		return "created"
	case NodeStatusInitialized:
		return "initialized"
	case NodeStatusOnline:
		return "online"
	case NodeStatusOffline:
		return "offline"
	case NodeStatusError:
		return "error"
	case NodeStatusDeleted:
		return "deleted"
	}
	return "unknown"
}

// NodeGPUAllocation describes a single allocation on a GPU
//
// swagger:model
type NodeGPUAllocation struct {
	// List of GPU IDs allocated to the container
	// Example: [0,1]
	GPUs []uint64 `json:"gpus" yaml:"gpus"`
	// Number of GPU Slots allocated to the container
	// Example: 1
	Slots int `json:"slots" yaml:"slots"`
	// Number of Encoder Slots allocated to the container
	// Example: 1
	EncoderSlots int `json:"encoder_slots" yaml:"encoder_slots"`
}

// NodeGPU describes a single GPU available on a node
//
// swagger:model
type NodeGPU struct {
	// ID of the GPU configured on the node
	// Example: 0
	ID uint64 `json:"id" yaml:"id"`
	// Type is the type of the GPU. Possible values are: nvidia, amd, intel
	// Example: nvidia
	Type string `json:"type" yaml:"type"`
	// PCI Bus Address used by the GPU
	// Example: 00:08.0
	PCIAddress string `json:"pci_address" yaml:"pci_address"`
	// PCI Bus Address used by the GPU
	// Example: D129
	RenderName string `json:"render_name" yaml:"render_name"`
	// Number of the GPU slots available
	// Example: 20
	Slots int `json:"slots" yaml:"slots"`
	// Number of the encoder slots available on the GPU
	// Example: 20
	EncoderSlots int `json:"encoder_slots" yaml:"encoder_slots"`
	// Map of current allocations and containers on the GPU
	Allocations map[string]NodeGPUAllocation `json:"allocations" yaml:"allocations"`
	// NUMA Node number for the GPU
	// Example: 0
	NUMANode uint64 `json:"numa_node" yaml:"numa_node"`
}

// NodeVPUAllocation describes a single allocation for a VPU
//
// swagger:model
type NodeVPUAllocation struct {
	// VPU IDs the allocation is for
	IDs []uint64 `json:"ids" yaml:"ids"`
	// Number of slots used by this allocation
	Slots int `json:"slots" yaml:"slots"`
}

// NodeVPU describes a single independent VPU available on a node
//
// swagger:model
type NodeVPU struct {
	// ID of the VPU
	ID uint64 `json:"id" yaml:"id"`
	// Type of the VPU. Valid values are: unknown, netint
	Type string `json:"type" yaml:"type"`
	// Model name of the VPU
	Model string `json:"model" yaml:"model"`
	// NUMA node the card sits on
	NUMANode uint64 `json:"numa_node" yaml:"numa_node"`
	// Number of slots available on the VPU
	Slots int `json:"slots" yaml:"slots"`
	// Map of current allocations on the VPU
	Allocations map[string]NodeVPUAllocation `json:"allocations" yaml:"allocations"`
}

// Node describes a single node of the underlying LXD cluster AMS manages
//
// swagger:model
type Node struct {
	// Name of the node
	// Example: lxd0
	Name string `json:"name" yaml:"name"`
	// Internal IP address of the node
	// Example: 10.0.0.1
	// swagger:strfmt ipv4
	Address string `json:"address" yaml:"address"`
	// Public IP address of the node
	// Example: 10.0.0.1
	// swagger:strfmt ipv4
	PublicAddress string `json:"public_address" yaml:"public_address"`
	// MTU for the configured network bridge on LXD
	// Example: 1500
	NetworkBridgeMTU int `json:"network_bridge_mtu" yaml:"network_bridge_mtu"`
	// Number of CPUs on the node
	// Example: 4
	CPUs int `json:"cpus" yaml:"cpus"`
	// CPU allocation rate for the node
	// Example: 4
	CPUAllocationRate float32 `json:"cpu_allocation_rate" yaml:"cpu_allocation_rate"`
	// Memory (in GB) of the LXD node
	// Example: 8GB
	Memory string `json:"memory" yaml:"memory"`
	// Memory allocation rate for the node
	// Example: 2
	MemoryAllocationRate float32 `json:"memory_allocation_rate" yaml:"memory_allocation_rate"`
	// Current status code of the node as an integer value
	// Example: 4
	StatusCode NodeStatus `json:"status_code" yaml:"status_code"`
	// Current status of the node
	// Example: online
	Status string `json:"status" yaml:"status"`
	// Flag to represent the master node for the AMS cluster
	// Example: true
	IsMaster bool `json:"is_master" yaml:"is_master"`
	// Disk size for the node
	// Example: true
	DiskSize string `json:"disk_size" yaml:"disk_size"`
	// Number of GPU slots present on the node
	// Example: 0
	GPUSlots int `json:"gpu_slots" yaml:"gpu_slots"`
	// Number of GPU encoder slots present on the node
	// Example: 0
	GPUEncoderSlots int `json:"gpu_encoder_slots" yaml:"gpu_encoder_slots"`
	// Tags attached to the node
	// Example: ["created_by=anbox", "gpu=nvidia"]
	Tags []string `json:"tags" yaml:"tags"`
	// Flag used to see if the node is available to schedule containers
	// Example: false
	Unschedulable bool `json:"unschedulable" yaml:"unschedulable"`
	// CPU architecture of the node
	// Example: aarch64
	Architecture string `json:"architecture,omitempty" yaml:"architecture,omitempty"`
	// Name of the storage pool configured for the node
	// Example: default
	StoragePool string `json:"storage_pool" yaml:"storage_pool"`
	// Flag used to control if AMS can manage the LXD node
	// Example: false
	Managed bool `json:"managed" yaml:"managed"`
	// GPU information for the node
	GPUs []NodeGPU `json:"gpus" yaml:"gpus"`
	// VPU information for the node
	VPUs []NodeVPU `json:"vpus" yaml:"vpus"`

	// DEPRECATED Flag in favour of `unschedulable` flag
	// Example: false
	DEPRECATEDUnschedulable bool `json:"unscheduable" yaml:"unscheduable"`
}

// NodesPost describes a request to create a new node on AMS
//
// swagger:model
type NodesPost struct {
	// Name of the node
	// Example: lxd0
	Name string `json:"name"`
	// Internal IP address of the node
	// Example: 10.0.0.1
	// swagger:strfmt ipv4
	Address string `json:"address"`
	// Public IP address of the node
	// Example: 10.0.0.1
	// swagger:strfmt ipv4
	PublicAddress string `json:"public_address"`
	// Trust password for the LXD instance
	// Example: sUp3rs3cr3t
	TrustPassword string `json:"trust_password"`
	// MTU for the configured network bridge on LXD
	// Example: 1500
	NetworkBridgeMTU int `json:"network_bridge_mtu"`
	// Storage device to use for configuring LXD storage pools
	// Example: /dev/sdb
	StorageDevice string `json:"storage_device"`
	// Number of CPUs on the node
	// Example: 4
	CPUs int `json:"cpus"`
	// CPU allocation rate for the node
	// Example: 4
	CPUAllocationRate float32 `json:"cpu_allocation_rate"`
	// Memory (in GB) of the node
	// Example: 8GB
	Memory string `json:"memory"`
	// Memory allocation rate for the node
	// Example: 2
	MemoryAllocationRate float32 `json:"memory_allocation_rate"`
	// Number of GPU slots to configure on the node
	// Example: 2
	GPUSlots int `json:"gpu_slots"`
	// Number of GPU encoder slots to configure on the node
	// Example: 4
	GPUEncoderSlots int `json:"gpu_encoder_slots" yaml:"gpu_encoder_slots"`
	// Tags to attach to the node
	// Example: ["created_by=anbox", "gpu=nvidia"]
	Tags []string `json:"tags" yaml:"tags"`
	// Flag used to control if AMS can manage the LXD node
	// Example: false
	Unmanaged bool `json:"unmanaged" yaml:"unmanaged"`
	// Name of the storage pool to use for configuring the LXD node
	// Example: default
	StoragePool string `json:"storage_pool" yaml:"storage_pool"`
	// Name of the network bridge to create on the LXD node
	// Example: amsbr0
	NetworkName string `json:"network_name" yaml:"network_name"`
	// CIDR of the subnet to configure for the network bridge on LXD
	// Example: 10.0.0.0/24
	// swagger:strfmt ipv4
	NetworkSubnet string `json:"network_subnet" yaml:"network_subnet"`
	// Trust token for the LXD instance
	// Example: csdflkj3lks
	TrustToken string `json:"trust_token"`

	// Name of the network ACL to create on the LXD node
	// Example: ams0
	// Deprecated: This field is no longer supported since 1.23
	DEPRECATEDNetworkACLName string `json:"network_acl_name" yaml:"network_acl_name"`
}

// NodeGPUPatch allows changing configuration for individual GPUs
//
// swagger:model
type NodeGPUPatch struct {
	// ID must match the ID of an existing GPU and is used for identifying
	// the GPU which should be patched

	// ID of the GPU configured on the node
	// Example: 0
	ID uint64 `json:"id" yaml:"id"`
	// Update the number of the GPU slots available on the Node
	// Example: 20
	Slots *int `json:"slots" yaml:"slots"`
	// Update the number of GPU encoder slots
	// Example: 4
	EncoderSlots *int `json:"encoder_slots" yaml:"encoder_slots"`
}

// NodePatch describes a request to update an existing node
//
// swagger:model
type NodePatch struct {
	// The public, reachable address of the node.
	// Example: 10.0.0.1
	// swagger:strfmt ipv4
	PublicAddress *string `json:"public_address"`
	// Number of CPUs dedicated to instances.
	// Example: 4
	CPUs *int `json:"cpus"`
	// CPU allocation rate used for over-committing resources
	// Example: 4
	// Extensions:
	// x-docs-ref: sec-over-committing
	CPUAllocationRate *float32 `json:"cpu_allocation_rate"`
	// Update the memory (in GB) for the node
	// Example: 2GB
	Memory *string `json:"memory"`
	// Memory allocation rate used for over-committing resources.
	// Example: 2
	// Extensions:
	// x-docs-ref: sec-over-committing
	MemoryAllocationRate *float32 `json:"memory_allocation_rate"`
	// Update the number of GPU slots to configure on the node
	// Example: 2
	GPUSlots *int `json:"gpu_slots"`
	// Number of GPU encoder slots available on the node.
	// `0` for nodes without GPU
	// `32` for nodes with NVIDIA GPU
	// `10` for nodes with AMD or Intel GPU
	// Example: 4
	// Extensions:
	// x-docs-ref: sec-gpu-slots
	GPUEncoderSlots *int `json:"gpu_encoder_slots" yaml:"gpu_encoder_slots"`
	// Tags to identify the node.
	// Example: ["created_by=anbox", "gpu=nvidia"]
	Tags *[]string `json:"tags" yaml:"tags"`
	// When set to `true`, the node cannot be scheduled, which prevents new instances from being launched on it.
	// Example: true
	Unschedulable *bool          `json:"unschedulable" yaml:"unschedulable"`
	GPUs          []NodeGPUPatch `json:"gpus" yaml:"gpus"`
	// The network subnet of the machine where the node runs.
	// Example: 10.0.0.1/24
	// swagger:strfmt ipv4
	Subnet *string `json:"subnet" yaml:"subnet"`

	// DEPRECATED Flag in favour of `unschedulable` flag
	// Example: false
	// Extensions:
	// x-deprecated-since: "1.20"
	DEPRECATEDUnschedulable *bool `json:"unscheduable" yaml:"unscheduable"`
}

// NodeDelete describes a request used to delete a node
//
// swagger:model
type NodeDelete struct {
	// Use this to force deletion of a node from AMS and LXD cluster
	// Example: true
	Force bool `json:"force"`
	// Use this to remove the node from the LXD cluster as well
	// Example: true
	KeepInCluster bool `json:"keep_in_cluster"`
}
