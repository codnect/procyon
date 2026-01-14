package http

// Controller represents an interface for mapping routes.
type Controller interface {
	MapRoutes(r Router)
}

type Router interface {
	Route(handler Handler) *Route
	RoutePath(path string, handler Handler) *Route

	Group(prefix string) *RouteGroup

	Get(pattern string, handler Handler) *Route
	Post(pattern string, handler Handler) *Route
	Delete(pattern string, handler Handler) *Route
	Put(pattern string, handler Handler) *Route
	Patch(pattern string, handler Handler) *Route
	Head(pattern string, handler Handler) *Route
}

type defaultRouter struct {
}

func (r *defaultRouter) Route(handler Handler) *Route {
	return nil
}

func (r *defaultRouter) RoutePath(path string, handler Handler) *Route {
	return nil
}

// Get creates a new GET route handler with the provided pattern and handler.
func (r *defaultRouter) Get(pattern string, handler Handler) *Route {
	return nil
}

// Post creates a new POST route handler with the provided pattern and handler.
func (r *defaultRouter) Post(pattern string, handler Handler) *Route {
	return nil
}

// Delete creates a new DELETE route handler with the provided pattern and handler.
func (r *defaultRouter) Delete(pattern string, handler Handler) *Route {
	return nil
}

// Put creates a new PUT route handler with the provided pattern and handler.
func (r *defaultRouter) Put(pattern string, handler Handler) *Route {
	return nil
}

// Patch creates a new PATCH route handler with the provided pattern and handler.
func (r *defaultRouter) Patch(pattern string, handler Handler) *Route {
	return nil
}

// Head creates a new HEAD route handler with the provided pattern and handler.
func (r *defaultRouter) Head(pattern string, handler Handler) *Route {
	return nil
}

// Group creates a new route group with the provided prefix.
func (r *defaultRouter) Group(prefix string) *RouteGroup {
	return nil
}

type Route struct {
}

func (r *Route) WithName(name string) *Route {
	return nil
}

func (r *Route) WithDescription(description string) *Route {
	return nil
}

func (r *Route) WithMethods(methods ...Method) *Route {
	return nil
}

func (r *Route) Produces(contentTypes ...string) *Route {
	return nil
}

func (r *Route) Consumes(contentTypes ...string) *Route {
	return nil
}

func (r *Route) UseFilter(filter Filter) *Route {
	return nil
}

func (r *Route) UseFilterFunc(filter FilterFunc) *Route {
	return nil
}

// WithTags adds the provided tags to the route handler.
func (r *Route) WithTags(tags ...string) *Route {
	return nil
}

// WithMetadata adds the provided metadata to the route handler.
func (r *Route) WithMetadata(metadata ...any) *Route {
	return nil
}

type RouteGroup struct {
}

// Get creates a new GET route handler with the provided pattern and handler.
func (r *RouteGroup) Get(pattern string, handler HandlerFunc) *Route {
	return nil
}

// Post creates a new POST route handler with the provided pattern and handler.
func (r *RouteGroup) Post(pattern string, handler HandlerFunc) *Route {
	return nil
}

// Delete creates a new DELETE route handler with the provided pattern and handler.
func (r *RouteGroup) Delete(pattern string, handler HandlerFunc) *Route {
	return nil
}

// Put creates a new PUT route handler with the provided pattern and handler.
func (r *RouteGroup) Put(pattern string, handler HandlerFunc) *Route {
	return nil
}

// Patch creates a new PATCH route handler with the provided pattern and handler.
func (r *RouteGroup) Patch(pattern string, handler HandlerFunc) *Route {
	return nil
}

// Head creates a new HEAD route handler with the provided pattern and handler.
func (r *RouteGroup) Head(pattern string, handler HandlerFunc) *Route {
	return nil
}

// Methods creates a new route handler with the provided pattern, methods and handler.
func (r *RouteGroup) Methods(pattern string, methods []Method, handler HandlerFunc) *Route {
	return nil
}

// Group creates a new route group with the provided prefix.
func (r *RouteGroup) Group(prefix string) *RouteGroup {
	return nil
}

func (r *RouteGroup) UseFilter(fn Filter) *RouteGroup {
	return nil
}

func (r *RouteGroup) UseFilterFunc(fn FilterFunc) *RouteGroup {
	return nil
}

// WithTags adds tags to the route group.
func (r *RouteGroup) WithTags(tags ...string) *RouteGroup {
	return nil
}
