// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2019 Canonical Ltd.  All rights reserved.

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	errs "github.com/anbox-cloud/ams-sdk/shared/errors"
	"github.com/anbox-cloud/ams-sdk/shared/rest/api"
	"github.com/anbox-cloud/ams-sdk/shared/rest/client"
)

// ListCertificates lists all available certificates the AMS service knows about
func (c *clientImpl) ListCertificates() ([]api.Certificate, error) {
	params := client.QueryParams{
		"recursion": "1",
	}
	var certs []api.Certificate
	_, err := c.QueryStruct("GET", client.APIPath("certificates"), params, nil, nil, "", &certs)
	return certs, err
}

// AddCertificate adds a new certificate to the service trust store
func (c *clientImpl) AddCertificate(details *api.CertificatesPost) (*api.Response, error) {
	if len(details.Certificate) == 0 {
		return nil, fmt.Errorf("No certificate specified")
	}

	b, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}

	resp, _, err := c.CallAPI("POST", client.APIPath("certificates"), nil, nil, bytes.NewReader(b), "")
	return resp, err
}

// DeleteCertificate deletes an existing trusted certificate by its fingerprint
func (c *clientImpl) DeleteCertificate(fingerprint string) error {
	if len(fingerprint) == 0 {
		return errs.NewInvalidArgument("fingerprint")
	}
	op, _, err := c.QueryOperation("DELETE", client.APIPath("certificates", fingerprint), nil, nil, nil, "")
	if err != nil {
		return err
	}
	return op.Wait(context.Background())
}
