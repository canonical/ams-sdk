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

type appUpdateCmd struct {
	examples.ConnectionCmd
	id          string
	packagePath string
}

func (command *appUpdateCmd) Parse() {
	flag.StringVar(&command.id, "id", "", "Application id")
	flag.StringVar(&command.packagePath, "package-path", "", "Path to the application package")

	command.ConnectionCmd.Parse()

	if len(command.id) == 0 || len(command.packagePath) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	cmd := &appUpdateCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	if err := updateApplication(c, cmd.id, cmd.packagePath); err != nil {
		log.Fatal(err)
	}
}

func updateApplication(c client.Client, id, packagePath string) error {
	// We can read from this channel the number of bytes already
	// sent to the server. This is useful to identify when all the
	// payload has been sent.
	sentBytesCh := make(chan float64)

	operation, err := c.UpdateApplicationWithPackage(id, packagePath, sentBytesCh)
	if err != nil {
		return err
	}

	// Simply print sent bytes to stdout
	closeCh := make(chan struct{})
	go func() {
		for {
			select {
			case n := <-sentBytesCh:
				log.Printf("Sent: %v", n)
			case <-closeCh:
				return
			}
		}
	}()

	// Wait for update operation to finish
	return operation.Wait(context.Background())
}
