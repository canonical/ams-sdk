// -*- Mode: Go; indent-tabs-mode: t -*-
/*
 * This file is part of AMS SDK
 * Copyright 2021 Canonical Ltd.
 *
 * This program is free software: you can redistribute it and/or modify it under
 * the terms of the GNU Lesser General Public License version 3, as published
 * by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranties of MERCHANTABILITY, SATISFACTORY
 * QUALITY, or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public
 * License for more details.
 *
 * You should have received a copy of the Lesser GNU General Public License along
 * with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
