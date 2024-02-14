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

package constants

import (
	"github.com/anbox-cloud/ams-sdk/pkg/units"
)

const (
	// MaxUserdataSize is the maximum size the userdata of a container can take. The size
	// of the userdata attached to a container has to be limited as a single element in our
	// data store can only have a maximum size of 1.5 MB and we store the userdata as part
	// of it.
	MaxUserdataSize = 10 * units.KB
)
