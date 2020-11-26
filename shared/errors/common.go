// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package errors

import (
	"fmt"
)

var (
	// ErrAlreadyRunning is returned when an object is already running when it was started again
	ErrAlreadyRunning = fmt.Errorf("Already running")
)

type content struct {
	What string
}

// IgnoreErrNotFound returns nil when the provided error is of type ErrNotFound
// and otherwise the given error
func IgnoreErrNotFound(err error) error {
	switch err.(type) {
	case ErrNotFound:
		return nil
	default:
		return err
	}
}
