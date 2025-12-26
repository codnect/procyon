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

package http

import "io"

// Handler represents an HTTP handler function.
type Handler func(Context)

// Middleware represents an HTTP middleware component.
type Middleware = Handler

// Context represents the request-scoped context passed to handlers and middleware.
type Context interface {
	// Request returns the incoming HTTP request abstraction.
	Request() Request
	// Response returns the outgoing HTTP response abstraction.
	Response() Response
	// Security returns the active SecurityContext.
	Security() SecurityContext
	// SetSecurity replaces the current security context.
	SetSecurity(SecurityContext)
	// Next executes the next middleware or handler in the chain.
	Next()
	// Params returns the matched route parameters.
	Params() map[string]string
	// Param retrieves a specific route parameter by name.
	Param(string) string
	// Set stores a value in the context.
	Set(string, any)
	// Get retrieves a stored value and a boolean indicating presence.
	Get(string) (any, bool)
}

// Request is an abstract HTTP request representation.
type Request interface {
	// Method returns the HTTP method of the request.
	Method() string
	// Path returns the request path.
	Path() string
	// Header returns the first value for the named header.
	Header(string) string
	// Body returns the raw request body reader.
	Body() io.ReadCloser
}

// Response is an abstract HTTP response representation.
type Response interface {
	// Status writes the HTTP status code.
	Status(int)
	// Header sets a header value on the response.
	Header(string, string)
	// Write writes the response body bytes.
	Write([]byte) (int, error)
}
