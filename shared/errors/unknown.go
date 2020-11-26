// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2019 Canonical Ltd.  All rights reserved.

package errors

import (
	"fmt"
)

// ErrUnknown describes the error which a parameter was unknown
type ErrUnknown struct {
	content
}

// ErrUnknown returns the error string
func (e ErrUnknown) Error() string {
	return fmt.Sprintf("%s is unknown", e.What)
}

// NewErrUnknown returns a new ErrUnknown struct
func NewErrUnknown(what string) ErrUnknown {
	return ErrUnknown{content{what}}
}
