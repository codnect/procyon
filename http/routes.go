package http

// Controller represents an interface for mapping routes.
type Controller interface {
	MapRoutes(r *RouteBuilder)
}

type RouteBuilder struct {
}

// Get creates a new GET route handler with the provided pattern and handler.
func (r *RouteBuilder) Get(pattern string, handler RequestDelegate) *RouteHandler {
	return nil
}

// Post creates a new POST route handler with the provided pattern and handler.
func (r *RouteBuilder) Post(pattern string, handler RequestDelegate) *RouteHandler {
	return nil
}

// Delete creates a new DELETE route handler with the provided pattern and handler.
func (r *RouteBuilder) Delete(pattern string, handler RequestDelegate) *RouteHandler {
	return nil
}

// Put creates a new PUT route handler with the provided pattern and handler.
func (r *RouteBuilder) Put(pattern string, handler RequestDelegate) *RouteHandler {
	return nil
}

// Patch creates a new PATCH route handler with the provided pattern and handler.
func (r *RouteBuilder) Patch(pattern string, handler RequestDelegate) *RouteHandler {
	return nil
}

// Head creates a new HEAD route handler with the provided pattern and handler.
func (r *RouteBuilder) Head(pattern string, handler RequestDelegate) *RouteHandler {
	return nil
}

// Methods creates a new route handler with the provided pattern, methods and handler.
func (r *RouteBuilder) Methods(pattern string, methods []Method, handler RequestDelegate) *RouteHandler {
	return nil
}

// Group creates a new route group with the provided prefix.
func (r *RouteBuilder) Group(prefix string) *RouteGroup {
	return nil
}

type RouteHandler struct {
}

func (r *RouteHandler) WithName(name string) *RouteHandler {
	return nil
}

func (r *RouteHandler) WithDescription(description string) *RouteHandler {
	return nil
}

func (r *RouteHandler) Use(middleware Middleware) *RouteHandler {
	return nil
}

func (r *RouteHandler) UseFunc(fn MiddlewareFunc) *RouteHandler {
	return nil
}

// WithTags adds the provided tags to the route handler.
func (r *RouteHandler) WithTags(tags ...string) *RouteHandler {
	return nil
}

// WithMetadata adds the provided metadata to the route handler.
func (r *RouteHandler) WithMetadata(metadata ...any) *RouteHandler {
	return nil
}

type RouteGroup struct {
}

// Get creates a new GET route handler with the provided pattern and handler.
func (r *RouteGroup) Get(pattern string, handler RequestDelegate) *RouteHandler {
	return nil
}

// Post creates a new POST route handler with the provided pattern and handler.
func (r *RouteGroup) Post(pattern string, handler RequestDelegate) *RouteHandler {
	return nil
}

// Delete creates a new DELETE route handler with the provided pattern and handler.
func (r *RouteGroup) Delete(pattern string, handler RequestDelegate) *RouteHandler {
	return nil
}

// Put creates a new PUT route handler with the provided pattern and handler.
func (r *RouteGroup) Put(pattern string, handler RequestDelegate) *RouteHandler {
	return nil
}

// Patch creates a new PATCH route handler with the provided pattern and handler.
func (r *RouteGroup) Patch(pattern string, handler RequestDelegate) *RouteHandler {
	return nil
}

// Head creates a new HEAD route handler with the provided pattern and handler.
func (r *RouteGroup) Head(pattern string, handler RequestDelegate) *RouteHandler {
	return nil
}

// Methods creates a new route handler with the provided pattern, methods and handler.
func (r *RouteGroup) Methods(pattern string, methods []Method, handler RequestDelegate) *RouteHandler {
	return nil
}

// Group creates a new route group with the provided prefix.
func (r *RouteGroup) Group(prefix string) *RouteGroup {
	return nil
}

// Use adds a middleware to the route group.
func (r *RouteGroup) Use(fn MiddlewareFunc) *RouteGroup {
	return nil
}

// WithTags adds tags to the route group.
func (r *RouteGroup) WithTags(tags ...string) *RouteGroup {
	return nil
}
