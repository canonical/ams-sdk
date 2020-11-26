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

type containerDeleteCmd struct {
	examples.ConnectionCmd
	id string
}

func (command *containerDeleteCmd) Parse() {
	flag.StringVar(&command.id, "id", "", "Container id")

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

	if err := deleteContainer(c, cmd.id); err != nil {
		log.Fatal(err)
	}
}

func deleteContainer(c client.Client, id string) error {
	operation, err := c.DeleteContainerByID(id)
	if err != nil {
		return err
	}

	// Wait for delete operation to finish
	return operation.Wait(context.Background())
}
