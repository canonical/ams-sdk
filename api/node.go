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
type NodeGPUAllocation struct {
	GPUs         []uint64 `json:"gpus" yaml:"gpus"`
	Slots        int      `json:"slots" yaml:"slots"`
	EncoderSlots int      `json:"encoder_slots" yaml:"encoder_slots"`
}

// NodeGPU describes a single GPU available on a node
type NodeGPU struct {
	ID           uint64                       `json:"id" yaml:"id"`
	PCIAddress   string                       `json:"pci_address" yaml:"pci_address"`
	RenderName   string                       `json:"render_name" yaml:"render_name"`
	Slots        int                          `json:"slots" yaml:"slots"`
	EncoderSlots int                          `json:"encoder_slots" yaml:"encoder_slots"`
	Allocations  map[string]NodeGPUAllocation `json:"allocations" yaml:"allocations"`
	NUMANode     uint64                       `json:"numa_node" yaml:"numa_node"`
}

// Node describes a single node of the underlying LXD cluster AMS manages
type Node struct {
	Name                 string     `json:"name" yaml:"name"`
	Address              string     `json:"address" yaml:"address"`
	PublicAddress        string     `json:"public_address" yaml:"public_address"`
	NetworkBridgeMTU     int        `json:"network_bridge_mtu" yaml:"network_bridge_mtu"`
	CPUs                 int        `json:"cpus" yaml:"cpus"`
	CPUAllocationRate    float32    `json:"cpu_allocation_rate" yaml:"cpu_allocation_rate"`
	Memory               string     `json:"memory" yaml:"memory"`
	MemoryAllocationRate float32    `json:"memory_allocation_rate" yaml:"memory_allocation_rate"`
	StatusCode           NodeStatus `json:"status_code" yaml:"status_code"`
	Status               string     `json:"status" yaml:"status"`
	IsMaster             bool       `json:"is_master" yaml:"is_master"`
	DiskSize             string     `json:"disk_size" yaml:"disk_size"`
	GPUSlots             int        `json:"gpu_slots" yaml:"gpu_slots"`
	GPUEncoderSlots      int        `json:"gpu_encoder_slots" yaml:"gpu_encoder_slots"`
	Tags                 []string   `json:"tags" yaml:"tags"`
	Unschedulable        bool       `json:"unscheduable" yaml:"unscheduable"`
	Architecture         string     `json:"architecture,omitempty" yaml:"architecture,omitempty"`
	StoragePool          string     `json:"storage_pool" yaml:"storage_pool"`
	Managed              bool       `json:"managed" yaml:"managed"`
	GPUs                 []NodeGPU  `json:"gpus" yaml:"gpus"`
}

// NodesPost describes a request to create a new node on AMS
type NodesPost struct {
	Name                 string   `json:"name"`
	Address              string   `json:"address"`
	PublicAddress        string   `json:"public_address"`
	TrustPassword        string   `json:"trust_password"`
	NetworkBridgeMTU     int      `json:"network_bridge_mtu"`
	StorageDevice        string   `json:"storage_device"`
	CPUs                 int      `json:"cpus"`
	CPUAllocationRate    float32  `json:"cpu_allocation_rate"`
	Memory               string   `json:"memory"`
	MemoryAllocationRate float32  `json:"memory_allocation_rate"`
	GPUSlots             int      `json:"gpu_slots"`
	GPUEncoderSlots      int      `json:"gpu_encoder_slots" yaml:"gpu_encoder_slots"`
	Tags                 []string `json:"tags" yaml:"tags"`
	Unmanaged            bool     `json:"unmanaged" yaml:"unmanaged"`
	StoragePool          string   `json:"storage_pool" yaml:"storage_pool"`
	NetworkName          string   `json:"network_name" yaml:"network_name"`
	NetworkSubnet        string   `json:"network_subnet" yaml:"network_subnet"`
	NetworkACLName       string   `json:"network_acl_name" yaml:"network_acl_name"`
}

// NodeGPUPatch allows changing configuration for individual GPUs
type NodeGPUPatch struct {
	// ID must match the ID of an existing GPU and is used for identifying
	// the GPU which should be patched
	ID           uint64 `json:"id" yaml:"id"`
	Slots        *int   `json:"slots" yaml:"slots"`
	EncoderSlots *int   `json:"encoder_slots" yaml:"encoder_slots"`
}

// NodePatch describes a request to update an existing node
type NodePatch struct {
	PublicAddress        *string        `json:"public_address"`
	CPUs                 *int           `json:"cpus"`
	CPUAllocationRate    *float32       `json:"cpu_allocation_rate"`
	Memory               *string        `json:"memory"`
	MemoryAllocationRate *float32       `json:"memory_allocation_rate"`
	GPUSlots             *int           `json:"gpu_slots"`
	GPUEncoderSlots      *int           `json:"gpu_encoder_slots" yaml:"gpu_encoder_slots"`
	Tags                 *[]string      `json:"tags" yaml:"tags"`
	Unschedulable        *bool          `json:"unscheduable" yaml:"unscheduable"`
	GPUs                 []NodeGPUPatch `json:"gpus" yaml:"gpus"`
	Subnet               *string        `json:"subnet" yaml:"subnet"`
}

// NodeDelete describes a request used to delete a node
type NodeDelete struct {
	Force         bool `json:"force"`
	KeepInCluster bool `json:"keep_in_cluster"`
}
