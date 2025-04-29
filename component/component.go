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
	"maps"
	"reflect"
	"slices"
	"sync"
	"unicode"
)

var (
	// components holds the globally registered component instances.
	components = make(map[string]*Component)

	// muComponents is the mutex used to guard access to the components map.
	muComponents = sync.RWMutex{}
)

// Component represents a registered component and holds its definition.
type Component struct {
	definition *Definition
	conditions []Condition
}

// Definition returns the component's definition metadata.
func (c *Component) Definition() *Definition {
	return c.definition
}

// Conditions returns the component's condition list.
func (c *Component) Conditions() []Condition {
	return slices.Clone(c.conditions)
}

// Registration lets you configure the registered component.
type Registration struct {
	component *Component
}

// Conditional attaches a runtime condition to the component.
// The condition is evaluated before the component is loaded into the container.
func (r *Registration) Conditional(cond Condition) *Registration {
	if cond != nil {
		r.component.conditions = append(r.component.conditions, cond)
	}

	return r
}

// Register registers a new component using the given constructor function and optional definition options.
// It panics if the component name already exists or if definition creation fails.
func Register(fn ConstructorFunc, opts ...DefinitionOption) *Registration {
	muComponents.Lock()
	defer muComponents.Unlock()

	def, err := MakeDefinition(fn, opts...)
	if err != nil {
		panic(err)
	}

	name := def.Name()

	if _, exists := components[name]; exists {
		panic(fmt.Errorf("component with name '%s' already exists", name))
	}

	component := &Component{
		definition: def,
	}

	components[name] = component
	return &Registration{
		component: component,
	}
}

// List returns all registered components as a slice.
func List() []*Component {
	muComponents.RLock()
	defer muComponents.RUnlock()
	return slices.Collect(maps.Values(components))
}

// ListOf returns all registered components whose type is assignable to the type T.
func ListOf[T any]() []*Component {
	muComponents.RLock()
	defer muComponents.RUnlock()

	targetType := reflect.TypeFor[T]()
	matches := make([]*Component, 0)

	for _, component := range components {
		sourceType := component.definition.Type()
		if convertibleTo(sourceType, targetType) {
			matches = append(matches, component)
		}
	}

	return matches
}

// convertibleTo checks if a source type can be converted to a target type.
// It unwraps pointer types and supports interface assignment compatibility.
func convertibleTo(sourceType reflect.Type, targetType reflect.Type) bool {
	if sourceType == targetType || (targetType.Kind() == reflect.Interface && sourceType.ConvertibleTo(targetType)) {
		return true
	} else if sourceType.Kind() == reflect.Pointer {
		return convertibleTo(sourceType.Elem(), targetType)
	}

	return false
}

// generateComponentName returns the name of the definition based on the return type of the constructor function
func generateComponentName(typ reflect.Type) string {
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	name := typ.Name()
	if name == "" {
		return ""
	}

	runes := []rune(name)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}
