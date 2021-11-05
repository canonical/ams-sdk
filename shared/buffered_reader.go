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

package shared

import "io"

// BufferedReader represents a reader with buf size
type BufferedReader struct {
	Reader io.Reader
	Size   chan float64
}

// BufferedReader implements io.Reader interface
func (r *BufferedReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	if r.Size != nil {
		r.Size <- float64(n)
	}
	return n, err
}
