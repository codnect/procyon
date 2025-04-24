// Copyright 2025 Codnect
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package component

import (
	"context"
	"time"
)

// ConditionContext provides runtime context and container access during condition evaluation.
type ConditionContext struct {
	ctx       context.Context
	container Container
}

// NewConditionContext creates a new ConditionContext with the given base context and container.
func NewConditionContext(ctx context.Context, container Container) ConditionContext {
	if ctx == nil {
		panic("nil context")
	}

	if container == nil {
		panic("nil container")
	}

	return ConditionContext{
		ctx:       ctx,
		container: container,
	}
}

// Deadline method returns the time when work done on behalf of
// this context should be canceled.
func (c ConditionContext) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

// Done method returns a channel that's closed when work done on behalf of
// this context should be canceled.
func (c ConditionContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

// Err method returns a non-nil error value after Done is closed.
func (c ConditionContext) Err() error {
	return c.ctx.Err()
}

// Value method returns the value associated with this context for key,
// or nil if no value is associated with key.
func (c ConditionContext) Value(key any) any {
	return c.ctx.Value(key)
}

// Container returns the container associated with this condition context.
func (c ConditionContext) Container() Container {
	return c.container
}

// Condition represents a rule that determines whether a component should be included at runtime.
// It is evaluated during the component loading phase.
type Condition interface {
	// Matches returns true if the condition is satisfied in the given context.
	Matches(ctx ConditionContext) bool
}

// ConditionEvaluator evaluates a set of conditions.
type ConditionEvaluator struct {
	container Container
}

// NewConditionEvaluator function creates a new ConditionEvaluator.
func NewConditionEvaluator(container Container) *ConditionEvaluator {
	if container == nil {
		panic("nil container")
	}

	return &ConditionEvaluator{
		container: container,
	}
}

// Evaluate returns true if all given conditions match.
func (e *ConditionEvaluator) Evaluate(ctx context.Context, conditions []Condition) bool {
	if len(conditions) == 0 {
		return true
	}

	conditionContext := NewConditionContext(ctx, e.container)

	for _, condition := range conditions {
		if !condition.Matches(conditionContext) {
			return false
		}
	}

	return true
}
