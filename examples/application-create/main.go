// -*- Mode: Go; indent-tabs-mode: t -*-
/*
 * This file is part of AMS SDK
 * Copyright 2021 Canonical Ltd.
 *
 * This program is free software: you can redistribute it and/or modify it under
 * the terms of the GNU Lesser General Public License version 3, as published
 * by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranties of MERCHANTABILITY, SATISFACTORY
 * QUALITY, or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public
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

	"github.com/anbox-cloud/ams-sdk/client"
	"github.com/anbox-cloud/ams-sdk/examples"
)

type appCreateCmd struct {
	examples.ConnectionCmd
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

	operation, err := c.CreateApplication(packagePath, sentBytesCh)
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

	// Wait for create operation to finish
	if err := operation.Wait(context.Background()); err != nil {
		return err
	}

	// Print to stdout the generated resources path
	examples.PrintCreated(operation.Get().Resources)
	return nil
}
