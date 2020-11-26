// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package main

import (
	"log"

	"github.com/anbox-cloud/ams-sdk/client"
	"github.com/anbox-cloud/ams-sdk/examples"
)

type imageListCmd struct {
	examples.ConnectionCmd
}

func main() {
	cmd := &imageListCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	if err := listImages(c); err != nil {
		log.Fatal(err)
	}
}

func listImages(c client.Client) error {
	images, err := c.ListImages()
	if err != nil {
		return err
	}

	return examples.DumpData(images)
}
