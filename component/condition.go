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
)

// Condition represents a rule that determines whether a component should be included at runtime.
// It is evaluated during the component loading phase.
type Condition interface {
	// Matches returns true if the condition is satisfied in the given context.
	Matches(ctx context.Context, container Container) bool
}

// conditionEvaluator evaluates a set of conditions.
type conditionEvaluator struct {
	container Container
}

// newConditionEvaluator creates a new conditionEvaluator.
func newConditionEvaluator(container Container) *conditionEvaluator {
	if container == nil {
		panic("nil container")
	}

	return &conditionEvaluator{
		container: container,
	}
}

// evaluate returns true if all given conditions match.
func (e *conditionEvaluator) evaluate(ctx context.Context, conditions []Condition) bool {
	if len(conditions) == 0 {
		return true
	}

	for _, condition := range conditions {
		if !condition.Matches(ctx, e.container) {
			return false
		}
	}

	return true
}
