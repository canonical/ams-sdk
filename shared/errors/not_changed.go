// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2019 Canonical Ltd.  All rights reserved.

package errors

import "fmt"

// ErrNotChanged describes an error when a value has changed
type ErrNotChanged struct {
	content
}

// Error returns the error string
func (e ErrNotChanged) Error() string {
	return fmt.Sprintf("%s not changed", e.What)
}

// NewErrNotChanged returns a new ErrNotChanged struct
func NewErrNotChanged(what string) ErrNotChanged {
	return ErrNotChanged{content{what}}
}
