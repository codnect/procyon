package container

import (
	"fmt"
	"github.com/procyon-projects/reflector"
)

type Option func(def *Definition) error

func Name(name string) Option {
	return func(def *Definition) error {
		def.name = name
		return nil
	}
}

func Optional[T any]() Option {
	return func(def *Definition) error {
		typ := reflector.TypeOf[T]()

		exists := false
		for _, input := range def.inputs {
			if input.reflectorType().Compare(typ) {
				input.optional = true
				exists = true
			}
		}

		if !exists {
			return fmt.Errorf("container: could not find any input of type %s", typ.Name())
		}

		return nil
	}
}

func OptionalAt(index int) Option {
	return func(def *Definition) error {
		if index < 0 {
			panic(fmt.Sprintf("container: index should be greater than or equal to zero, but got index %d", index))
		}

		if len(def.inputs) <= index {
			return fmt.Errorf("container: could not find any input at index %d", index)
		}

		def.inputs[index].optional = true
		return nil
	}
}

func Qualifier[T any](name string) Option {
	return func(def *Definition) error {
		typ := reflector.TypeOf[T]()

		exists := false
		for _, input := range def.inputs {
			if input.reflectorType().Compare(typ) {
				input.name = name
				exists = true
			}
		}

		if !exists {
			return fmt.Errorf("container: could not find any input of type %s", typ.Name())
		}

		return nil
	}
}

func QualifierAt(index int, name string) Option {
	return func(def *Definition) error {
		if index < 0 {
			panic(fmt.Sprintf("container: index should be greater than or equal to zero, but got index %d", index))
		}

		if len(def.inputs) <= index {
			return fmt.Errorf("container: could not find any input at index %d", index)
		}

		def.inputs[index].name = name
		return nil
	}
}

func Scoped(scope string) Option {
	return func(def *Definition) error {
		def.scope = scope
		return nil
	}
}
