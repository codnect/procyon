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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.codnect.io/procyon/component"
)

type AnyCondition struct {
	matched bool
}

func (c AnyCondition) Matches(ctx context.Context, container component.Container) bool {
	return c.matched
}

func TestComposite_Matches(t *testing.T) {
	testCases := []struct {
		name       string
		condition  *Composite
		wantResult bool
	}{
		{
			name:       "matches when all conditions match",
			condition:  All(AnyCondition{true}, AnyCondition{true}),
			wantResult: true,
		},
		{
			name:       "does not match when one condition does not match",
			condition:  All(AnyCondition{true}, AnyCondition{false}),
			wantResult: false,
		},
		{
			name:       "matches when empty",
			condition:  All(),
			wantResult: true,
		},
		{
			name:       "matches when one condition matches",
			condition:  Any(AnyCondition{false}, AnyCondition{true}),
			wantResult: true,
		},
		{
			name:       "does not match when no conditions match",
			condition:  Any(AnyCondition{false}, AnyCondition{false}),
			wantResult: false,
		},
		{
			name:       "does not match when empty",
			condition:  Any(),
			wantResult: false,
		},
		{
			name:       "none matches when no conditions match",
			condition:  None(AnyCondition{false}, AnyCondition{false}),
			wantResult: true,
		},
		{
			name:       "none does not match when one condition matches",
			condition:  None(AnyCondition{false}, AnyCondition{true}),
			wantResult: false,
		},
		{
			name:       "none matches when empty",
			condition:  None(),
			wantResult: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := component.NewStandardContainer()
			ctx := context.Background()

			// when
			result := tc.condition.Matches(ctx, container)

			// then
			assert.Equal(t, tc.wantResult, result)
		})
	}
}
