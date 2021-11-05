// -*- Mode: Go; indent-tabs-mode: t -*-
/*
 * This file is part of AMS SDK
 * Copyright 2021 Canonical Ltd.
 *
 * This program is free software: you can redistribute it and/or modify it under
 * the terms of the GNU Lesser General Public License version 3, as published
 * by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranties of MERCHANTABILITY, SATISFACTORY
 * QUALITY, or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public
 * License for more details.
 *
 * You should have received a copy of the Lesser GNU General Public License along
 * with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/anbox-cloud/ams-sdk/client"
	"github.com/anbox-cloud/ams-sdk/examples"
)

type containerListLogCmd struct {
	examples.ConnectionCmd
	containerID string
}

func (command *containerListLogCmd) Parse() {
	command.ConnectionCmd.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		fmt.Println("  <container_id>      string")
		os.Exit(1)
	}

	command.containerID = flag.Args()[0]
}

func main() {
	cmd := &containerListLogCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	if err := listContainerLogs(c, cmd.containerID); err != nil {
		log.Fatal(err)
	}
}

func listContainerLogs(c client.Client, id string) error {
	container, _, err := c.RetrieveContainerByID(id)
	if err != nil {
		return err
	}

	examples.DumpData(container.StoredLogs)
	return nil
}
