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

	stdhttp "net/http"
)

type standardContext struct {
	request  Request
	response Response
	security SecurityContext
	params   map[string]string
	storage  map[string]any
	handlers []Handler
	index    int
}

func newStandardContext(req Request, res Response, sec SecurityContext, params map[string]string, middleware []Middleware, handler Handler) *standardContext {
	handlers := make([]Handler, 0, len(middleware)+1)
	handlers = append(handlers, middleware...)
	handlers = append(handlers, handler)

	return &standardContext{
		request:  req,
		response: res,
		security: sec,
		params:   params,
		storage:  make(map[string]any),
		handlers: handlers,
		index:    -1,
	}
}

func (c *standardContext) Request() Request {
	return c.request
}

func (c *standardContext) Response() Response {
	return c.response
}

func (c *standardContext) Security() SecurityContext {
	return c.security
}

func (c *standardContext) SetSecurity(sec SecurityContext) {
	if sec != nil {
		c.security = sec
	}
}

func (c *standardContext) Next() {
	c.index++
	if c.index < len(c.handlers) {
		c.handlers[c.index](c)
	}
}

func (c *standardContext) Params() map[string]string {
	clone := make(map[string]string)
	for k, v := range c.params {
		clone[k] = v
	}
	return clone
}

func (c *standardContext) Param(name string) string {
	return c.params[name]
}

func (c *standardContext) Set(key string, value any) {
	c.storage[key] = value
}

func (c *standardContext) Get(key string) (any, bool) {
	val, ok := c.storage[key]
	return val, ok
}

// standardRequest adapts net/http.Request to the Request interface.
type standardRequest struct {
	req *stdhttp.Request
}

func (r *standardRequest) Method() string {
	return r.req.Method
}

func (r *standardRequest) Path() string {
	return r.req.URL.Path
}

func (r *standardRequest) Header(name string) string {
	return r.req.Header.Get(name)
}

func (r *standardRequest) Body() io.ReadCloser {
	return r.req.Body
}

// standardResponse adapts http.ResponseWriter to the Response interface.
type standardResponse struct {
	writer      stdhttp.ResponseWriter
	wroteHeader bool
}

func (r *standardResponse) Status(code int) {
	r.writer.WriteHeader(code)
	r.wroteHeader = true
}

func (r *standardResponse) Header(name, value string) {
	r.writer.Header().Set(name, value)
}

func (r *standardResponse) Write(body []byte) (int, error) {
	if !r.wroteHeader {
		r.Status(stdhttp.StatusOK)
	}
	return r.writer.Write(body)
}
