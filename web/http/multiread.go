package http

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
