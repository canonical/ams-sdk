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

package main

import (
	"flag"
	"log"
	"os"

	"github.com/anbox-cloud/ams-sdk/client"
	"github.com/anbox-cloud/ams-sdk/examples"
)

type nodeShowCmd struct {
	examples.ConnectionCmd
	name string
}

func (command *nodeShowCmd) Parse() {
	flag.StringVar(&command.name, "name", "", "Name of the node")

	command.ConnectionCmd.Parse()

	if len(command.name) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	cmd := &nodeShowCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	if err := showNode(c, cmd.name); err != nil {
		log.Fatal(err)
	}
}

func showNode(c client.Client, name string) error {
	node, _, err := c.RetrieveNodeByName(name)
	if err != nil {
		return err
	}

	var outputData struct {
		Name    string `yaml:"name" json:"name"`
		Network struct {
			Address   string `yaml:"address" json:"address"`
			BridgeMTU int    `yaml:"bridge-mtu" json:"bridge-mtu"`
		} `yaml:"network" json:"network"`
		Status string `yaml:"status" json:"status"`
		Config struct {
			PublicAddress        string  `yaml:"public-address" json:"public-address"`
			CPUs                 int     `yaml:"cpus" json:"cpus"`
			CPUAllocationRate    float32 `yaml:"cpu-allocation-rate" json:"cpu-allocation-rate"`
			Memory               string  `yaml:"memory" json:"memory"`
			MemoryAllocationRate float32 `yaml:"memory-allocation-rate" json:"memory-allocation-rate"`
		} `yaml:"config" json:"config"`
	}
	outputData.Name = node.Name
	outputData.Network.Address = node.Address
	outputData.Network.BridgeMTU = node.NetworkBridgeMTU
	outputData.Status = node.Status
	outputData.Config.PublicAddress = node.PublicAddress
	outputData.Config.CPUs = node.CPUs
	outputData.Config.CPUAllocationRate = node.CPUAllocationRate
	outputData.Config.Memory = node.Memory
	outputData.Config.MemoryAllocationRate = node.MemoryAllocationRate

	return examples.DumpData(outputData)
}
