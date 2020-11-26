// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/anbox-cloud/ams-sdk/client"
	"github.com/anbox-cloud/ams-sdk/examples"
)

type containerLogCmd struct {
	examples.ConnectionCmd
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
