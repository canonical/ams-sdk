// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package client

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/anbox-cloud/ams-sdk/shared/rest/api"
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
