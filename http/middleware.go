package http

// MiddlewareFunc defines a function type for middleware that can process HTTP requests.
type MiddlewareFunc func(ctx *Context, next RequestDelegate) error

func (f MiddlewareFunc) Invoke(ctx *Context, next RequestDelegate) error {
	return f(ctx, next)
}

// Middleware represents an interface for middleware that can process HTTP requests.
type Middleware interface {
	// Invoke processes the HTTP request within the given context and calls the next delegate in the chain.
	Invoke(ctx *Context, next RequestDelegate) error
}
