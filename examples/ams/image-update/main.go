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

type imageUpdateCmd struct {
	common.ConnectionCmd
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

	common.PrintCreated(operation.Get().Resources)
	return nil
}
