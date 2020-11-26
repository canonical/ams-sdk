// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package shared

import (
	"crypto/x509"
	"fmt"
	"net/http"
	"time"
)

const (
	// DefaultTransportTimeout stands for the default value for client requests to wait for a reply
	DefaultTransportTimeout = 30 * time.Second
)

// GetRemoteCertificate connects to the server and returns its certificate
func GetRemoteCertificate(address string) (*x509.Certificate, error) {
	// Setup a permissive TLS config
	tlsConfig, err := GetTLSConfig("", "", "", nil)
	if err != nil {
		return nil, err
	}

	tlsConfig.InsecureSkipVerify = true
	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
		Dial:            RFC3493Dialer,
	}

	// Connect
	client := &http.Client{
		Transport: tr,
		Timeout:   DefaultTransportTimeout,
	}
	resp, err := client.Get(address)
	if err != nil {
		return nil, err
	}

	// Retrieve the certificate
	if resp.TLS == nil || len(resp.TLS.PeerCertificates) == 0 {
		return nil, fmt.Errorf("Unable to read remote TLS certificate")
	}

	return resp.TLS.PeerCertificates[0], nil
}
