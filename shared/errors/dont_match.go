// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package errors

import (
	"fmt"
)

// ErrDontMatch error struct for a parameter whose value does not match the expected
type ErrDontMatch struct {
	content
	got      string
	expected string
}

// Error returns the error string
func (e ErrDontMatch) Error() string {
	return fmt.Sprintf("%v don't match, got %v but expected %v", e.What, e.got, e.expected)
}

// NewErrDontMatch returns a new ErrDontMatch struct
func NewErrDontMatch(what, got, expected string) ErrDontMatch {
	return ErrDontMatch{content{what}, got, expected}
}
