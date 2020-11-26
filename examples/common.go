// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

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
