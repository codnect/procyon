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

// DependencyRegistry stores instances by type for later retrieval during dependency resolution.
type DependencyRegistry interface {
	// RegisterDependency registers an instance of the specified type for dependency injection.
	RegisterDependency(typ reflect.Type, val any) error
}

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

	// DependencyRegistry registers non-component dependencies for injection.
	DependencyRegistry

	// ScopeRegistry manages custom scopes.
	ScopeRegistry

	// ProcessorRegistry manages lifecycle hooks.
	ProcessorRegistry
}

type HierarchicalContainer interface {
	Container

	// ParentContainer returns the parent container, or nil if there is no parent.
	ParentContainer() Container

	// SetParentContainer sets the parent container. It panics if the provided parent is nil
	// or if the container is already associated with a different parent.
	SetParentContainer(container Container)
}

// ContainerHolder is an interface that provides access to a Container.
type ContainerHolder interface {
	// Container returns the associated Container.
	Container() Container
}

// StandardContainer is the default implementation of the Container interface.
// It manages component definitions, singleton instances, custom scopes,
// and lifecycle processing.
type StandardContainer struct {
	parent        Container
	definitions   map[string]*Definition
	muDefinitions sync.RWMutex

	singletons             map[string]any
	singletonOrder         []string
	singletonState         *creationState
	typesOfSingletons      map[string]reflect.Type
	dependents             map[string]map[string]struct{}
	dependencies           map[string]map[string]struct{}
	muSingletons           sync.RWMutex
	resolvableDependencies map[reflect.Type]any
	muResolvableInstances  sync.RWMutex

	scopes   map[string]Scope
	muScopes sync.RWMutex

	beforeInitProcessors []BeforeInitProcessor
	afterInitProcessors  []AfterInitProcessor
	muProcessors         sync.RWMutex
}

// NewStandardContainer creates a StandardContainer.
func NewStandardContainer() *StandardContainer {
	return &StandardContainer{
		definitions:   make(map[string]*Definition),
		muDefinitions: sync.RWMutex{},

		singletons:        make(map[string]any),
		singletonOrder:    make([]string, 0),
		singletonState:    newCreationState(),
		typesOfSingletons: make(map[string]reflect.Type),
		dependents:        make(map[string]map[string]struct{}),
		dependencies:      make(map[string]map[string]struct{}),
		muSingletons:      sync.RWMutex{},

		scopes:   make(map[string]Scope),
		muScopes: sync.RWMutex{},

		resolvableDependencies: make(map[reflect.Type]any),
		muResolvableInstances:  sync.RWMutex{},

		beforeInitProcessors: []BeforeInitProcessor{},
		afterInitProcessors:  []AfterInitProcessor{},
		muProcessors:         sync.RWMutex{},
	}
}

// ParentContainer returns the parent container, or nil if there is no parent.
func (d *StandardContainer) ParentContainer() Container {
	return d.parent
}

// SetParentContainer sets the parent container. It panics if the provided parent is nil
// or if the container is already associated with a different parent.
func (d *StandardContainer) SetParentContainer(parent Container) {
	if parent == nil {
		panic("nil parent container")
	}

	if d.parent != nil && d.parent != parent {
		panic("already associated with a parent container")
	}

	d.parent = parent
}

// RegisterDefinition registers a new component definition.
// Returns an error if a definition with the same name already exists.
func (d *StandardContainer) RegisterDefinition(def *Definition) error {
	if def == nil {
		return errors.New("nil definition")
	}

	name := def.Name()

	d.muDefinitions.Lock()
	defer d.muDefinitions.Unlock()

	if _, dup := d.definitions[name]; dup {
		return fmt.Errorf("register definition %q: duplicate definition", name)
	}

	d.definitions[name] = def

	return nil
}

// UnregisterDefinition removes the component definition associated with the given name.
// Returns an error if the definition does not exist.
func (d *StandardContainer) UnregisterDefinition(name string) error {
	d.muDefinitions.Lock()
	defer d.muDefinitions.Unlock()

	if _, exists := d.definitions[name]; !exists {
		return fmt.Errorf("unregister definition %q: definition not found", name)
	}

	delete(d.definitions, name)

	return nil
}

// Definition retrieves the component definition associated with the given name.
// Returns the definition and a boolean indicating its existence.
func (d *StandardContainer) Definition(name string) (*Definition, bool) {
	d.muDefinitions.RLock()
	defer d.muDefinitions.RUnlock()

	if def, exists := d.definitions[name]; exists {
		return def, true
	}

	return nil, false
}

// ContainsDefinition checks whether a component definition with the specified name exists.
func (d *StandardContainer) ContainsDefinition(name string) bool {
	d.muDefinitions.RLock()
	defer d.muDefinitions.RUnlock()

	if _, exists := d.definitions[name]; exists {
		return true
	}

	return false
}

// Definitions return a slice of all registered component definitions.
func (d *StandardContainer) Definitions() []*Definition {
	d.muDefinitions.RLock()
	defer d.muDefinitions.RUnlock()

	if len(d.definitions) == 0 {
		return make([]*Definition, 0)
	}

	return slices.Collect(maps.Values(d.definitions))
}

// DefinitionsOf returns a slice of component definitions that are assignable to the specified type.
func (d *StandardContainer) DefinitionsOf(typ reflect.Type) []*Definition {
	if typ == nil {
		panic("nil definition type")
	}

	d.muDefinitions.RLock()
	defer d.muDefinitions.RUnlock()

	matches := make([]*Definition, 0)

	for _, def := range d.definitions {
		sourceType := def.Type()
		if convertibleTo(sourceType, typ) {
			matches = append(matches, def)
		}
	}

	return matches
}

// DefinitionNames returns a slice of all registered component definition names.
func (d *StandardContainer) DefinitionNames() []string {
	d.muDefinitions.RLock()
	defer d.muDefinitions.RUnlock()
	return slices.Collect(maps.Keys(d.definitions))
}

// DefinitionNamesOf returns a slice of component definition names that are assignable to the specified type.
func (d *StandardContainer) DefinitionNamesOf(typ reflect.Type) []string {
	if typ == nil {
		panic("nil definition type")
	}

	d.muDefinitions.RLock()
	defer d.muDefinitions.RUnlock()

	names := make([]string, 0)

	for name, def := range d.definitions {
		if convertibleTo(def.Type(), typ) {
			names = append(names, name)
		}
	}

	return names
}

// RegisterSingleton registers a singleton instance with the given name.
// Returns an error if a singleton instance with the same name already exists.
func (d *StandardContainer) RegisterSingleton(name string, instance any) error {
	if name == "" {
		return errors.New("empty instance name")
	}

	if instance == nil {
		return errors.New("nil instance")
	}

	d.muSingletons.Lock()
	defer d.muSingletons.Unlock()

	return d.registerSingleton(name, instance)
}

// ContainsSingleton checks whether a singleton with the specified name exists.
func (d *StandardContainer) ContainsSingleton(name string) bool {
	d.muSingletons.RLock()
	defer d.muSingletons.RUnlock()

	_, exists := d.singletons[name]
	return exists
}

// Singleton retrieves the singleton associated with the given name.
// Returns the instance and a boolean indicating its existence.
func (d *StandardContainer) Singleton(name string) (any, bool) {
	d.muSingletons.RLock()
	defer d.muSingletons.RUnlock()

	if singleton, exists := d.singletons[name]; exists {
		return singleton, true
	}

	return nil, false
}

// RemoveSingleton removes the singleton associated with the specified name.
func (d *StandardContainer) RemoveSingleton(name string) error {
	d.muSingletons.Lock()
	defer d.muSingletons.Unlock()

	if _, exists := d.singletons[name]; !exists {
		return fmt.Errorf("remove singleton %q: not found", name)
	}

	delete(d.singletons, name)
	delete(d.typesOfSingletons, name)

	return nil
}

// DestroySingletons destroys all registered singletons in reverse creation order.
// For each singleton, its dependents are recursively destroyed first.
func (d *StandardContainer) DestroySingletons() {
	d.muSingletons.Lock()
	defer d.muSingletons.Unlock()

	destroyed := make(map[string]struct{}, len(d.singletonOrder))

	for i := len(d.singletonOrder) - 1; i >= 0; i-- {
		d.destroySingleton(d.singletonOrder[i], destroyed)
	}

	clear(d.singletons)
	clear(d.typesOfSingletons)
	clear(d.dependents)
	clear(d.dependencies)
	d.singletonOrder = d.singletonOrder[:0]
}

// SingletonNames returns a slice of all registered singleton names.
func (d *StandardContainer) SingletonNames() []string {
	d.muSingletons.RLock()
	defer d.muSingletons.RUnlock()

	return slices.Collect(maps.Keys(d.singletons))
}

// destroySingleton recursively destroys a singleton and its dependents first.
// Must be called while holding muSingletons lock.
func (d *StandardContainer) destroySingleton(name string, destroyed map[string]struct{}) {
	if _, done := destroyed[name]; done {
		return
	}

	// destroy dependents first (components that depend on this one)
	for dependent := range d.dependents[name] {
		d.destroySingleton(dependent, destroyed)
	}

	destroyed[name] = struct{}{}

	singleton, exists := d.singletons[name]
	if !exists {
		return
	}

	if disposable, ok := singleton.(Disposable); ok {
		if err := disposable.Dispose(); err != nil {
			log.Warn("Failed to dispose singleton '{}'", name, err)
		}
	}
}

// CanResolve checks if a component with the given name is resolvable.
func (d *StandardContainer) CanResolve(name string) bool {
	if name == "" {
		return false
	}

	d.muSingletons.RLock()
	_, existsInSingletons := d.singletons[name]
	d.muSingletons.RUnlock()

	if existsInSingletons {
		return true
	}

	d.muDefinitions.RLock()
	_, existsInDefinitions := d.definitions[name]
	d.muDefinitions.RUnlock()

	if existsInDefinitions {
		return true
	}

	if d.parent != nil {
		return d.parent.CanResolve(name)
	}

	return false
}

// CanResolveType checks if a component of the given type is resolvable.
func (d *StandardContainer) CanResolveType(typ reflect.Type) bool {
	if typ == nil {
		return false
	}

	d.muSingletons.RLock()
	for _, singletonTyp := range d.typesOfSingletons {
		if convertibleTo(singletonTyp, typ) {
			d.muSingletons.RUnlock()
			return true
		}
	}
	d.muSingletons.RUnlock()

	d.muDefinitions.RLock()
	for _, def := range d.definitions {
		if convertibleTo(def.Type(), typ) {
			d.muDefinitions.RUnlock()
			return true
		}
	}
	d.muDefinitions.RUnlock()

	if d.parent != nil {
		return d.parent.CanResolveType(typ)
	}

	return false
}

// Resolve retrieves a component instance by its name.
func (d *StandardContainer) Resolve(ctx context.Context, name string) (any, error) {
	if ctx == nil {
		return nil, errors.New("nil context")
	}

	if name == "" {
		return nil, errors.New("empty instance name")
	}

	ctx = withCreationState(ctx)

	candidate, ok := d.Singleton(name)
	if ok {
		return candidate, nil
	}

	def, defExists := d.Definition(name)
	if !defExists {
		if d.parent != nil {
			return d.parent.Resolve(ctx, name)
		}

		return nil, fmt.Errorf("resolve %q: %w", name, ErrNotFound)
	}

	if def.IsSingleton() || def.IsPrototype() {
		instance, err := d.createInstance(ctx, def)
		if err != nil {
			return nil, fmt.Errorf("resolve %q: %w", name, err)
		}

		return instance, nil
	}

	scopeName := def.Scope()
	scope, scopeExists := d.Scope(scopeName)
	if !scopeExists {
		return nil, fmt.Errorf("resolve %q: scope %q not found", name, scopeName)
	}

	instance, err := scope.Resolve(ctx, name, func(ctx context.Context) (any, error) {
		instance, err := d.createInstance(ctx, def)
		if err != nil {
			return nil, fmt.Errorf("resolve %q: %w", name, err)
		}

		return instance, nil
	})

	if err != nil {
		return nil, fmt.Errorf("resolve %q: scope %q: %w", name, scopeName, err)
	}

	return instance, nil
}

// ResolveType retrieves an instance of the specified type.
func (d *StandardContainer) ResolveType(ctx context.Context, typ reflect.Type) (any, error) {
	if ctx == nil {
		return nil, errors.New("nil context")
	}

	if typ == nil {
		return nil, errors.New("nil instance type")
	}

	ctx = withCreationState(ctx)

	singletons := d.resolveSingletons(typ)
	definitions := d.DefinitionsOf(typ)
	nonSingletonDefs := make([]*Definition, 0)

	for _, def := range definitions {
		if !d.ContainsSingleton(def.Name()) {
			nonSingletonDefs = append(nonSingletonDefs, def)
		}
	}

	total := len(singletons) + len(nonSingletonDefs)
	if total > 1 {
		return nil, fmt.Errorf("resolve type %s: %w", typ, ErrAmbiguousMatch)
	}

	if len(singletons) > 0 {
		return singletons[0], nil
	}

	if len(definitions) > 0 {
		return d.Resolve(ctx, definitions[0].Name())
	}

	if d.parent != nil {
		return d.parent.ResolveType(ctx, typ)
	}

	return nil, fmt.Errorf("resolve type %s: %w", typ, ErrNotFound)
}

// ResolveAs retrieves a component by both name and expected type.
func (d *StandardContainer) ResolveAs(ctx context.Context, name string, typ reflect.Type) (any, error) {
	if ctx == nil {
		return nil, errors.New("nil context")
	}

	if name == "" {
		return nil, errors.New("empty instance name")
	}

	if typ == nil {
		return nil, errors.New("nil instance type")
	}

	ctx = withCreationState(ctx)

	instance, err := d.Resolve(ctx, name)
	if err != nil {
		return nil, err
	}

	instanceType := reflect.TypeOf(instance)
	if !convertibleTo(instanceType, typ) {
		return nil, fmt.Errorf("resolve %q: %s is not convertible to %s: %w", name, instanceType, typ, ErrTypeMismatch)
	}

	return instance, nil
}

// ResolveAll retrieves all instances assignable to the specified type.
func (d *StandardContainer) ResolveAll(ctx context.Context, typ reflect.Type) ([]any, error) {
	if ctx == nil {
		return nil, errors.New("nil context")
	}

	if typ == nil {
		return nil, errors.New("nil instance type")
	}

	ctx = withCreationState(ctx)

	resolvedNames := make(map[string]struct{})
	instances := make([]any, 0)

	for _, name := range d.DefinitionNamesOf(typ) {
		resolvedNames[name] = struct{}{}
		instance, err := d.Resolve(ctx, name)
		if err != nil {
			return nil, err
		}
		instances = append(instances, instance)
	}

	for _, name := range d.SingletonNames() {
		if _, already := resolvedNames[name]; already {
			continue
		}

		singleton, exists := d.Singleton(name)
		if !exists {
			continue
		}

		if convertibleTo(reflect.TypeOf(singleton), typ) {
			resolvedNames[name] = struct{}{}
			instances = append(instances, singleton)
		}
	}

	if d.parent != nil {
		for _, name := range d.parent.DefinitionNamesOf(typ) {
			if _, already := resolvedNames[name]; already {
				continue
			}

			resolvedNames[name] = struct{}{}
			instance, err := d.parent.Resolve(ctx, name)
			if err != nil {
				return nil, err
			}

			instances = append(instances, instance)
		}

		for _, name := range d.parent.SingletonNames() {
			if _, already := resolvedNames[name]; already {
				continue
			}

			singleton, exists := d.parent.Singleton(name)
			if !exists {
				continue
			}

			if convertibleTo(reflect.TypeOf(singleton), typ) {
				resolvedNames[name] = struct{}{}
				instances = append(instances, singleton)
			}
		}
	}

	return instances, nil
}

// RegisterDependency registers an instance of the specified type for dependency injection.
func (d *StandardContainer) RegisterDependency(typ reflect.Type, instance any) error {
	if typ == nil {
		return errors.New("nil instance type")
	}

	if instance == nil {
		return errors.New("nil instance")
	}

	d.muResolvableInstances.Lock()
	defer d.muResolvableInstances.Unlock()
	d.resolvableDependencies[typ] = instance

	return nil
}

// RegisterScope adds a new scope with the specified name to the registry.
func (d *StandardContainer) RegisterScope(name string, scope Scope) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("empty scope name")
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

	return fmt.Errorf("register scope %q: reserved scope", name)
}

// Scope retrieves the scope associated with the given name.
// Returns the scope and a boolean indicating its existence.
func (d *StandardContainer) Scope(name string) (Scope, bool) {
	if name == "" {
		return nil, false
	}

	d.muScopes.RLock()
	defer d.muScopes.RUnlock()

	if scope, ok := d.scopes[name]; ok {
		return scope, true
	}

	return nil, false
}

// UseBeforeInitProcessor registers an BeforeInitProcessor to be applied before initialization.
func (d *StandardContainer) UseBeforeInitProcessor(processor BeforeInitProcessor) error {
	if processor == nil {
		return errors.New("nil before-init processor")
	}

	d.muProcessors.Lock()
	defer d.muProcessors.Unlock()

	d.beforeInitProcessors = append(d.beforeInitProcessors, processor)
	return nil
}

// UseAfterInitProcessor registers an AfterInitProcessor to be applied after initialization.
func (d *StandardContainer) UseAfterInitProcessor(processor AfterInitProcessor) error {
	if processor == nil {
		return errors.New("nil after-init processor")
	}

	d.muProcessors.Lock()
	defer d.muProcessors.Unlock()

	d.afterInitProcessors = append(d.afterInitProcessors, processor)
	return nil
}

// registerSingleton registers a singleton instance with the given name.
func (d *StandardContainer) registerSingleton(name string, instance any) error {
	if _, dup := d.singletons[name]; dup {
		return fmt.Errorf("register singleton %q: duplicate instance", name)
	}

	d.singletons[name] = instance
	d.singletonOrder = append(d.singletonOrder, name)
	d.typesOfSingletons[name] = reflect.TypeOf(instance)
	return nil
}

// createInstance constructs a new instance of a component using its definition.
func (d *StandardContainer) createInstance(ctx context.Context, def *Definition) (any, error) {
	name := def.Name()

	state := creationStateFromContext(ctx)
	if def.IsSingleton() {
		state = d.singletonState
	}

	defer state.removeFromPreparation(name)
	if err := state.putToPreparation(name); err != nil {
		return nil, err
	}

	constructor := def.Constructor()

	args, err := d.resolveArguments(ctx, constructor.Args())
	if err != nil {
		return nil, fmt.Errorf("create %q (%s): %w", name, def.Type(), err)
	}

	var (
		instance any
	)

	instance, err = constructor.Invoke(args...)
	if err != nil {
		return nil, fmt.Errorf("invoke constructor %q (%s): %w", name, def.Type(), err)
	}

	instance, err = d.initialize(ctx, instance)
	if err != nil {
		return nil, fmt.Errorf("initialize %q (%s): %w", name, def.Type(), err)
	}

	if def.IsSingleton() {
		d.muSingletons.Lock()
		defer d.muSingletons.Unlock()
		err = d.registerSingleton(name, instance)

		if err != nil {
			return nil, err
		}

		d.registerDependencies(name, constructor.Args())
	}

	return instance, nil
}

// resolveArguments resolves the constructor arguments required to instantiate a component.
func (d *StandardContainer) resolveArguments(ctx context.Context, args []Arg) ([]any, error) {
	resolvedArgs := make([]any, 0, len(args))

	for idx, arg := range args {

		if arg.Type().Kind() == reflect.Slice {
			elemType := arg.Type().Elem()
			instances, err := d.ResolveAll(ctx, elemType)
			if err != nil {
				return nil, fmt.Errorf("unsatisfied dependency for argument %d (%s): %w", idx, arg.Type(), err)
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
			if arg.Name() != "" {
				return nil, fmt.Errorf("unsatisfied dependency for argument %d %q (%s): %w", idx, arg.Name(), arg.Type(), err)
			}

			return nil, fmt.Errorf("unsatisfied dependency for argument %d (%s): %w", idx, arg.Type(), err)
		}

		resolvedArgs = append(resolvedArgs, instance)
	}

	return resolvedArgs, nil
}

// resolveSingletons returns all singleton instances that are assignable to the specified type.
func (d *StandardContainer) resolveSingletons(typ reflect.Type) []any {
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

// registerDependencies registers the dependencies of a component based on its constructor arguments.
func (d *StandardContainer) registerDependencies(name string, args []Arg) {
	for _, arg := range args {
		if arg.Type().Kind() == reflect.Slice {
			for _, depDef := range d.DefinitionsOf(arg.Type().Elem()) {
				d.registerDependency(name, depDef.Name())
			}
		} else if arg.Name() != "" {
			d.registerDependency(name, arg.Name())
		} else {
			if defs := d.DefinitionsOf(arg.Type()); len(defs) == 1 {
				d.registerDependency(name, defs[0].Name())
			}
		}
	}
}

// registerDependency registers a dependency relationship between a component and its dependency.
func (d *StandardContainer) registerDependency(name, dependency string) {
	if d.dependents[dependency] == nil {
		d.dependents[dependency] = make(map[string]struct{})
	}

	d.dependents[dependency][name] = struct{}{}

	if d.dependencies[name] == nil {
		d.dependencies[name] = make(map[string]struct{})
	}

	d.dependencies[name][dependency] = struct{}{}
}

// initialize runs pre-processors, the Init method (if defined), and post-processors
// on the given instance, and it returns the fully initialized instance.
func (d *StandardContainer) initialize(ctx context.Context, instance any) (any, error) {
	result, err := d.applyBeforeInitProcessors(ctx, instance)
	if err != nil {
		return nil, fmt.Errorf("apply before-init processors: %w", err)
	}

	if initializer, ok := result.(Initializer); ok {
		err = initializer.Init(ctx)

		if err != nil {
			return nil, fmt.Errorf("invoke init: %w", err)
		}
	}

	result, err = d.applyAfterInitProcessors(ctx, result)
	if err != nil {
		return nil, fmt.Errorf("apply after-init processors: %w", err)
	}

	return result, nil
}

// applyBeforeInitProcessors executes all registered BeforeInitProcessor hooks on the instance.
// Returns the processed object or an error.
func (d *StandardContainer) applyBeforeInitProcessors(ctx context.Context, instance any) (any, error) {
	d.muProcessors.RLock()
	defer d.muProcessors.RUnlock()

	for _, processor := range d.beforeInitProcessors {
		result, err := processor.ProcessBeforeInit(ctx, instance)

		if err != nil {
			return nil, fmt.Errorf("before-init processor (%T): %w", processor, err)
		}

		if result == nil {
			return nil, fmt.Errorf("before-init processor (%T) returned nil", processor)
		}

		instance = result
	}

	return instance, nil
}

// applyAfterInitProcessors executes all registered AfterInitProcessor hooks on the instance.
// Returns the processed object or an error.
func (d *StandardContainer) applyAfterInitProcessors(ctx context.Context, instance any) (any, error) {
	d.muProcessors.RLock()
	defer d.muProcessors.RUnlock()

	for _, processor := range d.afterInitProcessors {
		result, err := processor.ProcessAfterInit(ctx, instance)

		if err != nil {
			return nil, fmt.Errorf("after-init processor (%T): %w", processor, err)
		}

		if result == nil {
			return nil, fmt.Errorf("after-init processor (%T) returned nil", processor)
		}

		instance = result
	}

	return instance, nil
}
