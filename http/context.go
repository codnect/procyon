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
	"context"
	"errors"
	"net/http"
	"time"
)

type serverContext interface {
	context.Context
	SetValue(key, value any)
	Request() *ServerRequest
	Response() *ServerResponse
}

// Context represents the context for an HTTP request and response.
type Context struct {
	req *ServerRequest
	res *ServerResponse

	values map[any]any
	err    error
}

// Deadline method returns the time when work done on behalf of
// this context should be canceled.
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.req.nativeReq.Context().Deadline()
}

// Done method returns a channel that's closed when work done on behalf of
// this context should be canceled.
func (c *Context) Done() <-chan struct{} {
	return c.req.nativeReq.Context().Done()
}

// Err returns the error associated with the context.
func (c *Context) Err() error {
	return errors.Join(c.req.nativeReq.Context().Err(), c.err)
}

// Value returns the value associated with the key in the context.
func (c *Context) Value(key any) any {
	val, ok := c.values[key]
	if !ok {
		return c.req.nativeReq.Context().Value(key)
	}

	return val
}

// SetValue sets a value in the context.
func (c *Context) SetValue(key, value any) {
	c.values[key] = value
}

// Request returns the server request associated with the context.
func (c *Context) Request() *ServerRequest {
	return c.req
}

// Response returns the server response associated with the context.
func (c *Context) Response() *ServerResponse {
	return c.res
}

// reset clears the context state and assigns a new HTTP request and response writer.
func (c *Context) reset(w http.ResponseWriter, r *http.Request) {
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

// EndpointContext represents a typed context for an HTTP endpoint handler.
// It wraps the base Context and provides access to parsed input data
// from path parameters, query strings, and request body.
type EndpointContext[I any] struct {
	ctx   *Context
	input I
}

// Deadline returns the time when work done on behalf of
// this context should be canceled.
func (e *EndpointContext[I]) Deadline() (deadline time.Time, ok bool) {
	return e.ctx.Deadline()
}

// Done returns a channel that's closed when work done on behalf of
// this context should be canceled.
func (e *EndpointContext[I]) Done() <-chan struct{} {
	return e.ctx.Done()
}

// Err returns the error associated with the context.
func (e *EndpointContext[I]) Err() error {
	return e.ctx.Err()
}

// Value returns the value associated with the key in the context.
func (e *EndpointContext[I]) Value(key any) any {
	return e.ctx.Value(key)
}

// SetValue sets a value in the context.
func (e *EndpointContext[I]) SetValue(key, value any) {
	e.ctx.SetValue(key, value)
}

// Request returns the server request associated with the context.
func (e *EndpointContext[I]) Request() *ServerRequest {
	return e.ctx.Request()
}

// Response returns the server response associated with the context.
func (e *EndpointContext[I]) Response() *ServerResponse {
	return e.ctx.Response()
}

// Input returns the parsed input data bound from path parameters,
// query strings, headers, and request body.
func (e *EndpointContext[I]) Input() I {
	return e.input
}

// NativeContext returns the underlying base Context.
func (e *EndpointContext[I]) NativeContext() *Context {
	return e.ctx
}

// setContext sets the underlying base Context.
func (e *EndpointContext[I]) setContext(ctx *Context) {
	e.ctx = ctx
}
