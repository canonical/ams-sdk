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

type imageUpdateCmd struct {
	examples.ConnectionCmd
	name        string
	packagePath string
}

func (command *imageUpdateCmd) Parse() {
	flag.StringVar(&command.name, "name", "", "Name of the image")
	flag.StringVar(&command.packagePath, "package-path", "", "Path to the image package")

	command.ConnectionCmd.Parse()

	if len(command.name) == 0 || len(command.packagePath) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	cmd := &imageUpdateCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	if _, err := os.Stat(cmd.packagePath); err != nil {
		if os.IsNotExist(err) {
			log.Fatal("Specified image does not exist")
		}
		log.Fatal(err)
	}

	if err := updateImage(c, cmd.name, cmd.packagePath); err != nil {
		log.Fatal(err)
	}
}

func updateImage(c client.Client, name, imagePath string) error {
	// We can read from this channel the number of bytes already
	// sent to the server. This is useful to identify when all the
	// payload has been sent.
	sentBytesCh := make(chan float64)

	operation, err := c.UpdateImage(name, imagePath, sentBytesCh)
	if err != nil {
		log.Fatal(err)
	}

	// Wait for update operation to finish
	if err := operation.Wait(context.Background()); err != nil {
		log.Fatal(err)
	}

	examples.PrintCreated(operation.Get().Resources)
	return nil
}
