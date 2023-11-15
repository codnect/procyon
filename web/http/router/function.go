package router

import (
	"github.com/procyon-projects/procyon/web/http"
	"github.com/procyon-projects/procyon/web/http/middleware"
)

type Function struct {
	pattern    string
	method     http.Method
	handler    http.Handler
	consumes   []string
	produces   []string
	attributes map[string]any

	middlewares []*middleware.Middleware

	group *Routes
}

func newFunction(method http.Method, pattern string, handler http.Handler, group *Routes) *Function {
	return &Function{
		pattern:     pattern,
		method:      method,
		handler:     handler,
		attributes:  map[string]any{},
		middlewares: []*middleware.Middleware{},
		group:       group,
	}
}

func (r *Function) Accept(mediaTypes ...string) *Function {
	r.consumes = append(r.consumes, mediaTypes...)
	return r
}

func (r *Function) ContentType(mediaTypes ...string) *Function {
	r.produces = append(r.produces, mediaTypes...)
	return r
}

func (r *Function) Attribute(name string, val any) *Function {
	r.attributes[name] = val
	return r
}

func (r *Function) AttributeMap(attrs map[string]any) *Function {
	for key, val := range attrs {
		r.attributes[key] = val
	}

	return r
}

func (r *Function) UseMiddleware(function middleware.Function, options ...middleware.Option) *Function {
	if function != nil {
		r.middlewares = append(r.middlewares, middleware.New(function, options...))
	}

	return r
}
