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

type imageDeleteCmd struct {
	examples.ConnectionCmd
	id      string
	version string
}

func (command *imageDeleteCmd) Parse() {
	flag.StringVar(&command.id, "id", "", "ID or name of the image")
	flag.StringVar(&command.version, "version", "", "Version of the image to delete")

	command.ConnectionCmd.Parse()

	if len(command.id) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	cmd := &imageDeleteCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	var version int
	var err error
	if len(cmd.version) > 0 {
		version, err = strconv.Atoi(cmd.version)
		if err == nil {
			err = deleteImageVersion(c, cmd.id, version)
		}
	} else {
		err = deleteImage(c, cmd.id)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func deleteImage(c client.Client, id string) error {
	operation, err := c.DeleteImageByIDOrName(id, false)
	if err != nil {
		return err
	}

	// Wait for delete operation to finish
	return operation.Wait(context.Background())
}

func deleteImageVersion(c client.Client, id string, version int) error {
	operation, err := c.DeleteImageVersion(id, version)
	if err != nil {
		return err
	}

	// Wait for delete operation to finish
	return operation.Wait(context.Background())
}
