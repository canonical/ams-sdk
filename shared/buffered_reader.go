// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

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
