package web

import "context"

type HandlerProcessor struct {
	mapping *HandlerMapping
}

func (p *HandlerProcessor) Process() {
	var x any

	if handler, ok := x.(Handler); ok {
		routeGroup := handler.Routes()
		if routeGroup != nil {
			p.registerRouterGroups(routeGroup)
		}
	}
}

func (p *HandlerProcessor) registerRouterGroups(group *RouterGroup) {
	//path := group.FullPath()

	for _, routerFunction := range group.Functions() {
		handlerChain := &HandlerChain{}

		index := 0

		for _, interceptor := range routerFunction.getBeforeInterceptors() {
			handlerChain.interceptors = append(handlerChain.interceptors, func(ctx context.Context) (bool, error) {
				return interceptor(ctx.(*Context))
			})
			index++
		}

		handlerChain.interceptors = append(handlerChain.interceptors, handlerChain.function.wrapperFunc)
		handlerChain.handlerIndex = index

		for _, interceptor := range routerFunction.getAfterInterceptors() {
			handlerChain.interceptors = append(handlerChain.interceptors, func(ctx context.Context) (bool, error) {
				err := interceptor(ctx.(*Context))
				return true, err
			})
			index++
		}

		handlerChain.afterCompletionIndex = index

		for _, interceptor := range routerFunction.getAfterCompletionInterceptors() {
			handlerChain.interceptors = append(handlerChain.interceptors, func(ctx context.Context) (bool, error) {
				err := interceptor(ctx.(*Context))
				return true, err
			})
		}
	}

	for _, nestedGroup := range group.NestedGroups() {
		p.registerRouterGroups(nestedGroup)
	}
}
