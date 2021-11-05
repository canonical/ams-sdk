
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

// ErrDontMatch error struct for a parameter whose value does not match the expected
type ErrDontMatch struct {
	content
	got      string
	expected string
}

// Error returns the error string
func (e ErrDontMatch) Error() string {
	return fmt.Sprintf("%v don't match, got %v but expected %v", e.What, e.got, e.expected)
}

// NewErrDontMatch returns a new ErrDontMatch struct
func NewErrDontMatch(what, got, expected string) ErrDontMatch {
	return ErrDontMatch{content{what}, got, expected}
}
