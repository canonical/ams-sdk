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
	"strconv"

	api "github.com/anbox-cloud/ams-sdk/api/ams"
	"github.com/anbox-cloud/ams-sdk/examples/ams/common"
	"github.com/anbox-cloud/ams-sdk/pkg/ams/client"
)

type imageDeleteCmd struct {
	common.ConnectionCmd
	id      string
	version string
}

func (command *imageDeleteCmd) Parse() {
	flag.StringVar(&command.id, "id", "", "ID or name of the image")
	flag.StringVar(&command.version, "version", "", "Version of the image to delete")

	command.ConnectionCmd.Parse()

	if len(command.id) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	cmd := &imageDeleteCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	var version int
	var err error
	if len(cmd.version) > 0 {
		version, err = strconv.Atoi(cmd.version)
		if err == nil {
			err = deleteImageVersion(c, cmd.id, version)
		}
	} else {
		err = deleteImage(c, cmd.id)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func deleteImage(c client.Client, id string) error {
	operation, err := c.DeleteImageByIDOrName(id, false, api.ImageTypeAny)
	if err != nil {
		return err
	}

	// Wait for delete operation to finish
	return operation.Wait(context.Background())
}

func deleteImageVersion(c client.Client, id string, version int) error {
	operation, err := c.DeleteImageVersion(id, version)
	if err != nil {
		return err
	}

	// Wait for delete operation to finish
	return operation.Wait(context.Background())
}
