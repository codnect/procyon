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
	"errors"
	"strings"
)

// ErrInvalidPattern indicates the provided pattern is malformed.
var ErrInvalidPattern = errors.New("invalid route pattern")

// Router routes HTTP requests using radix-tree and ant-style pattern matching.
type Router struct {
	root      *radixNode
	antRoutes []*antRoute
	notFound  Handler
}

// NewRouter creates a new Router instance.
func NewRouter() *Router {
	return &Router{
		root:     newRadixNode("/"),
		notFound: func(ctx Context) {},
	}
}

// Handle registers a route for the given method and pattern.
func (r *Router) Handle(method, pattern string, handler Handler, middleware ...Middleware) error {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return ErrInvalidPattern
	}

	if strings.ContainsAny(pattern, "*?{") {
		route, err := newAntRoute(method, pattern, handler, middleware)
		if err != nil {
			return err
		}
		r.antRoutes = append(r.antRoutes, route)
		return nil
	}

	r.root.addRoute(method, pattern, handler, middleware)
	return nil
}

// NotFound sets the handler executed when no routes match.
func (r *Router) NotFound(handler Handler) {
	if handler != nil {
		r.notFound = handler
	}
}

// Resolve resolves the incoming method and path into a handler and middleware chain.
func (r *Router) Resolve(method, path string) (Handler, []Middleware, map[string]string) {
	params, entry, ok := r.root.match(strings.ToUpper(method), path)
	if ok {
		return entry.handler, entry.middleware, params
	}

	for _, route := range r.antRoutes {
		if params, matched := route.match(strings.ToUpper(method), path); matched {
			return route.handler, route.middleware, params
		}
	}

	return r.notFound, nil, map[string]string{}
}
