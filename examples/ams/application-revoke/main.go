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

	"github.com/anbox-cloud/ams-sdk/examples/ams/common"
	"github.com/anbox-cloud/ams-sdk/pkg/ams/client"
)

type appRevokeCmd struct {
	common.ConnectionCmd
	id      string
	version string
}

func (command *appRevokeCmd) Parse() {
	flag.StringVar(&command.id, "id", "", "Application id")
	flag.StringVar(&command.version, "version", "", "Application version to revoke")

	command.ConnectionCmd.Parse()

	if len(command.id) == 0 || len(command.version) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	cmd := &appRevokeCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	version, err := strconv.Atoi(cmd.version)
	if err != nil {
		log.Fatal(err)
	}

	if err := revokeApplicationVersion(c, cmd.id, version); err != nil {
		log.Fatal(err)
	}
}

func revokeApplicationVersion(c client.Client, id string, version int) error {
	operation, err := c.RevokeApplicationVersion(id, version)
	if err != nil {
		return err
	}

	// Wait for the revoke operation to finish
	return operation.Wait(context.Background())
}
