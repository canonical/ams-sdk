// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package errors

import (
	"fmt"
)

// ErrNotAllowed error struct when an operation is not allowed
type ErrNotAllowed struct {
	content
}

// Error returns the error string
func (e ErrNotAllowed) Error() string {
	return fmt.Sprintf("%v not allowed", e.What)
}

// NewErrNotAllowed returns a new ErrNotAllowed struct
func NewErrNotAllowed(what string) ErrNotAllowed {
	return ErrNotAllowed{content{what}}
}
