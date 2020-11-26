// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package errors

import (
	"fmt"
)

// ErrMalformed describes the error when a malformed content was given
type ErrMalformed struct {
	content
}

// Error returns the error string
func (e ErrMalformed) Error() string {
	return fmt.Sprintf("%s is malformed", e.What)
}

// NewErrMalformed returns a new ErrMalformed struct
func NewErrMalformed(what string) ErrMalformed {
	return ErrMalformed{content{what}}
}
