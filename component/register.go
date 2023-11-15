package component

import (
	"fmt"
	"sync"
)

var (
	components   = make(map[string]*Component)
	muComponents = sync.RWMutex{}
)

func Register(constructor Constructor, options ...Option) *Component {
	defer muComponents.Unlock()
	muComponents.Lock()

	component, err := New(constructor, options...)

	if err != nil {
		panic(err)
	}

	componentName := component.Definition().Name()

	if _, exists := components[componentName]; exists {
		panic(fmt.Sprintf("component: component with name %s already exists", componentName))
	}

	components[componentName] = component
	return component
}

func RegisteredComponents() []*Component {
	defer muComponents.Unlock()
	muComponents.Lock()

	copyOfComponents := make([]*Component, 0)
	for _, component := range components {
		copyOfComponents = append(copyOfComponents, component)
	}

	return copyOfComponents
}
