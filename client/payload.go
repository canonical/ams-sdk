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
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/anbox-cloud/ams-sdk/shared"
	errs "github.com/anbox-cloud/ams-sdk/shared/errors"
)

// Uploader represents a uploader that tracks progress
type Uploader struct {
	r    io.Reader
	sent chan float64
}

// Read implements io.Reader interface
func (u *Uploader) Read(p []byte) (int, error) {
	n, err := u.r.Read(p)
	if u.sent != nil {
		u.sent <- float64(n)
	}
	return n, err
}

func preparePayload(filepath string) (io.Reader, string, error) {
	if !shared.PathExists(filepath) {
		return nil, "", errs.NewErrNotFound("payload")
	}

	f, err := os.Open(filepath)
	if err != nil {
		return nil, "", err
	}

	hasher := sha256.New()
	_, err = io.Copy(hasher, f)
	if err != nil {
		return nil, "", err
	}

	fingerprint := fmt.Sprintf("%x", hasher.Sum(nil))

	// move cursor to the beginning of the file after hashing
	_, err = f.Seek(0, 0)
	if err != nil {
		return nil, "", err
	}

	return f, fingerprint, nil
}
