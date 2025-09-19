// Copyright 2025 Codnect
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package io

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// URLResource struct represents a file-based resource.
type URLResource struct {
	url        *url.URL
	httpClient *http.Client
}

// NewURLResource function creates a new URLResource with the given URL.
func NewURLResource(url *url.URL) *URLResource {
	return &URLResource{
		url: url,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Name method returns the name of the URL resource.
// It extracts the name from the URL path.
func (u *URLResource) Name() string {
	path := u.url.Path
	if path == "" || path == "/" {
		return u.url.Host
	}

	segments := strings.Split(strings.TrimSuffix(path, "/"), "/")
	name := segments[len(segments)-1]
	return name
}

// Location method returns the location of the URL resource.
func (u *URLResource) Location() string {
	return u.url.String()
}

// Exists method checks if the URL resource exists by making a HEAD request.
func (u *URLResource) Exists() bool {
	resp, err := u.httpClient.Head(u.url.String())
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}

	if resp.Body != nil {
		_ = resp.Body.Close()
	}

	return resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices
}

// Reader method makes a GET request to the URL and returns an io.ReadCloser.
func (u *URLResource) Reader() (io.ReadCloser, error) {
	resp, err := u.httpClient.Get(u.url.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {

		if resp.Body != nil {
			_ = resp.Body.Close()
		}

		return nil, errors.New("resource does not exist")
	}

	return resp.Body, nil
}

// URL method returns the underlying URL of the resource.
func (u *URLResource) URL() *url.URL {
	return u.url
}
