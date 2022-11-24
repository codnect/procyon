package container

import (
	"fmt"
	"sync"
)

var (
	definitions   = make(map[string]*Definition)
	muDefinitions = sync.RWMutex{}
)

func Register(constructor Constructor, options ...Option) {
	defer muDefinitions.Unlock()
	muDefinitions.Lock()

	def, err := MakeDefinition(constructor, options...)

	if err != nil {
		panic(fmt.Sprintf("container: %s", err))
	}

	if _, exists := definitions[def.Name()]; exists {
		panic(fmt.Sprintf("container: definition with name %s already exists", def.Name()))
	}

	definitions[def.Name()] = def
}

func copyDefinitions() map[string]*Definition {
	defer muDefinitions.Unlock()
	muDefinitions.Lock()

	m := make(map[string]*Definition)

	for name, def := range definitions {
		m[name] = def
	}

	return m
}
