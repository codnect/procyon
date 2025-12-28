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

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// Response interface represents an HTTP response.
type Response interface {
	// Context returns the context associated with the response.
	Context() Context
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
	Reset()
}

type defaultResponse struct {
	responseWriter http.ResponseWriter
	ctx            Context
	writer         io.Writer

	headers        http.Header
	statusCode     Status
	writtenHeaders bool
	writerUsed     bool
}

func (r *defaultResponse) Context() Context {
	return r.ctx
}

func (r *defaultResponse) AddCookie(cookie *Cookie) {
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
		SameSite: http.SameSite(cookie.SameSite),
	}

	if v := stdCookie.String(); v != "" {
		r.headers.Add("Set-Cookie", v)
	}
}

func (r *defaultResponse) ContentLength() int {
	length := r.headers.Get("Content-Length")

	if length != "" {
		val, err := strconv.Atoi(length)
		if err != nil {
			return 0
		}

		return val
	}

	return 0
}

func (r *defaultResponse) SetContentLength(len int) {
	if r.writtenHeaders {
		return
	}

	r.headers.Add(HeaderContentLength, strconv.Itoa(len))
}

func (r *defaultResponse) CharacterEncoding() string {
	return ""
}

func (r *defaultResponse) SetCharacterEncoding(charset string) {

}

func (r *defaultResponse) ContentType() string {
	return r.headers.Get(HeaderContentType)
}

func (r *defaultResponse) SetContentType(contentType string) {
	if r.writtenHeaders {
		return
	}

	r.headers.Add(HeaderContentType, contentType)
}

func (r *defaultResponse) AddHeader(name string, value string) {
	if r.writtenHeaders {
		return
	}

	r.headers.Add(name, value)
}

func (r *defaultResponse) SetHeader(name string, value string) {
	if r.writtenHeaders {
		return
	}

	r.headers.Set(name, value)
}

func (r *defaultResponse) DeleteHeader(name string) {
	if r.writtenHeaders {
		return
	}

	r.headers.Del(name)
}

func (r *defaultResponse) Header(name string) (string, bool) {
	values := r.headers.Values(name)

	if len(values) != 0 {
		return values[0], true
	}

	return "", false
}

func (r *defaultResponse) HeaderNames() []string {
	headers := make([]string, 0, len(r.headers))

	for header := range r.headers {
		headers = append(headers, header)
	}

	return headers
}

func (r *defaultResponse) Headers(name string) []string {
	return r.headers.Values(name)
}

func (r *defaultResponse) Status() Status {
	return r.statusCode
}

func (r *defaultResponse) SetStatus(status Status) {
	if r.writtenHeaders {
		return
	}

	r.statusCode = status
}

func (r *defaultResponse) Redirect(location string, status Status) error {
	return nil
}

func (r *defaultResponse) Writer() io.Writer {
	r.writerUsed = true

	if r.writer != nil {
		return r.writer
	}

	return r.responseWriter
}

func (r *defaultResponse) Flush() error {
	r.writeHeaders()
	if !r.writerUsed {
		return nil
	}

	if r.writer == nil {
		r.responseWriter.WriteHeader(int(r.statusCode))
	}

	// flush data
	return nil
}

func (r *defaultResponse) IsCommitted() bool {
	return false
}

func (r *defaultResponse) Reset() {
	r.headers = http.Header{}
}

func (r *defaultResponse) writeHeaders() {
	if !r.writtenHeaders {
		for key, values := range r.headers {
			if len(values) == 1 {
				r.responseWriter.Header().Set(key, values[0])
			} else {
				for _, value := range values {
					r.responseWriter.Header().Add(key, value)
				}
			}
		}
		r.writtenHeaders = true
	}
}

// ResponseWrapper is a wrapper for the Response.
type ResponseWrapper struct {
	// ResponseWrapper is a wrapper for the Response interface.
	response Response
	// context is the context associated with the response.
	context Context
}

// Context returns the context associated with the response.
func (r ResponseWrapper) Context() Context {
	return r.context
}

// AddCookie adds a cookie to the response.
func (r ResponseWrapper) AddCookie(cookie *Cookie) {
	r.response.AddCookie(cookie)
}

// ContentLength returns the content length of the response.
func (r ResponseWrapper) ContentLength() int {
	return r.response.ContentLength()
}

// SetContentLength sets the content length of the response.
func (r ResponseWrapper) SetContentLength(len int) {
	r.response.SetContentLength(len)
}

// CharacterEncoding returns the character encoding of the response.
func (r ResponseWrapper) CharacterEncoding() string {
	return r.response.CharacterEncoding()
}

// SetCharacterEncoding sets the character encoding of the response.
func (r ResponseWrapper) SetCharacterEncoding(charset string) {
	r.response.SetCharacterEncoding(charset)
}

// ContentType returns the content type of the response.
func (r ResponseWrapper) ContentType() string {
	return r.response.ContentType()
}

// SetContentType sets the content type of the response.
func (r ResponseWrapper) SetContentType(contentType string) {
	r.response.SetContentType(contentType)
}

// AddHeader adds a header to the response.
// If the header is already set, it will be appended.
func (r ResponseWrapper) AddHeader(name string, value string) {
	r.response.AddHeader(name, value)
}

// SetHeader sets a header in the response.
// If the header is already set, it will be overwritten.
func (r ResponseWrapper) SetHeader(name string, value string) {
	r.response.SetHeader(name, value)
}

// DeleteHeader deletes a header from the response.
func (r ResponseWrapper) DeleteHeader(name string) {
	r.response.DeleteHeader(name)
}

// Header returns the value of a header in the response.
func (r ResponseWrapper) Header(name string) (string, bool) {
	return r.response.Header(name)
}

// HeaderNames returns the names of all headers in the response.
func (r ResponseWrapper) HeaderNames() []string {
	return r.response.HeaderNames()
}

// Headers returns all the values of a header in the response.
func (r ResponseWrapper) Headers(name string) []string {
	return r.response.Headers(name)
}

// Status returns the status of the response.
func (r ResponseWrapper) Status() Status {
	return r.response.Status()
}

// SetStatus sets the status of the response.
func (r ResponseWrapper) SetStatus(status Status) {
	r.response.SetStatus(status)
}

// Redirect redirects the response to a location with a status.
func (r ResponseWrapper) Redirect(location string, status Status) error {
	return r.response.Redirect(location, status)
}

// Writer returns the writer of the response.
func (r ResponseWrapper) Writer() io.Writer {
	return r.response.Writer()
}

// Flush flushes the response.
func (r ResponseWrapper) Flush() error {
	return r.response.Flush()
}

// IsCommitted checks if the response is committed.
func (r ResponseWrapper) IsCommitted() bool {
	return r.response.IsCommitted()
}

// Reset resets the response.
func (r ResponseWrapper) Reset() {
	r.response.Reset()
}
