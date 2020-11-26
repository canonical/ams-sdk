// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package errors

import (
	"fmt"
)

// ErrInvalidArgument describes the error when an invalid argument was given
type ErrInvalidArgument struct {
	content
}

// Error returns the error string
func (e ErrInvalidArgument) Error() string {
	return fmt.Sprintf("argument %s is invalid", e.What)
}

// NewInvalidArgument returns a new ErrInvalidArgument struct
func NewInvalidArgument(what string) ErrInvalidArgument {
	return ErrInvalidArgument{content{what}}
}
