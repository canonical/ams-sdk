// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2019 Canonical Ltd.  All rights reserved.

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
