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

	"github.com/anbox-cloud/ams-sdk/client"
	"github.com/anbox-cloud/ams-sdk/examples"
)

type imageAddCmd struct {
	examples.ConnectionCmd
	name        string
	packagePath string
}

func (command *imageAddCmd) Parse() {
	flag.StringVar(&command.name, "name", "", "Name of the image")
	flag.StringVar(&command.packagePath, "package-path", "", "Path to the image package")

	command.ConnectionCmd.Parse()

	if len(command.name) == 0 || len(command.packagePath) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	cmd := &imageAddCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	if _, err := os.Stat(cmd.packagePath); err != nil {
		if os.IsNotExist(err) {
			log.Fatal("Specified iamge does not exist")
		}
		log.Fatal(err)
	}

	if err := addImage(c, cmd.name, cmd.packagePath); err != nil {
		log.Fatal(err)
	}
}

func addImage(c client.Client, name, sourcePath string) error {
	// We can read from this channel the number of bytes already
	// sent to the server. This is useful to identify when all the
	// payload has been sent.
	sentBytesCh := make(chan float64)

	operation, err := c.AddImage(name, sourcePath, false, sentBytesCh)
	if err != nil {
		return err
	}

	// Wait for add operation to finish
	if err := operation.Wait(context.Background()); err != nil {
		return err
	}

	examples.PrintCreated(operation.Get().Resources)
	return nil
}
