// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/anbox-cloud/ams-sdk/client"
	"github.com/anbox-cloud/ams-sdk/examples"
	"github.com/anbox-cloud/ams-sdk/shared"
)

type imageShowCmd struct {
	examples.ConnectionCmd
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
	image, _, err := c.RetrieveImageByIDOrName(id)
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

	return examples.DumpData(outputData)
}
