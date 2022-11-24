package web

import (
	"fmt"
	"net/http"
	"procyon-test/web/mediatype"
)

type Router struct {
	mapping *HandlerMapping
}

func newRouter() *Router {
	return &Router{}
}

func (r *Router) Route(writer http.ResponseWriter, request *http.Request) {
	var (
		handlerChain   *HandlerChain
		exists         bool
		requestContext any
	)

	defer func() {
		if r := recover(); r != nil {

		}
		if handlerChain != nil && requestContext != nil {
			handlerChain.function.putToPool(requestContext)
		}
	}()

	handlerChain, exists = r.mapping.FindHandlerChain(request.URL.Path, HttpMethod(request.Method))

	if !exists {
		return
	}

	requestContext = handlerChain.function.getOrCreateContext()
	baseContext := handlerChain.function.toContext(requestContext)
	baseContext.reset(writer, request)

	for index, handler := range handlerChain.interceptors {
		if index == handlerChain.afterCompletionIndex {
			baseContext.writeResponse()
		}

		next, err := handler(baseContext)

		if !next {
			baseContext.writeResponse()
			return
		}

		if err != nil {
			baseContext.err = err
			baseContext.writeResponse()
			return
		}

		if index == handlerChain.handlerIndex && index == len(handlerChain.interceptors) {
			baseContext.writeResponse()
		}
	}
}

type RouterFunction struct {
	path       string
	fullPath   string
	method     HttpMethod
	handler    *HandlerFunction
	consumes   []mediatype.MediaType
	produces   []mediatype.MediaType
	attributes map[any]any

	beforeInterceptors          []func(ctx *Context) (bool, error)
	afterInterceptors           []func(ctx *Context) error
	afterCompletionInterceptors []func(ctx *Context) error

	group *RouterGroup
}

func newRouterFunction(method HttpMethod, basePath string, path string, handler *HandlerFunction, group *RouterGroup) *RouterFunction {
	fullPath := ""

	if basePath != "" && path != "" {
		fullPath = fmt.Sprintf("%s%s", basePath, path)
	}

	return &RouterFunction{
		fullPath:   fullPath,
		path:       path,
		method:     method,
		handler:    handler,
		attributes: map[any]any{},
		group:      group,
	}
}

func (f *RouterFunction) FullPath() string {
	return f.fullPath
}

func (f *RouterFunction) Path() string {
	return f.path
}

func (f *RouterFunction) Method() HttpMethod {
	return f.method
}

func (f *RouterFunction) Handler() *HandlerFunction {
	return f.handler
}

func (f *RouterFunction) Consumes() []mediatype.MediaType {
	consumes := make([]mediatype.MediaType, 0)

	if f.group != nil {
		groupConsumes := f.group.Consumes()
		if groupConsumes != nil {
			consumes = append(consumes, groupConsumes...)
		}
	}

	consumes = append(consumes, f.consumes...)

	return consumes
}

func (f *RouterFunction) Produces() []mediatype.MediaType {
	produces := make([]mediatype.MediaType, 0)

	if f.group != nil {
		groupProduces := f.group.Produces()
		if groupProduces != nil {
			produces = append(produces, groupProduces...)
		}
	}

	produces = append(produces, f.consumes...)

	return produces
}

func (f *RouterFunction) Attributes() map[any]any {
	attributes := make(map[any]any)

	if f.group != nil {
		groupAttributes := f.group.Attributes()

		for k, v := range groupAttributes {
			attributes[k] = v
		}
	}

	for k, v := range f.attributes {
		attributes[k] = v
	}

	return attributes
}

func (f *RouterFunction) getBeforeInterceptors() []func(ctx *Context) (bool, error) {
	interceptors := make([]func(ctx *Context) (bool, error), 0)

	if f.group != nil {
		groupInterceptors := f.group.getBeforeInterceptors()
		if groupInterceptors != nil {
			interceptors = append(interceptors, groupInterceptors...)
		}
	}

	interceptors = append(interceptors, f.beforeInterceptors...)
	return interceptors
}

func (f *RouterFunction) getAfterInterceptors() []func(ctx *Context) error {
	interceptors := make([]func(ctx *Context) error, 0)

	if f.group != nil {
		groupInterceptors := f.group.getAfterInterceptors()
		if groupInterceptors != nil {
			interceptors = append(interceptors, groupInterceptors...)
		}
	}

	interceptors = append(interceptors, f.afterInterceptors...)
	return interceptors
}

func (f *RouterFunction) getAfterCompletionInterceptors() []func(ctx *Context) error {
	interceptors := make([]func(ctx *Context) error, 0)

	if f.group != nil {
		groupInterceptors := f.group.getAfterCompletionInterceptors()
		if groupInterceptors != nil {
			interceptors = append(interceptors, groupInterceptors...)
		}
	}

	interceptors = append(interceptors, f.afterCompletionInterceptors...)
	return interceptors
}

func (f *RouterFunction) Consume(mediaTypes ...mediatype.MediaType) *RouterFunction {
	f.consumes = append(f.consumes, mediaTypes...)
	return f
}

func (f *RouterFunction) Produce(mediaTypes ...mediatype.MediaType) *RouterFunction {
	f.produces = append(f.produces, mediaTypes...)
	return f
}

func (f *RouterFunction) Attribute(name, val any) *RouterFunction {
	f.attributes[name] = val
	return f
}

func (f *RouterFunction) Before(interceptor func(ctx *Context) (bool, error)) *RouterFunction {
	f.beforeInterceptors = append(f.beforeInterceptors, interceptor)
	return f
}

func (f *RouterFunction) After(interceptor func(ctx *Context) error) *RouterFunction {
	f.afterInterceptors = append(f.afterInterceptors, interceptor)
	return f
}

func (f *RouterFunction) AfterCompletion(interceptor func(ctx *Context) error) *RouterFunction {
	f.afterCompletionInterceptors = append(f.afterCompletionInterceptors, interceptor)
	return f
}

type RouterGroup struct {
	path         string
	fullPath     string
	consumes     []mediatype.MediaType
	produces     []mediatype.MediaType
	nestedGroups []*RouterGroup
	attributes   map[any]any

	functions                   []*RouterFunction
	beforeInterceptors          []func(ctx *Context) (bool, error)
	afterInterceptors           []func(ctx *Context) error
	afterCompletionInterceptors []func(ctx *Context) error

	parentGroup *RouterGroup
}

func Routes(options ...RouteGroupOption) *RouterGroup {
	group := &RouterGroup{
		attributes: map[any]any{},
	}

	for _, option := range options {
		option(group)
	}

	group.fullPath = group.path

	return group
}

func (g *RouterGroup) Path() string {
	return g.path
}

func (g *RouterGroup) FullPath() string {
	if g.fullPath == "" {
		return g.fullPath
	}

	return g.fullPath
}

func (g *RouterGroup) Consumes() []mediatype.MediaType {
	consumes := make([]mediatype.MediaType, 0)

	if g.parentGroup != nil {
		parentConsumes := g.parentGroup.Consumes()
		if parentConsumes != nil {
			consumes = append(consumes, parentConsumes...)
		}
	}

	consumes = append(consumes, g.consumes...)

	return consumes
}

func (g *RouterGroup) Produces() []mediatype.MediaType {
	produces := make([]mediatype.MediaType, 0)

	if g.parentGroup != nil {
		parentProduces := g.parentGroup.Produces()
		if parentProduces != nil {
			produces = append(produces, parentProduces...)
		}
	}

	produces = append(produces, g.produces...)

	return produces
}

func (g *RouterGroup) NestedGroups() []*RouterGroup {
	return g.nestedGroups
}

func (g *RouterGroup) Attributes() map[any]any {
	attributes := make(map[any]any)

	if g.parentGroup != nil {
		parentAttributes := g.parentGroup.Attributes()

		for k, v := range parentAttributes {
			attributes[k] = v
		}
	}

	for k, v := range g.attributes {
		attributes[k] = v
	}

	return attributes
}

func (g *RouterGroup) Functions() []*RouterFunction {
	functions := make([]*RouterFunction, 0)

	for _, f := range g.functions {
		functions = append(functions, f)
	}

	return functions
}

func (g *RouterGroup) getBeforeInterceptors() []func(ctx *Context) (bool, error) {
	interceptors := make([]func(ctx *Context) (bool, error), 0)

	if g.parentGroup != nil {
		parentInterceptors := g.parentGroup.getBeforeInterceptors()
		if parentInterceptors != nil {
			interceptors = append(interceptors, parentInterceptors...)
		}
	}

	interceptors = append(interceptors, g.beforeInterceptors...)
	return interceptors
}

func (g *RouterGroup) getAfterInterceptors() []func(ctx *Context) error {
	interceptors := make([]func(ctx *Context) error, 0)

	if g.parentGroup != nil {
		parentInterceptors := g.parentGroup.getAfterInterceptors()
		if parentInterceptors != nil {
			interceptors = append(interceptors, parentInterceptors...)
		}
	}

	interceptors = append(interceptors, g.afterInterceptors...)
	return interceptors
}

func (g *RouterGroup) getAfterCompletionInterceptors() []func(ctx *Context) error {
	interceptors := make([]func(ctx *Context) error, 0)

	if g.parentGroup != nil {
		parentInterceptors := g.parentGroup.getAfterCompletionInterceptors()
		if parentInterceptors != nil {
			interceptors = append(interceptors, parentInterceptors...)
		}
	}

	interceptors = append(interceptors, g.afterCompletionInterceptors...)
	return interceptors
}

func (g *RouterGroup) Nest(options ...RouteGroupOption) *RouterGroup {
	nestedGroup := Routes(options...)
	nestedGroup.parentGroup = g
	nestedGroup.fullPath = fmt.Sprintf("%s%s", g.FullPath(), nestedGroup.path)

	g.nestedGroups = append(g.nestedGroups, nestedGroup)
	return nestedGroup
}

func (g *RouterGroup) Handler(method HttpMethod, handler *HandlerFunction) *RouterFunction {
	f := newRouterFunction(method, g.FullPath(), "", handler, g)
	g.functions = append(g.functions, f)
	return f
}

func (g *RouterGroup) GET(path string, handler *HandlerFunction) *RouterFunction {
	f := newRouterFunction(MethodGet, g.FullPath(), path, handler, g)
	g.functions = append(g.functions, f)
	return f
}

func (g *RouterGroup) POST(path string, handler *HandlerFunction) *RouterFunction {
	f := newRouterFunction(MethodPost, g.FullPath(), path, handler, g)
	g.functions = append(g.functions, f)
	return f
}

func (g *RouterGroup) DELETE(path string, handler *HandlerFunction) *RouterFunction {
	f := newRouterFunction(MethodDelete, g.FullPath(), path, handler, g)
	g.functions = append(g.functions, f)
	return f
}

func (g *RouterGroup) PATCH(path string, handler *HandlerFunction) *RouterFunction {
	f := newRouterFunction(MethodPatch, g.FullPath(), path, handler, g)
	g.functions = append(g.functions, f)
	return f
}

func (g *RouterGroup) HEAD(path string, handler *HandlerFunction) *RouterFunction {
	f := newRouterFunction(MethodHead, g.FullPath(), path, handler, g)
	g.functions = append(g.functions, f)
	return f
}
func (g *RouterGroup) OPTIONS(path string, handler *HandlerFunction) *RouterFunction {
	f := newRouterFunction(MethodOptions, g.FullPath(), path, handler, g)
	g.functions = append(g.functions, f)
	return f
}

func (g *RouterGroup) Consume(mediaTypes ...mediatype.MediaType) *RouterGroup {
	g.consumes = append(g.consumes, mediaTypes...)
	return g
}

func (g *RouterGroup) Produce(mediaTypes ...mediatype.MediaType) *RouterGroup {
	g.produces = append(g.produces, mediaTypes...)
	return g
}

func (g *RouterGroup) Attribute(name, val any) *RouterGroup {
	g.attributes[name] = val
	return g
}

func (g *RouterGroup) Before(interceptor func(ctx *Context) (bool, error)) *RouterGroup {
	g.beforeInterceptors = append(g.beforeInterceptors, interceptor)
	return g
}

func (g *RouterGroup) After(interceptor func(ctx *Context) error) *RouterGroup {
	g.afterInterceptors = append(g.afterInterceptors, interceptor)
	return g
}

func (g *RouterGroup) AfterCompletion(interceptor func(ctx *Context) error) *RouterGroup {
	g.afterCompletionInterceptors = append(g.afterCompletionInterceptors, interceptor)
	return g
}

type RouteGroupOption func(group *RouterGroup)

func Path(path string) RouteGroupOption {
	return func(group *RouterGroup) {
		group.path = path
	}
}
