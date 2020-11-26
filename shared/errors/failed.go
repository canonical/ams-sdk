// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package errors

import (
	"fmt"
)

// ErrFailed describes an error when an operation has failed
type ErrFailed struct {
	content
}

// Error returns the error string
func (e ErrFailed) Error() string {
	return fmt.Sprintf("%s failed", e.What)
}

// NewErrFailed returns a new ErrFailed struct
func NewErrFailed(what string) ErrFailed {
	return ErrFailed{content{what}}
}
