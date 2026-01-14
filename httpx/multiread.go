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

package httpx

import "io"

// cachedReader is a reader that caches the data read from it.
// This allows multiple reads from the same reader.
type cachedReader struct {
	// data is the cached data.
	data []byte
	// reader is the original reader.
	reader io.Reader
}

// newCachedReader creates a new instance of cachedReader.
// It ensures that the provided reader is not nil.
func newCachedReader(reader io.Reader) *cachedReader {
	if reader == nil {
		panic("nil reader")
	}

	return &cachedReader{
		data:   nil,
		reader: reader,
	}
}

// readAll reads all data from the original reader and caches it.
// This method is called before each read operation to ensure that the data is loaded.
func (c *cachedReader) readAll() (err error) {
	if c.data == nil {
		c.data, err = io.ReadAll(c.reader)
	}

	return err
}

// Read reads up to len(p) bytes into p from the cached data.
// It returns the number of bytes read and any error encountered.
func (c *cachedReader) Read(p []byte) (n int, err error) {
	err = c.readAll()
	if err != nil {
		return 0, err
	}

	n = copy(p, c.data)
	return
}

// MultiReadRequest is a wrapper for the Request that allows multiple reads.
type MultiReadRequest struct {
	// request is the original request.
	request Request
	// reader is the reader of the request body
	reader io.Reader
}

// NewMultiReadRequest creates a new instance of MultiReadRequest.
// It ensures that the provided defaultRequest is not nil.
func NewMultiReadRequest(request Request) *MultiReadRequest {
	if request == nil {
		panic("nil defaultRequest")
	}

	return &MultiReadRequest{
		request: request,
	}
}

// Context returns the context associated with the defaultRequest.
func (m *MultiReadRequest) Context() Context {
	return m.request.Context()
}

// Cookie returns the cookie with the specified name.
func (m *MultiReadRequest) Cookie(name string) (*Cookie, bool) {
	return m.request.Cookie(name)
}

// Cookies returns all the cookies associated with the defaultRequest.
func (m *MultiReadRequest) Cookies() []*Cookie {
	return m.request.Cookies()
}

// QueryParam returns the query parameter with the specified name.
func (m *MultiReadRequest) QueryParam(name string) (string, bool) {
	return m.request.QueryParam(name)
}

// QueryParamNames returns the names of all query parameters.
func (m *MultiReadRequest) QueryParamNames() []string {
	return m.request.QueryParamNames()
}

// QueryParams returns all the query parameters with the specified name.
func (m *MultiReadRequest) QueryParams(name string) []string {
	return m.request.QueryParams(name)
}

// QueryString returns the query string of the defaultRequest.
func (m *MultiReadRequest) QueryString() string {
	return m.request.QueryString()
}

// Header returns the header with the specified name.
func (m *MultiReadRequest) Header(name string) (string, bool) {
	return m.request.Header(name)
}

// HeaderNames returns the names of all headers.
func (m *MultiReadRequest) HeaderNames() []string {
	return m.request.HeaderNames()
}

// Headers returns all the headers with the specified name.
func (m *MultiReadRequest) Headers(name string) []string {
	return m.request.Headers(name)
}

// Path returns the path of the defaultRequest.
func (m *MultiReadRequest) Path() string {
	return m.request.Path()
}

// PathValue returns the value of the path parameter with the specified name.
func (m *MultiReadRequest) PathValue(name string) (string, bool) {
	return m.request.PathValue(name)
}

// Method returns the method of the defaultRequest.
func (m *MultiReadRequest) Method() Method {
	return m.request.Method()
}

// Reader returns the reader of the defaultRequest body.
// If the reader is nil, it creates a new cached reader for the defaultRequest.
func (m *MultiReadRequest) Reader() io.Reader {
	if m.reader == nil {
		m.reader = newCachedReader(m.request.Reader())
	}

	return m.reader
}

// Scheme returns the scheme of the defaultRequest.
func (m *MultiReadRequest) Scheme() string {
	return m.request.Scheme()
}

// IsSecure returns whether the defaultRequest is secure.
func (m *MultiReadRequest) IsSecure() bool {
	return m.request.IsSecure()
}

// MultiReadResponse is a wrapper for the Response that allows multiple reads.
type MultiReadResponse struct {
	// response is the original response.
	response Response
}

// NewMultiReadResponse creates a new instance of MultiReadResponse.
// It ensures that the provided response is not nil.
func NewMultiReadResponse(response Response) *MultiReadResponse {
	if response == nil {
		panic("nil response")
	}

	return &MultiReadResponse{
		response: response,
	}
}

// Context returns the context associated with the response.
func (m *MultiReadResponse) Context() Context {
	return m.response.Context()
}

// AddCookie adds a cookie to the response.
func (m *MultiReadResponse) AddCookie(cookie *Cookie) {
	m.response.AddCookie(cookie)
}

// ContentLength returns the content length of the response.
func (m *MultiReadResponse) ContentLength() int {
	return m.response.ContentLength()
}

// SetContentLength sets the content length of the response.
func (m *MultiReadResponse) SetContentLength(len int) {
	m.response.SetContentLength(len)
}

// CharacterEncoding returns the character encoding of the response.
func (m *MultiReadResponse) CharacterEncoding() string {
	return m.response.CharacterEncoding()
}

// SetCharacterEncoding sets the character encoding of the response.
func (m *MultiReadResponse) SetCharacterEncoding(charset string) {
	m.response.SetCharacterEncoding(charset)
}

// ContentType returns the content type of the response.
func (m *MultiReadResponse) ContentType() string {
	return m.response.ContentType()
}

// SetContentType sets the content type of the response.
func (m *MultiReadResponse) SetContentType(contentType string) {
	m.response.SetContentType(contentType)
}

// AddHeader adds a header to the response.
// If the header is already set, it will be appended.
func (m *MultiReadResponse) AddHeader(name string, value string) {
	m.response.AddHeader(name, value)
}

// SetHeader sets a header in the response.
// If the header is already set, it will be overwritten.
func (m *MultiReadResponse) SetHeader(name string, value string) {
	m.response.SetHeader(name, value)
}

// DeleteHeader deletes a header from the response.
func (m *MultiReadResponse) DeleteHeader(name string) {
	m.response.DeleteHeader(name)
}

// Header returns the value of a header in the response.
func (m *MultiReadResponse) Header(name string) (string, bool) {
	return m.response.Header(name)
}

// HeaderNames returns the names of all headers in the response.
func (m *MultiReadResponse) HeaderNames() []string {
	return m.HeaderNames()
}

// Headers returns all the values of a header in the response.
func (m *MultiReadResponse) Headers(name string) []string {
	return m.response.Headers(name)
}

// Status returns the status of the response.
func (m *MultiReadResponse) Status() Status {
	return m.response.Status()
}

// SetStatus sets the status of the response.
func (m *MultiReadResponse) SetStatus(status Status) {
	m.response.SetStatus(status)
}

// Redirect redirects the response to a location with a status.
func (m *MultiReadResponse) Redirect(location string, status Status) error {
	return m.response.Redirect(location, status)
}

// Writer returns the writer of the response.
func (m *MultiReadResponse) Writer() io.Writer {
	//TODO implement me
	panic("implement me")
}

// Flush flushes the response.
func (m *MultiReadResponse) Flush() error {
	//TODO implement me
	panic("implement me")
	return nil
}

// IsCommitted checks if the response is committed.
func (m *MultiReadResponse) IsCommitted() bool {
	//TODO implement me
	panic("implement me")
}

// Reset resets the response.
func (m *MultiReadResponse) Reset() {
	//TODO implement me
	panic("implement me")
}

// CopyBodyToResponse copies the cached body content to the response.
func (m *MultiReadResponse) CopyBodyToResponse() error {
	return nil
}
