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

package network

import "net"

// AllocatePorts asks the kernel for a set of free open ports that are
// not in use. The implementation gurantees that all ports are unique
// but does not give a gurantee that they haven't been taken by the time
// they are actually used.
func AllocatePorts(num int) ([]int, error) {
	var listeners []*net.TCPListener

	defer func() {
		for _, l := range listeners {
			l.Close()
		}
	}()

	var ports []int
	for n := 0; n < num; n++ {
		addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
		if err != nil {
			return nil, err
		}

		l, err := net.ListenTCP("tcp", addr)
		if err != nil {
			return nil, err
		}

		ports = append(ports, l.Addr().(*net.TCPAddr).Port)
		listeners = append(listeners, l)
	}

	return ports, nil
}

// AllocatePort asks the kernel for a free open port that is ready to use
func AllocatePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return -1, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return -1, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
