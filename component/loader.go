package component

import (
	"context"
	"fmt"
)

// Loader defines the interface responsible for loading registered components.
// It evaluates runtime conditions and only loads components that match.
type Loader interface {
	// Load evaluates each component's conditions and registers the eligible ones.
	// Returns an error if the loading process fails.
	Load(ctx context.Context) error
}

// ConditionalLoader loads component definitions into a container
// only if their associated runtime conditions are satisfied.
type ConditionalLoader struct {
	container  Container
	components []*Component
	evaluator  *ConditionEvaluator
}

// NewConditionalLoader creates a ConditionalLoader with the given container and components.
func NewConditionalLoader(container Container, components []*Component) *ConditionalLoader {
	if container == nil {
		panic("nil container")
	}

	return &ConditionalLoader{
		container:  container,
		components: components,
		evaluator:  NewConditionEvaluator(container),
	}
}

// Load evaluates the conditions of each component and registers its definition into the container
// only if all conditions are satisfied. Components that fail condition checks are skipped.
// Returns an error if any eligible component fails to register.
func (l *ConditionalLoader) Load(ctx context.Context) error {
	skipped := make([]*Component, 0)

	for _, comp := range l.components {
		if !l.evaluator.Evaluate(ctx, comp.Conditions()) {
			skipped = append(skipped, comp)
			continue
		}

		def := comp.definition
		if err := l.container.RegisterDefinition(def); err != nil {
			return fmt.Errorf("failed to register component %q: %w", def.Name(), err)
		}
	}

	if len(skipped) > 0 && len(skipped) < len(l.components) {
		loader := NewConditionalLoader(l.container, skipped)
		return loader.Load(ctx)
	}

	return nil
}
