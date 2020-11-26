// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package errors

import (
	"fmt"
)

// ErrNotSupported error struct for a not supported functionality
type ErrNotSupported struct {
	content
}

// Error returns the error string
func (e ErrNotSupported) Error() string {
	return fmt.Sprintf("%v not supported", e.What)
}

// NewErrNotSupported returns a new ErrNotSupported struct
func NewErrNotSupported(what string) ErrNotSupported {
	return ErrNotSupported{content{what}}
}
