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
	"fmt"
	"sync"
)

// The EventListener struct is used to interact with an event stream
type EventListener struct {
	c            *client
	chActive     chan bool
	disconnected bool
	err          error

	targets     []*EventTarget
	targetsLock sync.Mutex
}

// The EventTarget struct is returned to the caller of AddHandler and used in RemoveHandler
type EventTarget struct {
	function func(interface{})
	types    []string
}

// AddHandler adds a function to be called whenever an event is received
func (e *EventListener) AddHandler(types []string, function func(interface{})) (*EventTarget, error) {
	if function == nil {
		return nil, fmt.Errorf("A valid function must be provided")
	}

	// Handle locking
	e.targetsLock.Lock()
	defer e.targetsLock.Unlock()

	// Create a new target
	target := EventTarget{
		function: function,
		types:    types,
	}

	// And add it to the targets
	e.targets = append(e.targets, &target)

	return &target, nil
}

// RemoveHandler removes a function to be called whenever an event is received
func (e *EventListener) RemoveHandler(target *EventTarget) error {
	if target == nil {
		return fmt.Errorf("A valid event target must be provided")
	}

	// Handle locking
	e.targetsLock.Lock()
	defer e.targetsLock.Unlock()

	// Locate and remove the function from the list
	for i, entry := range e.targets {
		if entry == target {
			copy(e.targets[i:], e.targets[i+1:])
			e.targets[len(e.targets)-1] = nil
			e.targets = e.targets[:len(e.targets)-1]
			return nil
		}
	}

	return fmt.Errorf("Couldn't find this function and event types combination")
}

// Disconnect must be used once done listening for events
func (e *EventListener) Disconnect() {
	if e.disconnected {
		return
	}

	// Handle locking
	e.c.eventListenersLock.Lock()
	defer e.c.eventListenersLock.Unlock()

	// Locate and remove it from the global list
	for i, listener := range e.c.eventListeners {
		if listener == e {
			copy(e.c.eventListeners[i:], e.c.eventListeners[i+1:])
			e.c.eventListeners[len(e.c.eventListeners)-1] = nil
			e.c.eventListeners = e.c.eventListeners[:len(e.c.eventListeners)-1]
			break
		}
	}

	// Turn off the handler
	e.err = nil
	e.disconnected = true
	close(e.chActive)
}

// Wait hangs until the server disconnects the connection or Disconnect() is called
func (e *EventListener) Wait() error {
	<-e.chActive
	return e.err
}

// IsActive returns true if this listener is still connected, false otherwise.
func (e *EventListener) IsActive() bool {
	select {
	case <-e.chActive:
		return false // If the chActive channel is closed we got disconnected
	default:
		return true
	}
}
