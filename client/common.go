// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/anbox-cloud/ams-sdk/shared"
	"github.com/anbox-cloud/ams-sdk/shared/rest/client"
)

func (c *clientImpl) upload(httpOp, apiPath, packagePath string, details interface{}, sentBytes chan float64) (client.Operation, error) {
	f, fingerprint, err := preparePayload(packagePath)
	if err != nil {
		return nil, err
	}

	request := []byte{}
	if details != nil {
		request, err = json.Marshal(details)
		if err != nil {
			return nil, fmt.Errorf("Could not marshal request metadata: %v", err)
		}
	}

	header := http.Header{
		"Content-Type":      []string{"application/octet-stream"},
		"X-AMS-Fingerprint": []string{fingerprint},
		"X-AMS-Request":     []string{string(request)},
	}

	u := &shared.BufferedReader{Reader: f, Size: sentBytes}

	c.SetTransportTimeout(extendedTransportTimeout)
	op, _, err := c.QueryOperation(httpOp, apiPath, nil, header, u, "")
	c.SetTransportTimeout(client.DefaultTransportTimeout)
	return op, err
}

func (c *clientImpl) download(path string, params client.QueryParams, header http.Header, downloader func(header *http.Header, body io.ReadCloser) error) error {
	c.SetTransportTimeout(extendedTransportTimeout)
	err := c.DownloadFile(path, params, header, downloader)
	c.SetTransportTimeout(client.DefaultTransportTimeout)
	return err
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
