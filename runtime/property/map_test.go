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

func TestNewMapSource(t *testing.T) {
	testCases := []struct {
		name      string
		mapName   string
		keyValues map[string]any
		wantPanic error
	}{
		{
			name:      "empty map name",
			mapName:   "",
			wantPanic: errors.New("empty or blank name"),
		},
		{
			name:      "blank map name",
			mapName:   " ",
			wantPanic: errors.New("empty or blank name"),
		},
		{
			name:      "nil map",
			mapName:   "anyMapName",
			keyValues: nil,
			wantPanic: errors.New("nil map"),
		},
		{
			name:    "valid map source",
			mapName: "anyMapName",
			keyValues: map[string]any{
				"anyKey": "anyValue",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					NewMapSource(tc.mapName, tc.keyValues)
				})
				return
			}

			mapSource := NewMapSource(tc.mapName, tc.keyValues)

			// then
			require.NotNil(t, mapSource)
		})
	}
}

func TestMapSource_Name(t *testing.T) {
	// given
	mapSource := NewMapSource("anyMapName", map[string]any{})
	// when
	name := mapSource.Name()

	// then
	assert.Equal(t, "anyMapName", name)
}

func TestMapSource_Underlying(t *testing.T) {
	// given
	m := map[string]any{}
	mapSource := NewMapSource("anyMapName", m)

	// when
	underlying := mapSource.Underlying()

	// then
	assert.Equal(t, m, underlying)
}

func TestMapSource_ContainsProperty(t *testing.T) {
	testCases := []struct {
		name       string
		keyValues  map[string]any
		propName   string
		wantResult bool
	}{
		{
			name: "property does not exist",
			keyValues: map[string]any{
				"anyKey": "anyValue",
				"anyKeyWithSub": map[string]any{
					"anySubKey": "anySubValue",
				},
			},
			propName:   "anotherKey",
			wantResult: false,
		},
		{
			name: "property exists",
			keyValues: map[string]any{
				"anyKey": "anyValue",
				"anyKeyWithSub": map[string]any{
					"anySubKey": "anySubValue",
				},
			},
			propName:   "anyKey",
			wantResult: true,
		},
		{
			name: "nested-property exists",
			keyValues: map[string]any{
				"anyKey": "anyValue",
				"anyKeyWithSub": map[string]any{
					"anySubKey": "anySubValue",
				},
			},
			propName:   "anyKeyWithSub.anySubKey",
			wantResult: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			mapSource := NewMapSource("anyMapName", tc.keyValues)

			// when
			result := mapSource.ContainsProperty(tc.propName)

			// then
			assert.Equal(t, tc.wantResult, result)
		})
	}
}

func TestMapSource_Property(t *testing.T) {
	testCases := []struct {
		name       string
		keyValues  map[string]any
		propName   string
		wantExists bool
		wantValue  any
	}{
		{
			name: "property does not exist",
			keyValues: map[string]any{
				"anyKey": "anyValue",
				"anyKeyWithSub": map[string]any{
					"anySubKey": "anySubValue",
				},
			},
			propName:   "anotherKey",
			wantExists: false,
		},
		{
			name: "property exists",
			keyValues: map[string]any{
				"anyKey": "anyValue",
				"anyKeyWithSub": map[string]any{
					"anySubKey": "anySubValue",
				},
			},
			propName:   "anyKey",
			wantExists: true,
			wantValue:  "anyValue",
		},
		{
			name: "nested-property exists",
			keyValues: map[string]any{
				"anyKey": "anyValue",
				"anyKeyWithSub": map[string]any{
					"anySubKey": "anySubValue",
				},
			},
			propName:   "anyKeyWithSub.anySubKey",
			wantExists: true,
			wantValue:  "anySubValue",
		},
		{
			name: "property with value 'nil'",
			keyValues: map[string]any{
				"anyKey": nil,
			},
			propName:   "anyKey",
			wantExists: true,
			wantValue:  "",
		},
		{
			name: "array property",
			keyValues: map[string]any{
				"anyKey": "anyValue",
				"anyKeyWithSub": map[string]any{
					"anySubKey": "anySubValue",
				},
				"anyArrayKey": []any{
					"anyArrayValue",
					"anotherArrayValue",
				},
			},
			propName:   "anyArrayKey.0",
			wantExists: true,
			wantValue:  "anyArrayValue",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			mapSource := NewMapSource("anyMapName", tc.keyValues)

			// when
			val, exists := mapSource.Property(tc.propName)

			// then

			require.Equal(t, tc.wantExists, exists)

			if tc.wantValue != nil {
				assert.Equal(t, tc.wantValue, val)
			}

		})
	}
}

func TestMapSource_PropertyOrDefault(t *testing.T) {
	testCases := []struct {
		name         string
		keyValues    map[string]any
		propName     string
		defaultValue any
		wantValue    any
	}{
		{
			name: "property does not exist",
			keyValues: map[string]any{
				"anyKey": "anyValue",
				"anyKeyWithSub": map[string]any{
					"anySubKey": "anySubValue",
				},
			},
			propName:     "anotherKey",
			defaultValue: "anyDefaultValue",
			wantValue:    "anyDefaultValue",
		},
		{
			name: "property exists",
			keyValues: map[string]any{
				"anyKey": "anyValue",
				"anyKeyWithSub": map[string]any{
					"anySubKey": "anySubValue",
				},
			},
			propName:  "anyKey",
			wantValue: "anyValue",
		},
		{
			name: "nested-property exists",
			keyValues: map[string]any{
				"anyKey": "anyValue",
				"anyKeyWithSub": map[string]any{
					"anySubKey": "anySubValue",
				},
			},
			propName:  "anyKeyWithSub.anySubKey",
			wantValue: "anySubValue",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			mapSource := NewMapSource("anyMapName", tc.keyValues)

			// when
			val := mapSource.PropertyOrDefault(tc.propName, tc.defaultValue)

			// then
			assert.Equal(t, tc.wantValue, val)
		})
	}
}

func TestMapSource_PropertyNames(t *testing.T) {
	// given
	m := map[string]any{
		"anyKey":     "anyValue",
		"anotherKey": "anotherValue",
		"anyKeyWithSub": map[string]any{
			"anySubKey": "anySubValue",
		},
	}
	wantPropNames := []string{"anotherKey", "anyKey", "anyKeyWithSub.anySubKey"}
	mapSource := NewMapSource("anyMapName", m)

	// when
	propNames := mapSource.PropertyNames()

	// then
	assert.Len(t, propNames, len(wantPropNames))

	for _, wantProp := range wantPropNames {
		assert.Contains(t, propNames, wantProp)
	}
}
