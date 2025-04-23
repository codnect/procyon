package component

import "context"

// ConditionContext provides contextual information for condition evaluation.
// It extends context.Context and gives access to the current container.
type ConditionContext interface {
	context.Context

	// Container returns the container associated with the evaluation context.
	Container() Container
}

// Condition represents a rule that determines whether a component should be included at runtime.
// It is evaluated during the component loading phase.
type Condition interface {
	// Matches returns true if the condition is satisfied in the given context.
	Matches(ctx ConditionContext) bool
}

// ConditionEvaluator evaluates conditions attached to a component.
type ConditionEvaluator interface {
	// Evaluate returns true if all conditions of the given component match.
	Evaluate(ctx context.Context, component *Component) bool
}
