package mvc

import "github.com/procyon-projects/procyon/web/http"

type Context[T, E any] struct {
	ctx          http.Context
	modelAndView ModelAndView
}
