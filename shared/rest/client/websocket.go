// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package client

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

func (c *client) Websocket(resource string) (*websocket.Conn, error) {
	return c.dialWebsocket(c.composeWebsocketPath(resource))
}

// composeWebsocketPath returns websocket url related with rest client one
func (c *client) composeWebsocketPath(path string) string {
	host := c.serviceURL.Host

	scheme := "ws"
	if c.serviceURL.Scheme == "https" {
		scheme = "wss"
	}

	return fmt.Sprintf("%s://%s%s", scheme, host, path)
}

func (c *client) dialWebsocket(url string) (*websocket.Conn, error) {

	switch c.Doer.(type) {
	case *http.Client:
	default:
		return nil, errors.New("Client is not a valid http one")
	}

	t := c.Doer.(*http.Client).Transport.(*http.Transport)

	// Setup a new websocket dialer based on it
	dialer := websocket.Dialer{
		NetDial:         t.Dial,
		TLSClientConfig: t.TLSClientConfig,
		Proxy:           t.Proxy,
	}

	// Set the user agent
	headers := http.Header{}
	if c.httpUserAgent != "" {
		headers.Set("User-Agent", c.httpUserAgent)
	}

	// Establish the connection
	conn, _, err := dialer.Dial(url, headers)
	if err != nil {
		return nil, err
	}

	return conn, err
}
