package router

import (
	"github.com/procyon-projects/procyon/web/http"
	"github.com/procyon-projects/procyon/web/http/middleware"
)

type Routes struct {
	pattern      string
	accepts      []string
	contentTypes []string
	nestedRoutes []*Routes
	attributes   map[string]any

	handlers    []*Function
	middlewares []*middleware.Middleware

	parentGroup *Routes
}

func newRoutes() *Routes {
	return &Routes{
		accepts:      []string{},
		contentTypes: []string{},
		nestedRoutes: []*Routes{},
		attributes:   map[string]any{},
	}
}

func Route() *Routes {
	return newRoutes()
}

func RouteWith(pattern string) *Routes {
	routes := newRoutes()
	routes.pattern = pattern
	return routes
}

func (r *Routes) Nest(pattern string) *Routes {
	nestedGroup := newRoutes()
	nestedGroup.pattern = pattern
	nestedGroup.parentGroup = r
	r.nestedRoutes = append(r.nestedRoutes, nestedGroup)
	return nestedGroup
}

func (r *Routes) Route(method http.Method, handler http.Handler) *Function {
	routerFn := newFunction(method, "", handler, r)
	r.handlers = append(r.handlers, routerFn)
	return routerFn
}

func (r *Routes) GET(pattern string, handler http.Handler) *Function {
	routerFn := newFunction(http.MethodGet, pattern, handler, r)
	r.handlers = append(r.handlers, routerFn)
	return routerFn
}

func (r *Routes) POST(pattern string, handler http.Handler) *Function {
	routerFn := newFunction(http.MethodPost, pattern, handler, r)
	r.handlers = append(r.handlers, routerFn)
	return routerFn
}

func (r *Routes) DELETE(pattern string, handler http.Handler) *Function {
	routerFn := newFunction(http.MethodDelete, pattern, handler, r)
	r.handlers = append(r.handlers, routerFn)
	return routerFn
}

func (r *Routes) PATCH(pattern string, handler http.Handler) *Function {
	routerFn := newFunction(http.MethodPatch, pattern, handler, r)
	r.handlers = append(r.handlers, routerFn)
	return routerFn
}

func (r *Routes) HEAD(pattern string, handler http.Handler) *Function {
	routerFn := newFunction(http.MethodHead, pattern, handler, r)
	r.handlers = append(r.handlers, routerFn)
	return routerFn
}

func (r *Routes) OPTIONS(pattern string, handler http.Handler) *Function {
	routerFn := newFunction(http.MethodOptions, pattern, handler, r)
	r.handlers = append(r.handlers, routerFn)
	return routerFn
}

func (r *Routes) Accept(mediaTypes ...string) *Routes {
	r.accepts = append(r.accepts, mediaTypes...)
	return r
}

func (r *Routes) ContentType(mediaTypes ...string) *Routes {
	r.contentTypes = append(r.contentTypes, mediaTypes...)
	return r
}

func (r *Routes) Attribute(name string, val any) *Routes {
	r.attributes[name] = val
	return r
}

func (r *Routes) AttributeMap(attrs map[string]any) *Routes {
	for key, val := range attrs {
		r.attributes[key] = val
	}

	return r
}

func (r *Routes) UseMiddleware(function middleware.Function, options ...middleware.Option) *Routes {
	if function != nil {
		r.middlewares = append(r.middlewares, middleware.New(function, options...))
	}

	return r
}
