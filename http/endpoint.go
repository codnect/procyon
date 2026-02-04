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
	"path"
	"strings"
)

// Endpoint represents a fully described HTTP endpoint definition.
type Endpoint struct {
	// path is the route pattern associated with this endpoint
	// (e.g. "/users/{id}", "/health", "/files/**").
	path string

	// method is the HTTP method this endpoint responds to
	// (e.g. GET, POST, PUT).
	method Method

	// delegate is the request handler invoked when this endpoint
	// matches an incoming request.
	delegate RequestDelegate
}

// Path returns the route pattern of the endpoint.
func (e Endpoint) Path() string {
	return e.path
}

// Method returns the HTTP method associated with the endpoint.
func (e Endpoint) Method() Method {
	return e.method
}

// RequestDelegate returns the handler responsible for processing
// requests matched to this endpoint.
func (e Endpoint) RequestDelegate() RequestDelegate {
	return e.delegate
}

// EndpointDataSource provides a collection of endpoint definitions.
type EndpointDataSource interface {
	// Endpoints returns all available endpoint definitions.
	Endpoints() []*Endpoint
}

// EndpointMatcher matches an incoming request context
// against a set of endpoint definitions.
//
// The matcher is responsible only for selection logic;
// execution and metadata processing are handled elsewhere
// in the request pipeline.
type EndpointMatcher interface {
	// Match attempts to find a matching endpoint for the given context.
	// It returns the matched endpoint and true if a match is found,
	// or nil and false otherwise.
	Match(ctx *Context) (*Endpoint, bool)
}

// Endpoints interface represents a collection of HTTP routes.
// It provides methods to map handler functions to specific paths and HTTP methods.
type Endpoints interface {
	// MapAny maps a handler function to the specified path for all HTTP methods.
	MapAny(path string, handler Handler) *EndpointBuilder
	// MapMethods maps a handler function to the specified path for the given HTTP methods.
	MapMethods(path string, methods []Method, handler Handler) *EndpointBuilder
	// MapGet maps a handler function to the specified path for the GET HTTP method.
	MapGet(path string, handler Handler) *EndpointBuilder
	// MapPost maps a handler function to the specified path for the POST HTTP method.
	MapPost(path string, handler Handler) *EndpointBuilder
	// MapPut maps a handler function to the specified path for the PUT HTTP method.
	MapPut(path string, handler Handler) *EndpointBuilder
	// MapDelete maps a handler function to the specified path for the DELETE HTTP method.
	MapDelete(path string, handler Handler) *EndpointBuilder
	// MapPatch maps a handler function to the specified path for the PATCH HTTP method.
	MapPatch(path string, handler Handler) *EndpointBuilder
	// MapGroup creates a new EndpointGroup with the specified prefix.
	MapGroup(prefix string) *EndpointGroup
}

// EndpointConfigurer interface represents a type that can configure HTTP routes.
type EndpointConfigurer interface {
	// ConfigureEndpoints method configures the given routes.
	ConfigureEndpoints(endpoints Endpoints)
}

// EndpointBuilder represents a route definition bound to a path,
// a set of HTTP methods, and a handler.
type EndpointBuilder struct {
	path    string
	methods []Method
	handler Handler
}

// newEndpointBuilder creates a new EndpointBuilder with the given path, methods, and handler.
func newEndpointBuilder(path string, methods []Method, handler Handler) *EndpointBuilder {
	return &EndpointBuilder{
		path:    path,
		methods: methods,
		handler: handler,
	}
}

// EndpointGroup represents a group of routes with a common prefix.
type EndpointGroup struct {
	prefix   string
	routes   []*EndpointBuilder
	children []*EndpointGroup
}

func newEndpointGroup(prefix string) *EndpointGroup {
	if prefix == "" {
		prefix = "/"
	}

	return &EndpointGroup{
		prefix:   prefix,
		routes:   make([]*EndpointBuilder, 0),
		children: make([]*EndpointGroup, 0),
	}
}

// MapAny maps a handler function to the specified path for all HTTP methods within the group.
func (g *EndpointGroup) MapAny(path string, handler Handler) *EndpointBuilder {
	return g.MapMethods(path, nil, handler)
}

// MapMethods maps a handler function to the specified path for the given HTTP methods within the group.
func (g *EndpointGroup) MapMethods(path string, methods []Method, handler Handler) *EndpointBuilder {
	result := joinPaths(g.prefix, path)

	routeHandler := newEndpointBuilder(result, methods, handler)
	g.routes = append(g.routes, routeHandler)
	return routeHandler
}

// MapGet maps a handler function to the specified path for the GET HTTP method within the group.
func (g *EndpointGroup) MapGet(path string, handler Handler) *EndpointBuilder {
	return g.MapMethods(path, []Method{MethodGet}, handler)
}

// MapPost maps a handler function to the specified path for the POST HTTP method within the group.
func (g *EndpointGroup) MapPost(path string, handler Handler) *EndpointBuilder {
	return g.MapMethods(path, []Method{MethodPost}, handler)
}

// MapPut maps a handler function to the specified path for the PUT HTTP method within the group.
func (g *EndpointGroup) MapPut(path string, handler Handler) *EndpointBuilder {
	return g.MapMethods(path, []Method{MethodPut}, handler)
}

// MapDelete maps a handler function to the specified path for the DELETE HTTP method within the group.
func (g *EndpointGroup) MapDelete(path string, handler Handler) *EndpointBuilder {
	return g.MapMethods(path, []Method{MethodDelete}, handler)
}

// MapPatch maps a handler function to the specified path for the PATCH HTTP method within the group.
func (g *EndpointGroup) MapPatch(path string, handler Handler) *EndpointBuilder {
	return g.MapMethods(path, []Method{MethodPatch}, handler)
}

// MapGroup creates a new EndpointGroup with the specified prefix within the current group.
func (g *EndpointGroup) MapGroup(prefix string) *EndpointGroup {
	result := joinPaths(g.prefix, prefix)
	group := newEndpointGroup(result)
	g.children = append(g.children, group)
	return group
}

// joinPaths joins multiple path elements into a single path string,
// ensuring that there is exactly one '/' separator between elements
// and preserving leading and trailing slashes.
func joinPaths(elem ...string) string {
	var result string
	if !strings.HasPrefix(elem[0], "/") {
		elem[0] = "/" + elem[0]
		result = path.Join(elem...)[1:]
	} else {
		result = path.Join(elem...)
	}

	if strings.HasSuffix(elem[len(elem)-1], "/") && !strings.HasSuffix(result, "/") {
		result += "/"
	}

	return result
}
