// Copyright 2014 Mark Wolfe. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buildkite

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

// TokenAuthTransport manages injection of the API token for each request
type TokenAuthTransport struct {
	APIToken string
	Debug    bool
}

// RoundTrip invoked each time a request is made
func (t TokenAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.APIToken))
	ts := time.Now()
	res, err := http.DefaultTransport.RoundTrip(req)
	fmt.Printf("DEBUG uri = %s time = %s\n", req.RequestURI, time.Now().Sub(ts))
	return res, err
}

// Client builds a new http client.
func (t *TokenAuthTransport) Client() *http.Client {
	return &http.Client{Transport: t}
}

// NewTokenConfig configure authentication using an API token
func NewTokenConfig(apiToken string, debug bool) (*TokenAuthTransport, error) {
	if apiToken == "" {
		return nil, fmt.Errorf("Invalid token, empty string supplied")
	}
	return &TokenAuthTransport{APIToken: apiToken, Debug: debug}, nil
}

// BasicAuthTransport manages injection of the authorization header
type BasicAuthTransport struct {
	Username string
	Password string
}

// RoundTrip invoked each time a request is made
func (bat BasicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s",
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s",
			bat.Username, bat.Password)))))
	return http.DefaultTransport.RoundTrip(req)
}

// Client builds a new http client.
func (bat *BasicAuthTransport) Client() *http.Client {
	return &http.Client{Transport: bat}
}

// NewBasicConfig configure authentication using the supplied credentials
func NewBasicConfig(username string, password string) (*BasicAuthTransport, error) {
	if username == "" {
		return nil, fmt.Errorf("Invalid username, empty string supplied")
	}
	if password == "" {
		return nil, fmt.Errorf("Invalid password, empty string supplied")
	}
	return &BasicAuthTransport{username, password}, nil
}
