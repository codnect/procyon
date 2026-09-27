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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.codnect.io/procyon/component"
)

func TestResource_Matches(t *testing.T) {
	testCases := []struct {
		name         string
		ctx          context.Context
		preCondition func(ctx context.Context, container component.Container)
		condition    *Resource

		wantResult bool
	}{
		{
			name:       "does not match without resource resolver",
			ctx:        context.Background(),
			condition:  OnResource("resources/procyon.yaml"),
			wantResult: false,
		},
		{
			name: "does not match missing resource",
			ctx:  context.Background(),
			preCondition: func(ctx context.Context, container component.Container) {
				resource := new(AnyMockResource)
				resource.On("Exists").Return(false)

				resolver := new(AnyMockResourceResolver)
				resolver.
					On("Resolve", ctx, "resources/procyon.yaml").
					Return(resource, nil)

				err := container.RegisterSingleton("resourceResolver", resolver)
				require.NoError(t, err)
			},
			condition:  OnResource("resources/procyon.yaml"),
			wantResult: false,
		},
		{
			name: "does not match when resource resolution fails",
			ctx:  context.Background(),
			preCondition: func(ctx context.Context, container component.Container) {
				resolver := new(AnyMockResourceResolver)
				resolver.
					On("Resolve", ctx, "resources/procyon.yaml").
					Return(nil, errors.New("resolve resource error"))

				err := container.RegisterSingleton("resourceResolver", resolver)
				require.NoError(t, err)
			},
			condition:  OnResource("resources/procyon.yaml"),
			wantResult: false,
		},
		{
			name: "matches existing resource",
			ctx:  context.Background(),
			preCondition: func(ctx context.Context, container component.Container) {
				resource := new(AnyMockResource)
				resource.On("Exists").Return(true)

				resolver := new(AnyMockResourceResolver)
				resolver.
					On("Resolve", ctx, "resources/procyon.yaml").
					Return(resource, nil)

				err := container.RegisterSingleton("resourceResolver", resolver)
				require.NoError(t, err)
			},
			condition:  OnResource("resources/procyon.yaml"),
			wantResult: true,
		},
		{
			name: "matches all existing resources",
			ctx:  context.Background(),
			preCondition: func(ctx context.Context, container component.Container) {
				first := new(AnyMockResource)
				first.On("Exists").Return(true)

				second := new(AnyMockResource)
				second.On("Exists").Return(true)

				resolver := new(AnyMockResourceResolver)
				resolver.
					On("Resolve", ctx, "resources/procyon.yaml").
					Return(first, nil)
				resolver.
					On("Resolve", ctx, "resources/procyon-dev.yaml").
					Return(second, nil)

				err := container.RegisterSingleton("resourceResolver", resolver)
				require.NoError(t, err)
			},
			condition: OnResource(
				"resources/procyon.yaml",
				"resources/procyon-dev.yaml",
			),
			wantResult: true,
		},
		{
			name: "does not match when one resource is missing",
			ctx:  context.Background(),
			preCondition: func(ctx context.Context, container component.Container) {
				first := new(AnyMockResource)
				first.On("Exists").Return(true)

				second := new(AnyMockResource)
				second.On("Exists").Return(false)

				resolver := new(AnyMockResourceResolver)
				resolver.
					On("Resolve", ctx, "resources/procyon.yaml").
					Return(first, nil)
				resolver.
					On("Resolve", ctx, "resources/procyon-dev.yaml").
					Return(second, nil)

				err := container.RegisterSingleton("resourceResolver", resolver)
				require.NoError(t, err)
			},
			condition: OnResource(
				"resources/procyon.yaml",
				"resources/procyon-dev.yaml",
			),
			wantResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := component.NewStandardContainer()

			if tc.preCondition != nil {
				tc.preCondition(tc.ctx, container)
			}

			// when
			result := tc.condition.Matches(tc.ctx, container)

			// then
			assert.Equal(t, tc.wantResult, result)
		})
	}
}
