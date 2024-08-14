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
	"github.com/anbox-cloud/ams-sdk/pkg/ams/client"
)

type appUpdateCmd struct {
	common.ConnectionCmd
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
