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
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/anbox-cloud/ams-sdk/examples/ams/common"
	"github.com/anbox-cloud/ams-sdk/pkg/ams/client"
)

type containerLogCmd struct {
	common.ConnectionCmd
	id      string
	logName string
}

func (command *containerLogCmd) Parse() {
	flag.StringVar(&command.id, "id", "", "Container id")
	flag.StringVar(&command.logName, "log-name", "", "Log name to be fetched from the container")

	command.ConnectionCmd.Parse()

	if len(command.id) == 0 || len(command.logName) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	cmd := &containerLogCmd{}
	cmd.Parse()
	c := cmd.NewClient()

	if err := showContainerLog(c, cmd.id, cmd.logName); err != nil {
		log.Fatal(err)
	}
}

func showContainerLog(c client.Client, id, logName string) error {
	return c.RetrieveContainerLog(id, logName, func(header *http.Header, body io.ReadCloser) error {
		content, err := ioutil.ReadAll(body)
		if err != nil {
			return err
		}
		fmt.Printf(string(content))
		return nil
	})
}
