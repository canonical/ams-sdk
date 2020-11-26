// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package errors

import (
	"fmt"
)

// ErrNotExecutable represents a non executable file permission error structure
type ErrNotExecutable struct {
	content
}

// Error returns the error string
func (e ErrNotExecutable) Error() string {
	return fmt.Sprintf("%v not executable", e.What)
}

// NewErrNotExecutable returns a new ErrNotExecutable struct
func NewErrNotExecutable(what string) ErrNotExecutable {
	return ErrNotExecutable{content{what}}
}
