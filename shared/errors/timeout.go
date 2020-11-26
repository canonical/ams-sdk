// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package errors

import (
	"fmt"
)

// ErrTimeout describes an error when an operation timed out
type ErrTimeout struct {
	content
}

// Error returns the error string
func (e ErrTimeout) Error() string {
	return fmt.Sprintf("%s timed out", e.What)
}

// NewErrTimeout returns a new ErrAborted struct
func NewErrTimeout(what string) ErrTimeout {
	return ErrTimeout{content{what}}
}
