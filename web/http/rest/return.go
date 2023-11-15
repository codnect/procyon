package rest

import (
	"github.com/procyon-projects/procyon/web/http"
	"reflect"
)

var (
	responseEntityType = reflect.TypeOf((*ResponseEntity)(nil)).Elem()
)

type ResponseEntityHandler struct {
}

func (eh *ResponseEntityHandler) SupportsReturnType(returnType reflect.Type) bool {
	return returnType.ConvertibleTo(responseEntityType)
}

func (eh *ResponseEntityHandler) HandleReturnValue(ctx http.Context, returnValue any) error {

	switch responseEntity := returnValue.(type) {
	case ResponseEntity:
		return eh.handleResponseEntity(ctx, responseEntity)
	}
	return nil
}

func (eh *ResponseEntityHandler) handleResponseEntity(ctx http.Context, responseEntity ResponseEntity) error {
	return nil
}
