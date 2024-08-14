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
	"flag"
	"log"
	"os"
	"time"

	api "github.com/anbox-cloud/ams-sdk/api/ams"
	"github.com/anbox-cloud/ams-sdk/examples/ams/common"
	"github.com/anbox-cloud/ams-sdk/pkg/ams/client"
	"github.com/anbox-cloud/ams-sdk/pkg/ams/shared"
)

type imageShowCmd struct {
	common.ConnectionCmd
	id string
}

func (command *imageShowCmd) Parse() {
	flag.StringVar(&command.id, "id", "", "Id or name of image")

	command.ConnectionCmd.Parse()

	if len(command.id) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	cmd := &imageShowCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	if err := showImage(c, cmd.id); err != nil {
		log.Fatal(err)
	}
}

func showImage(c client.Client, id string) error {
	image, _, err := c.RetrieveImageByIDOrName(id, api.ImageTypeAny)
	if err != nil {
		return err
	}

	type imageVersion struct {
		Size      string `yaml:"size" json:"size"`
		CreatedAt string `yaml:"created-at" json:"created-at"`
		Status    string `yaml:"status" json:"status"`
	}

	var outputData struct {
		ID       string               `yaml:"id" json:"id"`
		Name     string               `yaml:"name" json:"name"`
		Status   string               `yaml:"status" json:"status"`
		Versions map[int]imageVersion `yaml:"versions" json:"versions"`
	}
	outputData.ID = image.ID
	outputData.Name = image.Name
	outputData.Status = image.Status
	outputData.Versions = make(map[int]imageVersion)

	for _, v := range image.Versions {
		t := time.Unix(v.CreatedAt, 0)
		outputData.Versions[v.Number] = imageVersion{
			Size:      shared.GetByteSizeString(v.Size, 2),
			CreatedAt: t.String(),
			Status:    v.Status,
		}
	}

	return common.DumpData(outputData)
}
