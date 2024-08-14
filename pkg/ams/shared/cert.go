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

package shared

import (
	"crypto/sha256"
	"crypto/tls"
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

	return GetRemoteCertificateWithTLSConfig(address, tlsConfig)
}

// GetRemoteCertificateWithTLSConfig connects to the server and returns its certificate
func GetRemoteCertificateWithTLSConfig(address string, tlsConfig *tls.Config) (*x509.Certificate, error) {
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
