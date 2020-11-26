// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package main

import (
	"log"

	"github.com/anbox-cloud/ams-sdk/client"
	"github.com/anbox-cloud/ams-sdk/examples"
)

type addonListCmd struct {
	examples.ConnectionCmd
}

func main() {
	cmd := &addonListCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	if err := listAddons(c); err != nil {
		log.Fatal(err)
	}
}

func listAddons(c client.Client) error {
	addons, err := c.ListAddons()
	if err != nil {
		return err
	}

	return examples.DumpData(addons)
}
