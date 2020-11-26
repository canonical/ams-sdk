// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

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
