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

package property

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewMultiSourceResolver(t *testing.T) {
	testCases := []struct {
		name       string
		sourceList *SourceList
		wantPanic  error
	}{
		{
			name:       "nil source list",
			sourceList: nil,
			wantPanic:  errors.New("nil sources"),
		},
		{
			name:       "with source list",
			sourceList: SourcesAsList(NewMapSource("anyMapName", map[string]any{"anyKey": "anyValue"})),
			wantPanic:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					NewMultiSourceResolver(tc.sourceList)
				})
				return
			}

			resolver := NewMultiSourceResolver(tc.sourceList)

			// then
			require.NotNil(t, resolver)
		})
	}
}

func TestMultiSourceResolver_ContainsProperty(t *testing.T) {
	testCases := []struct {
		name        string
		sources     []Source
		propertyKey string
		wantResult  bool
	}{
		{
			name: "property does not exist",
			sources: []Source{
				NewMapSource("anyMapName", map[string]any{"anyKey": "anyValue"}),
			},
			propertyKey: "anotherKey",
			wantResult:  false,
		},
		{
			name: "property exists",
			sources: []Source{
				NewMapSource("anyMapName", map[string]any{"anyKey": "anyValue"}),
			},
			propertyKey: "anyKey",
			wantResult:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			resolver := NewMultiSourceResolver(SourcesAsList(tc.sources...))

			// when
			result := resolver.ContainsProperty(tc.propertyKey)

			// then
			assert.Equal(t, tc.wantResult, result)
		})
	}
}

func TestMultiSourceResolver_Property(t *testing.T) {
	testCases := []struct {
		name        string
		sources     []Source
		propertyKey string
		wantExists  bool
		wantValue   any
	}{
		{
			name: "property does not exist",
			sources: []Source{
				NewMapSource("anyMapName", map[string]any{"anyKey": "anyValue"}),
			},
			propertyKey: "anotherKey",
			wantExists:  false,
		},
		{
			name: "property exists",
			sources: []Source{
				NewMapSource("anyMapName", map[string]any{"anyKey": "anyValue"}),
			},
			propertyKey: "anyKey",
			wantExists:  true,
			wantValue:   "anyValue",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			resolver := NewMultiSourceResolver(SourcesAsList(tc.sources...))

			// when
			val, exists := resolver.Property(tc.propertyKey)

			// then
			assert.Equal(t, tc.wantExists, exists)
			assert.Equal(t, tc.wantValue, val)
		})
	}
}

func TestMultiSourceResolver_PropertyOrDefault(t *testing.T) {
	testCases := []struct {
		name         string
		sources      []Source
		propertyKey  string
		defaultValue any
		wantExists   bool
		wantValue    any
	}{
		{
			name: "property does not exist",
			sources: []Source{
				NewMapSource("anyMapName", map[string]any{"anyKey": "anyValue"}),
			},
			propertyKey:  "anotherKey",
			defaultValue: "anotherValue",
			wantExists:   false,
			wantValue:    "anotherValue",
		},
		{
			name: "property exists",
			sources: []Source{
				NewMapSource("anyMapName", map[string]any{"anyKey": "anyValue"}),
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
			resolver := NewMultiSourceResolver(SourcesAsList(tc.sources...))

			// when
			val := resolver.PropertyOrDefault(tc.propertyKey, tc.defaultValue)

			// then
			assert.Equal(t, tc.wantValue, val)
		})
	}
}

func TestMultiSourceResolver_ResolvePlaceholders(t *testing.T) {
	testCases := []struct {
		name       string
		sources    []Source
		text       string
		wantResult string
	}{
		{
			name: "no placeholders resolved",
			sources: []Source{
				NewMapSource("anyMapName", map[string]any{}),
			},
			text:       "Server running on ${host}:${port} with profile ${environment}",
			wantResult: "Server running on ${host}:${port} with profile ${environment}",
		},
		{
			name: "empty placeholder",
			sources: []Source{
				NewMapSource("anyMapName", map[string]any{}),
			},
			text:       "Server running on ${}",
			wantResult: "Server running on ${}",
		},
		{
			name: "unterminated placeholder",
			sources: []Source{
				NewMapSource("anyMapName", map[string]any{}),
			},
			text:       "Server running on ${",
			wantResult: "Server running on ${",
		},
		{
			name: "num without curly braces",
			sources: []Source{
				NewMapSource("anyMapName", map[string]any{
					"123": "test",
				}),
			},
			text:       "Server running on 123",
			wantResult: "Server running on 123",
		},
		{
			name: "characters without curly braces",
			sources: []Source{
				NewMapSource("anyMapName", map[string]any{
					"AnyKey": "test",
				}),
			},
			text:       "Server running on $AnyKey",
			wantResult: "Server running on $AnyKey",
		},
		{
			name: "alpha num without curly braces",
			sources: []Source{
				NewMapSource("anyMapName", map[string]any{
					"Any123": "anyValue",
				}),
			},
			text:       "Server running on $Any123",
			wantResult: "Server running on $Any123",
		},
		{
			name: "partial placeholders resolved",
			sources: []Source{
				NewMapSource("anyMapName", map[string]any{
					"host":        "127.0.0.1",
					"environment": "dev",
				}),
			},
			text:       "Server running on ${host}:${port} with profile ${environment}",
			wantResult: "Server running on 127.0.0.1:${port} with profile dev",
		},
		{
			name: "all placeholders resolved",
			sources: []Source{
				NewMapSource("anyMapName", map[string]any{
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
			resolver := NewMultiSourceResolver(SourcesAsList(tc.sources...))

			// when
			result := resolver.ResolvePlaceholders(tc.text)

			// then
			assert.Equal(t, tc.wantResult, result)
		})
	}
}

func TestMultiSourceResolver_ResolveRequiredPlaceholders(t *testing.T) {
	testCases := []struct {
		name       string
		sources    []Source
		text       string
		wantResult string
		wantErr    error
	}{
		{
			name: "no placeholders resolved",
			sources: []Source{
				NewMapSource("anyMapName", map[string]any{}),
			},
			text:       "Server running on ${host}:${port} with profile ${environment}",
			wantResult: "Server running on ${host}:${port} with profile ${environment}",
			wantErr:    errors.New("cannot resolve placeholder '${host}'"),
		},
		{
			name: "empty placeholder",
			sources: []Source{
				NewMapSource("anyMapName", map[string]any{}),
			},
			text:       "Server running on ${}",
			wantResult: "Server running on ${}",
			wantErr:    errors.New("wrong placeholder format"),
		},
		{
			name: "unterminated placeholder",
			sources: []Source{
				NewMapSource("anyMapName", map[string]any{}),
			},
			text:       "Server running on ${",
			wantResult: "Server running on ${",
			wantErr:    errors.New("wrong placeholder format"),
		},
		{
			name: "num without curly braces",
			sources: []Source{
				NewMapSource("anyMapName", map[string]any{
					"123": "test",
				}),
			},
			text:       "Server running on 123",
			wantResult: "Server running on 123",
		},
		{
			name: "characters without curly braces",
			sources: []Source{
				NewMapSource("anyMapName", map[string]any{
					"AnyKey": "test",
				}),
			},
			text:       "Server running on $AnyKey",
			wantResult: "Server running on $AnyKey",
		},
		{
			name: "alpha num without curly braces",
			sources: []Source{
				NewMapSource("anyMapName", map[string]any{
					"Any123": "anyValue",
				}),
			},
			text:       "Server running on $Any123",
			wantResult: "Server running on $Any123",
		},
		{
			name: "partial placeholders resolved",
			sources: []Source{
				NewMapSource("anyMapName", map[string]any{
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
			sources: []Source{
				NewMapSource("anyMapName", map[string]any{
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
			resolver := NewMultiSourceResolver(SourcesAsList(tc.sources...))

			// when
			result, err := resolver.ResolveRequiredPlaceholders(tc.text)

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
