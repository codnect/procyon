package container

import (
	"errors"
	"fmt"
	"github.com/procyon-projects/reflector"
	"sync"
)

type InstanceRegistry struct {
	instances              map[string]any
	instancesInPreparation map[string]struct{}
	typesOfInstances       map[string]reflector.Type
	muInstances            sync.RWMutex
}

func NewInstanceRegistry() *InstanceRegistry {
	return &InstanceRegistry{
		instances:              make(map[string]any),
		instancesInPreparation: make(map[string]struct{}),
		typesOfInstances:       make(map[string]reflector.Type),
		muInstances:            sync.RWMutex{},
	}
}

func (c *InstanceRegistry) Add(name string, instance any) error {
	defer c.muInstances.Unlock()
	c.muInstances.Lock()

	if _, exists := c.instances[name]; exists {
		return fmt.Errorf("container: instance with name %s already exists", name)
	}

	c.instances[name] = instance
	c.typesOfInstances[name] = reflector.TypeOfAny(instance)
	return nil
}

func (c *InstanceRegistry) Find(name string) (any, bool) {
	defer c.muInstances.Unlock()
	c.muInstances.Lock()

	if instance, exists := c.instances[name]; exists {
		return instance, true
	}

	return nil, false
}

func (c *InstanceRegistry) Contains(name string) bool {
	defer c.muInstances.Unlock()
	c.muInstances.Lock()

	_, exists := c.instances[name]
	return exists
}

func (c *InstanceRegistry) InstanceNames() []string {
	defer c.muInstances.Unlock()
	c.muInstances.Lock()

	names := make([]string, 0)
	for name, _ := range c.instances {
		names = append(names, name)
	}

	return names
}

func (c *InstanceRegistry) FindByType(requiredType *Type) (any, error) {
	if requiredType == nil {
		return nil, errors.New("container: requiredType cannot be nil")
	}

	instances := c.FindAllByType(requiredType)
	if len(instances) > 1 {
		return nil, fmt.Errorf("container: instances cannot be distinguished for required type %s", requiredType.Name())
	}

	if len(instances) == 0 {
		return nil, fmt.Errorf("container: not found any instance of type %s", requiredType.Name())
	}

	return instances[0], nil
}

func (c *InstanceRegistry) FindAllByType(requiredType *Type) []any {
	defer c.muInstances.Unlock()
	c.muInstances.Lock()

	instances := make([]any, 0)

	for name, typ := range c.typesOfInstances {

		if typ.CanConvert(requiredType.typ) {
			instances = append(instances, c.instances[name])
		} else if reflector.IsPointer(typ) && !reflector.IsPointer(requiredType.typ) && !reflector.IsInterface(requiredType.typ) {
			ptrType := reflector.ToPointer(typ)

			if ptrType.Elem().CanConvert(requiredType.typ) {
				val, err := ptrType.Elem().Value()

				if err == nil {
					instances = append(instances, val)
				}
			}
		}

	}

	return instances
}

func (c *InstanceRegistry) OrElseGet(name string, supplier func() (any, error)) (any, error) {
	instance, ok := c.Find(name)

	if ok {
		return instance, nil
	}

	err := c.putToPreparation(name)

	if err != nil {
		return nil, err
	}

	defer func() {
		c.removeFromPreparation(name)
	}()

	instance, err = supplier()
	return instance, err
}

func (c *InstanceRegistry) Count() int {
	defer c.muInstances.Unlock()
	c.muInstances.Lock()

	return len(c.instances)
}

func (c *InstanceRegistry) putToPreparation(name string) error {
	defer c.muInstances.Unlock()
	c.muInstances.Lock()

	if _, ok := c.instancesInPreparation[name]; ok {
		return fmt.Errorf("container: instance with name %s is currently in preparation, maybe it has got circular dependency cycle", name)
	}

	c.instancesInPreparation[name] = struct{}{}
	return nil
}

func (c *InstanceRegistry) removeFromPreparation(name string) {
	defer c.muInstances.Unlock()
	c.muInstances.Lock()
	delete(c.instancesInPreparation, name)
}

type DefinitionRegistry struct {
	definitions   map[string]*Definition
	muDefinitions sync.RWMutex
}

func NewDefinitionRegistry() *DefinitionRegistry {
	return &DefinitionRegistry{
		definitions:   copyDefinitions(),
		muDefinitions: sync.RWMutex{},
	}
}

func (c *DefinitionRegistry) Add(def *Definition) error {
	if def == nil {
		return fmt.Errorf("container: definition should not be nil")
	}

	defer c.muDefinitions.Unlock()
	c.muDefinitions.Lock()

	if _, exists := c.definitions[def.Name()]; exists {
		return fmt.Errorf("container: definition with name %s already exists", def.Name())
	}

	c.definitions[def.Name()] = def

	return nil
}

func (c *DefinitionRegistry) Remove(name string) error {
	defer c.muDefinitions.Unlock()
	c.muDefinitions.Lock()

	if _, exists := c.definitions[name]; !exists {
		return fmt.Errorf("container: no found definition with name %s", name)
	}

	delete(c.definitions, name)
	return nil
}

func (c *DefinitionRegistry) Contains(name string) bool {
	defer c.muDefinitions.Unlock()
	c.muDefinitions.Lock()

	_, exists := c.definitions[name]
	return exists
}

func (c *DefinitionRegistry) Find(name string) (*Definition, bool) {
	defer c.muDefinitions.Unlock()
	c.muDefinitions.Lock()

	if def, exists := c.definitions[name]; exists {
		return def, true
	}

	return nil, false
}

func (c *DefinitionRegistry) Definitions() []*Definition {
	defer c.muDefinitions.Unlock()
	c.muDefinitions.Lock()

	defs := make([]*Definition, 0)
	for _, def := range c.definitions {
		defs = append(defs, def)
	}

	return defs
}

func (c *DefinitionRegistry) DefinitionNames() []string {
	defer c.muDefinitions.Unlock()
	c.muDefinitions.Lock()

	names := make([]string, 0)
	for name, _ := range c.definitions {
		names = append(names, name)
	}

	return names
}

func (c *DefinitionRegistry) DefinitionNamesByType(requiredType *Type) []string {
	defer c.muDefinitions.Unlock()
	c.muDefinitions.Lock()

	names := make([]string, 0)

	if requiredType == nil {
		return names
	}

	for name, def := range c.definitions {

		instanceType := def.reflectorType()

		if instanceType.CanConvert(requiredType.typ) {
			names = append(names, name)
		} else if reflector.IsPointer(instanceType) && !reflector.IsPointer(requiredType.typ) && !reflector.IsInterface(requiredType.typ) {
			ptrType := reflector.ToPointer(instanceType)

			if ptrType.Elem().CanConvert(requiredType.typ) {
				names = append(names, name)
			}
		}

	}

	return names
}

func (c *DefinitionRegistry) Count() int {
	defer c.muDefinitions.Unlock()
	c.muDefinitions.Lock()

	return len(c.definitions)
}
