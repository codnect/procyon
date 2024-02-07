package procyon

import (
	"codnect.io/procyon-core/component"
	"codnect.io/procyon-core/container"
	"codnect.io/reflector"
)

func getComponentsByType[T any](c container.Container, args ...any) ([]T, error) {
	customizers := make([]T, 0)
	rType := reflector.TypeOf[T]()

	registry := c.DefinitionRegistry()
	definitionNames := registry.DefinitionNamesByType(rType)

	for _, name := range definitionNames {
		definition, _ := registry.Find(name)

		if len(definition.Inputs()) != len(args) {
			continue
		}

		if !matchInputTypes(definition.Inputs(), args...) {
			continue
		}

		results, err := definition.Constructor().Invoke(args...)

		if err != nil {
			return nil, err
		}

		customizer := results[0].(T)
		customizers = append(customizers, customizer)
	}

	return customizers, nil
}

func matchInputTypes(inputs []*container.Input, args ...any) bool {
	for index, input := range inputs {
		rArg := reflector.TypeOfAny(args[index]).ReflectType()

		if !rArg.ConvertibleTo(input.Type().ReflectType()) {
			return false
		}
	}

	return true
}

func registerComponentDefinitions(container container.Container) error {
	for _, registeredComponent := range component.RegisteredComponents() {
		err := container.DefinitionRegistry().Register(registeredComponent.Definition())
		if err != nil {
			return err
		}
	}

	return nil
}
