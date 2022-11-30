package container

import (
	"errors"
	"fmt"
	"github.com/procyon-projects/reflector"
)

type Container struct {
	definitionRegistry *DefinitionRegistry
	instanceRegistry   *InstanceRegistry
}

func New() *Container {
	return &Container{
		definitionRegistry: NewDefinitionRegistry(copyDefinitions()),
		instanceRegistry:   NewInstanceRegistry(),
	}
}

func (c *Container) Start() error {
	return nil
}

func (c *Container) DefinitionRegistry() *DefinitionRegistry {
	return c.definitionRegistry
}

func (c *Container) InstanceRegistry() *InstanceRegistry {
	return c.instanceRegistry
}

func (c *Container) Get(name string) (any, error) {
	return c.getInstance(name, nil)
}

func (c *Container) GetByNameAndType(name string, typ *Type) (any, error) {
	return c.getInstance(name, typ)
}

func (c *Container) GetByNameAndArgs(name string, args ...any) (any, error) {
	return c.getInstance(name, nil, args...)
}

func (c *Container) GetByType(typ *Type) (any, error) {
	return c.getInstance("", typ)
}

func (c *Container) Contains(name string) bool {
	return c.instanceRegistry.Contains(name)
}

func (c *Container) IsShared(name string) bool {
	return false
}

func (c *Container) IsPrototype(name string) bool {

	return false
}

func (c *Container) getInstance(name string, requiredType *Type, args ...any) (any, error) {
	if name == "" && requiredType == nil {
		return nil, errors.New("either name or typ should be given")
	}

	if name == "" {
		candidate, err := c.instanceRegistry.FindByType(requiredType)

		if err != nil {
			return nil, err
		}

		return candidate, nil
	}

	def, ok := c.definitionRegistry.Find(name)

	if !ok {
		return nil, fmt.Errorf("not found definition with name %s", name)
	}

	if requiredType != nil && !c.match(def.reflectorType(), requiredType.typ) {
		return nil, fmt.Errorf("definition type with name %s does not match the required type", name)
	}

	if def.IsShared() {
		instance, err := c.instanceRegistry.OrElseGet(name, func() (any, error) {
			return c.createInstance(def, args)
		})

		return instance, err
	} else if def.IsPrototype() {
		return c.createInstance(def, args)
	}

	return nil, nil
}

func (c *Container) match(instanceType reflector.Type, requiredType reflector.Type) bool {
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

func (c *Container) createInstance(definition *Definition, args []any) (instance any, err error) {
	newFunc := definition.constructorFunc
	parameterCount := len(definition.Inputs())

	if parameterCount != 0 && len(args) == 0 {
		var resolvedArguments []any
		resolvedArguments, err = c.resolveInputs(definition.Inputs())

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
		return nil, fmt.Errorf("the number of provided arguments is wrong for definition %s", definition.Name())
	}

	return c.initializeInstance(definition.name, instance)
}

func (c *Container) resolveInputs(inputs []*Input) ([]any, error) {
	arguments := make([]any, len(inputs))
	for _, input := range inputs {
		instance, err := c.Get(input.Name())

		if !input.IsOptional() && err != nil {
			return nil, err
		}

		if input.IsOptional() && err != nil {
			if input.reflectorType().IsInstantiable() {
				var val reflector.Value
				val, err = input.reflectorType().Instantiate()

				if err != nil {
					return nil, err
				}

				instance = val.Val()
			}
		}

		arguments = append(arguments, instance)
	}
	return nil, nil
}

func (c *Container) initializeInstance(name string, instance any) (any, error) {
	return nil, nil
}
