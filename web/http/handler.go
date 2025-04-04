package http

type HandlerFunc func(ctx Context, next RequestDelegate) error

type HandlerChain []HandlerFunc

type Handler interface {
	HandleRequest(ctx Context) error
}

type RequestHandler struct {
	fn func(ctx Context) error
}

func Handle(fn func(ctx Context) error) Handler {
	return &RequestHandler{
		fn: fn,
	}
}

func (h *RequestHandler) HandleRequest(ctx Context) error {
	return h.fn(ctx)
}

// MiddlewareFunc represents a function that is used to process the request before it reaches the route handler.
type MiddlewareFunc func(ctx Context, next RequestDelegate) error

// MiddlewareOption represents an option that is used to configure the middleware.
type MiddlewareOption func(middleware *Middleware)

// Middleware represents a middleware that is used to process the request before it reaches the route handler.
type Middleware struct {
	fn    MiddlewareFunc
	order int
}

// NewMiddleware creates a new middleware with the provided function and options.
func NewMiddleware(fn MiddlewareFunc, options ...MiddlewareOption) *Middleware {
	if fn == nil {
		panic("nil middleware function")
	}

	middleware := &Middleware{
		fn: fn,
	}

	for _, option := range options {
		option(middleware)
	}

	return middleware
}

// Order returns the order of the middleware.
func (m *Middleware) Order() int {
	return m.order
}

// WithOrder creates a new middleware option with the provided order.
// The order is used to determine the order of the middleware in the middleware chain.
func WithOrder(order int) MiddlewareOption {
	return func(middleware *Middleware) {
		middleware.order = order
	}
}
