package http

type EndpointFilterDelegate func(ctx Context) (any, error)

// EndpointFilterFunc represents a function that filters HTTP requests at the endpoint level.
type EndpointFilterFunc func(ctx Context, next EndpointFilterDelegate) error

func (f EndpointFilterFunc) Filter(ctx Context, next EndpointFilterDelegate) error {
	return f(ctx, next)
}

// EndpointFilter represents an interface for filtering HTTP requests at the endpoint level.
type EndpointFilter interface {
	Filter(ctx Context, next EndpointFilterDelegate) error
}

// EndpointOption represents a function that configures an endpoint.
type EndpointOption func(endpoint *Endpoint)

// Endpoint represents a route that is used to handle an HTTP defaultRequest.
type Endpoint struct {
	pattern     string
	methods     []Method
	handler     RequestDelegate
	middlewares []*Middleware
	tags        []string
}

// NewEndpoint creates a new route with the provided pattern, handler and options.
func NewEndpoint(pattern string, handler RequestDelegate, options ...EndpointOption) *Endpoint {
	route := &Endpoint{
		pattern: pattern,
		handler: handler,
		tags:    []string{},
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
func (r *Endpoint) Pattern() string {
	return r.pattern
}

func (r *Endpoint) RequestDelegate() RequestDelegate {
	return r.handler
}

// Methods returns the methods of the route.
func (r *Endpoint) Methods() []Method {
	methods := make([]Method, len(r.methods))
	copy(methods, r.methods)
	return methods
}

// Metadata returns the metadata of the route.
func (r *Endpoint) Metadata() any {
	return nil
}

// Tags returns the tags of the route.
func (r *Endpoint) Tags() []string {
	return r.tags
}

// WithMethods creates a new route option with the provided methods.
func WithMethods(methods ...Method) EndpointOption {
	return func(route *Endpoint) {
		if len(methods) != 0 {
			route.methods = append(route.methods, methods...)
		}
	}
}

// WithMetadata creates a new route option with the provided metadata.
func WithMetadata(metadata any) EndpointOption {
	return func(route *Endpoint) {
		if metadata != nil {
			//handler.metadata[metadata.MetadataKey()] = metadata
		}
	}
}

// WithFilter creates a new route option with the provided middleware.
func WithFilter(filter EndpointFilter) EndpointOption {
	return func(route *Endpoint) {
		if filter != nil {
			route.middlewares = append(route.middlewares, nil)
		}
	}
}

// WithTags creates a new route option with the provided tags.
func WithTags(tags ...string) EndpointOption {
	return func(route *Endpoint) {
		if len(tags) != 0 {
			route.tags = append(route.tags, tags...)
		}
	}
}

// EndpointRegistry represents an interface for managing endpoints.
type EndpointRegistry interface {
	// Register registers the provided route.
	Register(endpoint *Endpoint) error
	// Unregister unregisters the provided route.
	Unregister(endpoint *Endpoint) error
	// Endpoints returns the list of registered endpoints.
	Endpoints() []*Endpoint
	// Match matches the route with the provided context.
	Match(ctx Context) (*Endpoint, bool)
}
