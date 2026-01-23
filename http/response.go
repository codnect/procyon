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
	"errors"
	"io"
	"net/http"
)

// ServerResponse represents an HTTP response being constructed for the client.
// It buffers status and headers until they are committed to the underlying writer.
type ServerResponse struct {
	ctx    *Context
	writer http.ResponseWriter

	headers Headers
	status  Status

	writtenHeaders bool
	writerUsed     bool
}

// Context returns the Context associated with this response.
func (r *ServerResponse) Context() *Context {
	return r.ctx
}

// AddCookie adds a Set-Cookie header to the response.
// If headers have already been written, the call is ignored.
// If the cookie path is empty, it defaults to "/".
func (r *ServerResponse) AddCookie(cookie *Cookie) {
	if r.writtenHeaders {
		return
	}

	if cookie.Path == "" {
		cookie.Path = "/"
	}

	if v := cookie.String(); v != "" {
		r.headers.Add("Set-Cookie", v)
	}
}

// AddHeader appends a header value without replacing existing ones.
// If headers have already been written, the call is ignored.
func (r *ServerResponse) AddHeader(name string, value string) {
	if r.writtenHeaders {
		return
	}

	r.headers.Add(name, value)
}

// SetHeader sets a header value, replacing any existing values.
// If headers have already been written, the call is ignored.
func (r *ServerResponse) SetHeader(name string, value string) {
	if r.writtenHeaders {
		return
	}

	r.headers.Set(name, value)
}

// DeleteHeader removes all values for the given header name.
// If headers have already been written, the call is ignored.
func (r *ServerResponse) DeleteHeader(name string) {
	if r.writtenHeaders {
		return
	}

	r.headers.Del(name)
}

// Header returns the first value for the given response header name.
// It returns ("", false) if name is empty or the header is not present.
func (r *ServerResponse) Header(name string) (string, bool) {
	if name == "" {
		return "", false
	}

	if v := r.headers.Get(name); v != "" {
		return v, true
	}

	return "", false
}

// HeaderValues returns all values for the given response header name.
// It returns nil if name is empty.
func (r *ServerResponse) HeaderValues(name string) []string {
	if name == "" {
		return nil
	}

	v := r.headers.Values(name)
	values := make([]string, len(v))
	copy(values, v)
	return values
}

// Status returns the HTTP status code for the response.
func (r *ServerResponse) Status() Status {
	return r.status
}

// SetStatus sets the HTTP status code for the response.
func (r *ServerResponse) SetStatus(status Status) {
	if r.writtenHeaders {
		return
	}

	r.status = status
}

// Writer returns the underlying io.Writer for writing the response body.
func (r *ServerResponse) Writer() io.Writer {
	r.writerUsed = true
	r.writeHeaders()
	return r.writer
}

// Flush sends any buffered data to the client.
func (r *ServerResponse) Flush() error {
	r.writeHeaders()

	if flusher, ok := r.writer.(http.Flusher); ok {
		flusher.Flush()
	}

	return nil
}

// IsCommitted returns true if the response headers have been written to the client.
func (r *ServerResponse) IsCommitted() bool {
	return r.writtenHeaders
}

// Redirect sets a redirect response with the given location and status code.
// It returns an error if the response has already been committed.
func (r *ServerResponse) Redirect(location string, status Status) error {
	if r.writtenHeaders {
		return errors.New("already committed")
	}

	r.status = status
	r.headers.Set("Location", location)
	r.writeHeaders()
	return nil
}

// Reset clears the response status and headers if not yet committed.
// It returns an error if the response has already been committed.
func (r *ServerResponse) Reset() error {
	if r.IsCommitted() {
		return errors.New("already committed")
	}

	r.status = StatusOK
	r.headers = http.Header{}
	return nil
}

// writeHeaders writes the buffered headers and status to the underlying writer.
func (r *ServerResponse) writeHeaders() {
	if r.writtenHeaders {
		return
	}

	if r.status == 0 {
		r.status = StatusOK
	}

	for key, values := range r.headers {
		for _, value := range values {
			r.writer.Header().Add(key, value)
		}
	}

	r.writer.WriteHeader(int(r.status))
	r.writtenHeaders = true
}
