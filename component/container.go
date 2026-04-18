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
	singletonOrder        []string
	singletonState        *creationState
	typesOfSingletons     map[string]reflect.Type
	dependents            map[string]map[string]struct{}
	dependencies          map[string]map[string]struct{}
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
		singletonOrder:    make([]string, 0),
		singletonState:    newCreationState(),
		typesOfSingletons: make(map[string]reflect.Type),
		dependents:        make(map[string]map[string]struct{}),
		dependencies:      make(map[string]map[string]struct{}),
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

	if _, dup := d.definitions[name]; dup {
		return fmt.Errorf("register definition %q: duplicate definition", name)
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
		return fmt.Errorf("unregister definition %q: definition not found", name)
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

// RegisterSingleton registers a singleton instance with the given name.
// Returns an error if a singleton instance with the same name already exists.
func (d *DefaultContainer) RegisterSingleton(name string, instance any) error {
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
func (d *DefaultContainer) ContainsSingleton(name string) bool {
	d.muSingletons.RLock()
	defer d.muSingletons.RUnlock()

	_, exists := d.singletons[name]
	return exists
}

// Singleton retrieves the singleton associated with the given name.
// Returns the instance and a boolean indicating its existence.
func (d *DefaultContainer) Singleton(name string) (any, bool) {
	d.muSingletons.RLock()
	defer d.muSingletons.RUnlock()

	if singleton, exists := d.singletons[name]; exists {
		return singleton, true
	}

	return nil, false
}

// RemoveSingleton removes the singleton associated with the specified name.
func (d *DefaultContainer) RemoveSingleton(name string) error {
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
func (d *DefaultContainer) DestroySingletons() {
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

// destroySingleton recursively destroys a singleton and its dependents first.
// Must be called while holding muSingletons lock.
func (d *DefaultContainer) destroySingleton(name string, destroyed map[string]struct{}) {
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
func (d *DefaultContainer) CanResolve(name string) bool {
	if name == "" {
		return false
	}

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
	if typ == nil {
		return false
	}

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
		return nil, fmt.Errorf("resolve %q: %w", name, ErrNotFound)
	}

	if def.IsSingleton() || def.IsPrototype() {
		return d.createInstance(ctx, def)
	}

	scopeName := def.Scope()
	scope, scopeExists := d.Scope(scopeName)
	if !scopeExists {
		return nil, fmt.Errorf("resolve %q: scope %q not found", name, scopeName)
	}

	instance, err := scope.Resolve(ctx, name, func(ctx context.Context) (any, error) {
		return d.createInstance(ctx, def)
	})

	if err != nil {
		return nil, fmt.Errorf("resolve %q: scope %q: %w", name, scopeName, err)
	}

	return instance, nil
}

// ResolveType retrieves an instance of the specified type.
func (d *DefaultContainer) ResolveType(ctx context.Context, typ reflect.Type) (any, error) {
	if ctx == nil {
		return nil, errors.New("nil context")
	}

	if typ == nil {
		return nil, errors.New("nil instance type")
	}

	ctx = withCreationState(ctx)

	resolvableCandidates := d.findResolvableCandidates(typ)
	if len(resolvableCandidates) > 1 {
		return nil, fmt.Errorf("resolve type %s: %w", typ, ErrAmbiguousMatch)
	} else if len(resolvableCandidates) == 1 {
		return resolvableCandidates[0], nil
	}

	singletons := d.resolveSingletons(typ)
	if len(singletons) > 1 {
		return nil, fmt.Errorf("resolve type %s: %w", typ, ErrAmbiguousMatch)
	} else if len(singletons) == 1 {
		return singletons[0], nil
	}

	definitions := d.DefinitionsOf(typ)
	if len(definitions) > 1 {
		return nil, fmt.Errorf("resolve type %s: %w", typ, ErrAmbiguousMatch)
	} else if len(definitions) == 1 {
		return d.Resolve(ctx, definitions[0].Name())
	}

	return nil, fmt.Errorf("resolve type %s: %w", typ, ErrNotFound)
}

// ResolveAs retrieves a component by both name and expected type.
func (d *DefaultContainer) ResolveAs(ctx context.Context, name string, typ reflect.Type) (any, error) {
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
func (d *DefaultContainer) ResolveAll(ctx context.Context, typ reflect.Type) ([]any, error) {
	if typ == nil {
		return nil, errors.New("nil instance type")
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

// RegisterResolvable registers type with the corresponding instance.
func (d *DefaultContainer) RegisterResolvable(typ reflect.Type, instance any) error {
	if typ == nil {
		return errors.New("nil instance type")
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
func (d *DefaultContainer) Scope(name string) (Scope, bool) {
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

// UsePreProcessor registers a PreProcessor to be applied before initialization.
func (d *DefaultContainer) UsePreProcessor(processor PreProcessor) error {
	if processor == nil {
		return errors.New("nil pre-processor")
	}

	d.muProcessors.Lock()
	defer d.muProcessors.Unlock()

	d.preProcessors = append(d.preProcessors, processor)
	return nil
}

// UsePostProcessor registers a PostProcessor to be applied after initialization.
func (d *DefaultContainer) UsePostProcessor(processor PostProcessor) error {
	if processor == nil {
		return errors.New("nil post-processor")
	}

	d.muProcessors.Lock()
	defer d.muProcessors.Unlock()

	d.postProcessors = append(d.postProcessors, processor)
	return nil
}

// registerSingleton registers a singleton instance with the given name.
func (d *DefaultContainer) registerSingleton(name string, instance any) error {
	if _, dup := d.singletons[name]; dup {
		return fmt.Errorf("register singleton %q: duplicate instance", name)
	}

	d.singletons[name] = instance
	d.singletonOrder = append(d.singletonOrder, name)
	d.typesOfSingletons[name] = reflect.TypeOf(instance)
	return nil
}

// createInstance constructs a new instance of a component using its definition.
func (d *DefaultContainer) createInstance(ctx context.Context, def *Definition) (any, error) {
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
func (d *DefaultContainer) resolveArguments(ctx context.Context, args []Arg) ([]any, error) {
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

// registerDependencies registers the dependencies of a component based on its constructor arguments.
func (d *DefaultContainer) registerDependencies(name string, args []Arg) {
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
func (d *DefaultContainer) registerDependency(name, dependency string) {
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
func (d *DefaultContainer) initialize(ctx context.Context, instance any) (any, error) {
	result, err := d.applyPreProcessors(ctx, instance)
	if err != nil {
		return nil, fmt.Errorf("apply pre-processors: %w", err)
	}

	if initializer, ok := result.(Initializer); ok {
		err = initializer.Init(ctx)

		if err != nil {
			return nil, fmt.Errorf("invoke init: %w", err)
		}
	}

	result, err = d.applyPostProcessors(ctx, result)
	if err != nil {
		return nil, fmt.Errorf("apply post-processors: %w", err)
	}

	return result, nil
}

// applyPreProcessors executes all registered PreProcessor hooks on the instance.
// Returns the processed object or an error.
func (d *DefaultContainer) applyPreProcessors(ctx context.Context, instance any) (any, error) {
	d.muProcessors.RLock()
	defer d.muProcessors.RUnlock()

	for _, processor := range d.preProcessors {
		result, err := processor.ProcessBeforeInit(ctx, instance)

		if err != nil {
			return nil, fmt.Errorf("pre-processor (%T): %w", processor, err)
		}

		if result == nil {
			return nil, fmt.Errorf("pre-processor (%T) returned nil", processor)
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
			return nil, fmt.Errorf("post-processor (%T): %w", processor, err)
		}

		if result == nil {
			return nil, fmt.Errorf("post-processor (%T) returned nil", processor)
		}

		instance = result
	}

	return instance, nil
}
