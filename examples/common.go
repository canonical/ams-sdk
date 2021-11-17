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

package examples

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

// PrintCreated prints out created resources
func PrintCreated(resources map[string][]string) {
	for k, r := range resources {
		fmt.Printf("%s:\n", k)
		for _, j := range r {
			fmt.Printf("  - %s\n", path.Base(j))
		}
	}
}

// DumpData prints out object in a human readable format
func DumpData(data interface{}) error {
	b, err := json.MarshalIndent(&data, "", "\t")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, string(b))
	return nil
}

// GetByteSizeString returns the size in a human readable format
func GetByteSizeString(input int64, precision uint) string {
	if input < 1024 {
		return fmt.Sprintf("%dB", input)
	}

	value := float64(input)

	for _, unit := range []string{"kB", "MB", "GB", "TB", "PB", "EB"} {
		value = value / 1024
		if value < 1024 {
			return fmt.Sprintf("%.*f%s", precision, value, unit)
		}
	}

	return fmt.Sprintf("%.*fEB", precision, value)
}
