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

type ServerResponse struct {
	responseWriter http.ResponseWriter
	ctx            Context
	writer         io.Writer

	headers        http.Header
	statusCode     Status
	writtenHeaders bool
	writerUsed     bool
}

func (r *ServerResponse) Context() Context {
	return r.ctx
}

func (r *ServerResponse) AddCookie(cookie *Cookie) {
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

func (r *ServerResponse) ContentLength() int {
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

func (r *ServerResponse) SetContentLength(len int) {
	if r.writtenHeaders {
		return
	}

	r.headers.Add(HeaderContentLength, strconv.Itoa(len))
}

func (r *ServerResponse) CharacterEncoding() string {
	return ""
}

func (r *ServerResponse) SetCharacterEncoding(charset string) {

}

func (r *ServerResponse) ContentType() string {
	return r.headers.Get(HeaderContentType)
}

func (r *ServerResponse) SetContentType(contentType string) {
	if r.writtenHeaders {
		return
	}

	r.headers.Add(HeaderContentType, contentType)
}

func (r *ServerResponse) AddHeader(name string, value string) {
	if r.writtenHeaders {
		return
	}

	r.headers.Add(name, value)
}

func (r *ServerResponse) SetHeader(name string, value string) {
	if r.writtenHeaders {
		return
	}

	r.headers.Set(name, value)
}

func (r *ServerResponse) DeleteHeader(name string) {
	if r.writtenHeaders {
		return
	}

	r.headers.Del(name)
}

func (r *ServerResponse) Header(name string) (string, bool) {
	values := r.headers.Values(name)

	if len(values) != 0 {
		return values[0], true
	}

	return "", false
}

func (r *ServerResponse) HeaderNames() []string {
	headers := make([]string, 0, len(r.headers))

	for header := range r.headers {
		headers = append(headers, header)
	}

	return headers
}

func (r *ServerResponse) Headers(name string) []string {
	return r.headers.Values(name)
}

func (r *ServerResponse) Status() Status {
	return r.statusCode
}

func (r *ServerResponse) SetStatus(status Status) {
	if r.writtenHeaders {
		return
	}

	r.statusCode = status
}

func (r *ServerResponse) Redirect(location string, status Status) error {
	return nil
}

func (r *ServerResponse) Writer() io.Writer {
	r.writerUsed = true

	if r.writer != nil {
		return r.writer
	}

	return r.responseWriter
}

func (r *ServerResponse) Flush() error {
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

func (r *ServerResponse) IsCommitted() bool {
	return false
}

func (r *ServerResponse) Reset() {
	r.headers = http.Header{}
}

func (r *ServerResponse) writeHeaders() {
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
