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

package packages

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
	"github.com/anbox-cloud/ams-sdk/shared"
	errs "github.com/anbox-cloud/ams-sdk/shared/errors"
)

const (
	// MaxHookExecutionTimeout represents the maximum hook execution timeout
	MaxHookExecutionTimeout = time.Minute * 15
)

// PackageType represents the file format for a pacakge
type PackageType int

const (
	// PackageTypeUnknown represents an unknown addon or application package format
	PackageTypeUnknown PackageType = iota - 1
	// PackageTypeTarBZ2 represents an addon or application package as a tarball
	PackageTypeTarBZ2
	// PackageTypeZip represents an addon or application package as a zip archive
	PackageTypeZip
)

// DetectPackageType is used to auto-determine the type of a package by looking
// at the magic bytes of a file
func DetectPackageType(filePath string) (PackageType, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return PackageTypeUnknown, err
	}
	defer file.Close()

	buff := make([]byte, 512)
	if _, err = file.Read(buff); err != nil {
		return PackageTypeUnknown, err
	}
	mimeType := http.DetectContentType(buff)
	switch mimeType {
	case "application/zip":
		return PackageTypeZip, nil
		// TODO: Add the mimetype required to detect tarballs. This is currently not
		// available in the standard functions in Go lib
	default:
		return PackageTypeUnknown, nil
	}

}

// IsTarball detects if the given package is a valid tarball
func IsTarball(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".tar.bz2") ||
		strings.HasSuffix(strings.ToLower(filePath), ".tbz2")
}

// IsZip detects if the given package is a valid zip archive
func IsZip(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".zip")
}

// ParseManifest parses the manifest structure from io reader
func ParseManifest(reader io.Reader, manifest interface{}) error {
	byteBuf := new(bytes.Buffer)
	if _, err := byteBuf.ReadFrom(reader); err != nil {
		return err
	}
	return yaml.Unmarshal(byteBuf.Bytes(), manifest)
}

// CreateTempPackage creates a temporary package file with the given contents
func CreateTempPackage(sources []string, packageType PackageType) (string, error) {
	srcDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	outputDir, err := ioutil.TempDir("", "app")
	if err != nil {
		return "", err
	}

	// The shared.DirCopy function that is used by this function to create a
	// package tarball does a file copy recursively. Hence we just need to keep
	// the top-level files and dirs for a tarball creation. Keeping the children
	// of one directory in the sources fails the function due to file existence.
	// it also results in performance degradation.
	topLevelFiles := map[string]bool{}
	for _, file := range sources {
		parts := strings.Split(file, "/")
		topLevelFiles[parts[0]] = true
	}

	// Copy all files to the output dir
	for file := range topLevelFiles {
		fi, err := os.Stat(file)
		if err != nil {
			return "", err
		}

		fileName := fi.Name()

		srcFile := file
		if !shared.PathExists(srcFile) {
			srcFile = filepath.Join(srcDir, fileName)
		}
		destFile := filepath.Join(outputDir, fileName)
		if fi.IsDir() {
			if err := shared.DirCopy(srcFile, destFile); err != nil {
				return "", fmt.Errorf("Failed to copy data %s to %s: %v",
					srcFile, destFile, err)
			}
		} else {
			if err := shared.FileCopy(srcFile, destFile); err != nil {
				return "", fmt.Errorf("Failed to copy data %s to %s: %v",
					srcFile, destFile, err)
			}
		}
	}

	var packagePath string
	switch packageType {
	case PackageTypeTarBZ2:
		packagePath = filepath.Join(outputDir, "application.tar.bz2")
		if err := shared.CreateBzip2Tarball(srcDir, packagePath, sources); err != nil {
			return "", fmt.Errorf("Failed to create tarball file")
		}
	case PackageTypeZip:
		fallthrough
	default:
		packagePath = filepath.Join(outputDir, "application.zip")
		if err := shared.CreateZip(srcDir, packagePath, sources); err != nil {
			return "", fmt.Errorf("Failed to create zip file")
		}
	}

	return packagePath, nil
}

// ValidateHookTimeout validates the given timeout for hooks
func ValidateHookTimeout(timeout string) error {
	d, err := time.ParseDuration(timeout)
	if err != nil || d < 0 {
		return errs.NewInvalidArgument("timeout")
	}

	if d > MaxHookExecutionTimeout {
		return fmt.Errorf("timeout exceeds the allowed maximum(%d mins)", MaxHookExecutionTimeout/time.Minute)
	}

	return nil
}
