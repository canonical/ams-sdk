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
