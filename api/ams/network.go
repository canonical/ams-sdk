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
