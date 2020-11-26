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

type appDeleteCmd struct {
	examples.ConnectionCmd
	id      string
	version string
}

func (command *appDeleteCmd) Parse() {
	flag.StringVar(&command.id, "id", "", "Application id")
	flag.StringVar(&command.version, "version", "", "Application id")

	command.ConnectionCmd.Parse()

	if len(command.id) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	cmd := &appDeleteCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	var version int
	var err error
	if len(cmd.version) > 0 {
		version, err = strconv.Atoi(cmd.version)
		if err == nil {
			err = deleteApplicationVersion(c, cmd.id, version)
		}
	} else {
		err = deleteApplication(c, cmd.id)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func deleteApplication(c client.Client, id string) error {
	operation, err := c.DeleteApplicationByID(id, false)
	if err != nil {
		return err
	}

	// Wait for delete operation to finish
	return operation.Wait(context.Background())
}

func deleteApplicationVersion(c client.Client, id string, version int) error {
	operation, err := c.DeleteApplicationVersion(id, version, false)
	if err != nil {
		return err
	}

	// Wait for delete operation to finish
	return operation.Wait(context.Background())
}
