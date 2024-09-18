package container

import (
	"fmt"
	"reflect"
	"strings"
)

// DefinitionOption is a function type that modifies a Definition.
type DefinitionOption func(definition *Definition) error

// Named sets the name of the object.
func Named(name string) DefinitionOption {
	return func(definition *Definition) error {
		if strings.TrimSpace(name) != "" {
			definition.name = name
		}

		return nil
	}
}

// Scoped sets the scope of the object.
func Scoped(scope string) DefinitionOption {
	return func(definition *Definition) error {
		if strings.TrimSpace(scope) == "" {
			definition.scope = SingletonScope
		} else {
			definition.scope = scope
		}

		return nil
	}
}

// Prioritized sets the priority of the object.
func Prioritized(priority int) DefinitionOption {
	return func(definition *Definition) error {
		definition.priority = priority
		return nil
	}
}

// Qualifier sets the name of the constructor's input parameter.
func Qualifier[T any](name string) DefinitionOption {
	return func(definition *Definition) error {
		typ := reflect.TypeFor[T]()
		objectConstructor := definition.constructor

		exists := false
		for index, arg := range objectConstructor.Arguments() {
			if arg.Type() == typ {
				objectConstructor.arguments[index].name = name
				exists = true
			}
		}

		if !exists {
			return fmt.Errorf("cannot find any input of type %s", typ.Name())
		}

		return nil
	}
}

// QualifierAt sets the name of the constructor's input parameter at the given index.
func QualifierAt(index int, name string) DefinitionOption {
	return func(definition *Definition) error {
		if index < 0 {
			panic(fmt.Sprintf("index should be greater than or equal to zero, but got index %d", index))
		}

		objectConstructor := definition.constructor
		if len(objectConstructor.Arguments()) <= index {
			return fmt.Errorf("cannot find any input at index %d", index)
		}

		objectConstructor.arguments[index].name = name
		return nil
	}
}

// Optional sets the constructor's input parameter as optional.
func Optional[T any]() DefinitionOption {
	return func(definition *Definition) error {
		typ := reflect.TypeFor[T]()
		objectConstructor := definition.constructor

		exists := false
		for index, arg := range objectConstructor.Arguments() {
			if arg.Type() == typ {
				objectConstructor.arguments[index].optional = true
				exists = true
			}
		}

		if !exists {
			return fmt.Errorf("cannot find any input of type %s", typ.Name())
		}

		return nil
	}
}

// OptionalAt sets the constructor's input parameter at the given index as optional.
func OptionalAt(index int) DefinitionOption {
	return func(definition *Definition) error {
		if index < 0 {
			panic(fmt.Sprintf("index should be greater than or equal to zero, but got index %d", index))
		}

		objectConstructor := definition.constructor
		if len(objectConstructor.Arguments()) <= index {
			return fmt.Errorf("cannot find any input at index %d", index)
		}

		objectConstructor.arguments[index].optional = true
		return nil
	}
}
