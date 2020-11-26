// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/anbox-cloud/ams-sdk/client"
	"github.com/anbox-cloud/ams-sdk/examples"
)

type appShowCmd struct {
	examples.ConnectionCmd
	id string
}

func (command *appShowCmd) Parse() {
	flag.StringVar(&command.id, "id", "", "Application id")

	command.ConnectionCmd.Parse()

	if len(command.id) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	cmd := &appShowCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	if err := showApplication(c, cmd.id); err != nil {
		log.Fatal(err)
	}
}

func showApplication(c client.Client, id string) error {
	app, _, err := c.RetrieveApplicationByID(id)
	if err != nil {
		return err
	}

	type appExtraData struct {
		Target      string `json:"target" yaml:"target"`
		Owner       string `json:"owner" yaml:"owner"`
		Permissions string `json:"permissions" yaml:"permissions"`
	}

	type appVersion struct {
		Image               string                  `yaml:"image" json:"image"`
		Published           bool                    `yaml:"published" json:"published"`
		Status              string                  `yaml:"status" json:"status"`
		Addons              []string                `yaml:"addons" json:"addons"`
		BootActivity        string                  `yaml:"boot-activity" json:"boot-activity"`
		RequiredPermissions []string                `yaml:"required-permissions" json:"required-permissions"`
		ExtraData           map[string]appExtraData `yaml:"extra-data" json:"extra-data"`
	}

	var outputData struct {
		ID        string `yaml:"id" json:"id"`
		Name      string `yaml:"name" json:"name"`
		Status    string `yaml:"status" json:"status"`
		Published bool   `yaml:"published" json:"published"`
		Config    struct {
			InstanceType string `yaml:"instance-type" json:"instance-type"`
			BootPakcage  string `yaml:"boot-package" json:"boot-package"`
		} `yaml:"config" json:"config"`
		Versions map[int]appVersion `yaml:"versions" json:"versions"`
	}
	outputData.ID = app.ID
	outputData.Name = app.Name
	outputData.Status = app.Status
	outputData.Published = app.Published
	outputData.Config.InstanceType = app.InstanceType
	outputData.Config.BootPakcage = app.BootPackage
	outputData.Versions = make(map[int]appVersion)
	for _, v := range app.Versions {
		newVersion := appVersion{
			Image:               fmt.Sprintf("%s (version %d)", app.ParentImageID, v.ParentImageVersion),
			Published:           v.Published,
			Status:              v.Status,
			BootActivity:        v.BootActivity,
			RequiredPermissions: v.RequiredPermissions,
		}

		newVersion.ExtraData = make(map[string]appExtraData)
		for name, value := range v.ExtraData {
			extraData := appExtraData{
				Target:      value.Target,
				Owner:       value.Owner,
				Permissions: value.Permissions,
			}

			newVersion.ExtraData[name] = extraData
		}

		for _, addon := range v.Addons {
			newVersion.Addons = append(newVersion.Addons, fmt.Sprintf("%s (version %d)", addon.Name, addon.Version))
		}

		outputData.Versions[v.Number] = newVersion
	}

	return examples.DumpData(outputData)
}
