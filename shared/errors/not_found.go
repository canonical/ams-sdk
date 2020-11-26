// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package errors

import (
	"fmt"
)

// ErrNotFound error struct for a not existing resource
type ErrNotFound struct {
	content
}

// Error returns the error string
func (e ErrNotFound) Error() string {
	return fmt.Sprintf("%v not found", e.What)
}

// NewErrNotFound returns a new ErrNotFound struct
func NewErrNotFound(what string) ErrNotFound {
	return ErrNotFound{content{what}}
}
