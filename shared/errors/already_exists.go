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

package errors

import (
	"fmt"
)

// ErrAlreadyExists describes an error when a resource already exists
type ErrAlreadyExists struct {
	content
}

// Error returns the error string
func (e ErrAlreadyExists) Error() string {
	return fmt.Sprintf("%s already exists", e.What)
}

// NewErrAlreadyExists returns a new ErrAlreadyExists struct
func NewErrAlreadyExists(what string) ErrAlreadyExists {
	return ErrAlreadyExists{content{what}}
}
