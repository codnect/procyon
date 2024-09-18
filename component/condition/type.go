package condition

import (
	"codnect.io/procyon/component/filter"
	"reflect"
)

// OnTypeCondition struct represents a condition that checks if a specific type exists.
type OnTypeCondition struct {
	typ reflect.Type // The type to check.
}

// OnType function creates a new OnTypeCondition.
func OnType[T any]() *OnTypeCondition {
	return &OnTypeCondition{
		typ: reflect.TypeFor[T](),
	}
}

// MatchesCondition method checks if the type exists.
// It retrieves the definitions and singletons from the container and checks if any of them match the type.
// If any match, it returns true. If none match, it returns false.
func (c *OnTypeCondition) MatchesCondition(ctx Context) bool {
	container := ctx.Container()
	if container == nil {
		return false
	}

	definitions := container.Definitions().List(filter.ByType(c.typ))
	singletons := container.Singletons().List(filter.ByType(c.typ))
	return len(definitions) != 0 || len(singletons) != 0
}

// OnMissingTypeCondition struct represents a condition that checks if a specific type does not exist.
type OnMissingTypeCondition struct {
	missingType reflect.Type // The type to check.
}

// OnMissingType function creates a new OnMissingTypeCondition.
func OnMissingType[T any]() *OnMissingTypeCondition {
	return &OnMissingTypeCondition{
		missingType: reflect.TypeFor[T](),
	}
}

// MatchesCondition method checks if the type does not exist.
// It retrieves the definitions and singletons from the container and checks if any of them match the type.
// If none match, it returns true. If any match, it returns false.
func (c *OnMissingTypeCondition) MatchesCondition(ctx Context) bool {
	container := ctx.Container()
	if container == nil {
		return false
	}

	definitions := container.Definitions().List(filter.ByType(c.missingType))
	singletons := container.Singletons().List(filter.ByType(c.missingType))
	return len(definitions) == 0 && len(singletons) == 0
}
