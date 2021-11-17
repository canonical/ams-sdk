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
	"log"

	"github.com/anbox-cloud/ams-sdk/client"
	"github.com/anbox-cloud/ams-sdk/examples"
)

type imageListCmd struct {
	examples.ConnectionCmd
}

func main() {
	cmd := &imageListCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	if err := listImages(c); err != nil {
		log.Fatal(err)
	}
}

func listImages(c client.Client) error {
	images, err := c.ListImages()
	if err != nil {
		return err
	}

	return examples.DumpData(images)
}
