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

package procyon

import (
	"testing"

	"codnect.io/procyon/runtime/config"
	"github.com/stretchr/testify/assert"
)

func TestEnvironment_ActiveProfiles(t *testing.T) {
	// given
	env := NewEnvironment()

	propSource := config.NewMapPropertySource("anyMapSource", map[string]any{
		"procyon.profiles.active": "dev",
	})

	env.PropertySources().PushBack(propSource)

	// when
	profiles := env.ActiveProfiles()

	// then
	assert.Equal(t, 1, len(profiles))
	assert.ElementsMatch(t, []string{"dev"}, profiles)
}

func TestEnvironment_DefaultProfiles(t *testing.T) {
	// given
	env := NewEnvironment()

	propSource := config.NewMapPropertySource("anyMapSource", map[string]any{
		"procyon.profiles.default": "dev",
	})

	env.PropertySources().PushBack(propSource)

	// when
	profiles := env.DefaultProfiles()

	// then
	assert.Equal(t, 1, len(profiles))
	assert.ElementsMatch(t, []string{"dev"}, profiles)
}

func TestEnvironment_IsProfileActive(t *testing.T) {
	testCases := []struct {
		name           string
		activeProfiles string
		profileToCheck string
		expected       bool
	}{
		{
			name:           "Profile is active",
			activeProfiles: "dev, test",
			profileToCheck: "dev",
			expected:       true,
		},
		{
			name:           "Profile is not active",
			activeProfiles: "dev",
			profileToCheck: "test",
			expected:       false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			env := NewEnvironment()

			propSource := config.NewMapPropertySource("anyMapSource", map[string]any{
				"procyon.profiles.active": tc.activeProfiles,
			})

			env.PropertySources().PushBack(propSource)

			// when
			isActive := env.IsProfileActive(tc.profileToCheck)

			// then
			assert.Equal(t, tc.expected, isActive)
		})
	}
}

func TestEnvironment_SetActiveProfiles(t *testing.T) {
	testCases := []struct {
		name           string
		activeProfiles []string
		expected       []string
	}{
		{
			name:           "Set single active profile",
			activeProfiles: []string{"dev"},
			expected:       []string{"dev"},
		},
		{
			name:           "Set multiple active profiles",
			activeProfiles: []string{"dev", "test"},
			expected:       []string{"dev", "test"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			env := NewEnvironment()
			err := env.SetActiveProfiles(tc.activeProfiles...)

			// when
			profiles := env.ActiveProfiles()

			// then
			assert.NoError(t, err)
			assert.ElementsMatch(t, tc.expected, profiles)
		})
	}
}

func TestEnvironment_AddActiveProfiles(t *testing.T) {
	testCases := []struct {
		name            string
		initialProfiles []string
		profilesToAdd   []string
		expected        []string
	}{
		{
			name:            "Add single active profile",
			initialProfiles: []string{"dev"},
			profilesToAdd:   []string{"test"},
			expected:        []string{"dev", "test"},
		},
		{
			name:            "Add multiple active profiles",
			initialProfiles: []string{"dev"},
			profilesToAdd:   []string{"test", "secure"},
			expected:        []string{"dev", "test", "secure"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			env := NewEnvironment()
			err := env.SetActiveProfiles(tc.initialProfiles...)
			assert.NoError(t, err)

			// when
			err = env.AddActiveProfiles(tc.profilesToAdd...)
			profiles := env.ActiveProfiles()

			// then
			assert.NoError(t, err)
			assert.ElementsMatch(t, tc.expected, profiles)
		})
	}
}

func TestEnvironment_SetDefaultProfiles(t *testing.T) {
	testCases := []struct {
		name            string
		defaultProfiles []string
		expected        []string
	}{
		{
			name:            "Set single default profile",
			defaultProfiles: []string{"dev"},
			expected:        []string{"dev"},
		},
		{
			name:            "Set multiple default profiles",
			defaultProfiles: []string{"dev", "test"},
			expected:        []string{"dev", "test"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			env := NewEnvironment()
			err := env.SetDefaultProfiles(tc.defaultProfiles...)
			assert.NoError(t, err)

			// when
			profiles := env.DefaultProfiles()

			// then
			assert.ElementsMatch(t, tc.expected, profiles)
		})
	}
}

func TestEnvironment_PropertySources(t *testing.T) {
	// given
	env := NewEnvironment()

	// when
	propertySources := env.PropertySources()

	// then
	assert.NotNil(t, propertySources)
	assert.Equal(t, env.propertySources, propertySources)
}

func TestEnvironment_PropertyResolver(t *testing.T) {
	// given
	env := NewEnvironment()

	// when
	propertyResolver := env.PropertyResolver()

	// then
	assert.NotNil(t, propertyResolver)
	assert.Equal(t, env.propertyResolver, propertyResolver)
}
