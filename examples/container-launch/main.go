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

	"github.com/anbox-cloud/ams-sdk/api"
	"github.com/anbox-cloud/ams-sdk/client"
	"github.com/anbox-cloud/ams-sdk/examples"
	"github.com/anbox-cloud/ams-sdk/shared/errors"
)

type containerLaunchCmd struct {
	examples.ConnectionCmd
	id           string
	version      string
	node         string
	instanceType string
	raw          bool
}

const (
	defaultRawContainerInstanceType = "a2.3"
)

func (command *containerLaunchCmd) Parse() {
	flag.StringVar(&command.id, "id", "", "Application id or image id(name) if launching a container from a raw image")
	flag.StringVar(&command.version, "version", "", "Application or image version to be used to launch a container")
	flag.StringVar(&command.node, "node", "", "In which node to launch the container")
	flag.StringVar(&command.instanceType, "instance-type", defaultRawContainerInstanceType, "Instance type to use when launching a container from an image instead of an application")
	flag.BoolVar(&command.raw, "raw", false, "Launched a container from a specific image instead of an application if it's set to true")

	command.ConnectionCmd.Parse()

	if len(command.id) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	cmd := &containerLaunchCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	details := &api.ContainersPost{}

	version := 0
	var err error
	if len(cmd.version) > 0 {
		version, err = strconv.Atoi(cmd.version)
		if err != nil {
			log.Fatal(err)
		}
	}

	if !cmd.raw && len(cmd.instanceType) > 0 {
		log.Fatal(errors.NewInvalidArgument("instance type"))
	}

	if cmd.raw {
		details.ImageID = cmd.id
		details.ImageVersion = &version
		details.InstanceType = cmd.instanceType

		if len(cmd.instanceType) == 0 {
			details.InstanceType = defaultRawContainerInstanceType
		}
	} else {
		details.ApplicationID = cmd.id
		details.ApplicationVersion = &version
	}
	details.Node = cmd.node

	if err := launchContainer(c, details); err != nil {
		log.Fatal(err)
	}
}

func launchContainer(c client.Client, details *api.ContainersPost) error {
	operation, err := c.LaunchContainer(details, false)
	if err != nil {
		return err
	}

	// Wait for launch operation to finish
	if err := operation.Wait(context.Background()); err != nil {
		return err
	}

	examples.PrintCreated(operation.Get().Resources)
	return nil
}
