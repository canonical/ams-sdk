// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package main

import (
	"log"

	"github.com/anbox-cloud/ams-sdk/client"
	"github.com/anbox-cloud/ams-sdk/examples"
)

type nodeListCmd struct {
	examples.ConnectionCmd
}

func main() {
	cmd := &nodeListCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	if err := listNodes(c); err != nil {
		log.Fatal(err)
	}
}

func listNodes(c client.Client) error {
	nodes, err := c.ListNodes()
	if err != nil {
		return err
	}

	return examples.DumpData(nodes)
}
