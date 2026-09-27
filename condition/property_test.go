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
	"go.codnect.io/procyon"
	"go.codnect.io/procyon/component"
	"go.codnect.io/procyon/runtime/config"
)

func TestProperty_Matches(t *testing.T) {
	testCases := []struct {
		name         string
		ctx          context.Context
		preCondition func(container component.Container)
		condition    *Property
		wantResult   bool
	}{
		{
			name:       "does not match without environment",
			ctx:        context.Background(),
			condition:  OnProperty("server.port"),
			wantResult: false,
		},
		{
			name: "does not match missing property",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()

				err := container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition:  OnProperty("server.port"),
			wantResult: false,
		},
		{
			name: "matches existing property",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()
				propSource := config.NewMapPropertySource("anyMapSource", map[string]any{
					"server.port": 8080,
				})

				env.PropertySources().PushFront(propSource)

				err := container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition:  OnProperty("server.port"),
			wantResult: true,
		},
		{
			name: "matches false property",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()
				propSource := config.NewMapPropertySource("anyMapSource", map[string]any{
					"metrics.enabled": false,
				})

				env.PropertySources().PushFront(propSource)

				err := container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition:  OnProperty("metrics.enabled"),
			wantResult: true,
		},
		{
			name: "matches equal property value",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()
				propSource := config.NewMapPropertySource("anyMapSource", map[string]any{
					"metrics.enabled": "true",
				})

				env.PropertySources().PushFront(propSource)

				err := container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition:  OnProperty("metrics.enabled").Equals("true"),
			wantResult: true,
		},
		{
			name: "does not match different property value",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()
				propSource := config.NewMapPropertySource("anyMapSource", map[string]any{
					"metrics.enabled": "false",
				})

				env.PropertySources().PushFront(propSource)

				err := container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition:  OnProperty("metrics.enabled").Equals("true"),
			wantResult: false,
		},
		{
			name: "matches property with prefix",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()
				propSource := config.NewMapPropertySource("anyMapSource", map[string]any{
					"metrics.enabled": true,
				})

				env.PropertySources().PushFront(propSource)

				err := container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition:  OnProperty("enabled").Prefix("metrics"),
			wantResult: true,
		},
		{
			name: "matches multiple properties",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()
				propSource := config.NewMapPropertySource("anyMapSource", map[string]any{
					"metrics.enabled": true,
					"report.enabled":  false,
				})

				env.PropertySources().PushFront(propSource)

				err := container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition: OnProperty(
				"metrics.enabled",
				"report.enabled",
			),
			wantResult: true,
		},
		{
			name: "does not match multiple properties with missing property",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()
				propSource := config.NewMapPropertySource("anyMapSource", map[string]any{
					"metrics.enabled": true,
				})

				env.PropertySources().PushFront(propSource)

				err := container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition: OnProperty(
				"metrics.enabled",
				"report.enabled",
			),
			wantResult: false,
		},
		{
			name: "matches non-string property value",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()
				propSource := config.NewMapPropertySource("anyMapSource", map[string]any{
					"server.port": 8080,
				})

				env.PropertySources().PushFront(propSource)

				err := container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition:  OnProperty("server.port").Equals("8080"),
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

func TestMissingProperty_Matches(t *testing.T) {
	testCases := []struct {
		name         string
		ctx          context.Context
		preCondition func(container component.Container)
		condition    *MissingProperty
		wantResult   bool
	}{
		{
			name:       "does not match without environment",
			ctx:        context.Background(),
			condition:  OnMissingProperty("server.port"),
			wantResult: false,
		},
		{
			name: "matches missing property",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()

				err := container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition:  OnMissingProperty("server.port"),
			wantResult: true,
		},
		{
			name: "does not match existing property",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()
				propSource := config.NewMapPropertySource("anyMapSource", map[string]any{
					"server.port": 8080,
				})

				env.PropertySources().PushFront(propSource)

				err := container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition:  OnMissingProperty("server.port"),
			wantResult: false,
		},
		{
			name: "does not match existing false property",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()
				propSource := config.NewMapPropertySource("anyMapSource", map[string]any{
					"metrics.enabled": false,
				})

				env.PropertySources().PushFront(propSource)

				err := container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition:  OnMissingProperty("metrics.enabled"),
			wantResult: false,
		},
		{
			name: "matches missing property with prefix",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()

				err := container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition:  OnMissingProperty("enabled").Prefix("metrics"),
			wantResult: true,
		},
		{
			name: "does not match existing property with prefix",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()
				propSource := config.NewMapPropertySource("anyMapSource", map[string]any{
					"metrics.enabled": true,
				})

				env.PropertySources().PushFront(propSource)

				err := container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition:  OnMissingProperty("enabled").Prefix("metrics"),
			wantResult: false,
		},
		{
			name: "matches multiple missing properties",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()

				err := container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition: OnMissingProperty(
				"metrics.enabled",
				"report.enabled",
			),
			wantResult: true,
		},
		{
			name: "does not match multiple properties with existing property",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()
				propSource := config.NewMapPropertySource("anyMapSource", map[string]any{
					"metrics.enabled": true,
				})

				env.PropertySources().PushFront(propSource)

				err := container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition: OnMissingProperty(
				"metrics.enabled",
				"report.enabled",
			),
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
