// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

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
