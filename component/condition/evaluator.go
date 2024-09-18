package condition

import (
	"codnect.io/procyon/component/container"
	"context"
)

// Evaluator struct represents an evaluator that evaluates component conditions.
type Evaluator struct {
	container container.Container
}

// NewEvaluator function creates a new condition Evaluator.
func NewEvaluator(container container.Container) Evaluator {
	if container == nil {
		panic("nil container")
	}

	return Evaluator{
		container: container,
	}
}

// Evaluate method evaluates the given conditions.
// If all conditions match, it returns true. If any condition does not match, it returns false.
func (e Evaluator) Evaluate(ctx context.Context, conditions []Condition) bool {
	if len(conditions) == 0 {
		return true
	}

	conditionContext := NewContext(ctx, e.container)

	for _, condition := range conditions {
		if !condition.MatchesCondition(conditionContext) {
			return false
		}
	}

	return true
}
