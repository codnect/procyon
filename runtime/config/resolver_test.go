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

package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultPropertyResolver(t *testing.T) {
	testCases := []struct {
		name            string
		propertySources *PropertySources
		wantPanic       error
	}{
		{
			name:            "nil property sources",
			propertySources: nil,
			wantPanic:       errors.New("nil property sources"),
		},
		{
			name:            "with source list",
			propertySources: NewPropertySources(NewMapPropertySource("anyMapName", map[string]any{"anyKey": "anyValue"})),
			wantPanic:       nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					NewDefaultPropertyResolver(tc.propertySources)
				})
				return
			}

			resolver := NewDefaultPropertyResolver(tc.propertySources)

			// then
			require.NotNil(t, resolver)
		})
	}
}

func TestDefaultPropertyResolver_Lookup(t *testing.T) {
	testCases := []struct {
		name            string
		propertySources []PropertySource
		propertyKey     string
		wantExists      bool
		wantValue       any
	}{
		{
			name: "property does not exist",
			propertySources: []PropertySource{
				NewMapPropertySource("anyMapName", map[string]any{"anyKey": "anyValue"}),
			},
			propertyKey: "anotherKey",
			wantExists:  false,
		},
		{
			name: "property exists",
			propertySources: []PropertySource{
				NewMapPropertySource("anyMapName", map[string]any{"anyKey": "anyValue"}),
			},
			propertyKey: "anyKey",
			wantExists:  true,
			wantValue:   "anyValue",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			propertySources := NewPropertySources(tc.propertySources...)
			resolver := NewDefaultPropertyResolver(propertySources)

			// when
			val, exists := resolver.Lookup(tc.propertyKey)

			// then
			assert.Equal(t, tc.wantExists, exists)
			assert.Equal(t, tc.wantValue, val)
		})
	}
}

func TestDefaultPropertyResolver_LookupOrDefault(t *testing.T) {
	testCases := []struct {
		name            string
		propertySources []PropertySource
		propertyKey     string
		defaultValue    any
		wantExists      bool
		wantValue       any
	}{
		{
			name: "property does not exist",
			propertySources: []PropertySource{
				NewMapPropertySource("anyMapName", map[string]any{"anyKey": "anyValue"}),
			},
			propertyKey:  "anotherKey",
			defaultValue: "anotherValue",
			wantExists:   false,
			wantValue:    "anotherValue",
		},
		{
			name: "property exists",
			propertySources: []PropertySource{
				NewMapPropertySource("anyMapName", map[string]any{"anyKey": "anyValue"}),
			},
			propertyKey:  "anyKey",
			defaultValue: "anotherValue",
			wantExists:   true,
			wantValue:    "anyValue",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			propertySources := NewPropertySources(tc.propertySources...)
			resolver := NewDefaultPropertyResolver(propertySources)

			// when
			val := resolver.LookupOrDefault(tc.propertyKey, tc.defaultValue)

			// then
			assert.Equal(t, tc.wantValue, val)
		})
	}
}

func TestDefaultPropertyResolver_Expand(t *testing.T) {
	testCases := []struct {
		name            string
		propertySources []PropertySource
		text            string
		wantResult      string
	}{
		{
			name: "no placeholders resolved",
			propertySources: []PropertySource{
				NewMapPropertySource("anyMapName", map[string]any{}),
			},
			text:       "Server running on ${host}:${port} with profile ${environment}",
			wantResult: "Server running on ${host}:${port} with profile ${environment}",
		},
		{
			name: "empty placeholder",
			propertySources: []PropertySource{
				NewMapPropertySource("anyMapName", map[string]any{}),
			},
			text:       "Server running on ${}",
			wantResult: "Server running on ${}",
		},
		{
			name: "unterminated placeholder",
			propertySources: []PropertySource{
				NewMapPropertySource("anyMapName", map[string]any{}),
			},
			text:       "Server running on ${",
			wantResult: "Server running on ${",
		},
		{
			name: "num without curly braces",
			propertySources: []PropertySource{
				NewMapPropertySource("anyMapName", map[string]any{
					"123": "test",
				}),
			},
			text:       "Server running on 123",
			wantResult: "Server running on 123",
		},
		{
			name: "characters without curly braces",
			propertySources: []PropertySource{
				NewMapPropertySource("anyMapName", map[string]any{
					"AnyKey": "test",
				}),
			},
			text:       "Server running on $AnyKey",
			wantResult: "Server running on $AnyKey",
		},
		{
			name: "alpha num without curly braces",
			propertySources: []PropertySource{
				NewMapPropertySource("anyMapName", map[string]any{
					"Any123": "anyValue",
				}),
			},
			text:       "Server running on $Any123",
			wantResult: "Server running on $Any123",
		},
		{
			name: "partial placeholders resolved",
			propertySources: []PropertySource{
				NewMapPropertySource("anyMapName", map[string]any{
					"host":        "127.0.0.1",
					"environment": "dev",
				}),
			},
			text:       "Server running on ${host}:${port} with profile ${environment}",
			wantResult: "Server running on 127.0.0.1:${port} with profile dev",
		},
		{
			name: "all placeholders resolved",
			propertySources: []PropertySource{
				NewMapPropertySource("anyMapName", map[string]any{
					"host":        "127.0.0.1",
					"port":        8090,
					"environment": "dev",
				}),
			},
			text:       "Server running on ${host}:${port} with profile ${environment}",
			wantResult: "Server running on 127.0.0.1:8090 with profile dev",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			propertySources := NewPropertySources(tc.propertySources...)
			resolver := NewDefaultPropertyResolver(propertySources)

			// when
			result := resolver.Expand(tc.text)

			// then
			assert.Equal(t, tc.wantResult, result)
		})
	}
}

func TestDefaultPropertyResolver_ExpandStrict(t *testing.T) {
	testCases := []struct {
		name            string
		propertySources []PropertySource
		text            string
		wantResult      string
		wantErr         error
	}{
		{
			name: "no placeholders resolved",
			propertySources: []PropertySource{
				NewMapPropertySource("anyMapName", map[string]any{}),
			},
			text:       "Server running on ${host}:${port} with profile ${environment}",
			wantResult: "Server running on ${host}:${port} with profile ${environment}",
			wantErr:    errors.New("cannot resolve placeholder '${host}'"),
		},
		{
			name: "empty placeholder",
			propertySources: []PropertySource{
				NewMapPropertySource("anyMapName", map[string]any{}),
			},
			text:       "Server running on ${}",
			wantResult: "Server running on ${}",
			wantErr:    errors.New("wrong placeholder format"),
		},
		{
			name: "unterminated placeholder",
			propertySources: []PropertySource{
				NewMapPropertySource("anyMapName", map[string]any{}),
			},
			text:       "Server running on ${",
			wantResult: "Server running on ${",
			wantErr:    errors.New("wrong placeholder format"),
		},
		{
			name: "num without curly braces",
			propertySources: []PropertySource{
				NewMapPropertySource("anyMapName", map[string]any{
					"123": "test",
				}),
			},
			text:       "Server running on 123",
			wantResult: "Server running on 123",
		},
		{
			name: "characters without curly braces",
			propertySources: []PropertySource{
				NewMapPropertySource("anyMapName", map[string]any{
					"AnyKey": "test",
				}),
			},
			text:       "Server running on $AnyKey",
			wantResult: "Server running on $AnyKey",
		},
		{
			name: "alpha num without curly braces",
			propertySources: []PropertySource{
				NewMapPropertySource("anyMapName", map[string]any{
					"Any123": "anyValue",
				}),
			},
			text:       "Server running on $Any123",
			wantResult: "Server running on $Any123",
		},
		{
			name: "partial placeholders resolved",
			propertySources: []PropertySource{
				NewMapPropertySource("anyMapName", map[string]any{
					"host":        "127.0.0.1",
					"environment": "dev",
				}),
			},
			text:       "Server running on ${host}:${port} with profile ${environment}",
			wantResult: "Server running on 127.0.0.1:${port} with profile dev",
			wantErr:    errors.New("cannot resolve placeholder '${port}'"),
		},
		{
			name: "all placeholders resolved",
			propertySources: []PropertySource{
				NewMapPropertySource("anyMapName", map[string]any{
					"host":        "127.0.0.1",
					"port":        8090,
					"environment": "dev",
				}),
			},
			text:       "Server running on ${host}:${port} with profile ${environment}",
			wantResult: "Server running on 127.0.0.1:8090 with profile dev",
			wantErr:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			propertySources := NewPropertySources(tc.propertySources...)
			resolver := NewDefaultPropertyResolver(propertySources)

			// when
			result, err := resolver.ExpandStrict(tc.text)

			// then
			if tc.wantErr != nil {
				require.Empty(t, result)

				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.wantResult, result)
		})
	}
}
