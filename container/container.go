package container

import (
	"context"
	"errors"
	"fmt"
	"github.com/procyon-projects/reflector"
	"reflect"
	"strings"
	"sync"
)

type Container interface {
	Start() error
	DefinitionRegistry() DefinitionRegistry
	SharedInstances() SharedInstances
	Hooks() Hooks
	Get(ctx context.Context, name string) (any, error)
	GetByNameAndType(ctx context.Context, name string, typ reflector.Type) (any, error)
	GetByNameAndArgs(ctx context.Context, name string, args ...any) (any, error)
	GetByType(ctx context.Context, typ reflector.Type) (any, error)
	GetInstancesByType(ctx context.Context, requiredType reflector.Type) ([]any, error)
	Contains(name string) bool
	IsShared(name string) bool
	IsPrototype(name string) bool
	RegisterResolvable(typ reflector.Type, object any) error
	ScopeNames() []string
	GetScope(scopeName string) (Scope, error)
}

type container struct {
	definitionRegistry DefinitionRegistry
	sharedInstances    SharedInstances
	hooks              Hooks

	scopes   map[string]Scope
	muScopes *sync.RWMutex

	resolvableDependencies   map[reflect.Type]any
	muResolvableDependencies *sync.RWMutex
}

func New() Container {
	return WithDefinitions(nil)
}

func WithDefinitions(definitions []*Definition) Container {
	return &container{
		definitionRegistry:       NewDefinitionRegistry(definitions),
		sharedInstances:          NewSharedInstances(),
		hooks:                    NewHooks(),
		scopes:                   map[string]Scope{},
		muScopes:                 &sync.RWMutex{},
		resolvableDependencies:   map[reflect.Type]any{},
		muResolvableDependencies: &sync.RWMutex{},
	}
}

func (c *container) Start() error {
	ctx := contextWithHolder(context.Background())

	for _, name := range c.definitionRegistry.DefinitionNames() {
		definition, ok := c.definitionRegistry.Find(name)
		if ok && definition.IsShared() {
			_, err := c.Get(ctx, definition.Name())

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *container) DefinitionRegistry() DefinitionRegistry {
	return c.definitionRegistry
}

func (c *container) SharedInstances() SharedInstances {
	return c.sharedInstances
}

func (c *container) Hooks() Hooks {
	return c.hooks
}

func (c *container) Get(ctx context.Context, name string) (any, error) {
	return c.getInstance(ctx, name, nil)
}

func (c *container) GetByNameAndType(ctx context.Context, name string, typ reflector.Type) (any, error) {
	return c.getInstance(ctx, name, typ)
}

func (c *container) GetByNameAndArgs(ctx context.Context, name string, args ...any) (any, error) {
	return c.getInstance(ctx, name, nil, args...)
}

func (c *container) GetByType(ctx context.Context, typ reflector.Type) (any, error) {
	return c.getInstance(ctx, "", typ)
}

func (c *container) GetInstancesByType(ctx context.Context, requiredType reflector.Type) ([]any, error) {
	if requiredType == nil {
		return nil, errors.New("container: requiredType cannot be nil")
	}

	instances, err := c.getInstances(ctx, reflector.ToSlice(reflector.TypeOf[[]any]()), requiredType)

	if err != nil {
		return nil, err
	}

	return instances.([]any), nil
}

func (c *container) Contains(name string) bool {
	return c.sharedInstances.Contains(name)
}

func (c *container) IsShared(name string) bool {
	def, ok := c.definitionRegistry.Find(name)
	return ok && def.IsShared()
}

func (c *container) IsPrototype(name string) bool {
	def, ok := c.definitionRegistry.Find(name)
	return ok && def.IsPrototype()
}

func (c *container) RegisterResolvable(typ reflector.Type, object any) error {
	if typ == nil {
		return errors.New("container: type should be passed")
	}

	if object == nil {
		return errors.New("container: object should be passed")
	}

	defer c.muResolvableDependencies.Unlock()
	c.muResolvableDependencies.Lock()

	c.resolvableDependencies[typ.ReflectType()] = object
	return nil
}

func (c *container) RegisterScope(scopeName string, scope Scope) error {
	if strings.TrimSpace(scopeName) == "" {
		panic("container: scopeName cannot be empty or blank")
	}

	if scope == nil {
		panic("container: scope cannot be nil")
	}

	if SharedScope != scopeName && PrototypeScope != scopeName {
		defer c.muScopes.Unlock()
		c.muScopes.Lock()
		c.scopes[scopeName] = scope
		return nil
	}

	return errors.New("container: cannot replace existing scopes 'shared' and 'prototype'")
}

func (c *container) ScopeNames() []string {
	defer c.muScopes.Unlock()
	c.muScopes.Lock()

	scopeNames := make([]string, 0)
	for scopeName, _ := range c.scopes {
		scopeNames = append(scopeNames, scopeName)
	}

	return scopeNames
}

func (c *container) GetScope(scopeName string) (Scope, error) {
	defer c.muScopes.Unlock()
	c.muScopes.Lock()
	if scope, ok := c.scopes[scopeName]; ok {
		return scope, nil
	}

	return nil, fmt.Errorf("container: no scope registered for scope name %s", scopeName)
}

func (c *container) getInstance(ctx context.Context, name string, requiredType reflector.Type, args ...any) (any, error) {
	if ctx == nil {
		return nil, errors.New("container: context should be passed")
	}

	ctx = contextWithHolder(ctx)

	if name == "" && requiredType == nil {
		return nil, errors.New("container: either name or requiredType should be given")
	}

	if name == "" {
		candidate, err := c.sharedInstances.FindByType(requiredType)

		if err == nil {
			return candidate, nil
		}

		names := c.definitionRegistry.DefinitionNamesByType(requiredType)

		if len(names) == 0 {
			return nil, &notFoundError{
				ErrorString: fmt.Sprintf("container: not found instance or definition with required type %s", requiredType.Name()),
			}
		} else if len(names) > 1 {
			return nil, fmt.Errorf("container: there is more than one definition for the required type %s, it cannot be distinguished", requiredType.Name())
		}

		name = names[0]
	}

	def, ok := c.definitionRegistry.Find(name)

	if !ok {
		return nil, &notFoundError{
			ErrorString: fmt.Sprintf("container: not found definition with name %s", name),
		}
	}

	if requiredType != nil && !c.match(def.Type(), requiredType) {
		return nil, fmt.Errorf("container: definition type with name %s does not match the required type", name)
	}

	if def.IsShared() {
		instance, err := c.sharedInstances.OrElseGet(name, func() (any, error) {
			return c.createInstance(ctx, def, args)
		})

		return instance, err
	} else if def.IsPrototype() {
		prototypeHolder := holderFromContext(ctx)
		err := prototypeHolder.beforeCreation(name)

		if err != nil {
			return nil, err
		}

		defer prototypeHolder.afterCreation(name)
		return c.createInstance(ctx, def, args)
	}

	if strings.TrimSpace(def.Scope()) == "" {
		return nil, fmt.Errorf("container: no scope name for required type %s", requiredType.Name())
	}

	scope, err := c.GetScope(def.Scope())

	if err != nil {
		return nil, err
	}

	return scope.GetObject(ctx, name, func(ctx context.Context) (any, error) {
		scopeHolder := holderFromContext(ctx)
		err = scopeHolder.beforeCreation(name)

		if err != nil {
			return nil, err
		}

		defer scopeHolder.afterCreation(name)
		return c.createInstance(ctx, def, args)
	})
}

func (c *container) match(instanceType reflector.Type, requiredType reflector.Type) bool {
	if instanceType.CanConvert(requiredType) {
		return true
	} else if reflector.IsPointer(instanceType) && !reflector.IsPointer(requiredType) && !reflector.IsInterface(requiredType) {
		ptrType := reflector.ToPointer(instanceType)

		if ptrType.Elem().CanConvert(requiredType) {
			return true
		}
	}

	return false
}

func (c *container) createInstance(ctx context.Context, definition *Definition, args []any) (instance any, err error) {
	newFunc := definition.constructorFunc
	parameterCount := len(definition.Inputs())

	if parameterCount != 0 && len(args) == 0 {
		var resolvedArguments []any
		resolvedArguments, err = c.resolveInputs(ctx, definition.Inputs())

		if err != nil {
			return nil, err
		}

		var results []any
		results, err = newFunc.Invoke(resolvedArguments...)

		if err != nil {
			return nil, err
		}

		instance = results[0]
	} else if (parameterCount == 0 && len(args) == 0) || (len(args) != 0 && parameterCount == len(args)) {
		var results []any
		results, err = newFunc.Invoke(args...)

		if err != nil {
			return nil, err
		}

		instance = results[0]
	} else {
		return nil, fmt.Errorf("container: the number of provided arguments is wrong for definition %s", definition.Name())
	}

	return c.initializeInstance(definition.name, instance)
}

func (c *container) getInstances(ctx context.Context, sliceType reflector.Slice, itemType reflector.Type) (any, error) {
	val, err := sliceType.Instantiate()

	if err != nil {
		return nil, err
	}

	var (
		instance any
		items    any
	)

	instances := c.sharedInstances.FindAllByType(itemType)

	sliceType = reflector.ToSlice(reflector.ToPointer(reflector.TypeOfAny(val.Val())).Elem())
	items, err = sliceType.Append(instances...)

	if err != nil {
		return nil, err
	}

	definitionNames := c.definitionRegistry.DefinitionNamesByType(itemType)

	for _, definitionName := range definitionNames {
		if c.sharedInstances.Contains(definitionName) {
			continue
		}

		instance, err = c.Get(ctx, definitionName)

		if err != nil {
			return nil, err
		}

		items, err = sliceType.Append(instance)
	}

	return items, nil
}

func (c *container) getResolvableInstance(typ reflector.Type) (any, bool) {
	defer c.muResolvableDependencies.Unlock()
	c.muResolvableDependencies.Lock()

	for resolvableType, instance := range c.resolvableDependencies {
		if resolvableType.ConvertibleTo(typ.ReflectType()) {
			return instance, true
		}
	}

	return nil, false
}

func (c *container) resolveInputs(ctx context.Context, inputs []*Input) ([]any, error) {
	arguments := make([]any, 0)
	for _, input := range inputs {

		if reflector.IsSlice(input.Type()) {
			sliceType := reflector.ToSlice(input.Type())
			instances, err := c.getInstances(ctx, sliceType, sliceType.Elem())

			if err != nil {
				return nil, err
			}

			arguments = append(arguments, instances)
			continue
		}

		var (
			instance any
			err      error
		)

		resolvableInstance, exists := c.getResolvableInstance(input.Type())
		if exists {
			arguments = append(arguments, resolvableInstance)
			continue
		}

		if input.Name() != "" {
			instance, err = c.Get(ctx, input.Name())
		} else {
			instance, err = c.GetByType(ctx, input.Type())
		}

		if err != nil {
			if notFoundErr := (*notFoundError)(nil); errors.As(err, &notFoundErr) && !reflector.IsPointer(input.Type()) && input.Type().IsInstantiable() {
				var val reflector.Value
				val, err = input.Type().Instantiate()

				if err != nil {
					return nil, err
				}

				instance = val.Elem()
				arguments = append(arguments, instance)
				continue
			}

			if !input.IsOptional() && err != nil {
				return nil, err
			} else if input.IsOptional() && err != nil {
				arguments = append(arguments, nil)
			}
		} else {
			arguments = append(arguments, instance)
		}
	}

	return arguments, nil
}

func (c *container) initializeInstance(name string, instance any) (any, error) {
	var (
		err error
	)
	hooks := c.Hooks().ToSlice()

	for _, hook := range hooks {
		if hook.OnPreInitialization != nil {
			instance, err = hook.OnPreInitialization(name, instance)
		}

		if err != nil {
			return nil, err
		}
	}

	if postConstructor, implements := instance.(PostConstructor); implements {
		err = postConstructor.PostConstruct()

		if err != nil {
			return nil, err
		}
	}

	for _, hook := range hooks {
		if hook.OnPostInitialization != nil {
			instance, err = hook.OnPostInitialization(name, instance)

			if err != nil {
				return nil, err
			}
		}
	}

	return instance, nil
}
