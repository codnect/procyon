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
	"fmt"
	"path"
	"strings"
)

// Endpoint represents a fully described HTTP endpoint definition.
type Endpoint struct {
	// method is the HTTP method this endpoint responds to
	// (e.g. GET, POST, PUT).
	method Method

	// path is the route pattern associated with this endpoint
	// (e.g. "/users/{id}", "/health", "/files/**").
	path string

	// delegate is the request handler invoked when this endpoint
	// matches an incoming request.
	delegate RequestDelegate
}

// NewEndpoint creates a new Endpoint with the specified method, path, and delegate.
func NewEndpoint(method Method, path string, delegate RequestDelegate) *Endpoint {
	if delegate == nil {
		panic("nil request delegate")
	}

	return &Endpoint{
		method:   method,
		path:     path,
		delegate: delegate,
	}
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

// EndpointSource provides access to endpoint definitions.
type EndpointSource interface {
	// Endpoints returns the available endpoint definitions.
	Endpoints() []*Endpoint
}

// EndpointRegistrar registers endpoint definitions.
type EndpointRegistrar interface {
	// Register registers the given endpoint.
	Register(endpoint *Endpoint) error
}

// DefaultEndpointRegistry is the default registry for storing endpoint definitions.
type DefaultEndpointRegistry struct {
	endpoints []*Endpoint
}

// NewDefaultEndpointRegistry creates a new empty DefaultEndpointRegistry.
func NewDefaultEndpointRegistry() *DefaultEndpointRegistry {
	return &DefaultEndpointRegistry{
		endpoints: make([]*Endpoint, 0),
	}
}

// Register registers the given endpoint.
func (r *DefaultEndpointRegistry) Register(endpoint *Endpoint) error {
	if endpoint == nil {
		return errors.New("nil endpoint")
	}

	r.endpoints = append(r.endpoints, endpoint)
	return nil
}

// Endpoints returns the registered endpoint definitions.
func (r *DefaultEndpointRegistry) Endpoints() []*Endpoint {
	return r.endpoints
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

// EndpointMapper represents a type that can map HTTP endpoints.
type EndpointMapper interface {
	// MapEndpoints maps HTTP endpoints to the given endpoint collection.
	MapEndpoints(endpoints Endpoints)
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
	if handler == nil {
		panic("nil handler")
	}

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

// endpointMappingProcessor processes EndpointMapper components and registers
// their mapped endpoints after component initialization.
type endpointMappingProcessor struct {
	endpointRegistrar EndpointRegistrar
	executors         ResultExecutorRegistry
}

// newEndpointMappingProcessor creates a new endpointMappingProcessor with the
// given endpoint registrar and result executor registry.
func newEndpointMappingProcessor(endpointRegistrar EndpointRegistrar,
	executorRegistry ResultExecutorRegistry) *endpointMappingProcessor {
	if endpointRegistrar == nil {
		panic("nil endpoint registrar")
	}

	if executorRegistry == nil {
		panic("nil result executor registry")
	}

	return &endpointMappingProcessor{
		endpointRegistrar: endpointRegistrar,
		executors:         executorRegistry,
	}
}

// ProcessAfterInit maps and registers endpoints from EndpointMapper components
// after their initialization has completed.
func (p *endpointMappingProcessor) ProcessAfterInit(_ context.Context, name string, instance any) (any, error) {
	mapper, ok := instance.(EndpointMapper)
	if !ok {
		return instance, nil
	}

	group := newEndpointGroup("/")
	mapper.MapEndpoints(group)

	if err := p.collectEndpoints(group); err != nil {
		return nil, fmt.Errorf("map endpoints for %q: %w", name, err)
	}

	return instance, nil
}

// collectEndpoints recursively collects endpoints from the given endpoint group
// and registers them with the endpoint registrar.
func (p *endpointMappingProcessor) collectEndpoints(group *EndpointGroup) error {
	for _, route := range group.routes {
		delegate := p.createRequestDelegate(route.handler)
		for _, method := range route.methods {
			endpoint := NewEndpoint(method, route.path, delegate)
			if err := p.endpointRegistrar.Register(endpoint); err != nil {
				return err
			}
		}
	}

	for _, child := range group.children {
		if err := p.collectEndpoints(child); err != nil {
			return err
		}
	}

	return nil
}

// createRequestDelegate creates a request delegate that invokes the given
// handler and executes its result using the appropriate result executor.
func (p *endpointMappingProcessor) createRequestDelegate(handler Handler) RequestDelegate {
	return func(ctx *Context) error {
		result, err := handler.Handle(ctx)

		if err != nil {
			return err
		}

		if result == nil {
			return nil
		}

		executor, ok := p.executors.Resolve(result)

		if !ok {
			return fmt.Errorf("no result executor for %T", result)
		}
		return executor.Execute(ctx, result)
	}
}

// joinPaths joins multiple path elements into a single path string,
// ensuring that there is exactly one '/' separator between elements
// and preserving leading and trailing slashes.
func joinPaths(elem ...string) string {
	last := elem[len(elem)-1]

	for i, e := range elem {
		if !strings.HasPrefix(e, "/") {
			elem[i] = "/" + e
		}
	}

	result := path.Join(elem...)

	if last != "" && strings.HasSuffix(last, "/") && !strings.HasSuffix(result, "/") {
		result += "/"
	}

	return result
}
