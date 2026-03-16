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

// Dispatcher interface represents a dispatcher that can process
// an HTTP request contained in the Context.
type Dispatcher interface {
	Dispatch(ctx *Context) error
}

type RequestDispatcher struct {
	delegate RequestDelegate
}

// NewRequestDispatcher creates a new dispatcher by building
// a middleware pipeline around the given EndpointMatcher.
//
// The pipeline is always structured as:
//
//	routing → [user middlewares] → endpoint
//
// User middlewares run after routing, so they can inspect the
// matched endpoint via ctx.Endpoint() before it executes.
func NewRequestDispatcher(endpointMatcher EndpointMatcher, middlewares ...Middleware) Dispatcher {
	pipeline := buildPipeline(endpointMatcher, middlewares...)
	return &RequestDispatcher{
		delegate: pipeline,
	}
}

// Dispatch executes the built pipeline for the given request context.
func (d *RequestDispatcher) Dispatch(ctx *Context) error {
	return d.delegate(ctx)
}

// buildPipeline chains all middleware into a single RequestDelegate.
// routing is always first, endpoint is always last, user middlewares
// are sandwiched in between.
func buildPipeline(endpointMatcher EndpointMatcher, userMiddlewares ...Middleware) RequestDelegate {
	middlewares := make([]Middleware, 0, len(userMiddlewares)+2)
	middlewares = append(middlewares, newRoutingMiddleware(endpointMatcher))
	middlewares = append(middlewares, userMiddlewares...)
	middlewares = append(middlewares, newEndpointMiddleware())

	return buildChain(middlewares)
}

// buildChain builds a RequestDelegate by folding middleware from right to left.
// The innermost delegate is a no-op terminal.
func buildChain(middlewares []Middleware) RequestDelegate {
	var next RequestDelegate = func(ctx *Context) error {
		return nil
	}

	for i := len(middlewares) - 1; i >= 0; i-- {
		middleware := middlewares[i]
		current := next
		next = func(ctx *Context) error {
			return middleware.Invoke(ctx, current)
		}
	}

	return next
}
