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
	"context"
	"net/http"
	"time"
)

// Context is an interface that extends the standard context.Context interface in Go.
// It provides additional methods that are specific to handling HTTP requests and responses.
type Context interface {
	context.Context
	// Endpoint returns the endpoint associated with the context.
	Endpoint() *Endpoint
	// IsCompleted checks if the HTTP defaultRequest has been completed.
	IsCompleted() bool
	// Abort aborts the HTTP defaultRequest.
	Abort()
	// IsAborted checks if the HTTP defaultRequest has been aborted.
	IsAborted() bool
	// Request returns the HTTP defaultRequest associated with the context.
	Request() Request
	// Response returns the HTTP response associated with the context.
	Response() Response
}

// defaultContext is the default implementation of the Context interface.
type defaultContext struct {
	request  defaultRequest
	response defaultResponse

	middlewares      []MiddlewareFunc
	nextHandlerIndex int

	err       error
	completed bool
	aborted   bool

	delegate defaultRequestDelegate
}

func (c *defaultContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *defaultContext) Done() <-chan struct{} {
	return nil
}

func (c *defaultContext) Err() error {
	return nil
}

// Value returns the value associated with this context for key, or nil if no value is associated with key.
// Successive calls to Value with the same key returns the same result.
func (c *defaultContext) Value(key any) any {
	if key == PathValuesContextKey {
		return c.request.pathValues
	}

	return nil
}

// IsCompleted checks if the HTTP defaultRequest has been completed.
func (c *defaultContext) IsCompleted() bool {
	return c.completed
}

// Abort aborts the HTTP defaultRequest.
func (c *defaultContext) Abort() {
	c.aborted = true
}

// IsAborted checks if the HTTP defaultRequest has been aborted.
func (c *defaultContext) IsAborted() bool {
	return c.aborted
}

// Request returns the HTTP defaultRequest associated with the context.
func (c *defaultContext) Request() Request {
	return &c.request
}

// Response returns the HTTP response associated with the context.
func (c *defaultContext) Response() Response {
	return &c.response
}

// Invoke invokes the handler chain.
func (c *defaultContext) Invoke(ctx Context) {
	if len(c.middlewares) == 0 {
		return
	}

	nextHandler := c.nextHandlerIndex

	if c.completed || c.aborted || len(c.middlewares) <= nextHandler {
		return
	}

	next := c.middlewares[nextHandler]
	nextHandler++
	err := next(ctx, c.delegate)

	if err != nil {
		c.err = err
	}

	if c.completed || c.aborted {
		return
	}

	if nextHandler != len(c.middlewares) {
		c.aborted = true
	} else {
		c.completed = true
	}
}

// reset resets the context with the specified writer and defaultRequest.
func (c *defaultContext) reset(writer http.ResponseWriter, request *http.Request) {
	c.request.req = request
	c.delegate.ctx = c

	c.err = nil
	c.completed = false
	c.aborted = false

	c.nextHandlerIndex = 0
}

// contextWrapper is an implementation of the Context interface.
// It wraps a parent Context, a Request, and a Response.
type contextWrapper struct {
	parent          Context
	requestWrapper  RequestWrapper
	responseWrapper ResponseWrapper
	key             any
	value           any
}

// NewContext creates a new instance of Context.
func NewContext(request Request, response Response) Context {
	if request == nil {
		panic("nil defaultRequest")
	}

	if response == nil {
		panic("nil response")
	}

	wrapper := &contextWrapper{
		requestWrapper: RequestWrapper{
			request: request,
		},
		responseWrapper: ResponseWrapper{
			response: response,
		},
	}

	wrapper.requestWrapper.context = wrapper
	wrapper.responseWrapper.context = wrapper
	return wrapper
}

// ContextWithValue creates a new instance of Context with a value.
func ContextWithValue(parent Context, key, val any) Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}

	if key == nil {
		panic("nil key")
	}

	wrapper := &contextWrapper{
		parent: parent,
		requestWrapper: RequestWrapper{
			request: parent.Request(),
		},
		responseWrapper: ResponseWrapper{
			response: parent.Response(),
		},
		key:   key,
		value: val,
	}

	wrapper.requestWrapper.context = wrapper
	wrapper.responseWrapper.context = wrapper
	return wrapper
}

// NewContextWithRequest creates a new instance of Context with a defaultRequest.
func NewContextWithRequest(parent Context, request Request) Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}

	if request == nil {
		panic("nil defaultRequest")
	}

	wrapper := &contextWrapper{
		parent: parent,
		requestWrapper: RequestWrapper{
			request: request,
		},
		responseWrapper: ResponseWrapper{
			response: parent.Response(),
		},
	}

	wrapper.requestWrapper.context = wrapper
	wrapper.responseWrapper.context = wrapper
	return wrapper
}

// NewContextWithResponse creates a new instance of Context with a response.
func NewContextWithResponse(parent Context, response Response) Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}

	if response == nil {
		panic("nil response")
	}

	wrapper := &contextWrapper{
		parent: parent,
		requestWrapper: RequestWrapper{
			request: parent.Request(),
		},
		responseWrapper: ResponseWrapper{
			response: response,
		},
	}

	wrapper.requestWrapper.context = wrapper
	wrapper.responseWrapper.context = wrapper
	return wrapper
}

// Deadline returns the time when work done on behalf of this context should be canceled.
func (c *contextWrapper) Deadline() (deadline time.Time, ok bool) {
	return c.parent.Deadline()
}

// Done returns a channel that's closed when work done on behalf of this context should be canceled.
// Done may return nil if this context can never be canceled.
func (c *contextWrapper) Done() <-chan struct{} {
	return c.parent.Done()
}

// Err returns a non-nil error value after Done is closed. Err returns Canceled if the context was canceled
// or DeadlineExceeded if the context's deadline passed.
func (c *contextWrapper) Err() error {
	return c.parent.Err()
}

// Value returns the value associated with this context for key, or nil if no value is associated with key.
// Successive calls to Value with the same key returns the same result.
func (c *contextWrapper) Value(key any) any {
	if c.key != nil && c.key == key {
		return c.value
	}

	return c.parent.Value(key)
}

// IsCompleted checks if the HTTP defaultRequest has been completed.
func (c *contextWrapper) IsCompleted() bool {
	return c.parent.IsCompleted()
}

// Abort aborts the HTTP defaultRequest.
func (c *contextWrapper) Abort() {
	c.parent.Abort()
}

// IsAborted checks if the HTTP defaultRequest has been aborted.
func (c *contextWrapper) IsAborted() bool {
	return c.parent.IsAborted()
}

// Request returns the HTTP defaultRequest associated with the context.
func (c *contextWrapper) Request() Request {
	return c.requestWrapper
}

// Response returns the HTTP response associated with the context.
func (c *contextWrapper) Response() Response {
	return c.responseWrapper
}
