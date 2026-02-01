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

// RequestDelegate represents the next middleware or handler in the pipeline.
type RequestDelegate func(ctx *Context) error

// ServerRequest wraps an incoming HTTP request and provides convenience accessors
// with small, per-request caches for cookies and query parameters.
type ServerRequest struct {
	ctx       *Context
	nativeReq *http.Request

	cookiesCache []*Cookie
	queryCache   url.Values
}

// Context returns the Context associated with the ServerRequest.
func (r *ServerRequest) Context() *Context {
	return r.ctx
}

// Cookie returns the named cookie provided in the request.
func (r *ServerRequest) Cookie(name string) (*Cookie, bool) {
	if name == "" {
		return nil, false
	}

	r.initCookieCache()

	for _, cookie := range r.cookiesCache {
		if cookie.Name == name {
			return cookie, true
		}
	}

	return nil, false
}

// CookieValues returns all cookies with the given name.
// It returns nil if name is empty.
func (r *ServerRequest) CookieValues(name string) []*Cookie {
	if name == "" {
		return nil
	}

	r.initCookieCache()

	cookies := make([]*Cookie, 0)
	for _, cookie := range r.cookiesCache {
		if cookie.Name == name {
			cookies = append(cookies, cookie)
		}
	}

	return cookies
}

// Cookies returns all cookies sent with the request.
func (r *ServerRequest) Cookies() []*Cookie {
	r.initCookieCache()

	cookies := make([]*Cookie, len(r.cookiesCache))
	copy(cookies, r.cookiesCache)

	return cookies
}

// Query returns the first value associated with the given query parameter name.
// It returns ("", false) if name is empty or the parameter is not present.
func (r *ServerRequest) Query(name string) (string, bool) {
	if name == "" {
		return "", false
	}

	r.initQueryCache()

	val, ok := r.queryCache[name]
	if ok && len(val) > 0 {
		return val[0], true
	}

	return "", false
}

// QueryValues returns all values associated with the given query parameter name.
// It returns nil if name is empty.
func (r *ServerRequest) QueryValues(name string) []string {
	if name == "" {
		return nil
	}

	r.initQueryCache()

	if v, ok := r.queryCache[name]; ok {
		out := make([]string, len(v))
		copy(out, v)
		return out
	}

	return nil
}

// QueryString returns the raw, unparsed query string (the part after '?').
func (r *ServerRequest) QueryString() string {
	return r.nativeReq.URL.RawQuery
}

// Header returns the first value for the given header name.
// It returns ("", false) if name is empty or the header is not present.
func (r *ServerRequest) Header(name string) (string, bool) {
	if name == "" {
		return "", false
	}

	if v := r.nativeReq.Header.Get(name); v != "" {
		return v, true
	}

	return "", false
}

// HeaderValues returns all values for the given header name.
// It returns nil if name is empty.
func (r *ServerRequest) HeaderValues(name string) []string {
	if name == "" {
		return nil
	}

	v := r.nativeReq.Header.Values(name)
	values := make([]string, len(v))
	copy(values, v)
	return values
}

// Path returns the request URL path (without the query string).
func (r *ServerRequest) Path() string {
	return r.nativeReq.URL.Path
}

// PathValue returns the value of a path parameter captured by the route matcher.
// If the parameter does not exist, it returns an empty string.
func (r *ServerRequest) PathValue(name string) string {
	// TODO: Implement path parameter extraction based on the routing mechanism.
	return ""
}

// Method returns the HTTP method (GET, POST, PUT, ...).
func (r *ServerRequest) Method() string {
	return r.nativeReq.Method
}

// Body returns the request body as an io.ReadCloser.
func (r *ServerRequest) Body() io.ReadCloser {
	return r.nativeReq.Body
}

// Scheme returns "https" when TLS is present; otherwise it returns "http".
func (r *ServerRequest) Scheme() string {
	if r.nativeReq.TLS != nil {
		return "https"
	}

	return "http"
}

// IsSecure reports whether the request was made over TLS (HTTPS).
func (r *ServerRequest) IsSecure() bool {
	return r.nativeReq.TLS != nil
}

// initCookieCache lazily caches request cookies on first access.
func (r *ServerRequest) initCookieCache() {
	if r.cookiesCache == nil {
		r.cookiesCache = r.nativeReq.Cookies()
	}
}

// initQueryCache lazily caches parsed query parameters on first access.
func (r *ServerRequest) initQueryCache() {
	if r.queryCache == nil {
		if r.nativeReq != nil && r.nativeReq.URL != nil {
			r.queryCache = r.nativeReq.URL.Query()
		} else {
			r.queryCache = url.Values{}
		}
	}
}
