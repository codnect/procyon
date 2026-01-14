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
	"net/url"
	"strconv"
	"strings"
)

// ServerResponse interface represents an HTTP response.
type ServerResponse interface {
	// Context returns the context associated with the response.
	Context() *Context
	// AddCookie adds a cookie to the response.
	AddCookie(cookie *Cookie)

	// ContentLength returns the content length of the response.
	ContentLength() int
	// SetContentLength sets the content length of the response.
	SetContentLength(len int)

	// CharacterEncoding returns the character encoding of the response.
	CharacterEncoding() string
	// SetCharacterEncoding sets the character encoding of the response.
	SetCharacterEncoding(charset string)

	// ContentType returns the content type of the response.
	ContentType() string
	// SetContentType sets the content type of the response.
	SetContentType(contentType string)

	// AddHeader adds a header to the response.
	// If the header is already set, it will be appended.
	AddHeader(name string, value string)
	// SetHeader sets a header in the response.
	// If the header is already set, it will be overwritten.
	SetHeader(name string, value string)
	// DeleteHeader deletes a header from the response.
	DeleteHeader(name string)
	// Header returns the value of a header in the response.
	Header(name string) (string, bool)
	// HeaderNames returns the names of all headers in the response.
	HeaderNames() []string
	// Headers returns all the values of a header in the response.
	Headers(name string) []string

	// Status returns the status of the response.
	Status() Status
	// SetStatus sets the status of the response.
	SetStatus(status Status)
	// Redirect redirects the response to a location with a status.
	Redirect(location string, status Status) error

	// Writer returns the writer of the response.
	Writer() io.Writer
	// Flush flushes the response.
	Flush() error
	// IsCommitted checks if the response is committed.
	IsCommitted() bool
	// Reset resets the response.
	Reset() error
}

type defaultServerResponse struct {
	responseWriter http.ResponseWriter
	ctx            *Context
	writer         io.Writer

	characterEncoding string
	headers           http.Header
	statusCode        Status
	writtenHeaders    bool
	writerUsed        bool
}

func (r *defaultServerResponse) Context() *Context {
	return r.ctx
}

func (r *defaultServerResponse) AddCookie(cookie *Cookie) {
	if r.writtenHeaders {
		return
	}

	path := cookie.Path
	if path == "" {
		path = "/"
	}

	stdCookie := &http.Cookie{
		Name:     cookie.Name,
		Value:    url.QueryEscape(cookie.Value),
		Path:     path,
		Domain:   cookie.Domain,
		Expires:  cookie.Expires,
		MaxAge:   cookie.MaxAge,
		Secure:   cookie.Secure,
		HttpOnly: cookie.HttpOnly,
		SameSite: cookie.SameSite,
	}

	if v := stdCookie.String(); v != "" {
		r.headers.Add(HeaderSetCookie, v)
	}
}

func (r *defaultServerResponse) ContentLength() int {
	length := r.headers.Get(HeaderContentLength)

	if length != "" {
		val, err := strconv.Atoi(length)
		if err != nil {
			return 0
		}

		return val
	}

	return 0
}

func (r *defaultServerResponse) SetContentLength(len int) {
	if r.writtenHeaders {
		return
	}

	r.headers.Set(HeaderContentLength, strconv.Itoa(len))
}

func (r *defaultServerResponse) CharacterEncoding() string {
	return r.characterEncoding
}

func (r *defaultServerResponse) SetCharacterEncoding(charset string) {
	r.characterEncoding = charset
}

func (r *defaultServerResponse) ContentType() string {
	return r.headers.Get(HeaderContentType)
}

func (r *defaultServerResponse) SetContentType(contentType string) {
	if r.writtenHeaders {
		return
	}

	r.headers.Set(HeaderContentType, contentType)
}

func (r *defaultServerResponse) AddHeader(name string, value string) {
	if r.writtenHeaders {
		return
	}

	r.headers.Add(name, value)
}

func (r *defaultServerResponse) SetHeader(name string, value string) {
	if r.writtenHeaders {
		return
	}

	r.headers.Set(name, value)
}

func (r *defaultServerResponse) DeleteHeader(name string) {
	if r.writtenHeaders {
		return
	}

	r.headers.Del(name)
}

func (r *defaultServerResponse) Header(name string) (string, bool) {
	values := r.headers.Values(name)

	if len(values) != 0 {
		return values[0], true
	}

	return "", false
}

func (r *defaultServerResponse) HeaderNames() []string {
	headers := make([]string, 0, len(r.headers))

	for header := range r.headers {
		headers = append(headers, header)
	}

	return headers
}

func (r *defaultServerResponse) Headers(name string) []string {
	return r.headers.Values(name)
}

func (r *defaultServerResponse) Status() Status {
	return r.statusCode
}

func (r *defaultServerResponse) SetStatus(status Status) {
	if r.writtenHeaders {
		return
	}

	r.statusCode = status
}

func (r *defaultServerResponse) Redirect(location string, status Status) error {
	if r.writtenHeaders {
		return errors.New("already committed")
	}

	r.headers.Set(HeaderLocation, location)
	r.statusCode = status
	return nil
}

func (r *defaultServerResponse) Writer() io.Writer {
	r.writerUsed = true
	r.writeHeaders()

	if r.writer != nil {
		return r.writer
	}

	return r.responseWriter
}

func (r *defaultServerResponse) Flush() error {
	r.writeHeaders()

	if flusher, ok := r.responseWriter.(http.Flusher); ok {
		flusher.Flush()
	}

	return nil
}

func (r *defaultServerResponse) IsCommitted() bool {
	return r.writtenHeaders
}

func (r *defaultServerResponse) Reset() error {
	if r.IsCommitted() {
		return errors.New("already committed")
	}

	r.statusCode = StatusOK
	r.headers = http.Header{}
	return nil
}

func (r *defaultServerResponse) writeHeaders() {
	if r.writtenHeaders {
		return
	}

	if r.characterEncoding != "" {
		ct := r.headers.Get(HeaderContentType)
		if ct != "" && !strings.Contains(ct, "charset") {
			r.headers.Set(HeaderContentType, ct+"; charset="+r.characterEncoding)
		}
	}

	for key, values := range r.headers {
		for _, value := range values {
			r.responseWriter.Header().Add(key, value)
		}
	}

	r.responseWriter.WriteHeader(int(r.statusCode))
	r.writtenHeaders = true
}
