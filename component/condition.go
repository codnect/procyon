package component

import "context"

// Condition represents a rule that determines whether a component should be included at runtime.
// It is evaluated during the component loading phase.
type Condition interface {
	// MatchesCondition returns true if the condition is satisfied in the given context.
	MatchesCondition(ctx ConditionContext) bool
}

// ConditionContext provides contextual information for condition evaluation.
// It extends context.Context and gives access to the current container.
type ConditionContext interface {
	context.Context

	// Container returns the container associated with the evaluation context.
	Container() Container
}

// ConditionEvaluator evaluates a set of conditions in the given context.
type ConditionEvaluator interface {
	// Evaluate returns true if all given conditions match the current context.
	Evaluate(ctx context.Context, conditions []Condition) bool
}
