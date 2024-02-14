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

package common

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/anbox-cloud/ams-sdk/internal/ams/client"
	"github.com/anbox-cloud/ams-sdk/internal/ams/shared"
)

// ConnectionCmd defines the options for an example to connect to the service
type ConnectionCmd struct {
	ClientCert string
	ClientKey  string
	ServiceURL string
}

// Parse parses command line arguments
func (c *ConnectionCmd) Parse() {
	flag.StringVar(&c.ClientCert, "cert", "", "Path to the file with the client certificate to use to connect to AMS")
	flag.StringVar(&c.ClientKey, "key", "", "Path to the file with the client key to use to connect to AMS")
	flag.StringVar(&c.ServiceURL, "url", "", "URL of the AMS server")

	flag.Parse()

	if err := c.Validate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Validate returns error if provided two way SSL parameters are invalid
func (c *ConnectionCmd) Validate() error {
	if len(c.ServiceURL) == 0 {
		return fmt.Errorf("Please provide a service URL")
	}

	if len(c.ClientCert) == 0 {
		return fmt.Errorf("Please provide a certificate path")
	}

	if len(c.ClientKey) == 0 {
		return fmt.Errorf("Please provide a certificate key path")
	}

	if _, err := os.Stat(c.ClientCert); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Specified client cert does not exist")
		}
		return err
	}

	if _, err := os.Stat(c.ClientKey); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Specified client cert key does not exist")
		}
		return err
	}

	return nil
}

// NewClient returns a REST client to connect to AMS
func (c *ConnectionCmd) NewClient() client.Client {
	// Server URL is accessible from client
	u, err := url.Parse(c.ServiceURL)
	if err != nil {
		log.Fatal(err)
	}

	serverCert, err := shared.GetRemoteCertificate(c.ServiceURL)
	if err != nil {
		log.Fatal(err)
	}

	// Asuming that client cert, client key and server cert files exist and are valid
	// certificate files.
	// Server must have client cert amongst trusted client certificates before
	// connecting
	tlsConfig, err := shared.GetTLSConfig(c.ClientCert, c.ClientKey, "", serverCert)
	if err != nil {
		log.Fatal(err)
	}

	amsClient, err := client.New(u, tlsConfig)
	if err != nil {
		log.Fatal(err)
	}

	return amsClient
}
