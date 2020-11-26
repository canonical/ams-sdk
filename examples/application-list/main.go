// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package main

import (
	"log"

	"github.com/anbox-cloud/ams-sdk/client"
	"github.com/anbox-cloud/ams-sdk/examples"
)

type appListCmd struct {
	examples.ConnectionCmd
}

func main() {
	cmd := &appListCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	if err := listApplications(c); err != nil {
		log.Fatal(err)
	}
}

func listApplications(c client.Client) error {
	apps, err := c.ListApplications()
	if err != nil {
		return err
	}

	return examples.DumpData(apps)
}
