// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package main

import (
	"log"

	"github.com/anbox-cloud/ams-sdk/client"
	"github.com/anbox-cloud/ams-sdk/examples"
)

type containerListCmd struct {
	examples.ConnectionCmd
}

func main() {
	cmd := &containerListCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	if err := listContainers(c); err != nil {
		log.Fatal(err)
	}
}

func listContainers(c client.Client) error {
	containers, err := c.ListContainers()
	if err != nil {
		return err
	}

	return examples.DumpData(containers)
}
