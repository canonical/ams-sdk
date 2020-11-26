// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package errors

import (
	"fmt"
)

// ErrRequired error struct for a required parameter
type ErrRequired struct {
	content
}

// Error returns the error string
func (e ErrRequired) Error() string {
	return fmt.Sprintf("%v is required", e.What)
}

// NewErrRequired returns a new ErrRequiredstruct
func NewErrRequired(what string) ErrRequired {
	return ErrRequired{content{what}}
}
