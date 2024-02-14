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

package client

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/anbox-cloud/ams-sdk/internal/ams/shared/rest/api"
)

type certificates struct {
	Client
}

// UpgradeToCertificatesClient wraps generic client to implement Operations interface operations
func UpgradeToCertificatesClient(c Client) Certificates {
	return &certificates{c}
}

func (c *certificates) ListCertificates() ([]api.Certificate, error) {
	certs := []api.Certificate{}
	_, err := c.QueryStruct("GET", APIPath("certificates"), nil, nil, nil, "", &certs)
	return certs, err
}

func (c *certificates) AddCertificate(base64PublicKey, trustPassword string) error {
	req := api.CertificatesPost{
		Certificate:   base64PublicKey,
		TrustPassword: trustPassword,
	}

	b, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("Could not marshal request body: %v", err)
	}

	_, _, err = c.CallAPI("POST", APIPath("certificates"), nil, nil, bytes.NewReader(b), "")
	return err
}

func (c *certificates) RetrieveCertificate(fingerprint string) (*api.Certificate, error) {
	cert := &api.Certificate{}
	_, err := c.QueryStruct("GET", APIPath("certificates", fingerprint), nil, nil, nil, "", &cert)
	return cert, err
}

func (c *certificates) DeleteCertificate(fingerprint string) (Operation, error) {
	op, _, err := c.QueryOperation("DELETE", APIPath("certificates", fingerprint), nil, nil, nil, "")
	return op, err
}
