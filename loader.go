package procyon

import (
	"codnect.io/procyon-core/component"
	"context"
)

type componentLoader struct {
	container component.Container
	evaluator component.ConditionEvaluator
}

func newComponentLoader(container component.Container) *componentLoader {
	return &componentLoader{
		container: container,
		evaluator: component.NewConditionEvaluator(container),
	}
}

func (l *componentLoader) loadDefinitions(ctx context.Context) error {
	components := component.List()

	for _, component := range components {
		if l.evaluator.Evaluate(ctx, component.Conditions()) {

			err := l.container.Definitions().Register(component.Definition())

			if err != nil {
				return err
			}
		}
	}

	return nil
}
