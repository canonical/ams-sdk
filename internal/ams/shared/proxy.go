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
//
// Taken from https://github.com/canonical/lxd/blob/master/shared/proxy.go under the
// terms of Apache-2.0

package shared

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

var (
	httpProxyEnv = &envOnce{
		names: []string{"HTTP_PROXY", "http_proxy"},
	}
	httpsProxyEnv = &envOnce{
		names: []string{"HTTPS_PROXY", "https_proxy"},
	}
	noProxyEnv = &envOnce{
		names: []string{"NO_PROXY", "no_proxy"},
	}
)

type envOnce struct {
	names []string
	once  sync.Once
	val   string
}

func (e *envOnce) Get() string {
	e.once.Do(e.init)
	return e.val
}

func (e *envOnce) init() {
	for _, n := range e.names {
		e.val = os.Getenv(n)
		if e.val != "" {
			return
		}
	}
}

// ProxyFromEnvironment This is basically the same as golang's ProxyFromEnvironment, except it
// doesn't fall back to http_proxy when https_proxy isn't around, which is
// incorrect behavior. It still respects HTTP_PROXY, HTTPS_PROXY, and NO_PROXY.
func ProxyFromEnvironment(req *http.Request) (*url.URL, error) {
	return ProxyFromConfig("", "", "")(req)
}

// ProxyFromConfig configures the proxy from the given values. If all values are empty a call
// is equals to ProxyFromEnvironment()
func ProxyFromConfig(httpsProxy string, httpProxy string, noProxy string) func(req *http.Request) (*url.URL, error) {
	return func(req *http.Request) (*url.URL, error) {
		var proxy, port string
		var err error

		switch req.URL.Scheme {
		case "https":
			proxy = httpsProxy
			if proxy == "" {
				proxy = httpsProxyEnv.Get()
			}
			port = ":443"
		case "http":
			proxy = httpProxy
			if proxy == "" {
				proxy = httpProxyEnv.Get()
			}
			port = ":80"
		default:
			return nil, fmt.Errorf("unknown scheme %s", req.URL.Scheme)
		}

		if proxy == "" {
			return nil, nil
		}

		addr := req.URL.Host
		if !hasPort(addr) {
			addr = addr + port
		}

		use, err := useProxy(addr, noProxy)
		if err != nil {
			return nil, err
		}
		if !use {
			return nil, nil
		}

		proxyURL, err := url.Parse(proxy)
		if err != nil || !strings.HasPrefix(proxyURL.Scheme, "http") {
			// proxy was bogus. Try prepending "http://" to it and
			// see if that parses correctly. If not, we fall
			// through and complain about the original one.
			if proxyURL, err := url.Parse("http://" + proxy); err == nil {
				return proxyURL, nil
			}
		}
		if err != nil {
			return nil, fmt.Errorf("invalid proxy address %q: %v", proxy, err)
		}
		return proxyURL, nil
	}
}

func hasPort(s string) bool {
	n := strings.LastIndex(s, ":")
	if n < 0 {
		return false
	}
	m := strings.LastIndex(s, "]")
	if m < 0 {
		if len(s) > n+1 {
			return s[n+1] != '/'
		}
		return false
	}
	return n > m
}

func useProxy(addr string, noProxy string) (bool, error) {
	if noProxy == "" {
		noProxy = noProxyEnv.Get()
	}

	if len(addr) == 0 {
		return true, nil
	}
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return false, nil
	}
	if host == "localhost" {
		return false, nil
	}
	if ip := net.ParseIP(host); ip != nil {
		if ip.IsLoopback() {
			return false, nil
		}
	}

	if noProxy == "*" {
		return false, nil
	}

	addr = strings.ToLower(strings.TrimSpace(addr))
	if hasPort(addr) {
		addr = addr[:strings.LastIndex(addr, ":")]
	}

	for _, p := range strings.Split(noProxy, ",") {
		p = strings.ToLower(strings.TrimSpace(p))
		if len(p) == 0 {
			continue
		}
		if hasPort(p) {
			p = p[:strings.LastIndex(p, ":")]
		}
		if addr == p {
			return false, nil
		}
		if p[0] == '.' && (strings.HasSuffix(addr, p) || addr == p[1:]) {
			// noProxy ".foo.com" matches "bar.foo.com" or "foo.com"
			return false, nil
		}
		if p[0] != '.' && strings.HasSuffix(addr, p) && addr[len(addr)-len(p)-1] == '.' {
			// noProxy "foo.com" matches "bar.foo.com"
			return false, nil
		}
	}
	return true, nil
}
