package component

import (
	"codnect.io/procyon/component/condition"
	"codnect.io/procyon/component/container"
	"codnect.io/procyon/component/filter"
	"fmt"
	"reflect"
	"sync"
)

// components is a map that stores all registered components.
// The key is the name of the component and the value is the Component instance.
var (
	components   = make(map[string]*Component)
	muComponents = sync.RWMutex{}
)

// Component struct represents a component in the application.
// It contains a component definition and a list of conditions that must be met for the component to be used.
type Component struct {
	definition        *container.Definition
	definitionOptions []container.DefinitionOption
	conditions        []condition.Condition
}

// Register function registers a new component with options.
func Register(constructorFunc container.ConstructorFunc, options ...Option) {
	defer muComponents.Unlock()
	muComponents.Lock()

	component := createComponent(constructorFunc, options...)
	componentName := component.Definition().Name()

	if _, exists := components[componentName]; exists {
		panic(fmt.Errorf("component with name '%s' already exists", componentName))
	}

	components[componentName] = component
}

// List function returns a list of components that match the provided filters.
func List(filters ...filter.Filter) []*Component {
	defer muComponents.Unlock()
	muComponents.Lock()

	filterOpts := filter.Of(filters...)
	componentList := make([]*Component, 0)

	for _, component := range components {
		definition := component.Definition()

		if filterOpts.Name != "" && filterOpts.Name != component.Definition().Name() {
			continue
		}

		if filterOpts.Type == nil {
			componentList = append(componentList, component)
			continue
		}

		if convertibleTo(definition.Type(), filterOpts.Type) {
			componentList = append(componentList, component)
		}
	}

	return componentList
}

// createComponent function creates a new component with options.
func createComponent(constructorFunc container.ConstructorFunc, options ...Option) *Component {
	component := &Component{
		definitionOptions: make([]container.DefinitionOption, 0),
		conditions:        make([]condition.Condition, 0),
	}

	err := applyComponentOptions(component, options)
	if err != nil {
		panic(err)
	}

	component.definition, err = container.MakeDefinition(constructorFunc, component.definitionOptions...)

	if err != nil {
		panic(err)
	}

	return component
}

// Definition function returns the Definition of the component.
func (c *Component) Definition() *container.Definition {
	return c.definition
}

// Conditions function returns a list of conditions that must be met.
func (c *Component) Conditions() []condition.Condition {
	copyOfConditions := make([]condition.Condition, 0)

	for _, condition := range c.conditions {
		copyOfConditions = append(copyOfConditions, condition)
	}

	return copyOfConditions
}

// applyComponentOptions applies the options to the component
func applyComponentOptions(component *Component, options []Option) error {
	for _, option := range options {
		err := option(component)
		if err != nil {
			return err
		}
	}
	return nil
}

// convertibleTo function checks if a source type can be converted to a target type.
func convertibleTo(sourceType reflect.Type, targetType reflect.Type) bool {
	if sourceType == targetType || (targetType.Kind() == reflect.Interface && sourceType.ConvertibleTo(targetType)) {
		return true
	} else if sourceType.Kind() == reflect.Pointer {
		return convertibleTo(sourceType.Elem(), targetType)
	}

	return false
}
