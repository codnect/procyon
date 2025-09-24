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
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestURLResource_Name(t *testing.T) {

	testCases := []struct {
		name string
		url  string

		wantName string
	}{
		{
			name:     "host only without trailing slash",
			url:      "http://localhost:8080",
			wantName: "localhost:8080",
		},
		{
			name:     "host only with trailing slash",
			url:      "http://localhost:8080/",
			wantName: "localhost:8080",
		},
		{
			name:     "path",
			url:      "http://localhost:8080/resources/procyon.yaml",
			wantName: "procyon.yaml",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			resourceUrl, urlErr := url.Parse(tc.url)
			require.NoError(t, urlErr)

			urlResource := NewURLResource(resourceUrl)

			// when
			name := urlResource.Name()

			// then
			require.Equal(t, tc.wantName, name)
		})
	}
}

func TestURLResource_Location(t *testing.T) {
	// given
	resourceUrl, urlErr := url.Parse("http://localhost:8080/resources/procyon.yaml")
	require.NoError(t, urlErr)

	urlResource := NewURLResource(resourceUrl)

	// when
	location := urlResource.Location()

	// then
	require.Equal(t, "http://localhost:8080/resources/procyon.yaml", location)
}

func TestURLResource_Exists(t *testing.T) {

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/resources/procyon.yaml" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	testServer := httptest.NewServer(handler)
	defer testServer.Close()

	testCases := []struct {
		name string
		url  string

		wantErr    error
		wantResult bool
	}{
		{
			name:       "resource exists",
			url:        fmt.Sprintf("%s/resources/procyon.yaml", testServer.URL),
			wantResult: true,
		},
		{
			name:       "resource does not exist",
			url:        fmt.Sprintf("%s/resources/unknown.yaml", testServer.URL),
			wantErr:    errors.New("resource does not exist"),
			wantResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			resourceUrl, urlErr := url.Parse(tc.url)
			require.NoError(t, urlErr)

			urlResource := NewURLResource(resourceUrl)

			// when
			exists := urlResource.Exists()

			// then
			require.Equal(t, tc.wantResult, exists)
		})
	}
}

func TestURLResource_Reader(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/resources/procyon.yaml" {
			w.Write([]byte("anyData"))
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	testServer := httptest.NewServer(handler)
	defer testServer.Close()

	testCases := []struct {
		name     string
		url      string
		wantErr  error
		wantData []byte
	}{
		{
			name:    "resource does not exists",
			url:     fmt.Sprintf("%s/resources/unknown.yaml", testServer.URL),
			wantErr: errors.New("resource does not exist"),
		},
		{
			name:     "valid url resource",
			url:      fmt.Sprintf("%s/resources/procyon.yaml", testServer.URL),
			wantData: []byte("anyData"),
		},
		{
			name: "no host",
			url:  "http://localhost:8080/resources/procyon.yaml",
			wantErr: &url.Error{
				Op:  "Get",
				URL: "http://localhost:8080/resources/procyon.yaml",
				Err: &net.OpError{
					Op:     "dial",
					Net:    "tcp",
					Source: nil,
					Addr: &net.TCPAddr{
						IP:   net.ParseIP("127.0.0.1"),
						Port: 8080,
					},
					Err: &os.SyscallError{
						Syscall: "connect",
						Err:     syscall.ECONNREFUSED,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			resourceUrl, urlErr := url.Parse(tc.url)
			require.NoError(t, urlErr)

			res := NewURLResource(resourceUrl)

			// when
			reader, err := res.Reader()

			// then
			if tc.wantErr != nil {
				require.Nil(t, reader)
				require.Equal(t, tc.wantErr.Error(), err.Error())
				return
			}

			data, _ := io.ReadAll(reader)
			require.Equal(t, tc.wantData, data)
		})
	}
}

func TestURLResource_URL(t *testing.T) {
	// given
	resourceUrl, urlErr := url.Parse("http://localhost:8080/resources/procyon.yaml")
	require.NoError(t, urlErr)

	urlResource := NewURLResource(resourceUrl)

	// when
	result := urlResource.URL()

	// then
	require.Equal(t, resourceUrl, result)
}
