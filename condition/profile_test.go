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
)

func TestProfile_Matches(t *testing.T) {
	testCases := []struct {
		name         string
		ctx          context.Context
		preCondition func(container component.Container)
		condition    *Profile
		wantResult   bool
	}{
		{
			name:       "does not match without environment",
			ctx:        context.Background(),
			condition:  OnProfile("dev"),
			wantResult: false,
		},
		{
			name: "does not match without specified profiles",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()

				err := container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition:  OnProfile(),
			wantResult: false,
		},
		{
			name: "matches default profile when no profile is active",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()

				err := container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition:  OnProfile("default"),
			wantResult: true,
		},
		{
			name: "does not match inactive profile",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()
				err := env.SetActiveProfiles("dev")
				assert.NoError(t, err)

				err = container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition:  OnProfile("secure"),
			wantResult: false,
		},
		{
			name: "matches active profile",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()
				err := env.SetActiveProfiles("dev")
				assert.NoError(t, err)

				err = container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition:  OnProfile("dev"),
			wantResult: true,
		},
		{
			name: "does not match default profile when profile is active",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()
				err := env.SetActiveProfiles("dev")
				assert.NoError(t, err)

				err = container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition:  OnProfile("default"),
			wantResult: false,
		},
		{
			name: "matches when one of multiple specified profiles is active",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()
				err := env.SetActiveProfiles("dev")
				assert.NoError(t, err)

				err = container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition:  OnProfile("secure", "dev"),
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

func TestMissingProfile_Matches(t *testing.T) {
	testCases := []struct {
		name         string
		ctx          context.Context
		preCondition func(container component.Container)
		condition    *MissingProfile
		wantResult   bool
	}{
		{
			name:       "does not match without environment",
			ctx:        context.Background(),
			condition:  OnMissingProfile("dev"),
			wantResult: false,
		},
		{
			name: "does not match without specified profiles",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()

				err := container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition:  OnMissingProfile(),
			wantResult: false,
		},
		{
			name: "does not match default profile when no profile is active",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()

				err := container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition:  OnMissingProfile("default"),
			wantResult: false,
		},
		{
			name: "matches inactive profile",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()
				err := env.SetActiveProfiles("dev")
				assert.NoError(t, err)

				err = container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition:  OnMissingProfile("secure"),
			wantResult: true,
		},
		{
			name: "does not match active profile",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()
				err := env.SetActiveProfiles("dev")
				assert.NoError(t, err)

				err = container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition:  OnMissingProfile("dev"),
			wantResult: false,
		},
		{
			name: "matches default profile when profile is active",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()
				err := env.SetActiveProfiles("dev")
				assert.NoError(t, err)

				err = container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition:  OnMissingProfile("default"),
			wantResult: true,
		},
		{
			name: "does not match when one of multiple specified profiles is active",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				env := procyon.NewEnvironment()
				err := env.SetActiveProfiles("dev")
				assert.NoError(t, err)

				err = container.RegisterSingleton("environment", env)
				assert.NoError(t, err)
			},
			condition:  OnMissingProfile("secure", "dev"),
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
