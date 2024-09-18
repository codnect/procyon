package container

import (
	"codnect.io/procyon/component/filter"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

// DefinitionRegistry interface provides methods for managing and accessing definitions.
type DefinitionRegistry interface {
	// Register a definition
	Register(definition *Definition) error
	// Remove a definition by name
	Remove(name string) error
	// Contains checks if a definition with the provided name exists
	Contains(name string) bool
	// Find returns the definition that matches the provided filters
	Find(filters ...filter.Filter) (*Definition, error)
	// FindFirst returns the first definition that matches the provided filters
	FindFirst(filters ...filter.Filter) (*Definition, bool)
	// List returns a list of definitions that match the provided filters
	List(filters ...filter.Filter) []*Definition
	// Names returns the names of all definitions
	Names() []string
	// Count returns the number of definitions
	Count() int
}

// Definition struct represents a definition of an object.
type Definition struct {
	name        string
	typ         reflect.Type
	scope       string
	priority    int
	constructor *Constructor
}

// MakeDefinition creates a new definition with the provided constructor function and options.
func MakeDefinition(constructorFunc ConstructorFunc, options ...DefinitionOption) (*Definition, error) {
	if constructorFunc == nil {
		return nil, fmt.Errorf("nil constructor function")
	}

	// Check if the constructor function is a function
	constructorType := reflect.TypeOf(constructorFunc)
	if constructorType.Kind() != reflect.Func {
		return nil, fmt.Errorf("constructor must be a function")
	}

	// Check if the constructor function returns only one result
	if constructorType.NumOut() != 1 {
		return nil, fmt.Errorf("constructor must only be a function returning one result")
	}

	// Get the return type of the constructor function
	returnType := constructorType.Out(0)

	// Get the name of the definition
	definitionName := getDefinitionName(returnType)

	// Create a new definition
	definition := createNewDefinition(definitionName, returnType, constructorType, constructorFunc)

	// Populate constructor arguments
	err := populateConstructorArguments(definition, constructorType)
	if err != nil {
		return nil, err
	}

	err = applyDefinitionOptions(definition, options)
	if err != nil {
		return nil, err
	}

	return definition, nil
}

// createNewDefinition creates a new definition
func createNewDefinition(definitionName string, returnType reflect.Type, constructorType reflect.Type, constructorFunc ConstructorFunc) *Definition {
	return &Definition{
		name:  definitionName,
		typ:   returnType,
		scope: SingletonScope,
		constructor: &Constructor{
			funcType:  constructorType,
			funcValue: reflect.ValueOf(constructorFunc),
			arguments: make([]ConstructorArgument, 0),
		},
	}
}

// Name returns the name of the object.
func (d *Definition) Name() string {
	return d.name
}

// Type returns the type of the object.
func (d *Definition) Type() reflect.Type {
	return d.typ
}

// Constructor returns the constructor of the object.
func (d *Definition) Constructor() *Constructor {
	return d.constructor
}

// Scope returns the scope of the object.
func (d *Definition) Scope() string {
	return d.scope
}

func (d *Definition) Priority() int {
	return d.priority
}

// IsSingleton checks if the definition is a singleton.
func (d *Definition) IsSingleton() bool {
	return d.scope == SingletonScope
}

// IsPrototype checks if the definition is a prototype.
func (d *Definition) IsPrototype() bool {
	return d.scope == PrototypeScope
}

// getDefinitionName returns the name of the definition based on the return type of the constructor function
func getDefinitionName(returnType reflect.Type) string {
	if returnType.Kind() == reflect.Pointer {
		return lowerCamelCase(returnType.Elem().Name())
	}
	return lowerCamelCase(returnType.Name())
}

// populateConstructorArguments populates constructor arguments for the given definition and constructor type.
func populateConstructorArguments(definition *Definition, constructorType reflect.Type) error {
	numIn := constructorType.NumIn()
	objectConstructor := definition.constructor

	for index := 0; index < numIn; index++ {
		argType := constructorType.In(index)

		arg := ConstructorArgument{
			index:    index,
			typ:      argType,
			optional: false,
		}

		objectConstructor.arguments = append(objectConstructor.arguments, arg)
	}

	return nil
}

// applyDefinitionOptions applies the options to the definition
func applyDefinitionOptions(definition *Definition, options []DefinitionOption) error {
	for _, option := range options {
		err := option(definition)
		if err != nil {
			return err
		}
	}
	return nil
}

// lowerCamelCase converts the first letter of the given string to lowercase.
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

// objectDefinitionRegistry struct represents a registry for object definitions.
type objectDefinitionRegistry struct {
	definitionMap map[string]*Definition
	muDefinitions *sync.RWMutex
}

// newObjectDefinitionRegistry creates a new object definition registry.
func newObjectDefinitionRegistry() *objectDefinitionRegistry {
	return &objectDefinitionRegistry{
		definitionMap: map[string]*Definition{},
		muDefinitions: &sync.RWMutex{},
	}
}

// Register adds a new definition to the registry.
func (r *objectDefinitionRegistry) Register(definition *Definition) error {
	if definition == nil {
		return fmt.Errorf("nil definition")
	}

	defer r.muDefinitions.Unlock()
	r.muDefinitions.Lock()

	if _, exists := r.definitionMap[definition.Name()]; exists {
		return ErrDefinitionAlreadyExists
	}

	r.definitionMap[definition.Name()] = definition

	return nil
}

// Remove removes a definition by name.
func (r *objectDefinitionRegistry) Remove(name string) error {
	defer r.muDefinitions.Unlock()
	r.muDefinitions.Lock()

	if _, exists := r.definitionMap[name]; !exists {
		return ErrDefinitionNotFound
	}

	delete(r.definitionMap, name)
	return nil
}

// Contains checks if a definition with the provided name exists.
func (r *objectDefinitionRegistry) Contains(name string) bool {
	defer r.muDefinitions.Unlock()
	r.muDefinitions.Lock()

	_, exists := r.definitionMap[name]
	return exists
}

// Find returns the definition that matches the provided filters.
func (r *objectDefinitionRegistry) Find(filters ...filter.Filter) (*Definition, error) {
	if len(filters) == 0 {
		return nil, ErrNoFilterProvided
	}

	definitionList := r.List(filters...)

	if len(definitionList) > 1 {
		return nil, ErrMultipleDefinitionsFound
	} else if len(definitionList) == 0 {
		return nil, ErrDefinitionNotFound
	}

	return definitionList[0], nil
}

// FindFirst returns the first definition that matches the provided filters.
func (r *objectDefinitionRegistry) FindFirst(filters ...filter.Filter) (*Definition, bool) {
	definitionList := r.List(filters...)

	if len(definitionList) == 0 {
		return nil, false
	}

	return definitionList[0], true
}

// List returns a list of definitions that match the provided filters.
func (r *objectDefinitionRegistry) List(filters ...filter.Filter) []*Definition {
	defer r.muDefinitions.Unlock()
	r.muDefinitions.Lock()

	filterOpts := filter.Of(filters...)
	definitionList := make([]*Definition, 0)

	for _, definition := range r.definitionMap {
		if filterOpts.Name != "" && filterOpts.Name != definition.Name() {
			continue
		}

		if filterOpts.Type == nil {
			definitionList = append(definitionList, definition)
			continue
		}

		if convertibleTo(definition.Type(), filterOpts.Type) {
			definitionList = append(definitionList, definition)
		}
	}

	return definitionList
}

// Names returns the names of all definitions.
func (r *objectDefinitionRegistry) Names() []string {
	defer r.muDefinitions.Unlock()
	r.muDefinitions.Lock()

	names := make([]string, 0)
	for name := range r.definitionMap {
		names = append(names, name)
	}

	return names
}

// Count returns the number of definitions.
func (r *objectDefinitionRegistry) Count() int {
	defer r.muDefinitions.Unlock()
	r.muDefinitions.Lock()

	return len(r.definitionMap)
}
