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

type addonAddCmd struct {
	examples.ConnectionCmd
	name      string
	addonPath string
}

func (command *addonAddCmd) Parse() {
	flag.StringVar(&command.name, "name", "", "Name of the addon")
	flag.StringVar(&command.addonPath, "path", "", "Path to the addon package")

	command.ConnectionCmd.Parse()

	if len(command.name) == 0 || len(command.addonPath) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	cmd := &addonAddCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	if _, err := os.Stat(cmd.addonPath); err != nil {
		if os.IsNotExist(err) {
			log.Fatal("Specified payload does not exist")
		}
		log.Fatal(err)
	}

	if err := addAddon(c, cmd.name, cmd.addonPath); err != nil {
		log.Fatal(err)
	}
}

func addAddon(c client.Client, name, sourcePath string) error {
	// We can read from this channel the number of bytes already
	// sent to the server. This is useful to identify when all the
	// payload has been sent.
	sentBytesCh := make(chan float64)

	operation, err := c.AddAddon(name, sourcePath, sentBytesCh)
	if err != nil {
		return err
	}

	// Wait for add operation to finish
	err = operation.Wait(context.Background())
	if err != nil {
		return err
	}

	examples.PrintCreated(operation.Get().Resources)
	return nil
}
