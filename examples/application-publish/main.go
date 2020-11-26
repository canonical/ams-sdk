// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/anbox-cloud/ams-sdk/client"
	"github.com/anbox-cloud/ams-sdk/examples"
)

type appPublishCmd struct {
	examples.ConnectionCmd
	id      string
	version string
}

func (command *appPublishCmd) Parse() {
	flag.StringVar(&command.id, "id", "", "Application id")
	flag.StringVar(&command.version, "version", "", "Application version to publish")

	command.ConnectionCmd.Parse()

	if len(command.id) == 0 || len(command.version) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	cmd := &appPublishCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	version, err := strconv.Atoi(cmd.version)
	if err != nil {
		log.Fatal(err)
	}

	if err := publishApplicationVersion(c, cmd.id, version); err != nil {
		log.Fatal(err)
	}
}

func publishApplicationVersion(c client.Client, id string, version int) error {
	operation, err := c.PublishApplicationVersion(id, version)
	if err != nil {
		return err
	}

	// Wait for publish operation to finish
	return operation.Wait(context.Background())
}
