package container

import (
	"codnect.io/reflector"
	"fmt"
	"sync"
)

type DefinitionRegistry interface {
	Register(def *Definition) error
	Remove(name string) error
	Contains(name string) bool
	Find(name string) (*Definition, bool)
	Definitions() []*Definition
	DefinitionNames() []string
	DefinitionNamesByType(requiredType reflector.Type) []string
	Count() int
}

type definitionRegistry struct {
	definitionMap map[string]*Definition
	muDefinitions sync.RWMutex
}

func NewDefinitionRegistry(defs []*Definition) DefinitionRegistry {
	registry := &definitionRegistry{
		definitionMap: make(map[string]*Definition),
		muDefinitions: sync.RWMutex{},
	}

	for _, def := range defs {
		registry.definitionMap[def.Name()] = def
	}

	return registry
}

func (c *definitionRegistry) Register(def *Definition) error {
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

func (c *definitionRegistry) Remove(name string) error {
	defer c.muDefinitions.Unlock()
	c.muDefinitions.Lock()

	if _, exists := c.definitionMap[name]; !exists {
		return fmt.Errorf("container: no found definition with name %s", name)
	}

	delete(c.definitionMap, name)
	return nil
}

func (c *definitionRegistry) Contains(name string) bool {
	defer c.muDefinitions.Unlock()
	c.muDefinitions.Lock()

	_, exists := c.definitionMap[name]
	return exists
}

func (c *definitionRegistry) Find(name string) (*Definition, bool) {
	defer c.muDefinitions.Unlock()
	c.muDefinitions.Lock()

	if def, exists := c.definitionMap[name]; exists {
		return def, true
	}

	return nil, false
}

func (c *definitionRegistry) Definitions() []*Definition {
	defer c.muDefinitions.Unlock()
	c.muDefinitions.Lock()

	defs := make([]*Definition, 0)
	for _, def := range c.definitionMap {
		defs = append(defs, def)
	}

	return defs
}

func (c *definitionRegistry) DefinitionNames() []string {
	defer c.muDefinitions.Unlock()
	c.muDefinitions.Lock()

	names := make([]string, 0)
	for name, _ := range c.definitionMap {
		names = append(names, name)
	}

	return names
}

func (c *definitionRegistry) DefinitionNamesByType(requiredType reflector.Type) []string {
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

func (c *definitionRegistry) Count() int {
	defer c.muDefinitions.Unlock()
	c.muDefinitions.Lock()

	return len(c.definitionMap)
}
