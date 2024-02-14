// -*- Mode: Go; indent-tabs-mode: t -*-
/*
 * This file is part of AMS SDK
 * Copyright 2021 Canonical Ltd.
 *
 * This program is free software: you can redistribute it and/or modify it under
 * the terms of the Lesser GNU General Public License version 3, as published
 * by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranties of MERCHANTABILITY, SATISFACTORY
 * QUALITY, or FITNESS FOR A PARTICULAR PURPOSE.  See the Lesser GNU General Public
 * License for more details.
 *
 * You should have received a copy of the Lesser GNU General Public License along
 * with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/anbox-cloud/ams-sdk/examples/ams/common"
	"github.com/anbox-cloud/ams-sdk/internal/ams/client"
)

type appCreateCmd struct {
	common.ConnectionCmd
	packagePath string
}

func (command *appCreateCmd) Parse() {
	flag.StringVar(&command.packagePath, "package-path", "", "Path to the application package")

	command.ConnectionCmd.Parse()

	if len(command.packagePath) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	cmd := &appCreateCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	if err := createApplication(c, cmd.packagePath); err != nil {
		log.Fatal(err)
	}
}

func createApplication(c client.Client, packagePath string) error {
	// We can read from this channel the number of bytes already
	// sent to the server. This is useful to identify when all the
	// payload has been sent.
	sentBytesCh := make(chan float64)
	closeCh := make(chan struct{})

	go func() {
		for {
			select {
			case n := <-sentBytesCh:
				// Simply print sent bytes to stdout
				log.Printf("Sent: %v", n)
			case <-closeCh:
				return
			}
		}
	}()

	operation, err := c.CreateApplication(packagePath, sentBytesCh)
	if err != nil {
		return err
	}

	// Wait for create operation to finish
	if err := operation.Wait(context.Background()); err != nil {
		return err
	}

	// Print to stdout the generated resources path
	common.PrintCreated(operation.Get().Resources)
	return nil
}
