// -*- Mode: Go; indent-tabs-mode: t -*-
/*
 * This file is part of AMS SDK
 * Copyright 2021 Canonical Ltd.
 *
 * This program is free software: you can redistribute it and/or modify it under
 * the terms of the GNU Lesser General Public License version 3, as published
 * by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranties of MERCHANTABILITY, SATISFACTORY
 * QUALITY, or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public
 * License for more details.
 *
 * You should have received a copy of the Lesser GNU General Public License along
 * with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
