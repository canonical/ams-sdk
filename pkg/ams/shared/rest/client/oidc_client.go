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
	"context"
	"net/http"
	"time"
)

type TokenProvider interface {
	GetAccessToken() (string, error)
	RefreshToken() error
	GetAccessTokenExpiry() (time.Time, error)
	Authenticate(context.Context) error
}

// oidcClient is a structure encapsulating an HTTP client, and attaches a token for each request.
type oidcClient struct {
	Client        *http.Client
	tokenProvider TokenProvider
}

// Do function executes an HTTP request using the oidcClient's http client, and manages authorization by refreshing or authenticating as needed.
// If the request fails with an HTTP Unauthorized status, it attempts to refresh the access token, or perform an OIDC authentication if refresh fails.
func (o *oidcClient) Do(req *http.Request) (*http.Response, error) {
	// Pre maturely refresh token in case the token has expired and a refresh
	// token is available
	expiry, err := o.tokenProvider.GetAccessTokenExpiry()
	if err != nil {
		return nil, err
	}
	if expiry.Before(time.Now()) {
		if err := o.tokenProvider.RefreshToken(); err != nil {
			return nil, err
		}
	}

	token, err := o.tokenProvider.GetAccessToken()
	if err != nil {
		return nil, err
	}
	// Set the new access token in the header.
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := o.Client.Do(req)
	if err != nil {
		return nil, err
	}

	// Return immediately if the error is not HTTP status unauthorized.
	if resp.StatusCode != http.StatusUnauthorized {
		return resp, nil
	}

	err = o.tokenProvider.RefreshToken()
	if err != nil {
		return nil, err
	}
	token, err = o.tokenProvider.GetAccessToken()
	if err != nil {
		return nil, err
	}

	// Set the new access token in the header.
	req.Header.Set("Authorization", "Bearer "+token)
	return o.Client.Do(req)
}
