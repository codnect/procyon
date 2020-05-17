package support

import "procyon/web"

type HandlerMethodRegistry struct {
	registerMap map[string][]*web.HandlerMethod
}

func NewHandlerMethodRegistry() HandlerMethodRegistry {
	return HandlerMethodRegistry{
		registerMap: make(map[string][]*web.HandlerMethod),
	}
}

func (registry HandlerMethodRegistry) Register(handlerMethod *web.HandlerMethod) {
	registry.RegisterGroup("", handlerMethod)
}

func (registry HandlerMethodRegistry) RegisterGroup(groupName string, handlerMethod *web.HandlerMethod) {
	if registry.registerMap[groupName] == nil {
		registry.registerMap[groupName] = make([]*web.HandlerMethod, 0)
	}
	registry.registerMap[groupName] = append(registry.registerMap[""], handlerMethod)
}
