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

package client

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"sync"
	"time"

	"github.com/anbox-cloud/ams-sdk/pkg/ams/shared/rest/api"
)

// Default value for client requests to wait for a reply
const (
	DefaultTransportTimeout = 1 * time.Minute
)

// Doer is the implementation of the Client engine
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// http client to use REST API
type client struct {
	Doer

	serviceURL *url.URL

	eventListeners     []*EventListener
	eventListenersLock *sync.Mutex

	// TODO for now this is not being filled anywhere
	httpUserAgent string
}

// QueryParams request query parameter
type QueryParams map[string]string

// NewTLSClient returns a REST client. Depending on provided addr parameter, it
// connects to a remote network server or through a unix socket
func NewTLSClient(addr any, tlsConfig *tls.Config) (Client, error) {
	if addr == nil {
		return nil, errors.New("Empty address given")
	}

	switch addr.(type) {
	case *url.URL:
		return newNetworkClient(addr.(*url.URL), tlsConfig)
	case string:
		return newUnixSocketClient(addr.(string))
	default:
		return nil, errors.New("Invalid address type given")
	}
}

// newNetworkClient returns a new REST client pointing to remote address received as parameter
// The connection is TLS enabled or not depending on the remote address schema.
// If TLS is enabled, a proper TLS config must be supplied as second parameter
func newNetworkClient(url *url.URL, tlsConfig *tls.Config) (Client, error) {
	if url == nil {
		return nil, fmt.Errorf("Invalid URL given")
	}

	if tlsConfig == nil {
		tlsConfig = &tls.Config{InsecureSkipVerify: true}
	}

	c := &client{
		serviceURL:         url,
		eventListenersLock: &sync.Mutex{},
		Doer: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
			Timeout: DefaultTransportTimeout,
		},
	}
	return c, nil
}

// NewOIDCClient constructs a new oidcClient, ensuring the token field is non-nil to prevent panics during authentication.
func NewOIDCClient(addr any, tlsConfig *tls.Config, tokenProvier TokenProvider) (Client, error) {
	if addr == nil {
		return nil, errors.New("empty address given")
	}
	url, ok := addr.(*url.URL)
	if !ok {
		return nil, fmt.Errorf("invalid URL given")
	}
	if tlsConfig == nil {
		tlsConfig = &tls.Config{InsecureSkipVerify: false}
	}
	client := client{
		Doer: &oidcClient{
			Client: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: tlsConfig,
				},
				Timeout: DefaultTransportTimeout,
			},
			tokenProvider: tokenProvier,
		},
		serviceURL:         url,
		eventListenersLock: &sync.Mutex{},
	}
	return &client, nil
}

// newUnixSocketClient returns a REST client pointing to local unix socket
func newUnixSocketClient(path string) (Client, error) {
	// Setup a Unix socket dialer
	unixDial := func(network, addr string) (net.Conn, error) {
		raddr, err := net.ResolveUnixAddr("unix", path)
		if err != nil {
			return nil, err
		}

		return net.DialUnix("unix", nil, raddr)
	}

	unixSocketServiceURL, err := url.Parse("http://unix")
	if err != nil {
		return nil, err
	}

	c := &client{
		Doer: &http.Client{
			Transport: &http.Transport{
				Dial:              unixDial,
				DisableKeepAlives: true,
			},
			Timeout: DefaultTransportTimeout,
		},
		serviceURL:         unixSocketServiceURL,
		eventListenersLock: &sync.Mutex{},
	}

	return c, nil
}

func extractErrorFromResponse(resp *http.Response) error {
	var errorResponse struct {
		Code    int    `json:"error_code"`
		Message string `json:"error"`
	}

	err := json.NewDecoder(resp.Body).Decode(&errorResponse)
	if err != nil {
		return err
	}

	return fmt.Errorf("%s", errorResponse.Message)
}

// ServiceURL returns the URL of the service the client is connected to
func (c *client) ServiceURL() string {
	return c.serviceURL.String()
}

// HTTPTransport returns the HTTP transport the client uses internally
func (c *client) HTTPTransport() *http.Transport {
	return c.Doer.(*http.Client).Transport.(*http.Transport)
}

// SetTimeout overwrites default timeout of the client with a new one
func (c *client) SetTransportTimeout(timeout time.Duration) {
	if c.Doer == nil {
		return
	}

	c.Doer.(*http.Client).Timeout = timeout
}

// QueryStruct sends a request to the server and stores response in a struct
func (c *client) QueryStruct(method, path string, params QueryParams, header http.Header, body io.Reader, etag string, target interface{}) (string, error) {
	resp, etag, err := c.CallAPI(method, path, params, header, body, etag)
	if err != nil {
		return "", err
	}

	err = resp.MetadataAsStruct(&target)
	return etag, err
}

func pretty(input interface{}) string {
	pretty, err := json.MarshalIndent(input, "\t", "\t")
	if err != nil {
		return fmt.Sprintf("%v", input)
	}

	return fmt.Sprintf("\n\t%s", pretty)
}

// QueryOperation sends a request to the server that will return an async response in an Operation object
// that allows additional logic like wait for completion or cancel it
func (c *client) QueryOperation(method, path string, params QueryParams, header http.Header, body io.Reader, etag string) (Operation, string, error) {
	// Attempt to setup an early event listener
	listener, err := c.GetEvents()
	if err != nil {
		listener = nil
	}

	resp, etag, err := c.CallAPI(method, path, params, header, body, etag)
	if err != nil {
		if listener != nil {
			listener.Disconnect()
		}
		return nil, "", err
	}

	apiOp, err := resp.MetadataAsOperation()
	if err != nil {
		if listener != nil {
			listener.Disconnect()
		}
		return nil, "", err
	}

	op := operation{
		Operation: *apiOp,
		c:         &operations{c},
		listener:  listener,
		chActive:  make(chan bool),
	}

	return &op, etag, nil
}

// CallAPI requests a REST api method with provided query params and body and returns related http response
func (c *client) CallAPI(method, path string, params QueryParams, header http.Header, body io.Reader, etag string) (*api.Response, string, error) {
	resp, err := c.performRequest(method, path, params, header, body, etag)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	return c.parseResponse(resp)
}

func (c *client) DownloadFile(path string, params QueryParams, header http.Header, downloader func(header *http.Header, body io.ReadCloser) error) error {
	resp, err := c.performRequest("GET", path, params, header, nil, "")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// NOTE: As a fileResposne is an inline response and different with
	// generic api.Response so that we can't simply parse the response
	// directly unless http status code is not StatusOK
	if resp.StatusCode != http.StatusOK {
		_, _, err := c.parseResponse(resp)
		return err
	}

	return downloader(&resp.Header, resp.Body)
}

func (c *client) performRequest(method, path string, params QueryParams, header http.Header, body io.Reader, etag string) (*http.Response, error) {
	u := c.serviceURL.ResolveReference(
		&url.URL{
			Path: path,
		},
	)

	r, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	v := r.URL.Query()
	for key, value := range params {
		v.Add(key, value)
	}
	r.URL.RawQuery = v.Encode()

	if len(c.httpUserAgent) > 0 {
		r.Header.Set("User-Agent", c.httpUserAgent)
	}

	if header != nil {
		for k, v := range header {
			for _, s := range v {
				r.Header.Add(k, s)
			}
		}
	}

	return c.Doer.Do(r)
}

// Internal functions
func (c *client) parseResponse(resp *http.Response) (*api.Response, string, error) {
	// Get the ETag
	etag := resp.Header.Get("ETag")

	decoder := json.NewDecoder(resp.Body)
	response := api.Response{}
	err := decoder.Decode(&response)

	// Not all API calls return a proper api.Response, in those cases we just print the status text
	if err != nil {
		return nil, "", errors.New(http.StatusText(resp.StatusCode))
	}

	if response.Type == api.ResponseTypeError {
		if response.Error != "" {
			return nil, "", errors.New(response.Error)
		}
		return nil, "", errors.New(http.StatusText(resp.StatusCode))
	}

	return &response, etag, nil
}

// APIPath prefixes API version to a path
func APIPath(path ...string) string {
	p := append([]string{"/", api.Version}, path...)
	return filepath.Join(p...)
}
