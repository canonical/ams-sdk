// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/anbox-cloud/ams-sdk/client"
	"github.com/anbox-cloud/ams-sdk/examples"
)

type nodeDeleteCmd struct {
	examples.ConnectionCmd
	name string
}

func (command *nodeDeleteCmd) Parse() {
	flag.StringVar(&command.name, "name", "", "Name of the node")

	command.ConnectionCmd.Parse()

	if len(command.name) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	cmd := &nodeDeleteCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	if err := deleteNode(c, cmd.name); err != nil {
		log.Fatal(err)
	}
}

func deleteNode(c client.Client, name string) error {
	operation, err := c.RemoveNode(name, false, false)
	if err != nil {
		return err
	}

	// Wait for delete operation to finish
	return operation.Wait(context.Background())
}
