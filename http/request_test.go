// Copyright 2026 Codnect
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

package http

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerRequest_Context(t *testing.T) {
	// given
	requestCtx := &Context{}
	serverRequest := ServerRequest{
		ctx: requestCtx,
	}

	// when
	ctx := serverRequest.Context()

	// then
	assert.Equal(t, requestCtx, ctx)
}

func TestServerRequest_Cookie(t *testing.T) {
	testCases := []struct {
		name            string
		cookies         []*Cookie
		cookieName      string
		wantCookieValue string
		wantExists      bool
	}{
		{
			name:            "empty cookie name",
			cookieName:      "",
			wantExists:      false,
			wantCookieValue: "",
		},
		{
			name: "cookie exists",
			cookies: []*Cookie{
				{
					Name:  "testCookie",
					Value: "testCookieValue",
				},
			},
			cookieName:      "testCookie",
			wantExists:      true,
			wantCookieValue: "testCookieValue",
		},
		{
			name:       "cookie not exists",
			cookies:    []*Cookie{},
			cookieName: "testCookie",
			wantExists: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			nativeReq := httptest.NewRequest(http.MethodGet, "/", nil)

			for _, cookie := range tc.cookies {
				nativeReq.AddCookie(cookie)
			}

			serverRequest := ServerRequest{
				nativeReq: nativeReq,
			}

			// when
			cookie, ok := serverRequest.Cookie(tc.cookieName)

			// then
			assert.Equal(t, tc.wantExists, ok)
			if tc.wantExists {
				assert.Equal(t, tc.wantCookieValue, cookie.Value)
			}
		})
	}
}

func TestServerRequest_CookieValues(t *testing.T) {
	testCases := []struct {
		name        string
		cookies     []*Cookie
		cookieName  string
		wantCookies []*Cookie
	}{
		{
			name:        "empty cookie name",
			wantCookies: []*Cookie{},
		},
		{
			name: "cookie exists",
			cookies: []*Cookie{
				{
					Name:  "testCookie",
					Value: "testCookieValue",
				},
			},
			cookieName: "testCookie",
			wantCookies: []*Cookie{
				{
					Name:  "testCookie",
					Value: "testCookieValue",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			nativeReq := httptest.NewRequest(http.MethodGet, "/", nil)
			for _, c := range tc.cookies {
				nativeReq.AddCookie(&http.Cookie{Name: c.Name, Value: c.Value})
			}
			serverRequest := ServerRequest{
				nativeReq: nativeReq,
			}

			// when
			values := serverRequest.CookieValues(tc.cookieName)

			// then
			assert.Equal(t, len(tc.wantCookies), len(values))
			for i, want := range tc.wantCookies {
				assert.Equal(t, want.Name, values[i].Name)
				assert.Equal(t, want.Value, values[i].Value)
			}
		})
	}
}

func TestServerRequest_Cookies(t *testing.T) {
	testCases := []struct {
		name        string
		cookies     []*Cookie
		cookieName  string
		wantCookies []*Cookie
	}{
		{
			name:        "no cookies",
			wantCookies: []*Cookie{},
		},
		{
			name: "cookie exists",
			cookies: []*Cookie{
				{
					Name:  "testCookie",
					Value: "testCookieValue",
				},
			},
			cookieName: "testCookie",
			wantCookies: []*Cookie{
				{
					Name:  "testCookie",
					Value: "testCookieValue",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			nativeReq := httptest.NewRequest(http.MethodGet, "/", nil)
			for _, c := range tc.cookies {
				nativeReq.AddCookie(&http.Cookie{Name: c.Name, Value: c.Value})
			}
			serverRequest := ServerRequest{
				nativeReq: nativeReq,
			}

			// when
			values := serverRequest.Cookies()

			// then
			assert.Equal(t, len(tc.wantCookies), len(values))
			for i, want := range tc.wantCookies {
				assert.Equal(t, want.Name, values[i].Name)
				assert.Equal(t, want.Value, values[i].Value)
			}
		})
	}
}

func TestServerRequest_Query(t *testing.T) {
	testCases := []struct {
		name           string
		location       string
		queryParamName string
		wantParamValue string
		wantExists     bool
	}{
		{
			name:       "empty query param name",
			location:   "/api/v1/users?name=test",
			wantExists: false,
		},
		{
			name:           "query param exists",
			location:       "/api/v1/users?name=test",
			queryParamName: "name",
			wantParamValue: "test",
			wantExists:     true,
		},
		{
			name:           "query param does not exist",
			location:       "/api/v1/users?name=test",
			queryParamName: "another",
			wantParamValue: "",
			wantExists:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			nativeReq := httptest.NewRequest(http.MethodGet, tc.location, nil)

			serverRequest := ServerRequest{
				nativeReq: nativeReq,
			}

			// when
			queryParamVal, ok := serverRequest.Query(tc.queryParamName)

			// then
			assert.Equal(t, tc.wantExists, ok)
			if tc.wantExists {
				assert.Equal(t, tc.wantParamValue, queryParamVal)
			}
		})
	}
}

func TestServerRequest_QueryValues(t *testing.T) {
	testCases := []struct {
		name            string
		location        string
		queryParamName  string
		wantParamValues []string
		wantExists      bool
	}{
		{
			name:       "empty query param name",
			location:   "/api/v1/users?name=test",
			wantExists: false,
		},
		{
			name:            "query param exists",
			location:        "/api/v1/users?name=test",
			queryParamName:  "name",
			wantParamValues: []string{"test"},
			wantExists:      true,
		},
		{
			name:            "multiple values for query param",
			location:        "/api/v1/users?name=test&name=another",
			queryParamName:  "name",
			wantParamValues: []string{"test", "another"},
		},
		{
			name:           "query param does not exist",
			location:       "/api/v1/users?name=test",
			queryParamName: "another",
			wantExists:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			nativeReq := httptest.NewRequest(http.MethodGet, tc.location, nil)

			serverRequest := ServerRequest{
				nativeReq: nativeReq,
			}

			// when
			values := serverRequest.QueryValues(tc.queryParamName)

			// then
			assert.ElementsMatch(t, tc.wantParamValues, values)
		})
	}
}

func TestServerRequest_QueryString(t *testing.T) {
	testCases := []struct {
		name            string
		location        string
		wantQueryString string
	}{
		{
			name:            "no query param",
			location:        "/api/v1/users",
			wantQueryString: "",
		},
		{
			name:            "query param exists",
			location:        "/api/v1/users?name=test&page=2",
			wantQueryString: "name=test&page=2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			nativeReq := httptest.NewRequest(http.MethodGet, tc.location, nil)
			serverRequest := ServerRequest{
				nativeReq: nativeReq,
			}

			// when
			queryString := serverRequest.QueryString()

			// then
			assert.Equal(t, tc.wantQueryString, queryString)
		})
	}
}

func TestServerRequest_Header(t *testing.T) {
	testCases := []struct {
		name            string
		headers         map[string][]string
		headerName      string
		wantHeaderValue string
		wantExists      bool
	}{
		{
			name:            "empty header name",
			headerName:      "",
			wantExists:      false,
			wantHeaderValue: "",
		},
		{
			name: "header exists",
			headers: map[string][]string{
				"testHeader": {"testHeaderValue"},
			},
			headerName:      "testHeader",
			wantExists:      true,
			wantHeaderValue: "testHeaderValue",
		},
		{
			name:       "header not exists",
			headers:    map[string][]string{},
			headerName: "testHeader",
			wantExists: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			nativeReq := httptest.NewRequest(http.MethodGet, "/", nil)

			for header, values := range tc.headers {
				for _, value := range values {
					nativeReq.Header.Add(header, value)
				}
			}

			serverRequest := ServerRequest{
				nativeReq: nativeReq,
			}

			// when
			header, ok := serverRequest.Header(tc.headerName)

			// then
			assert.Equal(t, tc.wantExists, ok)
			if tc.wantExists {
				assert.Equal(t, tc.wantHeaderValue, header)
			}
		})
	}
}

func TestServerRequest_HeaderValues(t *testing.T) {
	testCases := []struct {
		name             string
		headers          map[string][]string
		headerName       string
		wantHeaderValues []string
	}{
		{
			name:             "empty header name",
			headerName:       "",
			wantHeaderValues: []string{},
		},
		{
			name: "header exists",
			headers: map[string][]string{
				"testHeader": {"testHeaderValue"},
			},
			headerName:       "testHeader",
			wantHeaderValues: []string{"testHeaderValue"},
		},
		{
			name: "multiple header values for header name",
			headers: map[string][]string{
				"testHeader": {"testHeaderValue", "anotherHeaderValue"},
			},
			headerName:       "testHeader",
			wantHeaderValues: []string{"testHeaderValue", "anotherHeaderValue"},
		},
		{
			name:             "header not exists",
			headers:          map[string][]string{},
			headerName:       "testHeader",
			wantHeaderValues: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			nativeReq := httptest.NewRequest(http.MethodGet, "/", nil)

			for header, values := range tc.headers {
				for _, value := range values {
					nativeReq.Header.Add(header, value)
				}
			}

			serverRequest := ServerRequest{
				nativeReq: nativeReq,
			}

			// when
			values := serverRequest.HeaderValues(tc.headerName)

			// then
			assert.ElementsMatch(t, tc.wantHeaderValues, values)
		})
	}
}

func TestServerRequest_Path(t *testing.T) {
	// given
	nativeReq := httptest.NewRequest(http.MethodGet, "/api/v1/users?name=test", nil)
	serverRequest := ServerRequest{
		nativeReq: nativeReq,
	}

	// when
	path := serverRequest.Path()

	// then
	assert.Equal(t, "/api/v1/users", path)
}

func TestServerRequest_PathValues(t *testing.T) {
	testCases := []struct {
		name       string
		pathValues pathValues
		pathKey    string
		wantValue  string
	}{
		{
			name:       "empty path value",
			pathValues: pathValues{},
			pathKey:    "",
			wantValue:  "",
		},
		{
			name: "path exists",
			pathValues: pathValues{
				values: [16]pathValue{
					{
						name:  "testPathKey",
						value: "testValue",
					},
				},
				count: 1,
			},
			pathKey:   "testPathKey",
			wantValue: "testValue",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			serverRequest := ServerRequest{
				pathValues: tc.pathValues,
			}

			// when
			value := serverRequest.PathValue(tc.pathKey)

			// then
			assert.Equal(t, tc.wantValue, value)
		})
	}
}

func TestServerRequest_Method(t *testing.T) {
	// given
	serverRequest := ServerRequest{
		nativeReq: &http.Request{
			Method: "GET",
		},
	}

	// when
	method := serverRequest.Method()

	// then
	assert.Equal(t, MethodGet, method)
}

func TestServerRequest_Body(t *testing.T) {
	// given
	anyBody := io.NopCloser(bytes.NewBufferString("testBody"))
	serverRequest := ServerRequest{
		nativeReq: &http.Request{
			Body: anyBody,
		},
	}

	// when
	body := serverRequest.Body()

	// then
	assert.Equal(t, anyBody, body)
}

func TestServerRequest_Scheme(t *testing.T) {
	testCases := []struct {
		name       string
		nativeReq  *http.Request
		wantScheme string
	}{
		{
			name:       "http scheme",
			nativeReq:  &http.Request{},
			wantScheme: "http",
		},
		{
			name: "https scheme",
			nativeReq: &http.Request{
				TLS: &tls.ConnectionState{},
			},
			wantScheme: "https",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			serverRequest := &ServerRequest{
				nativeReq: tc.nativeReq,
			}

			// when
			scheme := serverRequest.Scheme()

			// then
			assert.Equal(t, tc.wantScheme, scheme)
		})
	}
}

func TestServerRequest_IsSecure(t *testing.T) {
	testCases := []struct {
		name       string
		nativeReq  *http.Request
		wantSecure bool
	}{
		{
			name:       "not secure",
			nativeReq:  &http.Request{},
			wantSecure: false,
		},
		{
			name: "secure",
			nativeReq: &http.Request{
				TLS: &tls.ConnectionState{},
			},
			wantSecure: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			serverRequest := &ServerRequest{
				nativeReq: tc.nativeReq,
			}

			// when
			isSecure := serverRequest.IsSecure()

			// then
			assert.Equal(t, tc.wantSecure, isSecure)
		})
	}
}
