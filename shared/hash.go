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
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"

	"golang.org/x/crypto/scrypt"
)

// Hash creates a hash of 192 chars from a value.
// A random salt of 32 bytes is used. The result is that salt plus the
// hashed result encoded to hexadecimal
func Hash(value string) (string, error) {
	// Nothing to do on unset
	if len(value) == 0 {
		return value, nil
	}

	buf := make([]byte, 32)
	// Generate a random salt of 32 bytes
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", err
	}

	hash, err := scrypt.Key([]byte(value), buf, 1<<14, 8, 1, 64)
	if err != nil {
		return "", err
	}

	// buf = salt + hash
	buf = append(buf, hash...)
	value = hex.EncodeToString(buf)

	return value, nil
}

// ValidateHash decodes hexadecimal hashed, extracts the salt from the first 32 bytes and
// hashes the value param using that salt. The result of that hash function should
// match the rest of the decoded secret buffer. Otherwise the validation fails
func ValidateHash(hashed, value string) error {
	if len(hashed) == 0 {
		return errors.New("No password is set")
	}

	buff, err := hex.DecodeString(hashed)
	if err != nil {
		return err
	}

	salt := buff[0:32]
	hash, err := scrypt.Key([]byte(value), salt, 1<<14, 8, 1, 64)
	if err != nil {
		return err
	}

	if !bytes.Equal(hash, buff[32:]) {
		return errors.New("Bad password provided")
	}

	return nil
}
