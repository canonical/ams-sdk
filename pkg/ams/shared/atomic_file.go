// -*- Mode: Go; indent-tabs-mode: t -*-
/*
 * This file is part of AMS SDK
 * Copyright 2021 Canonical Ltd.
 *
 * This program is free software: you can redistribute it and/or modify it under
 * the terms of the Lesser GNU General Public License version 3, as published
 * by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranties of MERCHANTABILITY, SATISFACTORY
 * QUALITY, or FITNESS FOR A PARTICULAR PURPOSE.  See the Lesser GNU General Public
 * License for more details.
 *
 * You should have received a copy of the Lesser GNU General Public License along
 * with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package shared

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// AtomicFile allows operation on a file to happen atomicly. Changes must be committed in
// the end before the final file will appear at its target location
type AtomicFile struct {
	*os.File

	targetPath string
	committed  bool
}

// WriteFileAtomic writes the given content to the target path in an atomic way
func WriteFileAtomic(targetPath string, content []byte, mode os.FileMode) error {
	f, err := NewAtomicFile(targetPath, mode)
	if err != nil {
		return err
	}
	defer f.Cancel()

	n, err := f.Write(content)
	if err != nil {
		return err
	}
	if n != len(content) {
		return fmt.Errorf("failed to write full content to file at path %s", targetPath)
	}

	return f.Commit()
}

// NewAtomicFile creates a new AtomicFile with the given target path and mode
func NewAtomicFile(targetPath string, mode os.FileMode) (*AtomicFile, error) {
	f, err := ioutil.TempFile(path.Dir(targetPath), path.Base(targetPath))
	if err != nil {
		return nil, err
	}

	if err := f.Chmod(mode); err != nil {
		f.Close()
		os.Remove(f.Name())
		return nil, err
	}

	af := AtomicFile{
		File:       f,
		targetPath: targetPath,
	}

	return &af, nil
}

// Commit commits changes made to the file and moves the temporary file to its
// final location
func (af *AtomicFile) Commit() error {
	if err := af.Sync(); err != nil {
		return err
	}

	if err := af.Close(); err != nil {
		return err
	}

	if err := os.Rename(af.Name(), af.targetPath); err != nil {
		return err
	}

	af.committed = true

	return nil
}

// Cancel all changes made. Is a no-op when Commit was already called
func (af *AtomicFile) Cancel() error {
	if af.committed {
		return fmt.Errorf("cannot canceled as already committed")
	}

	if err := af.Close(); err != nil {
		return err
	}

	if err := os.Remove(af.Name()); err != nil {
		return err
	}

	return nil
}
