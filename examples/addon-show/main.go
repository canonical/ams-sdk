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
	"flag"
	"log"
	"os"
	"time"

	"github.com/anbox-cloud/ams-sdk/client"
	"github.com/anbox-cloud/ams-sdk/examples"
)

type addonShowCmd struct {
	examples.ConnectionCmd
	name string
}

func (command *addonShowCmd) Parse() {
	flag.StringVar(&command.name, "name", "", "Name of the addon")
	command.ConnectionCmd.Parse()

	if len(command.name) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	cmd := &addonShowCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	if err := showAddon(c, cmd.name); err != nil {
		log.Fatal(err)
	}
}

func showAddon(c client.Client, name string) error {
	addon, _, err := c.RetrieveAddon(name)
	if err != nil {
		return err
	}

	type addonVersion struct {
		Size      string `yaml:"size" json:"size"`
		CreatedAt string `yaml:"created-at" json:"created-at"`
	}

	var outputData struct {
		Name     string               `yaml:"name" json:"name"`
		Versions map[int]addonVersion `yaml:"versions" json:"versions"`
	}
	outputData.Name = addon.Name
	outputData.Versions = make(map[int]addonVersion)

	for _, v := range addon.Versions {
		t := time.Unix(v.CreatedAt, 0)
		outputData.Versions[v.Number] = addonVersion{
			Size:      examples.GetByteSizeString(v.Size, 2),
			CreatedAt: t.String(),
		}
	}

	return examples.DumpData(outputData)
}
