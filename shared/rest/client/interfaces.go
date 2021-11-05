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

package client

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/anbox-cloud/ams-sdk/shared/rest/api"
)

// The Operation interface represents a currently running operation.
type Operation interface {
	AddHandler(function func(api.Operation)) (target *EventTarget, err error)
	Cancel() (err error)
	Get() (op api.Operation)
	GetWebsocket() (conn *websocket.Conn, err error)
	RemoveHandler(target *EventTarget) (err error)
	Refresh() (err error)
	Wait(ctx context.Context) (err error)
}

// The Operations interface represents operations exposed API methods
type Operations interface {
	ListOperationUUIDs() (uuids []string, err error)
	ListOperations() (operations []api.Operation, err error)
	RetrieveOperationByID(uuid string) (op *api.Operation, etag string, err error)
	WaitForOperationToFinish(uuid string, timeout time.Duration) (op *api.Operation, err error)
	GetOperationWebsocket(uuid string) (conn *websocket.Conn, err error)
	DeleteOperation(uuid string) (err error)
}

// The Certificates interface represents client certificates related API methods
type Certificates interface {
	ListCertificates() (certificates []api.Certificate, err error)
	AddCertificate(base64PublicKey, trustPassword string) (err error)
	RetrieveCertificate(fingerprint string) (certificate *api.Certificate, err error)
	DeleteCertificate(fingerprint string) (op Operation, err error)
}

// The Client interface represents all available REST client operations
type Client interface {
	ServiceURL() string
	HTTPTransport() *http.Transport

	SetTransportTimeout(timeout time.Duration)

	QueryStruct(method, path string, params QueryParams, header http.Header, body io.Reader, ETag string, target interface{}) (etag string, err error)
	QueryOperation(method, path string, params QueryParams, header http.Header, body io.Reader, ETag string) (operation Operation, etag string, err error)
	CallAPI(method, path string, params QueryParams, header http.Header, body io.Reader, ETag string) (response *api.Response, etag string, err error)
	DownloadFile(path string, params QueryParams, header http.Header, downloader func(header *http.Header, body io.ReadCloser) error) error

	Websocket(resource string) (conn *websocket.Conn, err error)

	// Event handling functions
	GetEvents() (listener *EventListener, err error)
}
