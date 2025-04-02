package http

// RouteGroupBuilder represents a builder that is used to build a route group.
type RouteGroupBuilder struct {
	prefix string

	routes      []*Route
	middlewares []*Middleware
	metadata    Metadata
	tags        []string
}

// NewRouteGroupBuilder creates a new route group builder.
func NewRouteGroupBuilder() *RouteGroupBuilder {
	return &RouteGroupBuilder{
		middlewares: []*Middleware{},
		metadata:    Metadata{},
		tags:        []string{},
	}
}

// MapGroup creates a new route group builder with the provided pattern.
func (r *RouteGroupBuilder) MapGroup(pattern string) *RouteGroupBuilder {
	return &RouteGroupBuilder{
		prefix:      pattern,
		middlewares: []*Middleware{},
		metadata:    Metadata{},
		tags:        []string{},
	}
}

// MapGet creates a new GET route builder with the provided pattern and handler.
func (r *RouteGroupBuilder) MapGet(pattern string, handler RequestHandler) *RouteBuilder {
	routeHandler := NewRoute(pattern, handler, WithMethods(MethodGet))
	r.routes = append(r.routes, routeHandler)
	return nil
}

// MapPost creates a new POST route builder with the provided pattern and handler.
func (r *RouteGroupBuilder) MapPost(pattern string, handler RequestHandler) *RouteBuilder {
	routeHandler := NewRoute(pattern, handler, WithMethods(MethodPost))
	r.routes = append(r.routes, routeHandler)
	return nil
}

// MapDelete creates a new DELETE route builder with the provided pattern and handler.
func (r *RouteGroupBuilder) MapDelete(pattern string, handler RequestHandler) *RouteBuilder {
	routeHandler := NewRoute(pattern, handler, WithMethods(MethodDelete))
	r.routes = append(r.routes, routeHandler)
	return nil
}

// MapPut creates a new PUT route builder with the provided pattern and handler.
func (r *RouteGroupBuilder) MapPut(pattern string, handler RequestHandler) *RouteBuilder {
	routeHandler := NewRoute(pattern, handler, WithMethods(MethodPut))
	r.routes = append(r.routes, routeHandler)
	return newRouteBuilder(routeHandler)
}

// MapPatch creates a new PATCH route builder with the provided pattern and handler.
func (r *RouteGroupBuilder) MapPatch(pattern string, handler RequestHandler) *RouteBuilder {
	routeHandler := NewRoute(pattern, handler, WithMethods(MethodPatch))
	r.routes = append(r.routes, routeHandler)
	return newRouteBuilder(routeHandler)
}

// MapHead creates a new HEAD route builder with the provided pattern and handler.
func (r *RouteGroupBuilder) MapHead(pattern string, handler RequestHandler) *RouteBuilder {
	routeHandler := NewRoute(pattern, handler, WithMethods(MethodHead))
	r.routes = append(r.routes, routeHandler)
	return newRouteBuilder(routeHandler)
}

// MapOptions creates a new OPTIONS route builder with the provided pattern and handler.
func (r *RouteGroupBuilder) MapOptions(pattern string, handler RequestHandler) *RouteBuilder {
	routeHandler := NewRoute(pattern, handler, WithMethods(MethodOptions))
	r.routes = append(r.routes, routeHandler)
	return newRouteBuilder(routeHandler)
}

// MapTrace creates a new TRACE route builder with the provided pattern and handler.
func (r *RouteGroupBuilder) MapTrace(pattern string, handler RequestHandler) *RouteBuilder {
	routeHandler := NewRoute(pattern, handler, WithMethods(MethodTrace))
	r.routes = append(r.routes, routeHandler)
	return newRouteBuilder(routeHandler)
}

// MapMethods creates a new route builder with the provided pattern, methods and handler.
func (r *RouteGroupBuilder) MapMethods(pattern string, methods []Method, handler RequestHandler) *RouteBuilder {
	route := NewRoute(pattern, handler, WithMethods(methods...))
	r.routes = append(r.routes, route)
	return nil
}

// Use adds a middleware to the route group.
func (r *RouteGroupBuilder) Use(middleware MiddlewareFunc, options ...MiddlewareOption) *RouteGroupBuilder {
	if middleware != nil {
		r.middlewares = append(r.middlewares, NewMiddleware(middleware, options...))
	}
	return r
}

// WithMetadata adds metadata to the route group.
func (r *RouteGroupBuilder) WithMetadata(metadata MetadataFunc) *RouteGroupBuilder {
	if metadata != nil {
		metadata(r.metadata)
	}
	return r
}

// WithTags adds tags to the route group.
func (r *RouteGroupBuilder) WithTags(tags ...string) *RouteGroupBuilder {
	r.tags = append(r.tags, tags...)
	return r
}
