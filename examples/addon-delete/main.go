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

type addonDeleteCmd struct {
	examples.ConnectionCmd
	name    string
	version string
}

func (command *addonDeleteCmd) Parse() {
	flag.StringVar(&command.name, "name", "", "Name of the addon")
	flag.StringVar(&command.version, "version", "", "Version of the addon")

	command.ConnectionCmd.Parse()

	if len(command.name) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	cmd := &addonDeleteCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	var version int
	var err error
	if len(cmd.version) > 0 {
		version, err = strconv.Atoi(cmd.version)
		if err == nil {
			err = deleteAddonVersion(c, cmd.name, version)
		}
	} else {
		err = deleteAddon(c, cmd.name)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func deleteAddon(c client.Client, name string) error {
	operation, err := c.DeleteAddon(name)
	if err != nil {
		return err
	}

	// Wait for delete operation to finish
	return operation.Wait(context.Background())
}

func deleteAddonVersion(c client.Client, name string, version int) error {
	operation, err := c.DeleteAddonVersion(name, version)
	if err != nil {
		return err
	}

	// Wait for delete operation to finish
	return operation.Wait(context.Background())
}
