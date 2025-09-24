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

package runtime

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvPropertySource_Name(t *testing.T) {
	// given
	envPropSource := NewEnvPropertySource()

	// when
	name := envPropSource.Name()

	// then
	assert.Equal(t, "systemEnvironment", name)
}

func TestEnvPropertySource_Origin(t *testing.T) {
	// given
	envPropSource := NewEnvPropertySource()

	// when
	origin := envPropSource.Origin()

	// then
	assert.NotNil(t, origin)
	assert.Equal(t, "os.Environ()", origin)
}

func TestEnvPropertySource_Value(t *testing.T) {
	testCases := []struct {
		name         string
		preCondition func()
		propName     string
		wantExists   bool
		wantValue    any
	}{
		{
			name:         "property does not exist",
			preCondition: func() {},
			propName:     "anotherKey",
			wantExists:   false,
		},
		{
			name: "lowercase property",
			preCondition: func() {
				os.Setenv("any_key", "anyValue")
			},
			propName:   "ANY_KEY",
			wantExists: true,
			wantValue:  "anyValue",
		},
		{
			name: "uppercase property",
			preCondition: func() {
				os.Setenv("ANY_KEY", "anyValue")
			},
			propName:   "ANY_KEY",
			wantExists: true,
			wantValue:  "anyValue",
		},
		{
			name: "no hyphen property",
			preCondition: func() {
				os.Setenv("ANY_KEY", "anyValue")
			},
			propName:   "ANY-KEY",
			wantExists: true,
			wantValue:  "anyValue",
		},
		{
			name: "no dot property",
			preCondition: func() {
				os.Setenv("ANY_KEY", "anyValue")
			},
			propName:   "ANY.KEY",
			wantExists: true,
			wantValue:  "anyValue",
		},
		{
			name: "no hyphen and no dot property",
			preCondition: func() {
				os.Setenv("ANY_KEY_PAIR", "anyValue")
			},
			propName:   "ANY.KEY-PAIR",
			wantExists: true,
			wantValue:  "anyValue",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			if tc.preCondition != nil {
				tc.preCondition()
			}

			envPropSource := NewEnvPropertySource()

			// when
			val, exists := envPropSource.Value(tc.propName)

			// then
			require.Equal(t, tc.wantExists, exists)

			if tc.wantValue != nil {
				assert.Equal(t, tc.wantValue, val)
			}

		})
	}
}

func TestEnvPropertySource_ValueOrDefault(t *testing.T) {

	testCases := []struct {
		name         string
		preCondition func()
		propName     string
		defaultValue string
		wantExists   bool
		wantValue    any
	}{
		{
			name:         "property does not exist",
			preCondition: func() {},
			propName:     "anyKey",
			defaultValue: "anyValue",
			wantValue:    "anyValue",
		},
		{
			name: "lowercase property",
			preCondition: func() {
				os.Setenv("any_key", "anyValue")
			},
			propName:  "ANY_KEY",
			wantValue: "anyValue",
		},
		{
			name: "uppercase property",
			preCondition: func() {
				os.Setenv("ANY_KEY", "anyValue")
			},
			propName:  "ANY_KEY",
			wantValue: "anyValue",
		},
		{
			name: "no hyphen property",
			preCondition: func() {
				os.Setenv("ANY_KEY", "anyValue")
			},
			propName:  "ANY-KEY",
			wantValue: "anyValue",
		},
		{
			name: "no dot property",
			preCondition: func() {
				os.Setenv("ANY_KEY", "anyValue")
			},
			propName:  "ANY.KEY",
			wantValue: "anyValue",
		},
		{
			name: "no hyphen and no dot property",
			preCondition: func() {
				os.Setenv("ANY_KEY_PAIR", "anyValue")
			},
			propName:  "ANY.KEY-PAIR",
			wantValue: "anyValue",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			if tc.preCondition != nil {
				tc.preCondition()
			}

			envPropSource := NewEnvPropertySource()

			// when
			val := envPropSource.ValueOrDefault(tc.propName, tc.defaultValue)

			// then
			assert.Equal(t, tc.wantValue, val)
		})
	}
}

func TestEnvPropertySource_PropertyNames(t *testing.T) {
	// given
	os.Setenv("ANY_KEY", "anyValue")

	envPropSource := NewEnvPropertySource()

	// when
	propNames := envPropSource.PropertyNames()

	// then
	assert.Contains(t, propNames, "ANY_KEY")
}
