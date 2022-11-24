package web

type HandlerMapping struct {
	routerTree *RouterTree
}

func newRouterMapping() *HandlerMapping {
	return &HandlerMapping{
		routerTree: newRouterTree(),
	}
}

func (m *HandlerMapping) FindHandlerChain(path string, method HttpMethod) (*HandlerChain, bool) {
	return nil, false
}
