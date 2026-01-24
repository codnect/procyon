package http

import (
	"path"
	"strings"
)

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

// Routes interface represents a collection of HTTP routes.
// It provides methods to map handler functions to specific paths and HTTP methods.
type Routes interface {
	// MapAny maps a handler function to the specified path for all HTTP methods.
	MapAny(path string, handler Handler) *RouteHandler
	// MapMethods maps a handler function to the specified path for the given HTTP methods.
	MapMethods(path string, methods []Method, handler Handler) *RouteHandler
	// MapGet maps a handler function to the specified path for the GET HTTP method.
	MapGet(path string, handler Handler) *RouteHandler
	// MapPost maps a handler function to the specified path for the POST HTTP method.
	MapPost(path string, handler Handler) *RouteHandler
	// MapPut maps a handler function to the specified path for the PUT HTTP method.
	MapPut(path string, handler Handler) *RouteHandler
	// MapDelete maps a handler function to the specified path for the DELETE HTTP method.
	MapDelete(path string, handler Handler) *RouteHandler
	// MapPatch maps a handler function to the specified path for the PATCH HTTP method.
	MapPatch(path string, handler Handler) *RouteHandler
	// MapGroup creates a new RouteGroup with the specified prefix.
	MapGroup(prefix string) *RouteGroup
}

// RouteConfigurer interface represents a type that can configure HTTP routes.
type RouteConfigurer interface {
	// ConfigureRoutes method configures the given routes.
	ConfigureRoutes(routes Routes)
}

// RouteHandler represents a route definition bound to a path,
// a set of HTTP methods, and a handler.
type RouteHandler struct {
	path    string
	methods []Method
	handler Handler
}

// newRouteHandler creates a new RouteHandler with the given path, methods, and handler.
func newRouteHandler(path string, methods []Method, handler Handler) *RouteHandler {
	return &RouteHandler{
		path:    path,
		methods: methods,
		handler: handler,
	}
}

// RouteGroup represents a group of routes with a common prefix.
type RouteGroup struct {
	prefix   string
	routes   []*RouteHandler
	children []*RouteGroup
}

func newRouteGroup(prefix string) *RouteGroup {
	if prefix == "" {
		prefix = "/"
	}

	return &RouteGroup{
		prefix:   prefix,
		routes:   make([]*RouteHandler, 0),
		children: make([]*RouteGroup, 0),
	}
}

// MapAny maps a handler function to the specified path for all HTTP methods within the group.
func (g *RouteGroup) MapAny(path string, handler Handler) *RouteHandler {
	return g.MapMethods(path, nil, handler)
}

// MapMethods maps a handler function to the specified path for the given HTTP methods within the group.
func (g *RouteGroup) MapMethods(path string, methods []Method, handler Handler) *RouteHandler {
	result := joinPaths(g.prefix, path)

	routeHandler := newRouteHandler(result, methods, handler)
	g.routes = append(g.routes, routeHandler)
	return routeHandler
}

// MapGet maps a handler function to the specified path for the GET HTTP method within the group.
func (g *RouteGroup) MapGet(path string, handler Handler) *RouteHandler {
	return g.MapMethods(path, []Method{MethodGet}, handler)
}

// MapPost maps a handler function to the specified path for the POST HTTP method within the group.
func (g *RouteGroup) MapPost(path string, handler Handler) *RouteHandler {
	return g.MapMethods(path, []Method{MethodPost}, handler)
}

// MapPut maps a handler function to the specified path for the PUT HTTP method within the group.
func (g *RouteGroup) MapPut(path string, handler Handler) *RouteHandler {
	return g.MapMethods(path, []Method{MethodPut}, handler)
}

// MapDelete maps a handler function to the specified path for the DELETE HTTP method within the group.
func (g *RouteGroup) MapDelete(path string, handler Handler) *RouteHandler {
	return g.MapMethods(path, []Method{MethodDelete}, handler)
}

// MapPatch maps a handler function to the specified path for the PATCH HTTP method within the group.
func (g *RouteGroup) MapPatch(path string, handler Handler) *RouteHandler {
	return g.MapMethods(path, []Method{MethodPatch}, handler)
}

// MapGroup creates a new RouteGroup with the specified prefix within the current group.
func (g *RouteGroup) MapGroup(prefix string) *RouteGroup {
	result := joinPaths(g.prefix, prefix)
	group := newRouteGroup(result)
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
