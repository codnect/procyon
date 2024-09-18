package component

import (
	"codnect.io/procyon/component/condition"
	"codnect.io/procyon/component/container"
	"context"
)

// Loader is a struct that represents a component loader.
// It uses a container to register components and an evaluator to evaluate conditions.
type Loader struct {
	container container.Container
	evaluator condition.Evaluator
}

// NewLoader function creates a new Loader instance with the provided container.
func NewLoader(container container.Container) *Loader {
	return &Loader{
		container: container,
		evaluator: condition.NewEvaluator(container),
	}
}

// LoadComponents method loads components into the container.
// It evaluates the conditions of each component and registers the component if the conditions are met.
// If the conditions are not met, the component is skipped.
// The method returns an error if the registration of a component fails.
func (l *Loader) LoadComponents(ctx context.Context, components []*Component) error {
	skippedComponents := make([]*Component, 0)

	for _, component := range components {
		if l.evaluator.Evaluate(ctx, component.Conditions()) {
			err := l.container.Definitions().Register(component.Definition())

			if err != nil {
				return err
			}
		} else {
			skippedComponents = append(skippedComponents, component)
		}
	}

	if len(components) == len(skippedComponents) {
		return nil
	}

	return l.LoadComponents(ctx, skippedComponents)
}
