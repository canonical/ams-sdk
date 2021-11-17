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

// ErrInvalidLength describes the error when an argument has an invalid length
type ErrInvalidLength struct {
	content
}

// Error returns the error string
func (e ErrInvalidLength) Error() string {
	return fmt.Sprintf("length of %s is invalid", e.What)
}

// NewErrInvalidLength returns a new ErrInvalidLength struct
func NewErrInvalidLength(what string) ErrInvalidLength {
	return ErrInvalidLength{content{what}}
}
