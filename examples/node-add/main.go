// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/anbox-cloud/ams-sdk/api"
	"github.com/anbox-cloud/ams-sdk/client"
	"github.com/anbox-cloud/ams-sdk/examples"
)

type nodeAddCmd struct {
	examples.ConnectionCmd
	name             string
	address          string
	trustPassword    string
	storageDevice    string
	networkBridgeMTU int
}

func (command *nodeAddCmd) Parse() {
	flag.StringVar(&command.name, "name", "", "Name of the node")
	flag.StringVar(&command.address, "address", "", "Address where the node is accessible")
	flag.StringVar(&command.trustPassword, "trust-password", "", "Trust password for the remote LXD node")
	flag.StringVar(&command.storageDevice, "storage-device", "", "Storage device LXD node should use")
	flag.IntVar(&command.networkBridgeMTU, "network-bridge-mtu", 1500, "MTU of the network bridge configured for LXD")

	command.ConnectionCmd.Parse()

	if len(command.name) == 0 || len(command.address) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	cmd := &nodeAddCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	node := &api.NodesPost{
		Name:             cmd.name,
		TrustPassword:    cmd.trustPassword,
		Address:          cmd.address,
		StorageDevice:    cmd.storageDevice,
		NetworkBridgeMTU: cmd.networkBridgeMTU,
	}

	if err := addNode(c, node); err != nil {
		log.Fatal(err)
	}
}

func addNode(c client.Client, node *api.NodesPost) error {
	operation, err := c.AddNode(node)
	if err != nil {
		return err
	}

	// Wait for add operation to finish
	err = operation.Wait(context.Background())
	if err != nil {
		return err
	}

	examples.PrintCreated(operation.Get().Resources)
	return nil
}
