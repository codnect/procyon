package http

type HandlerFunction func(ctx Context, next RequestDelegate) error

type Handler interface {
	Invoke(ctx Context) (any, error)
}

type requestHandler struct {
	fn func(ctx Context) error
}

func Handle(fn func(ctx Context) error) Handler {
	return &requestHandler{
		fn: fn,
	}
}

func (h *requestHandler) Invoke(ctx Context) (any, error) {
	err := h.fn(ctx)
	return nil, err
}
