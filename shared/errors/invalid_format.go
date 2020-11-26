// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package errors

import (
	"fmt"
)

// ErrInvalidFormat describes the error when an invalid format for an argument was given
type ErrInvalidFormat struct {
	content
}

// Error returns the error string
func (e ErrInvalidFormat) Error() string {
	return fmt.Sprintf("%s invalid format", e.What)
}

// NewErrInvalidFormat returns a new ErrInvalidFormat struct
func NewErrInvalidFormat(what string) ErrInvalidFormat {
	return ErrInvalidFormat{content{what}}
}
