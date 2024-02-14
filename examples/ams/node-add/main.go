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
	"context"
	"flag"
	"log"
	"os"

	api "github.com/anbox-cloud/ams-sdk/api/ams"
	"github.com/anbox-cloud/ams-sdk/examples/ams/common"
	"github.com/anbox-cloud/ams-sdk/internal/ams/client"
)

type nodeAddCmd struct {
	common.ConnectionCmd
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

	common.PrintCreated(operation.Get().Resources)
	return nil
}
