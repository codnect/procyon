package web

import (
	"net/http"
	"procyon-test/web/mediatype"
	"time"
)

type Context struct {
	err    error
	values map[any]any

	request        *http.Request
	responseWriter http.ResponseWriter

	queryParamBinder   *ValueBinder
	pathVariableBinder *ValueBinder

	response  Response
	completed bool
}

func (c *Context) writeResponse() {
	if !c.completed {
		for key, values := range c.response.headers {

			if len(values) == 1 {
				c.responseWriter.Header().Add(key, values[0])
			} else {
				for _, val := range values {
					c.responseWriter.Header().Add(key, val)
				}
			}
		}

		c.responseWriter.WriteHeader(int(c.response.status))
		c.responseWriter.Write(nil)
	}
}

func (c *Context) reset(writer http.ResponseWriter, request *http.Request) {
	c.responseWriter = writer
	c.request = request
	c.completed = false

	for k := range c.values {
		delete(c.values, k)
	}

	c.response.viewName = ""
	c.response.entity = nil
	c.response.status = StatusOK
	c.response.contentType = ""

	for k := range c.response.headers {
		delete(c.response.headers, k)
	}
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *Context) Done() <-chan struct{} {
	return nil
}

func (c *Context) Err() error {
	return c.err
}

func (c *Context) Value(key any) any {
	return c.values[key]
}

func (c *Context) Put(key, value any) {
	c.values[key] = value
}

func (c *Context) Request() *http.Request {
	return c.request
}

func (c *Context) ResponseWriter() http.ResponseWriter {
	return c.responseWriter
}

func (c *Context) Path() string {
	return c.request.URL.Path
}

func (c *Context) Method() HttpMethod {
	return HttpMethod(c.request.Method)
}

func (c *Context) Bind(dest any) error {
	return nil
}

func (c *Context) QueryParameterBinder() *ValueBinder {
	return c.queryParamBinder
}

func (c *Context) PathVariableBinder() *ValueBinder {
	return c.pathVariableBinder
}

func (c *Context) Response() *Response {
	return &c.response
}

func (c *Context) SetViewName(name string) {
	c.response.viewName = name
}

func (c *Context) SetEntity(entity any) {
	c.response.entity = entity
}

func (c *Context) SetStatus(status HttpStatus) {
	c.response.status = status
}

func (c *Context) SetContentType(contentType mediatype.MediaType) {
	c.response.contentType = contentType
}
