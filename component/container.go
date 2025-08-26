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
	"context"
	"errors"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"
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

	// ResolvableRegistry registers non-component dependencies for injection.
	ResolvableRegistry

	// ScopeRegistry manages custom scopes.
	ScopeRegistry

	// LifecycleManager manages lifecycle hooks.
	LifecycleManager
}

// ContainerCapable is an interface that indicates the ability to provide
// access to Container.
type ContainerCapable interface {
	// Container returns the associated Container.
	Container() Container
}

// DefaultContainer is the default implementation of the Container interface.
// It manages component definitions, singleton instances, custom scopes,
// and lifecycle processing.
type DefaultContainer struct {
	definitions   map[string]*Definition
	muDefinitions sync.RWMutex

	singletons            map[string]any
	singletonState        *creationState
	typesOfSingletons     map[string]reflect.Type
	muSingletons          sync.RWMutex
	resolvableInstances   map[reflect.Type]any
	muResolvableInstances sync.RWMutex

	scopes   map[string]Scope
	muScopes sync.RWMutex

	preProcessors  []PreProcessor
	postProcessors []PostProcessor
	muProcessors   sync.RWMutex
}

// NewDefaultContainer creates a DefaultContainer.
func NewDefaultContainer() *DefaultContainer {
	return &DefaultContainer{
		definitions:   make(map[string]*Definition),
		muDefinitions: sync.RWMutex{},

		singletons:        make(map[string]any),
		singletonState:    newCreationState(),
		typesOfSingletons: make(map[string]reflect.Type),
		muSingletons:      sync.RWMutex{},

		scopes:   make(map[string]Scope),
		muScopes: sync.RWMutex{},

		resolvableInstances:   make(map[reflect.Type]any),
		muResolvableInstances: sync.RWMutex{},

		preProcessors:  []PreProcessor{},
		postProcessors: []PostProcessor{},
		muProcessors:   sync.RWMutex{},
	}
}

// RegisterDefinition registers a new component definition.
// Returns an error if a definition with the same name already exists.
func (d *DefaultContainer) RegisterDefinition(def *Definition) error {
	if def == nil {
		return errors.New("nil definition")
	}

	name := def.Name()

	d.muDefinitions.Lock()
	defer d.muDefinitions.Unlock()

	if _, exists := d.definitions[name]; exists {
		return ErrDefinitionAlreadyExists
	}

	d.definitions[name] = def

	return nil
}

// UnregisterDefinition removes the component definition associated with the given name.
// Returns an error if the definition does not exist.
func (d *DefaultContainer) UnregisterDefinition(name string) error {
	d.muDefinitions.Lock()
	defer d.muDefinitions.Unlock()

	if _, exists := d.definitions[name]; !exists {
		return ErrDefinitionNotFound
	}

	delete(d.definitions, name)

	return nil
}

// Definition retrieves the component definition associated with the given name.
// Returns the definition and a boolean indicating its existence.
func (d *DefaultContainer) Definition(name string) (*Definition, bool) {
	d.muDefinitions.RLock()
	defer d.muDefinitions.RUnlock()

	if def, exists := d.definitions[name]; exists {
		return def, true
	}

	return nil, false
}

// ContainsDefinition checks whether a component definition with the specified name exists.
func (d *DefaultContainer) ContainsDefinition(name string) bool {
	d.muDefinitions.RLock()
	defer d.muDefinitions.RUnlock()

	if _, exists := d.definitions[name]; exists {
		return true
	}

	return false
}

// Definitions return a slice of all registered component definitions.
func (d *DefaultContainer) Definitions() []*Definition {
	d.muDefinitions.RLock()
	defer d.muDefinitions.RUnlock()

	if len(d.definitions) == 0 {
		return make([]*Definition, 0)
	}

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

// RegisterSingleton registers a singleton instance with the given name.
// Returns an error if a singleton with the same name already exists.
func (d *DefaultContainer) RegisterSingleton(name string, instance any) error {
	if name == "" {
		return errors.New("empty name")
	}

	if instance == nil {
		return errors.New("nil instance")
	}

	d.muSingletons.Lock()
	defer d.muSingletons.Unlock()

	if _, exists := d.singletons[name]; exists {
		return ErrInstanceAlreadyExists
	}

	d.singletons[name] = instance
	d.typesOfSingletons[name] = reflect.TypeOf(instance)
	return nil
}

// ContainsSingleton checks whether a singleton with the specified name exists.
func (d *DefaultContainer) ContainsSingleton(name string) bool {
	d.muSingletons.Lock()
	defer d.muSingletons.Unlock()

	_, exists := d.singletons[name]
	return exists
}

// Singleton retrieves the singleton instance associated with the given name.
// Returns the instance and a boolean indicating its existence.
func (d *DefaultContainer) Singleton(name string) (any, bool) {
	d.muSingletons.Lock()
	defer d.muSingletons.Unlock()

	if singleton, exists := d.singletons[name]; exists {
		return singleton, true
	}

	return nil, false
}

// RemoveSingleton removes the singleton instance associated with the specified name.
func (d *DefaultContainer) RemoveSingleton(name string) error {
	d.muSingletons.Lock()
	defer d.muSingletons.Unlock()

	if _, exists := d.singletons[name]; !exists {
		return ErrInstanceNotFound
	}

	delete(d.singletons, name)
	delete(d.typesOfSingletons, name)

	return nil
}

// CanResolve checks if a component with the given name is resolvable.
func (d *DefaultContainer) CanResolve(name string) bool {
	d.muSingletons.RLock()
	defer d.muSingletons.RUnlock()

	if _, exists := d.singletons[name]; exists {
		return true
	}

	d.muDefinitions.RLock()
	defer d.muDefinitions.RUnlock()

	if _, exists := d.definitions[name]; exists {
		return true
	}

	return false
}

// CanResolveType checks if a component of the given type is resolvable.
func (d *DefaultContainer) CanResolveType(typ reflect.Type) bool {
	d.muSingletons.RLock()
	defer d.muSingletons.RUnlock()

	for _, singletonTyp := range d.typesOfSingletons {
		if convertibleTo(singletonTyp, typ) {
			return true
		}
	}

	d.muDefinitions.RLock()
	defer d.muDefinitions.RUnlock()

	for _, def := range d.definitions {
		if convertibleTo(def.Type(), typ) {
			return true
		}
	}

	return false
}

// Resolve retrieves a component instance by its name.
func (d *DefaultContainer) Resolve(ctx context.Context, name string) (any, error) {
	if name == "" {
		return nil, errors.New("empty name")
	}

	ctx = withCreationState(ctx)

	candidate, ok := d.Singleton(name)
	if ok {
		return candidate, nil
	}

	def, defExists := d.Definition(name)
	if !defExists {
		return nil, ErrDefinitionNotFound
	}

	state := creationStateFromContext(ctx)

	if def.IsSingleton() {
		return d.createSingleton(ctx, def)
	} else if def.IsPrototype() {
		defer state.removeFromPreparation(name)
		if err := state.putToPreparation(name); err != nil {
			return nil, err
		}

		instance, err := d.createInstance(ctx, def)

		if err != nil {
			return nil, err
		}

		return instance, nil
	}

	scope, scopeExists := d.Scope(def.Scope())
	if !scopeExists {
		return nil, ErrScopeNotFound
	}

	return scope.Resolve(ctx, name, func(ctx context.Context) (any, error) {
		defer state.removeFromPreparation(name)
		if err := state.putToPreparation(name); err != nil {
			return nil, err
		}

		return d.createInstance(ctx, def)
	})
}

// ResolveType retrieves an instance of the specified type.
func (d *DefaultContainer) ResolveType(ctx context.Context, typ reflect.Type) (any, error) {
	if typ == nil {
		return nil, errors.New("nil type")
	}

	ctx = withCreationState(ctx)

	resolvableCandidates := d.findResolvableCandidates(typ)
	if len(resolvableCandidates) > 1 {
		return nil, errors.New("multiple instance found")
	} else if len(resolvableCandidates) == 1 {
		return resolvableCandidates[0], nil
	}

	singletons := d.resolveSingletons(typ)
	if len(singletons) > 1 {
		return nil, errors.New("multiple singletons found")
	} else if len(singletons) == 1 {
		return singletons[0], nil
	}

	definitions := d.DefinitionsOf(typ)
	if len(definitions) > 1 {
		return nil, errors.New("multiple definitions found")
	} else if len(definitions) == 1 {
		return d.Resolve(ctx, definitions[0].Name())
	}

	return nil, ErrInstanceNotFound
}

// ResolveAs retrieves a component by both name and expected type.
func (d *DefaultContainer) ResolveAs(ctx context.Context, name string, typ reflect.Type) (any, error) {
	if name == "" {
		return nil, errors.New("empty name")
	}

	if typ == nil {
		return nil, errors.New("nil type")
	}

	ctx = withCreationState(ctx)

	instance, err := d.Resolve(ctx, name)
	if err != nil {
		return nil, err
	}

	instanceType := reflect.TypeOf(instance)
	if !convertibleTo(instanceType, typ) {
		return nil, fmt.Errorf("component %q is not assignable to %s", name, typ)
	}

	return instance, nil
}

// ResolveAll retrieves all instances assignable to the specified type.
func (d *DefaultContainer) ResolveAll(ctx context.Context, typ reflect.Type) ([]any, error) {
	if typ == nil {
		return nil, errors.New("nil type")
	}

	ctx = withCreationState(ctx)

	instances := d.findResolvableCandidates(typ)

	for _, def := range d.DefinitionsOf(typ) {
		instance, err := d.Resolve(ctx, def.Name())

		if err != nil {
			return nil, err
		}

		instances = append(instances, instance)
	}

	return instances, nil
}

// RegisterResolvable registers type with the corresponding value.
func (d *DefaultContainer) RegisterResolvable(typ reflect.Type, instance any) error {
	if typ == nil {
		return errors.New("nil type")
	}

	if instance == nil {
		return errors.New("nil instance")
	}

	d.muResolvableInstances.Lock()
	defer d.muResolvableInstances.Unlock()
	d.resolvableInstances[typ] = instance

	return nil
}

// RegisterScope adds a new scope with the specified name to the registry.
func (d *DefaultContainer) RegisterScope(name string, scope Scope) error {
	if strings.TrimSpace(name) == "" {
		return ErrInvalidScopeName
	}

	if scope == nil {
		return errors.New("nil scope")
	}

	if SingletonScope != name && PrototypeScope != name {
		d.muScopes.Lock()
		defer d.muScopes.Unlock()

		d.scopes[name] = scope
		return nil
	}

	return ErrScopeReplacementNotAllowed
}

// Scope retrieves the scope associated with the given name.
// Returns the scope and a boolean indicating its existence.
func (d *DefaultContainer) Scope(name string) (Scope, bool) {
	d.muScopes.Lock()
	defer d.muScopes.Unlock()

	if scope, ok := d.scopes[name]; ok {
		return scope, true
	}

	return nil, false
}

// UsePreProcessor registers a PreProcessor to be applied before initialization.
func (d *DefaultContainer) UsePreProcessor(processor PreProcessor) error {
	if processor == nil {
		return errors.New("nil processor")
	}

	d.muProcessors.Lock()
	defer d.muProcessors.Unlock()

	d.preProcessors = append(d.preProcessors, processor)
	return nil
}

// UsePostProcessor registers a PostProcessor to be applied after initialization.
func (d *DefaultContainer) UsePostProcessor(processor PostProcessor) error {
	if processor == nil {
		return errors.New("nil processor")
	}

	d.muProcessors.Lock()
	defer d.muProcessors.Unlock()

	d.postProcessors = append(d.postProcessors, processor)
	return nil
}

// createInstance constructs a new instance of a component using its definition.
func (d *DefaultContainer) createInstance(ctx context.Context, def *Definition) (any, error) {
	constructor := def.Constructor()

	args, err := d.resolveArguments(ctx, constructor.Args())
	if err != nil {
		return nil, err
	}

	var instance any
	instance, err = constructor.Invoke(args...)
	if err != nil {
		return nil, err
	}

	instance, err = d.initialize(ctx, instance)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

// createSingleton constructs a singleton instance and registers it in the container.
// It also handles circular dependency protection via preparation state.
func (d *DefaultContainer) createSingleton(ctx context.Context, def *Definition) (any, error) {
	name := def.Name()

	defer d.singletonState.removeFromPreparation(name)
	if err := d.singletonState.putToPreparation(name); err != nil {
		return nil, err
	}

	instance, err := d.createInstance(ctx, def)
	if err != nil {
		return nil, err
	}

	err = d.RegisterSingleton(name, instance)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

// resolveArguments resolves the constructor arguments required to instantiate a component.
func (d *DefaultContainer) resolveArguments(ctx context.Context, args []Arg) ([]any, error) {
	resolvedArgs := make([]any, 0, len(args))

	for _, arg := range args {

		if arg.Type().Kind() == reflect.Slice {
			elemType := arg.Type().Elem()
			instances, err := d.ResolveAll(ctx, elemType)
			if err != nil {
				return nil, err
			}

			sliceVal := reflect.MakeSlice(arg.Type(), len(instances), len(instances))
			for i, inst := range instances {
				sliceVal.Index(i).Set(reflect.ValueOf(inst))
			}

			resolvedArgs = append(resolvedArgs, sliceVal.Interface())
			continue
		}

		var (
			instance any
			err      error
		)

		if arg.Name() != "" {
			instance, err = d.Resolve(ctx, arg.Name())
		} else {
			instance, err = d.ResolveType(ctx, arg.Type())
		}

		if err != nil {
			return nil, err
		}

		resolvedArgs = append(resolvedArgs, instance)
	}

	return resolvedArgs, nil
}

func (d *DefaultContainer) findResolvableCandidates(typ reflect.Type) []any {
	d.muResolvableInstances.RLock()
	defer d.muResolvableInstances.RUnlock()

	candidates := make([]any, 0)

	for candidateType, candidate := range d.resolvableInstances {
		if convertibleTo(candidateType, typ) {
			candidates = append(candidates, candidate)
		}
	}

	return candidates
}

// resolveSingletons returns all singleton instances that are assignable to the specified type.
func (d *DefaultContainer) resolveSingletons(typ reflect.Type) []any {
	d.muSingletons.RLock()
	defer d.muSingletons.RUnlock()

	singletons := make([]any, 0)

	for singletonName, singletonType := range d.typesOfSingletons {
		if convertibleTo(singletonType, typ) {
			singletons = append(singletons, d.singletons[singletonName])
		}
	}

	return singletons
}

// initialize runs pre-processors, the Init method (if defined), and post-processors
// on the given instance, and it returns the fully initialized instance.
func (d *DefaultContainer) initialize(ctx context.Context, instance any) (any, error) {
	result, err := d.applyPreProcessors(ctx, instance)
	if err != nil {
		return nil, err
	}

	if initializer, ok := instance.(Initializer); ok {
		err = initializer.Init(ctx)

		if err != nil {
			return nil, err
		}
	}

	return d.applyPostProcessors(ctx, result)
}

// applyPreProcessors executes all registered PreProcessor hooks on the instance.
// Returns the processed object or an error.
func (d *DefaultContainer) applyPreProcessors(ctx context.Context, instance any) (any, error) {
	d.muProcessors.RLock()
	defer d.muProcessors.RUnlock()

	for _, processor := range d.preProcessors {
		result, err := processor.ProcessBeforeInit(ctx, instance)

		if err != nil {
			return nil, err
		}

		if result == nil {
			return nil, errors.New("nil processor result")
		}

		instance = result
	}

	return instance, nil
}

// applyPostProcessors executes all registered PostProcessor hooks on the instance.
// Returns the processed object or an error.
func (d *DefaultContainer) applyPostProcessors(ctx context.Context, instance any) (any, error) {
	d.muProcessors.RLock()
	defer d.muProcessors.RUnlock()

	for _, processor := range d.postProcessors {
		result, err := processor.ProcessAfterInit(ctx, instance)

		if err != nil {
			return nil, err
		}

		if result == nil {
			return nil, errors.New("nil processor result")
		}

		instance = result
	}

	return instance, nil
}
