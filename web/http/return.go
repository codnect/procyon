package http

import (
	"reflect"
)

type ReturnValueHandler interface {
	SupportsReturnType(returnType reflect.Type) bool
	HandleReturnValue(ctx Context, returnValue any) error
}
