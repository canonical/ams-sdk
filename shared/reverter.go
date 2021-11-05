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

import (
	"context"
	"log"
	"time"
)

const (
	fullRevertTimeout = time.Second * 5
)

// RevertFunc describes a function used to revert an operation
type RevertFunc func(ctx context.Context) error

// Reverter provides functionality to automatically revert a set of executed operations. It
// can be used this way:
//
// r := shared.NewReverter()
// defer r.Finish()
//
// doOperation()
// r.Add(func(ctx context.Context) error {
//   revertOperation()
//   return nil
// })
//
// if err := doOtherOperation(); err != nil {
//   return err
// }
//
// r.Defuse()
type Reverter struct {
	needRevert bool
	reverters  []RevertFunc
}

// Add adds a new revert function to the reverter which will be called when Finish() is
// called unless the reverter gets defused.
func (r *Reverter) Add(f ...RevertFunc) {
	r.reverters = append(r.reverters, f...)
}

// Defuse defuses the reverter. If defused none of the added revert functions will be
// called when Finish() is invoked.
func (r *Reverter) Defuse() {
	r.needRevert = false
}

// Finish invokes all added revert functions if the reverter was not defused.
func (r *Reverter) Finish() {
	if !r.needRevert {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), fullRevertTimeout)
	defer cancel()
	// Walk in reverse order through our revertes and call them all
	for n := range r.reverters {
		revert := r.reverters[len(r.reverters)-n-1]
		if err := revert(ctx); err != nil {
			log.Printf("Failed to revert: %v", err)
		}
	}
}

// NewReverter constructs a new reverter.
func NewReverter() *Reverter {
	return &Reverter{needRevert: true}
}
