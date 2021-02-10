// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2021 Canonical Ltd.  All rights reserved.

package shared

import "sync/atomic"

const (
	stateUnlocked uint32 = iota
	stateLocked
)

// Locker represents a locker which acquires the lock only if it is
// not held by another goroutine at the time of invocation.
type Locker struct {
	val uint32
}

// NewLocker returns a new Locker which can be used by a goroutine to
// hold the internal lock and perform task which can not be invoked
// by other goroutine
func NewLocker() *Locker {
	return &Locker{val: stateUnlocked}
}

// TryLock attempts to acquire the lock immediately, return true if locking succeeds
// otherwise return false.
func (l *Locker) TryLock() bool {
	return atomic.CompareAndSwapUint32(&l.val, stateUnlocked, stateLocked)
}

// UnLock attempts to release the lock immediately.
func (l *Locker) UnLock() {
	atomic.StoreUint32(&l.val, stateUnlocked)
}
