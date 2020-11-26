// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package errors

// ErrAborted describes an error when an operation was aborted
type ErrAborted struct {
	content
}

// Error returns the error string
func (e ErrAborted) Error() string {
	return e.What
}

// NewErrAborted returns a new ErrAborted struct
func NewErrAborted(what string) ErrAborted {
	return ErrAborted{content{what}}
}
