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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/anbox-cloud/ams-sdk/pkg/ams/packages"
	"github.com/anbox-cloud/ams-sdk/pkg/ams/shared"
	errs "github.com/anbox-cloud/ams-sdk/pkg/ams/shared/errors"
	"github.com/anbox-cloud/ams-sdk/pkg/ams/shared/rest/client"
	"github.com/gorilla/websocket"
)

func (c *clientImpl) upload(httpOp, apiPath string, params client.QueryParams, packagePath string, details interface{}, sentBytes chan float64) (client.Operation, error) {
	hasZipSupport, err := c.HasExtension("zip_archive_support")
	if err != nil {
		return nil, err
	}
	if !hasZipSupport {
		if packages.IsZip(packagePath) {
			return nil, errs.NewErrNotSupported("api extension \"zip_archive_support\"")
		}
		pkgType, err := packages.DetectPackageType(packagePath)
		if err != nil {
			return nil, err
		}
		if pkgType == packages.PackageTypeZip {
			return nil, errs.NewErrNotSupported("api extension \"zip_archive_support\"")
		}
	}
	f, fingerprint, err := preparePayload(packagePath)
	if err != nil {
		return nil, err
	}
	request := []byte{}
	if details != nil {
		request, err = json.Marshal(details)
		if err != nil {
			return nil, fmt.Errorf("could not marshal request metadata: %v", err)
		}
	}

	header := http.Header{
		"Content-Type":      []string{"application/octet-stream"},
		"X-AMS-Fingerprint": []string{fingerprint},
		"X-AMS-Request":     []string{string(request)},
	}

	u := &shared.BufferedReader{Reader: f, Size: sentBytes}

	c.SetTransportTimeout(extendedTransportTimeout)
	op, _, err := c.QueryOperation(httpOp, apiPath, params, header, u, "")
	c.SetTransportTimeout(client.DefaultTransportTimeout)
	return op, err
}

func (c *clientImpl) download(path string, params client.QueryParams, header http.Header, downloader func(header *http.Header, body io.ReadCloser) error) error {
	c.SetTransportTimeout(extendedTransportTimeout)
	err := c.DownloadFile(path, params, header, downloader)
	c.SetTransportTimeout(client.DefaultTransportTimeout)
	return err
}

func convertFiltersToParams(filters []string) (client.QueryParams, error) {
	apiFilters := client.QueryParams{}
	for _, filter := range filters {
		parts := strings.SplitN(filter, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid filter '%s'", filter)
		}
		apiFilters[parts[0]] = parts[1]
	}
	return apiFilters, nil
}

func (c *clientImpl) rawWebsocket(url string) (*websocket.Conn, error) {
	httpTransport := c.HTTPTransport()

	dialer := websocket.Dialer{
		NetDial:         httpTransport.Dial,
		TLSClientConfig: httpTransport.TLSClientConfig,
		Proxy:           httpTransport.Proxy,
	}

	headers := http.Header{}

	conn, _, err := dialer.Dial(url, headers)
	if err != nil {
		return nil, err
	}

	return conn, err
}

func (c *clientImpl) websocket(path string) (*websocket.Conn, error) {
	// Generate the URL
	var url string
	serviceURL := c.ServiceURL()
	if strings.HasPrefix(serviceURL, "https://") {
		url = fmt.Sprintf("wss://%s/1.0%s", strings.TrimPrefix(serviceURL, "https://"), path)
	} else {
		url = fmt.Sprintf("ws://%s/1.0%s", strings.TrimPrefix(serviceURL, "http://"), path)
	}

	return c.rawWebsocket(url)
}

func (c *clientImpl) getOperationWebsocket(uuid string, secret string) (*websocket.Conn, error) {
	path := fmt.Sprintf("/operations/%s/websocket", url.QueryEscape(uuid))
	if secret != "" {
		path = fmt.Sprintf("%s?secret=%s", path, url.QueryEscape(secret))
	}

	return c.websocket(path)
}
