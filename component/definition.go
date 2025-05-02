// Copyright 2025 Codnect
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package component

import (
	"fmt"
	"reflect"
)

// DefinitionOption is a functional option used to configure a Definition.
type DefinitionOption func(def *Definition) error

// DefinitionRegistry defines methods for managing component definitions within the system.
type DefinitionRegistry interface {
	// RegisterDefinition registers a new component definition.
	// Returns an error if a definition with the same name already exists.
	RegisterDefinition(def *Definition) error

	// UnregisterDefinition removes the component definition associated with the given name.
	// Returns an error if the definition does not exist.
	UnregisterDefinition(name string) error

	// Definition retrieves the component definition associated with the given name.
	// Returns the definition and a boolean indicating its existence.
	Definition(name string) (*Definition, bool)

	// ContainsDefinition checks whether a component definition with the specified name exists.
	ContainsDefinition(name string) bool

	// Definitions returns a slice of all registered component definitions.
	Definitions() []*Definition

	// DefinitionsOf returns a slice of component definitions that are assignable to the specified type.
	DefinitionsOf(typ reflect.Type) []*Definition
}

// Definition represents the metadata and constructor for a component.
type Definition struct {
	name        string
	scope       string
	constructor Constructor
}

// Name returns the name of the definition.
func (d *Definition) Name() string {
	return d.name
}

// Scope returns the scope of the definition (e.g. singleton or prototype).
func (d *Definition) Scope() string {
	return d.scope
}

// IsSingleton returns true if the definition is a singleton-scoped component.
func (d *Definition) IsSingleton() bool {
	return d.scope == SingletonScope
}

// IsPrototype returns true if the definition is a prototype-scoped component.
func (d *Definition) IsPrototype() bool {
	return d.scope == PrototypeScope
}

// Type returns the reflect.Type of the component the definition produces.
func (d *Definition) Type() reflect.Type {
	return d.constructor.OutType()
}

// Constructor returns the constructor metadata used to build the component.
func (d *Definition) Constructor() Constructor {
	return d.constructor
}

// MakeDefinition creates a new definition with the provided constructor function and options.
func MakeDefinition(fn ConstructorFunc, opts ...DefinitionOption) (*Definition, error) {
	constructor, err := createConstructor(fn)
	if err != nil {
		return nil, err
	}

	// Get the return type of the constructor function
	outType := constructor.OutType()

	// Generate component name
	componentName := generateComponentName(outType)

	// Create a new definition
	def := &Definition{
		name:        componentName,
		scope:       SingletonScope,
		constructor: constructor,
	}

	err = applyDefinitionOpts(def, opts)
	if err != nil {
		return nil, err
	}

	return def, nil
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
