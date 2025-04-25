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
	"errors"
	"maps"
	"reflect"
	"slices"
	"sync"
)

// Container provides a unified interface for managing components,
// including definition registration, instance resolving, custom scopes,
// lifecycle management, and manual bindings.
type Container interface {
	// DefinitionRegistry provides access to component definitions and their metadata.
	DefinitionRegistry

	// SingletonRegistry manages singleton instances of components.
	SingletonRegistry

	// Resolver resolves component instances by type or name.
	Resolver

	// Binder allows manual binding of instances to specific types.
	Binder

	// ScopeRegistry manages custom scopes.
	ScopeRegistry

	// LifecycleManager manages lifecycle hooks.
	LifecycleManager
}

// DefinitionRegistry provides methods for registering and unregistering component definitions.
type DefaultContainer struct {
	definitions   map[string]*Definition
	muDefinitions sync.RWMutex
}

// ErrDefinitionAlreadyExists is returned when trying to register a definition that already exists.
func NewDefaultContainer() *DefaultContainer {
	return &DefaultContainer{
		definitions:   make(map[string]*Definition),
		muDefinitions: sync.RWMutex{},
	}
}

// ErrDefinitionAlreadyExists is returned when trying to register a definition that already exists.
func (d *DefaultContainer) RegisterDefinition(def *Definition) error {
	if def == nil {
		return errors.New("nil definition")
	}

	name := def.Name()

	defer d.muDefinitions.Unlock()
	d.muDefinitions.Lock()

	if _, exists := d.definitions[name]; exists {
		return ErrDefinitionAlreadyExists
	}

	d.definitions[name] = def

	return nil
}

// ErrDefinitionNotFound is returned when trying to unregister a definition that does not exist.
func (d *DefaultContainer) UnregisterDefinition(name string) error {
	defer d.muDefinitions.Unlock()
	d.muDefinitions.Lock()

	if _, exists := d.definitions[name]; !exists {
		return ErrDefinitionNotFound
	}

	delete(d.definitions, name)

	return nil
}

// Definition retrieves the component definition associated with the given name.
func (d *DefaultContainer) Definition(name string) (*Definition, bool) {
	defer d.muDefinitions.RUnlock()
	d.muDefinitions.RLock()

	if def, exists := d.definitions[name]; exists {
		return def, true
	}

	return nil, false
}

// ContainsDefinition checks whether a component definition with the specified name exists.
func (d *DefaultContainer) ContainsDefinition(name string) bool {
	defer d.muDefinitions.RUnlock()
	d.muDefinitions.RLock()

	if _, exists := d.definitions[name]; exists {
		return true
	}

	return false
}

// Definitions returns a slice of all registered component definitions.
func (d *DefaultContainer) Definitions() []*Definition {
	defer d.muDefinitions.RUnlock()
	d.muDefinitions.RLock()

	return slices.Collect(maps.Values(d.definitions))
}

// DefinitionsOf returns a slice of component definitions that are assignable to the specified type.
func (d *DefaultContainer) DefinitionsOf(typ reflect.Type) []*Definition {
	muComponents.RLock()
	defer muComponents.RUnlock()

	matches := make([]*Definition, 0)

	for _, def := range d.definitions {
		sourceType := def.Type()
		if convertibleTo(sourceType, typ) {
			matches = append(matches, def)
		}
	}

	return matches
}
