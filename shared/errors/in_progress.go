// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package errors

import (
	"fmt"
)

// ErrInProgress describes an error when an operation is already in progress
type ErrInProgress struct {
	content
}

// Error returns the error string
func (e ErrInProgress) Error() string {
	return fmt.Sprintf("%s already in progress", e.What)
}

// NewErrInProgress returns a new ErrInProgress struct
func NewErrInProgress(what string) ErrInProgress {
	return ErrInProgress{content{what}}
}
