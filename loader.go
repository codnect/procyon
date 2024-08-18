package procyon

import (
	"codnect.io/procyon-core/component"
	"context"
)

type componentDefinitionLoader struct {
	container  component.Container
	evaluator  component.ConditionEvaluator
	components []*component.Component
}

func newComponentDefinitionLoader(container component.Container) *componentDefinitionLoader {
	return &componentDefinitionLoader{
		container:  container,
		evaluator:  component.NewConditionEvaluator(container),
		components: component.List(),
	}
}

func (l *componentDefinitionLoader) load(ctx context.Context) error {
	skippedComponents := make([]*component.Component, 0)

	for _, component := range l.components {
		if l.evaluator.Evaluate(ctx, component.Conditions()) {
			err := l.container.Definitions().Register(component.Definition())

			if err != nil {
				return err
			}
		} else {
			skippedComponents = append(skippedComponents, component)
		}
	}

	if len(l.components) == len(skippedComponents) {
		return nil
	}

	l.components = skippedComponents
	return l.load(ctx)
}
