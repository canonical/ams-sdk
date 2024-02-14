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

package errors

import (
	"fmt"
)

// ErrNotExecutable represents a non executable file permission error structure
type ErrNotExecutable struct {
	content
}

// Error returns the error string
func (e ErrNotExecutable) Error() string {
	return fmt.Sprintf("%v not executable", e.What)
}

// NewErrNotExecutable returns a new ErrNotExecutable struct
func NewErrNotExecutable(what string) ErrNotExecutable {
	return ErrNotExecutable{content{what}}
}
