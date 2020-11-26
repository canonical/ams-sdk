// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package errors

import (
	"fmt"
)

// ErrAlreadyExists describes an error when a resource already exists
type ErrAlreadyExists struct {
	content
}

// Error returns the error string
func (e ErrAlreadyExists) Error() string {
	return fmt.Sprintf("%s already exists", e.What)
}

// NewErrAlreadyExists returns a new ErrAlreadyExists struct
func NewErrAlreadyExists(what string) ErrAlreadyExists {
	return ErrAlreadyExists{content{what}}
}
