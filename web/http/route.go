package http

import "codnect.io/procyon/metadata"

// Controller represents a controller that is used to map routes.
type Controller interface {
	MapRoutes(routes *RouteGroupBuilder)
}

// RouteOption represents an option that is used to configure the route.
type RouteOption func(handler *Route)

// Route represents a route that is used to handle an HTTP request.
type Route struct {
	pattern     string
	methods     []Method
	handler     RequestHandler
	middlewares []*Middleware
	metadata    metadata.Collection
	tags        []string
}

// NewRoute creates a new route with the provided pattern, handler and options.
func NewRoute(pattern string, handler RequestHandler, options ...RouteOption) *Route {
	route := &Route{
		pattern:  pattern,
		handler:  handler,
		metadata: metadata.Collection{},
		tags:     []string{},
	}

	for _, option := range options {
		option(route)
	}

	if len(route.methods) == 0 {
		route.methods = []Method{
			MethodGet,
		}
	}

	return route
}

// Pattern returns the pattern of the route.
func (r *Route) Pattern() string {
	return r.pattern
}

// Methods returns the methods of the route.
func (r *Route) Methods() []Method {
	methods := make([]Method, len(r.methods))
	copy(methods, r.methods)
	return methods
}

// HandlerChain returns the handler chain of the route.
func (r *Route) HandlerChain() HandlerChain {
	return nil
}

// Metadata returns the metadata of the route.
func (r *Route) Metadata() metadata.Collection {
	return nil
}

// Tags returns the tags of the route.
func (r *Route) Tags() []string {
	return r.tags
}

// WithMethods creates a new route option with the provided methods.
func WithMethods(methods ...Method) RouteOption {
	return func(handler *Route) {
		if len(methods) != 0 {
			handler.methods = append(handler.methods, methods...)
		}
	}
}

// WithMetadata creates a new route option with the provided metadata.
func WithMetadata(metadata metadata.Metadata) RouteOption {
	return func(handler *Route) {
		if metadata != nil {
			//handler.metadata[metadata.MetadataKey()] = metadata
		}
	}
}

// WithMiddleware creates a new route option with the provided middleware.
func WithMiddleware(middleware *Middleware) RouteOption {
	return func(handler *Route) {
		if middleware != nil {
			handler.middlewares = append(handler.middlewares, middleware)
		}
	}
}

// WithTags creates a new route option with the provided tags.
func WithTags(tags ...string) RouteOption {
	return func(handler *Route) {
		if len(tags) != 0 {
			handler.tags = append(handler.tags, tags...)
		}
	}
}

// RouteBuilder represents a builder that is used to configure a route.
type RouteBuilder struct {
	route *Route
}

// newRouteBuilder creates a new route builder with the provided route.
func newRouteBuilder(route *Route) *RouteBuilder {
	return &RouteBuilder{
		route: route,
	}
}

// Use adds the provided middleware to the route.
func (r *RouteBuilder) Use(middleware MiddlewareFunc, options ...MiddlewareOption) *RouteBuilder {
	if middleware != nil {
		r.route.middlewares = append(r.route.middlewares, NewMiddleware(middleware, options...))
	}
	return r
}

// WithMetadata adds the provided metadata to the route.
func (r *RouteBuilder) WithMetadata(metadata metadata.Metadata) *RouteBuilder {
	if metadata != nil {
		r.route.metadata[metadata.MetadataKey()] = metadata
	}
	return r
}

// WithTags adds the provided tags to the route.
func (r *RouteBuilder) WithTags(tags ...string) *RouteBuilder {
	r.route.tags = append(r.route.tags, tags...)
	return r
}
