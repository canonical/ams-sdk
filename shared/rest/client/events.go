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

	"github.com/anbox-cloud/ams-sdk/shared"
)

// Event handling functions

// GetEvents connects to the monitoring interface
func (c *client) GetEvents() (*EventListener, error) {
	// Prevent anything else from interacting with the listeners
	c.eventListenersLock.Lock()
	defer c.eventListenersLock.Unlock()

	// Setup a new listener
	listener := EventListener{
		c:        c,
		chActive: make(chan bool),
	}

	if c.eventListeners != nil {
		// There is an existing Go routine setup, so just add another target
		c.eventListeners = append(c.eventListeners, &listener)
		return &listener, nil
	}

	// Initialize the list if needed
	c.eventListeners = []*EventListener{}

	// Setup a new connection with the server
	resource := APIPath("events")
	conn, err := c.dialWebsocket(c.composeWebsocketPath(resource))
	if err != nil {
		return nil, err
	}

	// Add the listener
	c.eventListeners = append(c.eventListeners, &listener)

	// And spawn the listener
	go func() {
		for {
			c.eventListenersLock.Lock()
			if len(c.eventListeners) == 0 {
				// We don't need the connection anymore, disconnect
				conn.Close()

				c.eventListeners = nil
				c.eventListenersLock.Unlock()
				break
			}
			c.eventListenersLock.Unlock()

			_, data, err := conn.ReadMessage()
			if err != nil {
				// Prevent anything else from interacting with the listeners
				c.eventListenersLock.Lock()
				defer c.eventListenersLock.Unlock()

				// Tell all the current listeners about the failure
				for _, listener := range c.eventListeners {
					listener.err = err
					listener.disconnected = true
					close(listener.chActive)
				}

				// And remove them all from the list
				c.eventListeners = []*EventListener{}
				return
			}

			// Attempt to unpack the message
			message := make(map[string]interface{})
			err = json.Unmarshal(data, &message)
			if err != nil {
				continue
			}

			// Extract the message type
			_, ok := message["type"]
			if !ok {
				continue
			}
			messageType := message["type"].(string)

			// Send the message to all handlers
			c.eventListenersLock.Lock()
			for _, listener := range c.eventListeners {
				listener.targetsLock.Lock()
				for _, target := range listener.targets {
					if target.types != nil &&
						!shared.StringInSlice(messageType, target.types) &&
						!shared.StringInSlice("all", target.types) {
						continue
					}

					go target.function(message)
				}
				listener.targetsLock.Unlock()
			}
			c.eventListenersLock.Unlock()
		}
	}()

	return &listener, nil
}
