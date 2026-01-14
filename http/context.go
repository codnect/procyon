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

type Context struct {
	parent   *Context
	request  *defaultServerRequest
	response *defaultServerResponse

	nextHandlerIndex int

	err       error
	completed bool
	aborted   bool

	values map[any]any
}

func NewContext(request ServerRequest, response ServerResponse) *Context {
	if request == nil {
		panic("nil request")
	}

	if response == nil {
		panic("nil response")
	}

	return &Context{}
}

func ContextWithRequest(parent *Context, request ServerRequest) *Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}

	if request == nil {
		panic("nil request")
	}

	return &Context{
		parent: parent,
	}
}

func ContextWithResponse(parent *Context, response ServerResponse) *Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}

	if response == nil {
		panic("nil response")
	}

	return &Context{
		parent: parent,
	}
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	if c.request == nil {
		return
	}

	return c.request.nativeReq.Context().Deadline()
}

func (c *Context) Done() <-chan struct{} {
	if c.request == nil {
		return nil
	}

	return c.request.nativeReq.Context().Done()
}

func (c *Context) Err() error {
	ctxErr := c.request.nativeReq.Context().Err()
	if c.err != nil && ctxErr != nil {
		return errors.Join(c.err, ctxErr)
	}

	if c.err != nil {
		return c.err
	}

	return ctxErr
}

func (c *Context) Value(key any) any {
	val, ok := c.values[key]
	if !ok {
		return c.request.nativeReq.Context().Value(key)
	}

	return val
}

// Endpoint returns the endpoint associated with the context.
func (c *Context) Endpoint() *Endpoint {
	return nil
}

// IsCompleted checks if the HTTP defaultServerRequest has been completed.
func (c *Context) IsCompleted() bool {
	if c.parent != nil {
		return c.parent.IsCompleted()
	}

	return c.completed
}

// Abort aborts the HTTP defaultServerRequest.
func (c *Context) Abort() {
	if c.parent != nil {
		c.parent.Abort()
		return
	}

	c.aborted = true
}

// IsAborted checks if the HTTP defaultServerRequest has been aborted.
func (c *Context) IsAborted() bool {
	if c.parent == nil {
		return c.aborted
	}

	return c.parent.IsAborted()
}

// Request returns the HTTP defaultServerRequest associated with the context.
func (c *Context) Request() ServerRequest {
	if c.parent != nil && c.request == nil {
		return c.parent.request
	}

	return c.request
}

// Response returns the HTTP response associated with the context.
func (c *Context) Response() ServerResponse {
	if c.parent != nil && c.response == nil {
		return c.parent.response
	}

	return c.response
}

// reset resets the context with the specified writer and defaultRequest.
func (c *Context) reset(writer http.ResponseWriter, request *http.Request) {
	c.request.nativeReq = request
	//c.delegate.ctx = c

	c.err = nil
	c.completed = false
	c.aborted = false

	c.nextHandlerIndex = 0
}
