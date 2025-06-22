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
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type AnyCondition struct {
	matches bool
}

func (a AnyCondition) Matches(ctx ConditionContext) bool {
	return a.matches
}

func TestNewConditionContext(t *testing.T) {
	testCases := []struct {
		name      string
		ctx       context.Context
		container Container
		wantPanic error
	}{
		{
			name:      "nil context",
			ctx:       nil,
			container: NewDefaultContainer(),
			wantPanic: errors.New("nil context"),
		},
		{
			name:      "nil container",
			ctx:       context.Background(),
			container: nil,
			wantPanic: errors.New("nil container"),
		},
		{
			name:      "valid context and container",
			ctx:       context.Background(),
			container: NewDefaultContainer(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					newConditionContext(tc.ctx, tc.container)
				})
				return
			}

			condCtx := newConditionContext(tc.ctx, tc.container)

			// then
			require.NotNil(t, condCtx)
		})
	}
}

func TestConditionContext_Container(t *testing.T) {
	// given
	ctx := context.Background()
	container := NewDefaultContainer()
	conditionCtx := newConditionContext(ctx, container)

	// when
	result := conditionCtx.Container()

	// then
	assert.Equal(t, container, result)
}

func TestConditionContext_Err(t *testing.T) {
	// given
	ctx, cancelFn := context.WithCancel(context.Background())

	// cancel context
	cancelFn()

	container := NewDefaultContainer()
	conditionCtx := newConditionContext(ctx, container)

	// when
	err := conditionCtx.Err()

	// then
	assert.Equal(t, context.Canceled, err)
}

func TestConditionContext_Value(t *testing.T) {
	// given
	ctx := context.WithValue(context.Background(), "anyKey", "anyValue")
	container := NewDefaultContainer()
	conditionCtx := newConditionContext(ctx, container)

	// when
	result := conditionCtx.Value("anyKey")

	// then
	assert.Equal(t, "anyValue", result)
}

func TestConditionContext_Done(t *testing.T) {
	// given
	ctx, cancelFn := context.WithCancel(context.Background())

	// cancel context
	cancelFn()

	container := NewDefaultContainer()
	conditionCtx := newConditionContext(ctx, container)

	// when
	<-conditionCtx.Done()

	// then
	assert.Equal(t, context.Canceled, ctx.Err())
}

func TestConditionContext_Deadline(t *testing.T) {
	// given
	deadline := time.Now()
	ctx, _ := context.WithDeadline(context.Background(), deadline)

	container := NewDefaultContainer()
	conditionCtx := newConditionContext(ctx, container)

	// when
	result, ok := conditionCtx.Deadline()

	// then
	assert.True(t, ok)
	assert.Equal(t, deadline, result)
}

func TestNewConditionEvaluator(t *testing.T) {
	testCases := []struct {
		name      string
		ctx       context.Context
		container Container
		wantPanic error
	}{
		{
			name:      "nil container",
			ctx:       context.Background(),
			container: nil,
			wantPanic: errors.New("nil container"),
		},
		{
			name:      "valid container",
			ctx:       context.Background(),
			container: NewDefaultContainer(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					newConditionEvaluator(tc.container)
				})
				return
			}

			evaluator := newConditionEvaluator(tc.container)
			require.NotNil(t, evaluator)
		})
	}
}

func TestConditionEvaluator_Evaluate(t *testing.T) {
	testCases := []struct {
		name       string
		ctx        context.Context
		conditions []Condition
		container  Container
		wantResult bool
	}{
		{
			name:       "nil conditions",
			ctx:        context.Background(),
			container:  NewDefaultContainer(),
			conditions: nil,
			wantResult: true,
		},
		{
			name:       "no condition",
			ctx:        context.Background(),
			container:  NewDefaultContainer(),
			conditions: []Condition{},
			wantResult: true,
		},
		{
			name:      "all conditions match",
			ctx:       context.Background(),
			container: NewDefaultContainer(),
			conditions: []Condition{
				AnyCondition{
					matches: true,
				},
				AnyCondition{
					matches: true,
				},
			},
			wantResult: true,
		},
		{
			name:      "all conditions do not match",
			ctx:       context.Background(),
			container: NewDefaultContainer(),
			conditions: []Condition{
				AnyCondition{
					matches: false,
				},
				AnyCondition{
					matches: false,
				},
			},
			wantResult: false,
		},
		{
			name:      "at least one condition does not match",
			ctx:       context.Background(),
			container: NewDefaultContainer(),
			conditions: []Condition{
				AnyCondition{
					matches: true,
				},
				AnyCondition{
					matches: false,
				},
			},
			wantResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			evaluator := newConditionEvaluator(tc.container)

			// when
			result := evaluator.evaluate(tc.ctx, tc.conditions)

			// then
			assert.Equal(t, tc.wantResult, result)
		})
	}
}
