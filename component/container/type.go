package container

import (
	"reflect"
)

// convertibleTo checks if a source type can be converted to a target type.
func convertibleTo(sourceType reflect.Type, targetType reflect.Type) bool {
	if sourceType == targetType || (targetType.Kind() == reflect.Interface && sourceType.ConvertibleTo(targetType)) {
		return true
	} else if sourceType.Kind() == reflect.Pointer {
		return convertibleTo(sourceType.Elem(), targetType)
	}

	return false
}
