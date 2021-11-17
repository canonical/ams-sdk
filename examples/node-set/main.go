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
	"errors"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/anbox-cloud/ams-sdk/api"
	"github.com/anbox-cloud/ams-sdk/client"
	"github.com/anbox-cloud/ams-sdk/examples"
	"github.com/anbox-cloud/ams-sdk/shared"
)

type nodeSetCmd struct {
	examples.ConnectionCmd
	name  string
	key   string
	value string
}

func (command *nodeSetCmd) Parse() {
	flag.StringVar(&command.name, "name", "", "Name of the node")
	flag.StringVar(&command.key, "key", "", "Key of the property to set")
	flag.StringVar(&command.value, "value", "", "Value of the property to set")

	command.ConnectionCmd.Parse()

	if len(command.name) == 0 || len(command.key) == 0 ||
		len(command.value) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	cmd := &nodeSetCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	if err := updateNode(c, cmd.name, cmd.key, cmd.value); err != nil {
		log.Fatal(err)
	}
}

func updateNode(c client.Client, node, key, value string) error {
	details := &api.NodePatch{}
	switch key {
	case "public-address":
		details.PublicAddress = &value
	case "cpus":
		v, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		details.CPUs = &v
	case "cpu-allocation-rate":
		v, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}
		v32 := float32(v)
		details.CPUAllocationRate = &v32
	case "memory":
		_, err := shared.ParseByteSizeString(value)
		if err != nil {
			return err
		}
		details.Memory = &value
	case "memory-allocation-rate":
		v, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}
		v32 := float32(v)
		details.MemoryAllocationRate = &v32
	default:
		return errors.New("Unknown configuration item")
	}

	operation, err := c.UpdateNode(node, details)
	if err != nil {
		return err
	}

	// Wait for add operation to finish
	return operation.Wait(context.Background())
}
