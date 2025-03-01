package http

import (
	"io"
	"net/http"
	"net/url"
)

// Request interface represents an HTTP request.
type Request interface {
	// Context returns the context associated with the request.
	Context() Context

	// Cookie returns the cookie with the specified name.
	Cookie(name string) (*Cookie, bool)
	// Cookies returns all the cookies associated with the request.
	Cookies() []*Cookie

	// QueryParam returns the query parameter with the specified name.
	QueryParam(name string) (string, bool)
	// QueryParamNames returns the names of all query parameters.
	QueryParamNames() []string
	// QueryParams returns all the query parameters with the specified name.
	QueryParams(name string) []string
	// QueryString returns the query string of the request.
	QueryString() string

	// Header returns the header with the specified name.
	Header(name string) (string, bool)
	// HeaderNames returns the names of all headers.
	HeaderNames() []string
	// Headers returns all the headers with the specified name.
	Headers(name string) []string

	// Path returns the path of the request.
	Path() string
	// PathValue returns the value of the path parameter with the specified name.
	PathValue(name string) (string, bool)
	// Method returns the method of the request.
	Method() Method
	// Reader returns the reader of the request body.
	Reader() io.Reader
	// Scheme returns the scheme of the request.
	Scheme() string
	// IsSecure returns whether the request is secure.
	IsSecure() bool
}

// requestWrapper is a wrapper for the Request.
type requestWrapper struct {
	// request is the original request.
	request Request
	// context is the context associated with the request.
	context Context
}

// Context returns the context associated with the request.
func (r requestWrapper) Context() Context {
	return r.context
}

// Cookie returns the cookie with the specified name.
func (r requestWrapper) Cookie(name string) (*Cookie, bool) {
	return r.request.Cookie(name)
}

// Cookies returns all the cookies associated with the request.
func (r requestWrapper) Cookies() []*Cookie {
	return r.request.Cookies()
}

// QueryParam returns the query parameter with the specified name.
func (r requestWrapper) QueryParam(name string) (string, bool) {
	return r.request.QueryParam(name)
}

// QueryParamNames returns the names of all query parameters.
func (r requestWrapper) QueryParamNames() []string {
	return r.request.QueryParamNames()
}

// QueryParams returns all the query parameters with the specified name.
func (r requestWrapper) QueryParams(name string) []string {
	return r.request.QueryParams(name)
}

// QueryString returns the query string of the request.
func (r requestWrapper) QueryString() string {
	return r.request.QueryString()
}

// Header returns the header with the specified name.
func (r requestWrapper) Header(name string) (string, bool) {
	return r.request.Header(name)
}

// HeaderNames returns the names of all headers.
func (r requestWrapper) HeaderNames() []string {
	return r.request.HeaderNames()
}

// Headers returns all the headers with the specified name.
func (r requestWrapper) Headers(name string) []string {
	return r.request.Headers(name)
}

// Path returns the path of the request.
func (r requestWrapper) Path() string {
	return r.request.Path()
}

func (r requestWrapper) PathValue(name string) (string, bool) {
	return r.request.PathValue(name)
}

// Method returns the method of the request.
func (r requestWrapper) Method() Method {
	return r.request.Method()
}

// Reader returns the reader of the request body.
func (r requestWrapper) Reader() io.Reader {
	return r.request.Reader()
}

// Scheme returns the scheme of the request.
func (r requestWrapper) Scheme() string {
	return r.request.Scheme()
}

// IsSecure returns whether the request is secure.
func (r requestWrapper) IsSecure() bool {
	return r.request.IsSecure()
}

// MultiReadRequest is a wrapper for the Request that allows multiple reads.
type MultiReadRequest struct {
	// request is the original request.
	request Request
	// reader is the reader of the request body
	reader io.Reader
}

// NewMultiReadRequest creates a new instance of MultiReadRequest.
// It ensures that the provided request is not nil.
func NewMultiReadRequest(request Request) *MultiReadRequest {
	if request == nil {
		panic("nil request")
	}

	return &MultiReadRequest{
		request: request,
	}
}

// Context returns the context associated with the request.
func (m *MultiReadRequest) Context() Context {
	return m.request.Context()
}

// Cookie returns the cookie with the specified name.
func (m *MultiReadRequest) Cookie(name string) (*Cookie, bool) {
	return m.request.Cookie(name)
}

// Cookies returns all the cookies associated with the request.
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

// QueryString returns the query string of the request.
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

// Path returns the path of the request.
func (m *MultiReadRequest) Path() string {
	return m.request.Path()
}

// PathValue returns the value of the path parameter with the specified name.
func (m *MultiReadRequest) PathValue(name string) (string, bool) {
	return m.request.PathValue(name)
}

// Method returns the method of the request.
func (m *MultiReadRequest) Method() Method {
	return m.request.Method()
}

// Reader returns the reader of the request body.
// If the reader is nil, it creates a new cached reader for the request.
func (m *MultiReadRequest) Reader() io.Reader {
	if m.reader == nil {
		m.reader = newCachedReader(m.request.Reader())
	}

	return m.reader
}

// Scheme returns the scheme of the request.
func (m *MultiReadRequest) Scheme() string {
	return m.request.Scheme()
}

// IsSecure returns whether the request is secure.
func (m *MultiReadRequest) IsSecure() bool {
	return m.request.IsSecure()
}

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

type defaultServerRequest struct {
	req    *http.Request
	ctx    Context
	reader io.Reader

	queryCache   url.Values
	cookiesCache []*Cookie
	pathValues   PathValues
}

func (r *defaultServerRequest) Context() Context {
	return r.ctx
}

func (r *defaultServerRequest) initCookieCache() {
	if r.cookiesCache == nil {
		r.cookiesCache = parseCookies(r.req.Header)
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
		if r.req != nil && r.req.URL != nil {
			r.queryCache = r.req.URL.Query()
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
	//TODO implement me
	panic("implement me")
}

func (r *defaultServerRequest) QueryParams(name string) []string {
	r.initQueryCache()

	queryParams := make([]string, 0, len(r.queryCache))

	for queryParam := range r.queryCache {
		queryParams = append(queryParams, queryParam)
	}

	return queryParams
}

func (r *defaultServerRequest) QueryString() string {
	r.initQueryCache()
	return r.req.URL.RawQuery
}

func (r *defaultServerRequest) Header(name string) (string, bool) {
	values := r.req.Header.Values(name)

	if len(values) != 0 {
		return values[0], true
	}

	return "", false
}

func (r *defaultServerRequest) HeaderNames() []string {
	headers := make([]string, 0, len(r.req.Header))

	for header := range r.req.Header {
		headers = append(headers, header)
	}

	return headers
}

func (r *defaultServerRequest) Headers(name string) []string {
	return r.req.Header.Values(name)
}

func (r *defaultServerRequest) Path() string {
	return r.req.URL.Path
}

func (r *defaultServerRequest) PathValue(name string) (string, bool) {
	return r.pathValues.Value(name)
}

func (r *defaultServerRequest) Method() Method {
	return Method(r.req.Method)
}

func (r *defaultServerRequest) Reader() io.Reader {
	if r.reader != nil {
		return r.reader
	}

	return r.req.Body
}

func (r *defaultServerRequest) Scheme() string {
	return ""
}

func (r *defaultServerRequest) IsSecure() bool {
	return false
}
