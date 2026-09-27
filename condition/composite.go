// Copyright 2026 Codnect
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

package condition

import (
	"context"

	"go.codnect.io/procyon/component"
)

// compositeMode defines how the conditions in a Composite are evaluated.
type compositeMode uint8

const (
	compositeAll compositeMode = iota
	compositeAny
	compositeNone
)

// Composite combines multiple conditions using a logical operation.
type Composite struct {
	mode       compositeMode
	conditions []component.Condition
}

// All creates a composite condition that matches when all the given
// conditions match.
//
// An empty All condition matches by definition.
func All(conditions ...component.Condition) *Composite {
	return &Composite{
		mode:       compositeAll,
		conditions: conditions,
	}
}

// Any creates a composite condition that matches when at least one of the
// given conditions matches.
//
// An empty Any condition does not match.
func Any(conditions ...component.Condition) *Composite {
	return &Composite{
		mode:       compositeAny,
		conditions: conditions,
	}
}

// None creates a composite condition that matches when none of the given
// conditions match.
//
// An empty None condition matches by definition.
func None(conditions ...component.Condition) *Composite {
	return &Composite{
		mode:       compositeNone,
		conditions: conditions,
	}
}

// Matches reports whether the composite condition matches according to its
// logical operation.
func (c *Composite) Matches(ctx context.Context, container component.Container) bool {
	for _, condition := range c.conditions {
		matched := condition.Matches(ctx, container)

		switch c.mode {
		case compositeAll:
			if !matched {
				return false
			}
		case compositeAny:
			if matched {
				return true
			}
		case compositeNone:
			if matched {
				return false
			}
		}
	}

	return c.mode != compositeAny
}
