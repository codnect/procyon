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

import "context"

type handlerContext[T any] interface {
	*T
	context.Context

	SetValue(key, value any)
	Request() *ServerRequest
	Response() *ServerResponse
	private()
}

// Handler represents an HTTP request handler that processes
// an incoming request and returns a Result or an error.
type Handler interface {
	Handle(ctx *Context) (Result, error)
}

// HandlerFunc is an adapter type that allows a plain function
// with the signature func(*Context) (Result, error) to act as a Handler.
type HandlerFunc func(ctx *Context) (Result, error)

// Handle executes the underlying handler function.
func (f HandlerFunc) Handle(ctx *Context) (Result, error) {
	return f(ctx)
}

// handlerAdapter wraps a function that takes a typed context and returns only an error.
// Used internally when the handler doesn't need to return a Result.
type handlerAdapter[T any, C handlerContext[T]] struct {
	fn func(C) error
}

// Handle creates a typed context, binds the base Context to it,
// and executes the handler function.
func (h *handlerAdapter[T, C]) Handle(ctx *Context) (Result, error) {
	endpointCtx := C(new(T))
	if provider, ok := any(endpointCtx).(interface{ setContext(*Context) }); ok {
		provider.setContext(ctx)
	}
	return nil, h.fn(endpointCtx)
}

// handlerResultAdapter wraps a function that takes a typed context and returns a Result.
// Used internally for handlers that return structured responses.
type handlerResultAdapter[T any, C handlerContext[T], R Result] struct {
	fn func(C) (R, error)
}

// Handle creates a typed context, binds the base Context to it,
// and executes the handler function, returning the Result.
func (h *handlerResultAdapter[T, C, R]) Handle(ctx *Context) (Result, error) {
	endpointCtx := C(new(T))
	if provider, ok := any(endpointCtx).(interface{ setContext(*Context) }); ok {
		provider.setContext(ctx)
	}
	return h.fn(endpointCtx)
}

// Handle creates a Handler from a function that returns only an error.
// Type parameters are inferred from the function signature.
func Handle[T any, C handlerContext[T]](fn func(C) error) Handler {

	if _, ok := any((*T)(nil)).(*Context); ok {
		return HandlerFunc(func(ctx *Context) (Result, error) {
			return nil, fn(any(ctx).(C))
		})
	}
	return &handlerAdapter[T, C]{fn: fn}
}

// HandleResult creates a Handler from a function that returns a Result.
// Type parameters are inferred from the function signature.
func HandleResult[T any, C handlerContext[T], R Result](fn func(C) (R, error)) Handler {
	if _, ok := any((*T)(nil)).(*Context); ok {
		return HandlerFunc(func(ctx *Context) (Result, error) {
			return fn(any(ctx).(C))
		})
	}

	return &handlerResultAdapter[T, C, R]{fn: fn}
}
