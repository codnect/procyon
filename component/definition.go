package component

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

// DefinitionOption is a functional option used to configure a Definition.
type DefinitionOption func(def *Definition) error

// Definition represents the metadata and constructor for a component.
type Definition struct {
	name        string
	scope       string
	typ         reflect.Type
	constructor Constructor
	attrs       map[string]any
}

// Name returns the name of the definition.
func (d *Definition) Name() string {
	return d.name
}

// Scope returns the scope of the definition (e.g. singleton or prototype).
func (d *Definition) Scope() string {
	return d.scope
}

// IsSingleton returns true if the definition is a singleton scoped component.
func (d *Definition) IsSingleton() bool {
	return d.scope == SingletonScope
}

// IsPrototype returns true if the definition is a prototype scoped component.
func (d *Definition) IsPrototype() bool {
	return d.scope == PrototypeScope
}

// Type returns the reflect.Type of the component the definition produces.
func (d *Definition) Type() reflect.Type {
	return d.typ
}

// Constructor returns the constructor metadata used to build the component.
func (d *Definition) Constructor() Constructor {
	return d.constructor
}

// MakeDefinition creates a new definition with the provided constructor function and options.
func MakeDefinition(fn ConstructorFunc, opts ...DefinitionOption) (*Definition, error) {
	if fn == nil {
		return nil, fmt.Errorf("nil constructor")
	}

	// Check if the constructor function is a function
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		return nil, fmt.Errorf("constructor must be a function")
	}

	// Check if the constructor function returns only one result
	if fnType.NumOut() != 1 {
		return nil, fmt.Errorf("constructor must only be a function returning one result")
	}

	// Get the return type of the constructor function
	outType := fnType.Out(0)

	// Get the name of the definition
	definitionName := generateDefinitionName(outType)

	// Create a new definition
	definition := &Definition{
		name:  definitionName,
		typ:   outType,
		scope: SingletonScope,
		constructor: Constructor{
			funcType:  outType,
			funcValue: reflect.ValueOf(fn),
			args:      make([]Arg, 0),
		},
	}

	// Populate constructor arguments
	err := populateConstructorArgs(definition, fnType)
	if err != nil {
		return nil, err
	}

	err = applyDefinitionOpts(definition, opts)
	if err != nil {
		return nil, err
	}

	return definition, nil
}

// applyDefinitionOptions applies the options to the definition
func applyDefinitionOpts(def *Definition, opts []DefinitionOption) error {
	for _, opt := range opts {
		err := opt(def)
		if err != nil {
			return err
		}
	}
	return nil
}

// populateConstructorArguments populates constructor arguments for the given definition and constructor type.
func populateConstructorArgs(def *Definition, typ reflect.Type) error {
	numIn := typ.NumIn()
	constructor := def.constructor

	for index := 0; index < numIn; index++ {
		argType := typ.In(index)

		arg := Arg{
			index: index,
			typ:   argType,
		}

		constructor.args = append(constructor.args, arg)
	}

	return nil
}

// WithName sets the name of the component definition.
func WithName(name string) DefinitionOption {
	return func(def *Definition) error {
		def.name = name
		return nil
	}
}

// WithScope sets a custom scope string for the component definition.
func WithScope(scope string) DefinitionOption {
	return func(def *Definition) error {
		def.scope = scope
		return nil
	}
}

// AsSingleton marks the component definition to use singleton scope.
func AsSingleton() DefinitionOption {
	return func(def *Definition) error {
		def.scope = SingletonScope
		return nil
	}
}

// AsPrototype marks the component definition to use prototype scope.
func AsPrototype() DefinitionOption {
	return func(def *Definition) error {
		def.scope = PrototypeScope
		return nil
	}
}

// WithQualifierFor sets a named qualifier for the constructor argument
// that matches the given type T.
func WithQualifierFor[T any](name string) DefinitionOption {
	return func(def *Definition) error {
		typ := reflect.TypeFor[T]()
		objectConstructor := def.constructor

		exists := false
		for index, arg := range objectConstructor.Args() {
			if arg.Type() == typ {
				objectConstructor.args[index].name = name
				exists = true
			}
		}

		if !exists {
			return fmt.Errorf("cannot find any input of type %s", typ.Name())
		}

		return nil
	}
}

// WithQualifierAt sets a named qualifier for the constructor argument
// at the specified index.
func WithQualifierAt(index int, name string) DefinitionOption {
	return func(def *Definition) error {
		if index < 0 {
			panic(fmt.Sprintf("index should be greater than or equal to zero, but got index %d", index))
		}

		objectConstructor := def.constructor
		if len(objectConstructor.Args()) <= index {
			return fmt.Errorf("cannot find any input at index %d", index)
		}

		objectConstructor.args[index].name = name
		return nil
	}
}

// generateDefinitionName returns the name of the definition based on the return type of the constructor function
func generateDefinitionName(returnType reflect.Type) string {
	if returnType.Kind() == reflect.Pointer {
		return lowerCamelCase(returnType.Elem().Name())
	}
	return lowerCamelCase(returnType.Name())
}

func lowerCamelCase(str string) string {
	isFirst := true

	return strings.Map(func(r rune) rune {
		if isFirst {
			isFirst = false
			return unicode.ToLower(r)
		}

		return r
	}, str)
}
