// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/anbox-cloud/ams-sdk/api"
	"github.com/anbox-cloud/ams-sdk/client"
	"github.com/anbox-cloud/ams-sdk/examples"
)

type containerShowCmd struct {
	examples.ConnectionCmd
	id string
}

func (command *containerShowCmd) Parse() {
	flag.StringVar(&command.id, "id", "", "Container id")

	command.ConnectionCmd.Parse()

	if len(command.id) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	cmd := &containerShowCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	if err := showContainer(c, cmd.id); err != nil {
		log.Fatal(err)
	}
}

func showContainer(c client.Client, id string) error {
	container, _, err := c.RetrieveContainerByID(id)
	if err != nil {
		return err
	}

	var outputData struct {
		ID        string    `yaml:"id" json:"id"`
		Name      string    `yaml:"name" json:"name"`
		Status    string    `yaml:"status" json:"status"`
		Node      string    `yaml:"node" json:"node"`
		CreatedAt time.Time `yaml:"created_at" json:"created_at"`

		Application struct {
			ID      string `yaml:"id" json:"id"`
			Version int    `yaml:"version" json:"version"`
		} `yaml:"application" json:"application"`

		Network struct {
			Address       string                 `yaml:"address" json:"address"`
			PublicAddress string                 `yaml:"public_address" json:"public_address"`
			Services      []api.ContainerService `yaml:"services" json:"services"`
		} `yaml:"network" json:"network"`

		StoredLogs   []string `json:"stored_logs" yaml:"stored_logs"`
		ErrorMessage string   `json:"error_message" yaml:"error_message"`
	}
	outputData.ID = container.ID
	outputData.Name = container.Name
	outputData.Status = container.Status
	outputData.Node = container.Node
	outputData.CreatedAt = time.Unix(container.CreatedAt, 0)
	outputData.Application.ID = container.AppID
	outputData.Application.Version = container.AppVersion
	outputData.Network.Address = container.Address
	outputData.Network.PublicAddress = container.PublicAddress
	outputData.Network.Services = container.Services
	outputData.StoredLogs = container.StoredLogs
	outputData.ErrorMessage = container.ErrorMessage

	return examples.DumpData(outputData)
}
