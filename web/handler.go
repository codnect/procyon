package web

import (
	"context"
	"github.com/procyon-projects/reflector"
	"reflect"
	"strconv"
	"sync"
)

type Handler interface {
	Routes() *RouterGroup
}

type Variable struct {
	Index    int
	Name     string
	Required bool
	Type     reflect.Type
	Default  any
}

type Variables struct {
	Body            Variable
	QueryParameters []Variable
	PathVariables   []Variable
	Headers         []Variable
	FormFields      []Variable
}

type HandlerFunction struct {
	pool *sync.Pool

	wrapperFunc  func(ctx context.Context) (bool, error)
	toContextFun func(ctx any) *Context

	variables *Variables
}

func (f *HandlerFunction) getOrCreateContext() any {
	return f.pool.Get()
}

func (f *HandlerFunction) putToPool(ctx any) {
	f.pool.Put(ctx)
}

func (f *HandlerFunction) toContext(ctx any) *Context {
	return f.toContextFun(ctx)
}

func (f *HandlerFunction) invoke(ctx context.Context) (bool, error) {
	return f.wrapperFunc(ctx)
}

type RestHandlerFun[E, R any] interface {
	func(ctx *RestContext[E, R]) error
}

type MvcHandlerFun[E, R any] interface {
	func(ctx *MvcContext[E, R]) error
}

func collectVariables[E, R any]() *Variables {
	var (
		zeroRequest E
	)

	variables := &Variables{}

	originalRequestType := reflect.TypeOf(zeroRequest)
	requestType := originalRequestType
	if requestType.Kind() == reflect.Pointer {
		requestType = requestType.Elem()
	}

	if requestType.Kind() == reflect.Struct {

		for i := 0; i < requestType.NumField(); i++ {
			field := requestType.Field(i)
			typeTag, ok := field.Tag.Lookup("type")

			if ok {
				variable := &Variable{
					Index: i,
					Type:  field.Type,
				}

				nameTag, ok := field.Tag.Lookup("name")

				if ok {
					variable.Name = nameTag
				}

				defaultTag, ok := field.Tag.Lookup("default")

				if ok {
					variable.Default = defaultTag
				}

				requiredTag, ok := field.Tag.Lookup("required")

				if ok {
					required, err := strconv.ParseBool(requiredTag)
					if err == nil {
						variable.Required = required
					}
				}

				switch typeTag {
				case "body":
					variable.Required = true
					variable.Name = ""
					variables.Body = *variable
				case "path":
					variable.Required = true
					variable.Default = nil
					variables.PathVariables = append(variables.PathVariables, *variable)
				case "query":
					variables.QueryParameters = append(variables.QueryParameters, *variable)
				case "header":
					variables.Headers = append(variables.Headers, *variable)
				case "form":
					variables.FormFields = append(variables.FormFields, *variable)
				}
			}

		}

		if len(variables.QueryParameters) == 0 && len(variables.PathVariables) == 0 && len(variables.FormFields) == 0 && len(variables.Headers) == 0 {
			variables.Body.Type = originalRequestType
		}
	}

	return variables
}

func RestHandler[E, R any, F RestHandlerFun[E, R]](fn F) *HandlerFunction {
	m := reflector.TypeOf[R]()
	if m == nil {

	}
	name := m.PackagePath()
	if name == "" {

	}

	handler := &HandlerFunction{
		pool: &sync.Pool{},
	}

	handler.pool.New = func() any {
		restContext := &RestContext[E, R]{}
		restContext.bodyBuilder.response = restContext.ctx.Response()
		restContext.headersBuilder.response = restContext.ctx.Response()
		return restContext
	}

	handler.toContextFun = func(ctx any) *Context {
		return &ctx.(*RestContext[E, R]).ctx
	}

	switch typedFunction := any(fn).(type) {
	case func(ctx *RestContext[E, R]) error:
		handler.wrapperFunc = func(ctx context.Context) (bool, error) {
			return true, typedFunction(ctx.(*RestContext[E, R]))
		}
	}

	handler.variables = collectVariables[E, R]()

	return handler
}

func MvcHandler[E, R any, F MvcHandlerFun[E, R]](fn F) *HandlerFunction {
	handler := &HandlerFunction{
		pool: &sync.Pool{},
	}

	handler.pool.New = func() any {
		mvcContext := &MvcContext[E, R]{}
		return mvcContext
	}

	handler.toContextFun = func(ctx any) *Context {
		return nil
	}

	switch typedFunction := any(fn).(type) {
	case func(ctx *MvcContext[E, R]) error:
		handler.wrapperFunc = func(ctx context.Context) (bool, error) {
			return true, typedFunction(ctx.(*MvcContext[E, R]))
		}
	}

	collectVariables[E, R]()

	return handler
}
