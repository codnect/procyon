package mvc

import "codnect.io/procyon/web/http"

type requestHandler[T, E any] struct {
	fn Function[T, E]
}

func Handle[T, E any](fn Function[T, E]) http.Handler {
	return &requestHandler[T, E]{
		fn: fn,
	}
}

func (h *requestHandler[T, E]) Invoke(ctx http.Context) (any, error) {
	mvcContext := &Context[T, E]{
		ctx: ctx,
	}

	err := h.fn(mvcContext)
	return mvcContext.modelAndView, err
}
