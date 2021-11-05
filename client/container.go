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
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/anbox-cloud/ams-sdk/api"
	"github.com/anbox-cloud/ams-sdk/shared"
	errs "github.com/anbox-cloud/ams-sdk/shared/errors"
	"github.com/anbox-cloud/ams-sdk/shared/rest/client"
)

// The ContainerExecArgs struct is used to pass additional options during a
// container shell session
type ContainerExecArgs struct {
	Stdin    io.ReadCloser
	Stdout   io.WriteCloser
	Stderr   io.WriteCloser
	Control  func(conn *websocket.Conn)
	DataDone chan bool
}

// ListContainers lists all available containers the AMS service currently manages
func (c *clientImpl) ListContainers() ([]api.Container, error) {
	containers := []api.Container{}
	params := client.QueryParams{
		"recursion": "1",
	}
	_, err := c.QueryStruct("GET", client.APIPath("containers"), params, nil, nil, "", &containers)
	return containers, err
}

// LaunchContainer launches a single new container on the AMS endpoint
func (c *clientImpl) LaunchContainer(details *api.ContainersPost) (client.Operation, error) {
	b, err := json.Marshal(details)

	if err != nil {
		return nil, err
	}

	op, _, err := c.QueryOperation("POST", client.APIPath("containers"), nil, nil, bytes.NewReader(b), "")
	return op, err
}

// RetrieveContainerByID retrieves a single container by its ID
func (c *clientImpl) RetrieveContainerByID(id string) (*api.Container, string, error) {
	if len(id) == 0 {
		return nil, "", errs.NewInvalidArgument("id")
	}
	container := &api.Container{}
	etag, err := c.QueryStruct("GET", client.APIPath("containers", id), nil, nil, nil, "", container)
	return container, etag, err
}

// DeleteContainerByID deletes a single container specified by its id
func (c *clientImpl) DeleteContainerByID(id string, force bool) (client.Operation, error) {
	if len(id) == 0 {
		return nil, errs.NewInvalidArgument("id")
	}

	details := api.ContainerDelete{
		Force: force,
	}
	b, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}

	op, _, err := c.QueryOperation("DELETE", client.APIPath("containers", id), nil, nil, bytes.NewReader(b), "")
	return op, err
}

// RetrieveContainerLog retrieves a specific log file of a container
func (c *clientImpl) RetrieveContainerLog(id, name string, downloader func(header *http.Header, body io.ReadCloser) error) error {
	if len(id) == 0 {
		return errs.NewInvalidArgument("id")
	}
	if len(name) == 0 {
		return errs.NewInvalidArgument("name")
	}

	if !c.HasExtension("container_logs") {
		return errs.NewErrNotSupported("api extension \"container_logs\"")
	}

	return c.download(client.APIPath("containers", id, "logs", name), nil, nil, downloader)
}

// ExecuteContainer requests that AMS opens a shell inside a container
func (c *clientImpl) ExecuteContainer(id string, details *api.ContainerExecPost, args *ContainerExecArgs) (client.Operation, error) {
	if len(id) == 0 {
		return nil, errs.NewInvalidArgument("id")
	}

	if !c.HasExtension("container_exec") {
		return nil, errs.NewErrNotSupported("api extension \"container_exec\"")
	}

	b, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}

	op, _, err := c.QueryOperation("POST", client.APIPath("containers", id, "exec"), nil, nil, bytes.NewReader(b), "")
	if err != nil {
		return nil, err
	}

	if args != nil {
		opAPI := op.Get()

		fds := map[string]string{}
		value, ok := opAPI.Metadata["fds"]
		if ok {
			values := value.(map[string]interface{})
			for k, v := range values {
				fds[k] = v.(string)
			}
		}

		if args.Control != nil && fds["control"] != "" {
			conn, err := c.getOperationWebsocket(opAPI.ID, fds["control"])
			if err != nil {
				return nil, err
			}
			go args.Control(conn)
		}

		if details.Interactive {
			// Handle interactive sections
			if args.Stdin != nil && args.Stdout != nil {
				// Connect to the websocket
				conn, err := c.getOperationWebsocket(opAPI.ID, fds["0"])
				if err != nil {
					return nil, err
				}

				// And attach stdin and stdout to it
				go func() {
					shared.WebsocketSendStream(conn, args.Stdin, -1)
					<-shared.WebsocketRecvStream(args.Stdout, conn)
					conn.Close()

					if args.DataDone != nil {
						close(args.DataDone)
					}
				}()
			} else {
				if args.DataDone != nil {
					close(args.DataDone)
				}
			}
		} else {
			dones := map[int]chan bool{}
			conns := []*websocket.Conn{}

			// Handle stdin
			if fds["0"] != "" {
				conn, err := c.getOperationWebsocket(opAPI.ID, fds["0"])
				if err != nil {
					return nil, err
				}

				conns = append(conns, conn)
				dones[0] = shared.WebsocketSendStream(conn, args.Stdin, -1)
			}

			// Handle stdout
			if fds["1"] != "" {
				conn, err := c.getOperationWebsocket(opAPI.ID, fds["1"])
				if err != nil {
					return nil, err
				}

				conns = append(conns, conn)
				dones[1] = shared.WebsocketRecvStream(args.Stdout, conn)
			}

			// Handle stderr
			if fds["2"] != "" {
				conn, err := c.getOperationWebsocket(opAPI.ID, fds["2"])
				if err != nil {
					return nil, err
				}

				conns = append(conns, conn)
				dones[2] = shared.WebsocketRecvStream(args.Stderr, conn)
			}

			// Wait for everything to be done
			go func() {
				for i, chDone := range dones {
					// Skip stdin, dealing with it separately below
					if i == 0 {
						continue
					}
					<-chDone
				}

				if fds["0"] != "" {
					args.Stdin.Close()
				}

				for _, conn := range conns {
					conn.Close()
				}

				if args.DataDone != nil {
					close(args.DataDone)
				}
			}()
		}
	}

	return op, nil
}
