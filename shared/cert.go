// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package shared

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"time"
)

const (
	// DefaultTransportTimeout stands for the default value for client requests to wait for a reply
	DefaultTransportTimeout = 30 * time.Second
)

// CertFingerprint returns the sha256 fingerprint of the given x509.Certificate
func CertFingerprint(cert *x509.Certificate) string {
	return fmt.Sprintf("%x", sha256.Sum256(cert.Raw))
}

// CertFingerprintStr decodes the given string to x509.Certificate and then generates a sha256 fingerprint
func CertFingerprintStr(c string) (string, error) {
	pemCertificate, _ := pem.Decode([]byte(c))
	if pemCertificate == nil {
		return "", fmt.Errorf("invalid certificate")
	}

	cert, err := x509.ParseCertificate(pemCertificate.Bytes)
	if err != nil {
		return "", err
	}

	return CertFingerprint(cert), nil
}

// GetRemoteCertificate connects to the server and returns its certificate
func GetRemoteCertificate(address string) (*x509.Certificate, error) {
	// Setup a permissive TLS config
	tlsConfig, err := GetTLSConfig("", "", "", nil)
	if err != nil {
		return nil, err
	}

	tlsConfig.InsecureSkipVerify = true
	tr := &http.Transport{
		TLSClientConfig:   tlsConfig,
		Dial:              RFC3493Dialer,
		DisableKeepAlives: true,
		IdleConnTimeout:   30 * time.Second,
		MaxIdleConns:      1,
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   DefaultTransportTimeout,
	}
	resp, err := client.Get(address)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Retrieve the certificate
	if resp.TLS == nil || len(resp.TLS.PeerCertificates) == 0 {
		return nil, fmt.Errorf("unable to read remote TLS certificate")
	}

	return resp.TLS.PeerCertificates[0], nil
}
