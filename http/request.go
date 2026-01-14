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
	"io"
	"net/http"
	"net/url"
)

type Void struct{}

// HandlerFunc represents a delegate that can handle an HTTP defaultServerRequest.
// It is used in the middleware chain to invoke the next handler.
type HandlerFunc func(ctx *Context) (Result, error)

func (f HandlerFunc) Handle(ctx *Context) (Result, error) {
	return f(ctx)
}

type Handler interface {
	Handle(ctx *Context) (Result, error)
}

type RequestDelegate func(ctx *Context) error
type FilterDelegate func(ctx *Context) (Result, error)

// ServerRequest interface represents an HTTP defaultServerRequest.
type ServerRequest interface {
	// Context returns the context associated with the defaultServerRequest.
	Context() *Context

	// Cookie returns the cookie with the specified name.
	Cookie(name string) (*Cookie, bool)
	// Cookies returns all the cookies associated with the defaultServerRequest.
	Cookies() []*Cookie

	// QueryParam returns the query parameter with the specified name.
	QueryParam(name string) (string, bool)
	// QueryParamNames returns the names of all query parameters.
	QueryParamNames() []string
	// QueryParams returns all the query parameters with the specified name.
	QueryParams(name string) []string
	// QueryString returns the query string of the defaultServerRequest.
	QueryString() string

	// Header returns the header with the specified name.
	Header(name string) (string, bool)
	// HeaderNames returns the names of all headers.
	HeaderNames() []string
	// Headers returns all the headers with the specified name.
	Headers(name string) []string

	// Path returns the path of the defaultServerRequest.
	Path() string
	// PathValue returns the value of the path parameter with the specified name.
	PathValue(name string) (string, bool)
	// Method returns the method of the defaultServerRequest.
	Method() Method
	// Body returns the reader of the body.
	Body() io.Reader
	// Scheme returns the scheme of the defaultServerRequest.
	Scheme() string
	// IsSecure returns whether the defaultServerRequest is secure.
	IsSecure() bool
}

type defaultServerRequest struct {
	nativeReq *http.Request
	ctx       *Context
	body      io.Reader

	pathValues   PathValues
	queryCache   url.Values
	cookiesCache []*Cookie
}

func (r *defaultServerRequest) Context() *Context {
	return r.ctx
}

func (r *defaultServerRequest) initCookieCache() {
	if r.cookiesCache == nil {
		r.cookiesCache = r.nativeReq.Cookies()
	}
}

func (r *defaultServerRequest) Cookie(name string) (*Cookie, bool) {
	r.initCookieCache()

	for _, cookie := range r.cookiesCache {
		if cookie.Name == name {
			return cookie, true
		}
	}

	return nil, false
}

func (r *defaultServerRequest) Cookies() []*Cookie {
	r.initCookieCache()
	return r.cookiesCache
}

func (r *defaultServerRequest) initQueryCache() {
	if r.queryCache == nil {
		if r.nativeReq != nil && r.nativeReq.URL != nil {
			r.queryCache = r.nativeReq.URL.Query()
		} else {
			r.queryCache = url.Values{}
		}
	}
}

func (r *defaultServerRequest) QueryParam(name string) (string, bool) {
	r.initQueryCache()

	values, ok := r.queryCache[name]
	if ok {
		return values[0], true
	}

	return "", false
}

func (r *defaultServerRequest) QueryParamNames() []string {
	r.initQueryCache()

	names := make([]string, 0, len(r.queryCache))

	for name := range r.queryCache {
		names = append(names, name)
	}

	return names
}

func (r *defaultServerRequest) QueryParams(name string) []string {
	r.initQueryCache()
	values, ok := r.queryCache[name]
	if ok {
		return values
	}

	return nil
}

func (r *defaultServerRequest) QueryString() string {
	return r.nativeReq.URL.RawQuery
}

func (r *defaultServerRequest) Header(name string) (string, bool) {
	values := r.nativeReq.Header.Values(name)

	if len(values) != 0 {
		return values[0], true
	}

	return "", false
}

func (r *defaultServerRequest) HeaderNames() []string {
	headers := make([]string, 0, len(r.nativeReq.Header))

	for header := range r.nativeReq.Header {
		headers = append(headers, header)
	}

	return headers
}

func (r *defaultServerRequest) Headers(name string) []string {
	return r.nativeReq.Header.Values(name)
}

func (r *defaultServerRequest) Path() string {
	return r.nativeReq.URL.Path
}

func (r *defaultServerRequest) PathValue(name string) (string, bool) {
	return r.pathValues.Value(name)
}

func (r *defaultServerRequest) Method() Method {
	return Method(r.nativeReq.Method)
}

func (r *defaultServerRequest) Body() io.Reader {
	if r.body != nil {
		return r.body
	}

	return r.nativeReq.Body
}

func (r *defaultServerRequest) Scheme() string {
	if r.nativeReq.TLS != nil {
		return "https"
	}

	return "http"
}

func (r *defaultServerRequest) IsSecure() bool {
	return r.nativeReq.TLS != nil
}
