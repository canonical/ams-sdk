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
	"os"
	"strings"
)

// ContentList represents a list of content
type ContentList []string

// Has checks if the given content is included in the content list
func (c *ContentList) Has(content string) bool {
	for _, contentPath := range *c {
		// Dir path contains a trailing separator, which needs to be removed before comparing
		if strings.TrimRight(contentPath, string(os.PathSeparator)) == content {
			return true
		}
	}
	return false
}

// Add adds new contents into the list
func (c *ContentList) Add(content ...string) {
	*c = append(*c, content...)
}
