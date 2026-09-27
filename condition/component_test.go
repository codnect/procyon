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

type AnyInterface interface {
	AnyMethod()
}

type AnyComponent struct {
}

func (a AnyComponent) AnyMethod() {}

func TestComponent_Matches(t *testing.T) {
	testCases := []struct {
		name         string
		ctx          context.Context
		preCondition func(container component.Container)
		condition    *Component
		wantResult   bool
	}{
		{
			name:       "does not match empty name",
			ctx:        context.Background(),
			condition:  OnComponentName(""),
			wantResult: false,
		},
		{
			name:       "does not match missing component",
			ctx:        context.Background(),
			condition:  OnComponent[AnyComponent](),
			wantResult: false,
		},
		{
			name: "matches concrete component",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				err := container.RegisterSingleton("anyComponent", AnyComponent{})
				assert.NoError(t, err)
			},
			condition:  OnComponent[AnyComponent](),
			wantResult: true,
		},
		{
			name:       "does not match missing named component",
			ctx:        context.Background(),
			condition:  OnComponentName("anyComponent"),
			wantResult: false,
		},
		{
			name: "matches named component",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				err := container.RegisterSingleton("anyComponent", AnyComponent{})
				assert.NoError(t, err)
			},
			condition:  OnComponentName("anyComponent"),
			wantResult: true,
		},
		{
			name: "matches interface component",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				err := container.RegisterSingleton("anyComponent", AnyComponent{})
				assert.NoError(t, err)
			},
			condition:  OnComponent[AnyInterface](),
			wantResult: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := component.NewStandardContainer()

			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			// when
			result := tc.condition.Matches(tc.ctx, container)

			// then
			assert.Equal(t, tc.wantResult, result)
		})
	}
}

func TestMissingComponent_Matches(t *testing.T) {
	testCases := []struct {
		name         string
		ctx          context.Context
		preCondition func(container component.Container)
		condition    *MissingComponent
		wantResult   bool
	}{
		{
			name:       "does not match empty name",
			ctx:        context.Background(),
			condition:  OnMissingComponentName(""),
			wantResult: false,
		},
		{
			name:       "matches missing component",
			ctx:        context.Background(),
			condition:  OnMissingComponent[AnyComponent](),
			wantResult: true,
		},
		{
			name: "does not match existing concrete component",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				err := container.RegisterSingleton("anyComponent", AnyComponent{})
				assert.NoError(t, err)
			},
			condition:  OnMissingComponent[AnyComponent](),
			wantResult: false,
		},
		{
			name:       "matches missing named component",
			ctx:        context.Background(),
			condition:  OnMissingComponentName("anyComponent"),
			wantResult: true,
		},
		{
			name: "does not match existing named component",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				err := container.RegisterSingleton("anyComponent", AnyComponent{})
				assert.NoError(t, err)
			},
			condition:  OnMissingComponentName("anyComponent"),
			wantResult: false,
		},
		{
			name: "does not match existing interface component",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				err := container.RegisterSingleton("anyComponent", AnyComponent{})
				assert.NoError(t, err)
			},
			condition:  OnMissingComponent[AnyInterface](),
			wantResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := component.NewStandardContainer()

			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			// when
			result := tc.condition.Matches(tc.ctx, container)

			// then
			assert.Equal(t, tc.wantResult, result)
		})
	}
}
