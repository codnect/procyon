package mvc

import "codnect.io/procyon/web/http"

type Context[T, E any] struct {
	ctx          http.Context
	modelAndView ModelAndView
}
