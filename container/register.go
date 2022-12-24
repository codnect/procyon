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
		panic(err)
	}

	if _, exists := definitions[def.Name()]; exists {
		panic(fmt.Sprintf("container: definition with name %s already exists", def.Name()))
	}

	definitions[def.Name()] = def
}

func RegisteredDefinitions() []*Definition {
	defer muDefinitions.Unlock()
	muDefinitions.Lock()

	copyOfDefinitions := make([]*Definition, 0)

	for _, def := range definitions {
		copyOfDefinitions = append(copyOfDefinitions, def)
	}

	return copyOfDefinitions
}
