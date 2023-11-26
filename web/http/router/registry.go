package router

import (
	"codnect.io/procyon/web/http"
	"codnect.io/procyon/web/http/middleware"
)

type MappingRegistry struct {
}

/*type SimpleMappingRegistry struct {
	//routes         []*tree
	returnValueHandlers []http.ReturnValueHandler
}*/

func newMappingRegistry(returnValueHandlers []http.ReturnValueHandler) *MappingRegistry {
	return nil
}

func (r *MappingRegistry) GetHandler(ctx http.Context) (http.HandlerChain, bool) {
	return nil, false
}
func (r *MappingRegistry) Handlers() map[*Mapping]any {
	return nil
}

func (r *MappingRegistry) Register(mapping *Mapping, handler http.Handler, middlewares ...*middleware.Middleware) {
}

func (r *MappingRegistry) Unregister(mapping *Mapping) {

}

/*
func NewSimpleHandlerRegistry() *SimpleHandlerRegistry {
	mapping := &SimpleHandlerRegistry{
		routes: make([]*tree, 9),
	}

	mapping.createMappingTree(MethodGet)
	mapping.createMappingTree(MethodHead)
	mapping.createMappingTree(MethodPost)
	mapping.createMappingTree(MethodPut)
	mapping.createMappingTree(MethodPatch)
	mapping.createMappingTree(MethodDelete)
	mapping.createMappingTree(MethodConnect)
	mapping.createMappingTree(MethodOptions)
	mapping.createMappingTree(MethodTrace)
	return mapping
}

func (hm *SimpleHandlerRegistry) createMappingTree(method Method) {
	mappingTree := &tree{}
	hm.routes[method.IntValue()] = mappingTree
}

func (hm *SimpleHandlerRegistry) GetHandler(ctx *Context) *HandlerChain {
	request := ctx.Request()
	path := request.Path()
	method := request.Method()
	mappingTree := hm.routes[method.IntValue()]

	if mappingTree == nil {
		return nil
	}

	if mappingTree.staticRoutes != nil {
		if _, ok := mappingTree.staticRoutes[path]; ok {
			return nil
		}
	}

	chain := mappingTree.getHandlerChain(ctx)

	if chain == nil {
		return &HandlerChain{}
	}

	return chain
}

func (hm *SimpleHandlerRegistry) Handlers() map[*RequestMapping]any {
	return nil
}

func (hm *SimpleHandlerRegistry) Register(mapping *RequestMapping, handler Function, middlewares ...*middleware.Middleware) {
	methods := mapping.Methods()
	for _, method := range methods {
		mappingTree := hm.routes[method.IntValue()]

		if mappingTree.root == nil {
			mappingTree.root = &treeNode{}
		}

		mappingTree.addMapping(mapping, &HandlerChain{value: mapping.pattern,
			pathVariableNameMap:  map[string]int{},
			pathVariableIndexMap: map[int]string{},
		})
	}
}

func (hm *SimpleHandlerRegistry) Unregister(mapping *RequestMapping) {

}

*/
