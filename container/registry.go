package container

import (
	"fmt"
	"github.com/procyon-projects/reflector"
	"sync"
)

type DefinitionRegistry struct {
	definitionMap map[string]*Definition
	muDefinitions sync.RWMutex
}

func NewDefinitionRegistry(defs []*Definition) *DefinitionRegistry {
	registry := &DefinitionRegistry{
		definitionMap: make(map[string]*Definition),
		muDefinitions: sync.RWMutex{},
	}

	for _, def := range defs {
		registry.definitionMap[def.Name()] = def
	}

	return registry
}

func (c *DefinitionRegistry) Add(def *Definition) error {
	if def == nil {
		return fmt.Errorf("container: definition should not be nil")
	}

	defer c.muDefinitions.Unlock()
	c.muDefinitions.Lock()

	if _, exists := c.definitionMap[def.Name()]; exists {
		return fmt.Errorf("container: definition with name %s already exists", def.Name())
	}

	c.definitionMap[def.Name()] = def

	return nil
}

func (c *DefinitionRegistry) Remove(name string) error {
	defer c.muDefinitions.Unlock()
	c.muDefinitions.Lock()

	if _, exists := c.definitionMap[name]; !exists {
		return fmt.Errorf("container: no found definition with name %s", name)
	}

	delete(c.definitionMap, name)
	return nil
}

func (c *DefinitionRegistry) Contains(name string) bool {
	defer c.muDefinitions.Unlock()
	c.muDefinitions.Lock()

	_, exists := c.definitionMap[name]
	return exists
}

func (c *DefinitionRegistry) Find(name string) (*Definition, bool) {
	defer c.muDefinitions.Unlock()
	c.muDefinitions.Lock()

	if def, exists := c.definitionMap[name]; exists {
		return def, true
	}

	return nil, false
}

func (c *DefinitionRegistry) Definitions() []*Definition {
	defer c.muDefinitions.Unlock()
	c.muDefinitions.Lock()

	defs := make([]*Definition, 0)
	for _, def := range c.definitionMap {
		defs = append(defs, def)
	}

	return defs
}

func (c *DefinitionRegistry) DefinitionNames() []string {
	defer c.muDefinitions.Unlock()
	c.muDefinitions.Lock()

	names := make([]string, 0)
	for name, _ := range c.definitionMap {
		names = append(names, name)
	}

	return names
}

func (c *DefinitionRegistry) DefinitionNamesByType(requiredType reflector.Type) []string {
	defer c.muDefinitions.Unlock()
	c.muDefinitions.Lock()

	names := make([]string, 0)

	if requiredType == nil {
		return names
	}

	for name, def := range c.definitionMap {

		instanceType := def.Type()

		if instanceType.CanConvert(requiredType) {
			names = append(names, name)
		} else if reflector.IsPointer(instanceType) && !reflector.IsPointer(requiredType) && !reflector.IsInterface(requiredType) {
			ptrType := reflector.ToPointer(instanceType)

			if ptrType.Elem().CanConvert(requiredType) {
				names = append(names, name)
			}
		}

	}

	return names
}

func (c *DefinitionRegistry) Count() int {
	defer c.muDefinitions.Unlock()
	c.muDefinitions.Lock()

	return len(c.definitionMap)
}
