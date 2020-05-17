package web

type HandlerMethod struct {
	Path        string
	Method      HttpMethod
	HandlerFunc func(...interface{})
}

func NewHandlerMethod(path string, handler func(...interface{})) *HandlerMethod {
	return &HandlerMethod{}
}

func (handlerMethod *HandlerMethod) WithGet() *HandlerMethod {
	handlerMethod.Method = HttpMethodGet
	return handlerMethod
}

func (handlerMethod *HandlerMethod) WithHead() *HandlerMethod {
	handlerMethod.Method = HttpMethodHead
	return handlerMethod
}

func (handlerMethod *HandlerMethod) WithPost() *HandlerMethod {
	handlerMethod.Method = HttpMethodPost
	return handlerMethod
}

func (handlerMethod *HandlerMethod) WithPut() *HandlerMethod {
	handlerMethod.Method = HttpMethodPut
	return handlerMethod
}

func (handlerMethod *HandlerMethod) WithPatch() *HandlerMethod {
	handlerMethod.Method = HttpMethodPatch
	return handlerMethod
}
