// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2020 Canonical Ltd.  All rights reserved.
// +build docs

package api

func SwaggerModels() []string {
	return []string{
		swaggerModelOperation,
		swaggerModelStatusCode,
	}
}
