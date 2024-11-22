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
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	api "github.com/anbox-cloud/ams-sdk/api/ams"
	errs "github.com/anbox-cloud/ams-sdk/pkg/ams/shared/errors"
	"github.com/anbox-cloud/ams-sdk/pkg/ams/shared/rest/client"
	"github.com/anbox-cloud/ams-sdk/pkg/network"
	"github.com/gorilla/websocket"
)

// The InstanceExecArgs struct is used to pass additional options during a
// instance shell session
type InstanceExecArgs struct {
	Stdin    io.ReadCloser
	Stdout   io.WriteCloser
	Stderr   io.WriteCloser
	Control  func(conn *websocket.Conn)
	DataDone chan bool
}

// ListInstancesWithFilters lists all available instances the AMS service currently manages
func (c *clientImpl) ListInstancesWithFilters(filters []string) ([]api.Instance, error) {
	instances := []api.Instance{}
	if !c.hasInstanceSupport {
		containers, err := c.ListContainersWithFilters(filters)
		if err != nil {
			return nil, err
		}
		for _, c := range containers {
			inst, err := api.MapContainerToInstance(&c)
			if err != nil {
				return nil, err
			}
			instances = append(instances, inst)
		}
	} else {
		params, err := convertFiltersToParams(filters)
		if err != nil {
			return nil, err
		}
		params["recursion"] = "1"
		_, err = c.QueryStruct("GET", client.APIPath("instances"), params, nil, nil, "", &instances)
		if err != nil {
			return nil, err
		}
	}
	return instances, nil
}

// ListInstances lists all available instances the AMS service currently manages
func (c *clientImpl) ListInstances() ([]api.Instance, error) {
	instances := []api.Instance{}
	if !c.hasInstanceSupport {
		containers, err := c.ListContainers()
		if err != nil {
			return nil, err
		}
		for _, c := range containers {
			inst, err := api.MapContainerToInstance(&c)
			if err != nil {
				return nil, err
			}
			instances = append(instances, inst)
		}
	} else {
		params := client.QueryParams{
			"recursion": "1",
		}
		_, err := c.QueryStruct("GET", client.APIPath("instances"), params, nil, nil, "", &instances)
		if err != nil {
			return nil, err
		}
	}
	return instances, nil
}

// LaunchInstance launches a single new instance on the AMS endpoint
func (c *clientImpl) LaunchInstance(details *api.InstancesPost, noWait bool) (client.Operation, error) {
	if !c.hasInstanceSupport {
		d := api.ContainersPost{
			ApplicationID:      details.ApplicationID,
			ApplicationVersion: details.ApplicationVersion,
			ImageID:            details.ImageID,
			ImageVersion:       details.ImageVersion,
			Node:               details.Node,
			Userdata:           details.Userdata,
			Addons:             details.Addons,
			Services:           details.Services,
			CPUs:               details.Resources.CPUs,
			DiskSize:           details.Resources.DiskSize,
			Memory:             details.Resources.Memory,
			GPUSlots:           details.Resources.GPUSlots,
			VPUSlots:           details.Resources.VPUSlots,
			Tags:               details.Tags,
		}
		d.Config.Platform = details.Config.Platform
		d.Config.BootPackage = details.Config.BootPackage
		d.Config.BootActivity = details.Config.BootActivity
		d.Config.MetricsServer = details.Config.MetricsServer
		d.Config.DisableWatchdog = details.Config.DisableWatchdog
		d.Config.Features = details.Config.Features
		d.Config.DevMode = details.Config.DevMode
		d.NoStart = details.NoStart

		if len(d.ImageID) > 0 && (d.CPUs == nil || d.Memory == nil || d.DiskSize == nil) {
			return nil, errs.NewInvalidArgument("resources")
		}

		return c.LaunchContainer(&d, noWait)
	}

	hasVMSupport, err := c.HasExtension("vm_support")
	if err != nil {
		return nil, err
	}

	if details.Type == api.InstanceTypeVM && !hasVMSupport {
		return nil, errs.NewErrNotSupported("VM")
	}

	b, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}
	params := client.QueryParams{"no_wait": strconv.FormatBool(noWait)}
	op, _, err := c.QueryOperation("POST", client.APIPath("instances"), params, nil, bytes.NewReader(b), "")
	return op, err
}

// RetrieveInstanceByID retrieves a single instance by its ID
func (c *clientImpl) RetrieveInstanceByID(id string) (*api.Instance, string, error) {
	if len(id) == 0 {
		return nil, "", errs.NewInvalidArgument("id")
	}
	if !c.hasInstanceSupport {
		container, etag, err := c.RetrieveContainerByID(id)
		if err != nil {
			return nil, "", err
		}
		instance, err := api.MapContainerToInstance(container)
		if err != nil {
			return nil, "", err
		}
		return &instance, etag, err
	}

	instance := &api.Instance{}
	etag, err := c.QueryStruct("GET", client.APIPath("instances", id), nil, nil, nil, "", instance)
	return instance, etag, err
}

// UpdateInstanceByID updates an existing instance specified by its id
func (c *clientImpl) UpdateInstanceByID(id string, details *api.InstancePatch, noWait bool) (client.Operation, error) {
	if len(id) == 0 {
		return nil, errs.NewInvalidArgument("id")
	}

	if !c.hasInstanceSupport {
		return c.UpdateContainerByID(id, &api.ContainerPatch{
			DesiredStatus: details.DesiredStatus,
		}, noWait)
	}

	b, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}

	params := client.QueryParams{"no_wait": strconv.FormatBool(noWait)}
	op, _, err := c.QueryOperation("PATCH", client.APIPath("instances", id), params, nil, bytes.NewReader(b), "")
	return op, err
}

// DeleteInstanceByID deletes a single instance specified by its id
func (c *clientImpl) DeleteInstanceByID(id string, force bool) (client.Operation, error) {
	if len(id) == 0 {
		return nil, errs.NewInvalidArgument("id")
	}

	if !c.hasInstanceSupport {
		return c.DeleteContainerByID(id, force)
	}

	details := api.InstanceDelete{
		Force: force,
	}
	b, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}

	op, _, err := c.QueryOperation("DELETE", client.APIPath("instances", id), nil, nil, bytes.NewReader(b), "")
	return op, err
}

// DeleteInstances deletes multiple instances in a bulk operation
func (c *clientImpl) DeleteInstances(ids []string, force bool) (client.Operation, error) {
	if len(ids) == 0 {
		return nil, errs.NewInvalidArgument("ids")
	}

	if !c.hasInstanceSupport {
		return c.DeleteContainers(ids, force)
	}

	details := api.InstancesDelete{
		IDs:   ids,
		Force: force,
	}
	b, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}

	op, _, err := c.QueryOperation("DELETE", client.APIPath("instances"), nil, nil, bytes.NewReader(b), "")
	return op, err
}

// RetrieveInstanceLog retrieves a specific log file of an instance
func (c *clientImpl) RetrieveInstanceLog(id, name string, downloader func(header *http.Header, body io.ReadCloser) error) error {
	if len(id) == 0 {
		return errs.NewInvalidArgument("id")
	}
	if len(name) == 0 {
		return errs.NewInvalidArgument("name")
	}

	if !c.hasInstanceSupport {
		return c.RetrieveContainerLog(id, name, downloader)
	}

	return c.download(client.APIPath("instances", id, "logs", name), nil, nil, downloader)
}

// ExecuteInstance requests that AMS opens a shell inside an instance
func (c *clientImpl) ExecuteInstance(id string, details *api.InstanceExecPost, args *InstanceExecArgs) (client.Operation, error) {
	if len(id) == 0 {
		return nil, errs.NewInvalidArgument("id")
	}

	if !c.hasInstanceSupport {
		return c.ExecuteContainer(id, &api.ContainerExecPost{
			Command:     details.Command,
			Environment: details.Environment,
			Interactive: details.Interactive,
			Width:       details.Width,
			Height:      details.Height,
		}, &ContainerExecArgs{
			Stdin:    args.Stdin,
			Stdout:   args.Stdout,
			Stderr:   args.Stderr,
			Control:  args.Control,
			DataDone: args.DataDone,
		})
	}

	b, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}

	op, _, err := c.QueryOperation("POST", client.APIPath("instances", id, "exec"), nil, nil, bytes.NewReader(b), "")
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
					network.WebsocketSendStream(conn, args.Stdin, -1)
					<-network.WebsocketRecvStream(args.Stdout, conn)
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
			waitConns := 0

			// Handle stdin
			if fds["0"] != "" {
				conn, err := c.getOperationWebsocket(opAPI.ID, fds["0"])
				if err != nil {
					return nil, err
				}

				conns = append(conns, conn)
				dones[0] = network.WebsocketSendStream(conn, args.Stdin, -1)
			}

			// Handle stdout
			if fds["1"] != "" {
				conn, err := c.getOperationWebsocket(opAPI.ID, fds["1"])
				if err != nil {
					return nil, err
				}

				conns = append(conns, conn)
				dones[1] = network.WebsocketRecvStream(args.Stdout, conn)
				waitConns++
			}

			// Handle stderr
			if fds["2"] != "" {
				conn, err := c.getOperationWebsocket(opAPI.ID, fds["2"])
				if err != nil {
					return nil, err
				}

				conns = append(conns, conn)
				dones[2] = network.WebsocketRecvStream(args.Stderr, conn)
				waitConns++
			}

			// Wait for everything to be done
			go func() {
				for {
					select {
					case <-dones[0]:
						// Handle stdin finish, but don't wait for it
						dones[0] = nil
						_ = conns[0].Close()
					case <-dones[1]:
						dones[1] = nil
						_ = conns[1].Close()
						waitConns--
					case <-dones[2]:
						dones[2] = nil
						_ = conns[2].Close()
						waitConns--
					}

					if waitConns <= 0 {
						// Close stdin websocket if defined and not already closed.
						if dones[0] != nil {
							conns[0].Close()
						}

						break
					}
				}

				if args.DataDone != nil {
					close(args.DataDone)
				}
			}()
		}
	}

	return op, nil
}
