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

// ErrUnknown describes the error which a parameter was unknown
type ErrUnknown struct {
	content
}

// ErrUnknown returns the error string
func (e ErrUnknown) Error() string {
	return fmt.Sprintf("%s is unknown", e.What)
}

// NewErrUnknown returns a new ErrUnknown struct
func NewErrUnknown(what string) ErrUnknown {
	return ErrUnknown{content{what}}
}
