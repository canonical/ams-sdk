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

	"github.com/anbox-cloud/ams-sdk/examples/ams/common"
	"github.com/anbox-cloud/ams-sdk/pkg/ams/client"
)

type containerDeleteCmd struct {
	common.ConnectionCmd
	id    string
	force bool
}

func (command *containerDeleteCmd) Parse() {
	flag.StringVar(&command.id, "id", "", "Container id")
	flag.BoolVar(&command.force, "force", false, "Force the removal of the container")

	command.ConnectionCmd.Parse()

	if len(command.id) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	cmd := &containerDeleteCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	if err := deleteContainer(c, cmd.id, cmd.force); err != nil {
		log.Fatal(err)
	}
}

func deleteContainer(c client.Client, id string, force bool) error {
	operation, err := c.DeleteContainerByID(id, force)
	if err != nil {
		return err
	}

	// Wait for delete operation to finish
	return operation.Wait(context.Background())
}
