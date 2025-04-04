package http

import (
	"io"
	"net/http"
	"net/url"
)

const (
	// PathValuesContextKey is the key for path values in the context
	PathValuesContextKey = "PathValues"
)

// PathValues represents the values of the path parameters.
type PathValues map[string]string

// Put adds a new path parameter with the provided name and value.
func (p PathValues) Put(name string, value string) {
	p[name] = value
}

// Value returns the value of the path parameter with the provided name.
func (p PathValues) Value(name string) (string, bool) {
	if val, ok := p[name]; ok {
		return val, true
	}

	return "", false
}

// Clear removes all path values.
func (p PathValues) Clear() {
	clear(p)
}

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

type RequestDelegate interface {
	Invoke(ctx Context)
}

type ServerRequestDelegate struct {
	ctx *ServerContext
}

func (d ServerRequestDelegate) Invoke(ctx Context) {
	d.ctx.Invoke(ctx)
}

/*
type contextDelegate struct {
	ctx *Context
}

func (cd *contextDelegate) Invoke(ctx *Context) {
	if cd.ctx.IsCompleted() || cd.ctx.IsAborted() || len(cd.ctx.HandlerChain.functions) <= cd.ctx.handlerIndex {
		return
	}

	next := cd.ctx.HandlerChain.functions[cd.ctx.handlerIndex]
	cd.ctx.handlerIndex++

	err := next(ctx, cd)

	if err != nil {
		cd.ctx.err = err
	}

	if cd.ctx.IsCompleted() || cd.ctx.IsAborted() {
		return
	}

	if cd.ctx.handlerIndex != len(cd.ctx.HandlerChain.functions) {
		cd.ctx.Abort()
	}
}
*/

// RequestWrapper is a wrapper for the Request.
type RequestWrapper struct {
	// request is the original request.
	request Request
	// context is the context associated with the request.
	context Context
}

// Context returns the context associated with the request.
func (r RequestWrapper) Context() Context {
	return r.context
}

// Cookie returns the cookie with the specified name.
func (r RequestWrapper) Cookie(name string) (*Cookie, bool) {
	return r.request.Cookie(name)
}

// Cookies returns all the cookies associated with the request.
func (r RequestWrapper) Cookies() []*Cookie {
	return r.request.Cookies()
}

// QueryParam returns the query parameter with the specified name.
func (r RequestWrapper) QueryParam(name string) (string, bool) {
	return r.request.QueryParam(name)
}

// QueryParamNames returns the names of all query parameters.
func (r RequestWrapper) QueryParamNames() []string {
	return r.request.QueryParamNames()
}

// QueryParams returns all the query parameters with the specified name.
func (r RequestWrapper) QueryParams(name string) []string {
	return r.request.QueryParams(name)
}

// QueryString returns the query string of the request.
func (r RequestWrapper) QueryString() string {
	return r.request.QueryString()
}

// Header returns the header with the specified name.
func (r RequestWrapper) Header(name string) (string, bool) {
	return r.request.Header(name)
}

// HeaderNames returns the names of all headers.
func (r RequestWrapper) HeaderNames() []string {
	return r.request.HeaderNames()
}

// Headers returns all the headers with the specified name.
func (r RequestWrapper) Headers(name string) []string {
	return r.request.Headers(name)
}

// Path returns the path of the request.
func (r RequestWrapper) Path() string {
	return r.request.Path()
}

func (r RequestWrapper) PathValue(name string) (string, bool) {
	return r.request.PathValue(name)
}

// Method returns the method of the request.
func (r RequestWrapper) Method() Method {
	return r.request.Method()
}

// Reader returns the reader of the request body.
func (r RequestWrapper) Reader() io.Reader {
	return r.request.Reader()
}

// Scheme returns the scheme of the request.
func (r RequestWrapper) Scheme() string {
	return r.request.Scheme()
}

// IsSecure returns whether the request is secure.
func (r RequestWrapper) IsSecure() bool {
	return r.request.IsSecure()
}

type ServerRequest struct {
	req    *http.Request
	ctx    Context
	reader io.Reader

	queryCache   url.Values
	cookiesCache []*Cookie
	pathValues   PathValues
}

func (r *ServerRequest) Context() Context {
	return r.ctx
}

func (r *ServerRequest) initCookieCache() {
	if r.cookiesCache == nil {
		r.cookiesCache = r.req.Cookies()
	}
}

func (r *ServerRequest) Cookie(name string) (*Cookie, bool) {
	r.initCookieCache()

	for _, cookie := range r.cookiesCache {
		if cookie.Name == name {
			return cookie, true
		}
	}

	return nil, false
}

func (r *ServerRequest) Cookies() []*Cookie {
	r.initCookieCache()
	return r.cookiesCache
}

func (r *ServerRequest) initQueryCache() {
	if r.queryCache == nil {
		if r.req != nil && r.req.URL != nil {
			r.queryCache = r.req.URL.Query()
		}
	}
}

func (r *ServerRequest) QueryParam(name string) (string, bool) {
	r.initQueryCache()

	values, ok := r.queryCache[name]
	if ok {
		return values[0], true
	}

	return "", false
}

func (r *ServerRequest) QueryParamNames() []string {
	//TODO implement me
	panic("implement me")
}

func (r *ServerRequest) QueryParams(name string) []string {
	r.initQueryCache()

	queryParams := make([]string, 0, len(r.queryCache))

	for queryParam := range r.queryCache {
		queryParams = append(queryParams, queryParam)
	}

	return queryParams
}

func (r *ServerRequest) QueryString() string {
	r.initQueryCache()
	return r.req.URL.RawQuery
}

func (r *ServerRequest) Header(name string) (string, bool) {
	values := r.req.Header.Values(name)

	if len(values) != 0 {
		return values[0], true
	}

	return "", false
}

func (r *ServerRequest) HeaderNames() []string {
	headers := make([]string, 0, len(r.req.Header))

	for header := range r.req.Header {
		headers = append(headers, header)
	}

	return headers
}

func (r *ServerRequest) Headers(name string) []string {
	return r.req.Header.Values(name)
}

func (r *ServerRequest) Path() string {
	return r.req.URL.Path
}

func (r *ServerRequest) PathValue(name string) (string, bool) {
	return r.pathValues.Value(name)
}

func (r *ServerRequest) Method() Method {
	return Method(r.req.Method)
}

func (r *ServerRequest) Reader() io.Reader {
	if r.reader != nil {
		return r.reader
	}

	return r.req.Body
}

func (r *ServerRequest) Scheme() string {
	return ""
}

func (r *ServerRequest) IsSecure() bool {
	return false
}
