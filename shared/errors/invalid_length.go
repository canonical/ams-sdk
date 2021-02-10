// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package errors

import (
	"fmt"
)

// ErrInvalidLength describes the error when an argument has an invalid length
type ErrInvalidLength struct {
	content
}

// Error returns the error string
func (e ErrInvalidLength) Error() string {
	return fmt.Sprintf("length of %s is invalid", e.What)
}

// NewErrInvalidLength returns a new ErrInvalidLength struct
func NewErrInvalidLength(what string) ErrInvalidLength {
	return ErrInvalidLength{content{what}}
}
