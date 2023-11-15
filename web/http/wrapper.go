package http

import "reflect"

type HandlerWrapper struct {
	handler             Handler
	returnValueHandlers []ReturnValueHandler
}

func WrapHandler(handler Handler) *HandlerWrapper {
	return &HandlerWrapper{
		handler: handler,
	}
}

func (w *HandlerWrapper) SetReturnValueHandlers(handlers []ReturnValueHandler) {
	w.returnValueHandlers = handlers
}

func (w *HandlerWrapper) Unwrap() Handler {
	return w.handler
}

func (w *HandlerWrapper) Invoke(ctx Context, next RequestDelegate) error {
	returnValue, err := w.handler.Invoke(ctx)

	if err != nil {
		return err
	}

	returnType := reflect.TypeOf(returnValue)

	for _, handler := range w.returnValueHandlers {
		if handler.SupportsReturnType(returnType) {
			return handler.HandleReturnValue(ctx, next)
		}
	}

	return nil
}
