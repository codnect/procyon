package http

type HandlerFunc func(ctx Context, next RequestDelegate) error

type HandlerChain []HandlerFunc

type RequestHandler interface {
	HandleRequest(ctx Context) error
}

type requestHandler struct {
	fn func(ctx Context) error
}

func Handler(fn func(ctx Context) error) RequestHandler {
	return &requestHandler{
		fn: fn,
	}
}

func (h *requestHandler) HandleRequest(ctx Context) error {
	return h.fn(ctx)
}
