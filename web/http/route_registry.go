package http

import (
	"errors"
	"fmt"
)

type RouteRegistry struct {
	tree []*routeTree
}

func NewRouteRegistry() *RouteRegistry {
	registry := &RouteRegistry{
		make([]*routeTree, 9),
	}

	methods := []Method{
		MethodGet,
		MethodHead,
		MethodPost,
		MethodPut,
		MethodPatch,
		MethodDelete,
		MethodConnect,
		MethodOptions,
		MethodTrace,
	}

	for _, method := range methods {
		registry.tree[method.IntValue()] = &routeTree{
			staticRoutes: make(map[string]*Route, 0),
		}
	}

	return registry
}

func (r *RouteRegistry) Register(route *Route) error {
	methods := route.Methods()

	if len(methods) == 0 {
		return errors.New("route must have at least one method")
	}

	for _, method := range methods {
		intValue := method.IntValue()
		if intValue == -1 {
			return fmt.Errorf("invalid method: %s", method)
		}

		methodTree := r.tree[intValue]

		if methodTree.children == nil {
			methodTree.children = &routeNode{}
		}

		methodTree.add(route)
	}

	return nil
}

func (r *RouteRegistry) Find(ctx Context) (*Route, bool) {
	request := ctx.Request()
	path := request.Path()
	method := request.Method()

	intValue := method.IntValue()
	if intValue < 0 || intValue >= len(r.tree) {
		return nil, false
	}

	methodTree := r.tree[intValue]

	if route, ok := methodTree.staticRoutes[path]; ok {
		return route, true
	}

	route := methodTree.match(ctx)
	return route, true
}

func (r *RouteRegistry) List() []*Route {
	return nil
}

func (r *RouteRegistry) Unregister(route *Route) error {
	return nil
}
