package http

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
	handler     Handler
	middlewares []*Middleware
	metadata    Metadata
	tags        []string
}

// NewRoute creates a new route with the provided pattern, handler and options.
func NewRoute(pattern string, handler Handler, options ...RouteOption) *Route {
	route := &Route{
		pattern:  pattern,
		handler:  handler,
		metadata: Metadata{},
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
func (r *Route) Metadata() Metadata {
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
func WithMetadata(metadata MetadataFunc) RouteOption {
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
func (r *RouteBuilder) WithMetadata(metadata MetadataFunc) *RouteBuilder {
	if metadata != nil {
		//r.route.metadata[metadata.MetadataKey()] = metadata
	}
	return r
}

// WithTags adds the provided tags to the route.
func (r *RouteBuilder) WithTags(tags ...string) *RouteBuilder {
	r.route.tags = append(r.route.tags, tags...)
	return r
}

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
func (r *RouteGroupBuilder) MapGet(pattern string, handler Handler) *RouteBuilder {
	routeHandler := NewRoute(pattern, handler, WithMethods(MethodGet))
	r.routes = append(r.routes, routeHandler)
	return nil
}

// MapPost creates a new POST route builder with the provided pattern and handler.
func (r *RouteGroupBuilder) MapPost(pattern string, handler Handler) *RouteBuilder {
	routeHandler := NewRoute(pattern, handler, WithMethods(MethodPost))
	r.routes = append(r.routes, routeHandler)
	return nil
}

// MapDelete creates a new DELETE route builder with the provided pattern and handler.
func (r *RouteGroupBuilder) MapDelete(pattern string, handler Handler) *RouteBuilder {
	routeHandler := NewRoute(pattern, handler, WithMethods(MethodDelete))
	r.routes = append(r.routes, routeHandler)
	return nil
}

// MapPut creates a new PUT route builder with the provided pattern and handler.
func (r *RouteGroupBuilder) MapPut(pattern string, handler Handler) *RouteBuilder {
	routeHandler := NewRoute(pattern, handler, WithMethods(MethodPut))
	r.routes = append(r.routes, routeHandler)
	return newRouteBuilder(routeHandler)
}

// MapPatch creates a new PATCH route builder with the provided pattern and handler.
func (r *RouteGroupBuilder) MapPatch(pattern string, handler Handler) *RouteBuilder {
	routeHandler := NewRoute(pattern, handler, WithMethods(MethodPatch))
	r.routes = append(r.routes, routeHandler)
	return newRouteBuilder(routeHandler)
}

// MapHead creates a new HEAD route builder with the provided pattern and handler.
func (r *RouteGroupBuilder) MapHead(pattern string, handler Handler) *RouteBuilder {
	routeHandler := NewRoute(pattern, handler, WithMethods(MethodHead))
	r.routes = append(r.routes, routeHandler)
	return newRouteBuilder(routeHandler)
}

// MapOptions creates a new OPTIONS route builder with the provided pattern and handler.
func (r *RouteGroupBuilder) MapOptions(pattern string, handler Handler) *RouteBuilder {
	routeHandler := NewRoute(pattern, handler, WithMethods(MethodOptions))
	r.routes = append(r.routes, routeHandler)
	return newRouteBuilder(routeHandler)
}

// MapTrace creates a new TRACE route builder with the provided pattern and handler.
func (r *RouteGroupBuilder) MapTrace(pattern string, handler Handler) *RouteBuilder {
	routeHandler := NewRoute(pattern, handler, WithMethods(MethodTrace))
	r.routes = append(r.routes, routeHandler)
	return newRouteBuilder(routeHandler)
}

// MapMethods creates a new route builder with the provided pattern, methods and handler.
func (r *RouteGroupBuilder) MapMethods(pattern string, methods []Method, handler Handler) *RouteBuilder {
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

type RouteRegistry struct {
	//tree []*routeTree
}

func NewRouteRegistry() *RouteRegistry {

	/*registry := &RouteRegistry{
		make([]*routeTree, 8),
	}

	methods := []Method{
		MethodGet,
		MethodHead,
		MethodPost,
		MethodPut,
		MethodPatch,
		MethodDelete,
		MethodOptions,
		MethodTrace,
	}

	for _, method := range methods {
		registry.tree[httpMethodToInt(method)] = &routeTree{
			staticRoutes: make(map[string]*Route, 0),
		}
	}

	return registry*/
	return nil
}

func (r *RouteRegistry) Register(route *Route) error {

	/*methods := route.Methods()

	if len(methods) == 0 {
		return errors.New("route must have at least one method")
	}

	for _, method := range methods {
		intValue := httpMethodToInt(method)
		if intValue == -1 {
			return fmt.Errorf("invalid method: %s", method)
		}

		methodTree := r.tree[intValue]

		if methodTree.children == nil {
			methodTree.children = &routeNode{}
		}

		methodTree.addRoute(route)
	}
	*/
	return nil
}

func (r *RouteRegistry) Find(ctx Context) (*Route, bool) {
	/*request := ctx.Request()
	path := request.Path()

	intValue := httpMethodToInt(request.Method())
	if intValue < 0 || intValue >= len(r.tree) {
		return nil, false
	}

	methodTree := r.tree[intValue]

	if route, ok := methodTree.staticRoutes[path]; ok {
		return route, true
	}

	route := methodTree.findMatchingRoute(ctx)
	return route, true*/
	return nil, false
}

// List returns a list of all registered routes.
func (r *RouteRegistry) List() []*Route {
	return nil
}

func (r *RouteRegistry) Unregister(route *Route) error {
	return nil
}
