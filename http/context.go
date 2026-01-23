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
	"net/http"
	"time"
)

// Context represents the context for an HTTP request and response.
type Context struct {
	parent *Context
	req    *ServerRequest
	res    *ServerResponse

	values map[any]any
	err    error
}

// Deadline method returns the time when work done on behalf of
// this context should be canceled.
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	if c.parent != nil {
		return c.parent.Deadline()
	}

	return c.req.nativeReq.Context().Deadline()
}

// Done method returns a channel that's closed when work done on behalf of
// this context should be canceled.
func (c *Context) Done() <-chan struct{} {
	if c.parent != nil {
		return c.parent.Done()
	}

	return c.req.nativeReq.Context().Done()
}

// Err returns the error associated with the context.
func (c *Context) Err() error {
	if c.parent != nil {
		return c.parent.Err()
	}

	return errors.Join(c.req.nativeReq.Context().Err(), c.err)
}

// Value returns the value associated with the key in the context.
func (c *Context) Value(key any) any {
	if c.parent != nil {
		return c.parent.Value(key)
	}

	val, ok := c.values[key]
	if !ok {
		return c.req.nativeReq.Context().Value(key)
	}

	return val
}

// SetValue sets a value in the context.
func (c *Context) SetValue(key, value any) {
	if c.parent != nil {
		c.parent.SetValue(key, value)
		return
	}

	c.values[key] = value
}

// Request returns the server request associated with the context.
func (c *Context) Request() *ServerRequest {
	if c.parent != nil {
		return c.parent.Request()
	}

	return c.req
}

// Response returns the server response associated with the context.
func (c *Context) Response() *ServerResponse {
	if c.parent != nil {
		return c.parent.Response()
	}

	return c.res
}

// reset clears the context state and assigns a new HTTP request and response writer.
func (c *Context) reset(w http.ResponseWriter, r *http.Request) {
	c.parent = nil
	c.err = nil
	clear(c.values)

	c.req.nativeReq = r
	c.req.cookiesCache = nil
	c.req.queryCache = nil

	c.res.writer = w
	c.res.status = StatusOK
	c.res.writtenHeaders = false
	c.res.writerUsed = false
	if c.res.headers == nil {
		c.res.headers = make(http.Header)
	} else {
		clear(c.res.headers)
	}
}
