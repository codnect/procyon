package procyon

import (
	"codnect.io/procyon-core/component"
	"codnect.io/procyon-core/component/filter"
	"context"
)

type definitionLoader struct {
	container  component.Container
	evaluator  component.ConditionEvaluator
	components []*component.Component
}

func newDefinitionLoader(container component.Container) *definitionLoader {
	return &definitionLoader{
		container:  container,
		evaluator:  component.NewConditionEvaluator(container),
		components: component.List(),
	}
}

func (l *definitionLoader) load(ctx context.Context) error {
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

func initializeSingletons(ctx context.Context, container component.Container) error {

	for _, definition := range container.Definitions().List() {
		_, err := container.GetObject(ctx, filter.ByName(definition.Name()))

		if err != nil {
			return err
		}
	}

	return nil
}
