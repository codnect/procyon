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

func TestPropertySources_Has(t *testing.T) {
	testCases := []struct {
		name       string
		sourceName string
		sources    []PropertySource
		wantResult bool
	}{
		{
			name:       "source does not exist",
			sourceName: "anotherMapSource",
			sources: []PropertySource{
				NewMapPropertySource("anyMapSource", make(map[string]any)),
			},
			wantResult: false,
		},
		{
			name:       "source exists",
			sourceName: "anyMapSource",
			sources: []PropertySource{
				NewMapPropertySource("anyMapSource", make(map[string]any)),
				NewMapPropertySource("anotherMapSource", make(map[string]any)),
			},
			wantResult: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			propertySources := NewPropertySources(tc.sources...)

			// when
			result := propertySources.Has(tc.sourceName)

			// then
			assert.Equal(t, tc.wantResult, result)
		})
	}
}

func TestPropertySources_Find(t *testing.T) {
	anyMapSource := NewMapPropertySource("anyMapSource", make(map[string]any))
	anotherMapSource := NewMapPropertySource("anotherMapSource", make(map[string]any))

	testCases := []struct {
		name       string
		sourceName string
		sources    []PropertySource
		wantExists bool
		wantSource PropertySource
	}{
		{
			name:       "source does not exist",
			sourceName: "anotherMapSource",
			sources: []PropertySource{
				anyMapSource,
			},
			wantExists: false,
		},
		{
			name:       "source exists",
			sourceName: "anyMapSource",
			sources: []PropertySource{
				anyMapSource,
				anotherMapSource,
			},
			wantExists: true,
			wantSource: anyMapSource,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			propertySources := NewPropertySources(tc.sources...)

			// when
			source, exists := propertySources.Get(tc.sourceName)

			// then
			assert.Equal(t, tc.wantExists, exists)
			assert.Equal(t, tc.wantSource, source)
		})
	}
}

func TestPropertySources_AddFirst(t *testing.T) {
	anyMapSource := NewMapPropertySource("anyMapSource", make(map[string]any))
	anotherMapSource := NewMapPropertySource("anotherMapSource", make(map[string]any))
	otherMapSource := NewMapPropertySource("anyMapSource", make(map[string]any))

	testCases := []struct {
		name        string
		sources     []PropertySource
		source      PropertySource
		wantPanic   error
		wantSources []PropertySource
	}{
		{
			name: "nil property source",
			sources: []PropertySource{
				anyMapSource,
			},
			source:    nil,
			wantPanic: errors.New("nil property source"),
		},
		{
			name: "any source",
			sources: []PropertySource{
				anyMapSource,
			},
			source:      anotherMapSource,
			wantSources: []PropertySource{anotherMapSource, anyMapSource},
		},
		{
			name:        "with empty sources",
			sources:     []PropertySource{},
			source:      anotherMapSource,
			wantSources: []PropertySource{anotherMapSource},
		},
		{
			name: "existing source",
			sources: []PropertySource{
				anyMapSource,
				anotherMapSource,
			},
			source:      otherMapSource,
			wantSources: []PropertySource{otherMapSource, anotherMapSource},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			propertySources := NewPropertySources(tc.sources...)

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					propertySources.PushFront(tc.source)
				})
				return
			}

			propertySources.PushFront(tc.source)

			// then
			assert.Equal(t, tc.wantSources, propertySources.items)
		})
	}
}

func TestPropertySources_AddLast(t *testing.T) {
	anyMapSource := NewMapPropertySource("anyMapSource", make(map[string]any))
	anotherMapSource := NewMapPropertySource("anotherMapSource", make(map[string]any))
	otherMapSource := NewMapPropertySource("anyMapSource", make(map[string]any))

	testCases := []struct {
		name        string
		sources     []PropertySource
		source      PropertySource
		wantPanic   error
		wantSources []PropertySource
	}{
		{
			name: "nil property source",
			sources: []PropertySource{
				anyMapSource,
			},
			source:    nil,
			wantPanic: errors.New("nil property source"),
		},
		{
			name: "any source",
			sources: []PropertySource{
				anyMapSource,
			},
			source:      anotherMapSource,
			wantSources: []PropertySource{anyMapSource, anotherMapSource},
		},
		{
			name: "existing source",
			sources: []PropertySource{
				anyMapSource,
				anotherMapSource,
			},
			source:      otherMapSource,
			wantSources: []PropertySource{anotherMapSource, otherMapSource},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			propertySources := NewPropertySources(tc.sources...)

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					propertySources.PushBack(tc.source)
				})
				return
			}

			propertySources.PushBack(tc.source)

			// then
			assert.Equal(t, tc.wantSources, propertySources.items)
		})
	}
}

func TestPropertySources_Insert(t *testing.T) {
	anyMapSource := NewMapPropertySource("anyMapSource", make(map[string]any))
	anotherMapSource := NewMapPropertySource("anotherMapSource", make(map[string]any))
	otherMapSource := NewMapPropertySource("otherMapSource", make(map[string]any))

	testCases := []struct {
		name        string
		sources     []PropertySource
		index       int
		source      PropertySource
		wantPanic   error
		wantSources []PropertySource
	}{
		{
			name: "negative index",
			sources: []PropertySource{
				anyMapSource,
				anotherMapSource,
			},
			index:     -5,
			wantPanic: errors.New("negative index"),
		},
		{
			name: "nil property source",
			sources: []PropertySource{
				anyMapSource,
				anotherMapSource,
			},
			index:     1,
			wantPanic: errors.New("nil property source"),
		},
		{
			name: "add last",
			sources: []PropertySource{
				anyMapSource,
			},
			source:      anotherMapSource,
			index:       1,
			wantSources: []PropertySource{anyMapSource, anotherMapSource},
		},
		{
			name: "add first",
			sources: []PropertySource{
				anyMapSource,
			},
			source:      anotherMapSource,
			index:       0,
			wantSources: []PropertySource{anotherMapSource, anyMapSource},
		},
		{
			name: "add at index",
			sources: []PropertySource{
				anyMapSource,
				anotherMapSource,
			},
			source:      otherMapSource,
			index:       1,
			wantSources: []PropertySource{anyMapSource, otherMapSource, anotherMapSource},
		},
		{
			name: "existing source",
			sources: []PropertySource{
				anotherMapSource,
				anyMapSource,
				otherMapSource,
			},
			source:      anyMapSource,
			index:       2,
			wantSources: []PropertySource{anotherMapSource, otherMapSource, anyMapSource},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			propertySources := NewPropertySources(tc.sources...)

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					propertySources.Insert(tc.index, tc.source)
				})
				return
			}

			propertySources.Insert(tc.index, tc.source)

			// then
			assert.Equal(t, tc.wantSources, propertySources.items)
		})
	}
}

func TestPropertySources_Remove(t *testing.T) {
	anyMapSource := NewMapPropertySource("anyMapSource", make(map[string]any))
	anotherMapSource := NewMapPropertySource("anotherMapSource", make(map[string]any))

	testCases := []struct {
		name       string
		sourceName string
		sources    []PropertySource
		wantSource PropertySource
		wantLen    int
	}{
		{
			name:       "source does not exist",
			sourceName: "anotherMapSource",
			sources: []PropertySource{
				anyMapSource,
			},
			wantSource: nil,
			wantLen:    1,
		},
		{
			name:       "source exists",
			sourceName: "anyMapSource",
			sources: []PropertySource{
				anyMapSource,
				anotherMapSource,
			},
			wantSource: anyMapSource,
			wantLen:    1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			propertySources := NewPropertySources(tc.sources...)

			// when
			source := propertySources.Remove(tc.sourceName)

			// then
			assert.Equal(t, tc.wantSource, source)
			assert.Equal(t, tc.wantLen, len(propertySources.items))
		})
	}
}

func TestPropertySources_Replace(t *testing.T) {
	anyMapSource := NewMapPropertySource("anyMapSource", make(map[string]any))
	anotherMapSource := NewMapPropertySource("anotherMapSource", make(map[string]any))
	otherMapSource := NewMapPropertySource("anyMapSource", make(map[string]any))

	testCases := []struct {
		name        string
		sources     []PropertySource
		sourceName  string
		source      PropertySource
		wantPanic   error
		wantSources []PropertySource
	}{
		{
			name:       "nil property source",
			sourceName: "anySourceName",
			sources: []PropertySource{
				anyMapSource,
				anotherMapSource,
			},
			source:    nil,
			wantPanic: errors.New("nil property source"),
		},
		{
			name:       "replace with existing",
			sourceName: "anyMapSource",
			sources: []PropertySource{
				anyMapSource,
				anotherMapSource,
			},
			source:      otherMapSource,
			wantSources: []PropertySource{otherMapSource, anyMapSource},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			propertySources := NewPropertySources(tc.wantSources...)

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					propertySources.Replace(tc.sourceName, tc.source)
				})
				return
			}

			propertySources.Replace(tc.sourceName, tc.source)

			// then
			assert.Equal(t, tc.wantSources, propertySources.items)
		})
	}
}

func TestPropertySources_Len(t *testing.T) {
	// given
	anyMapSource := NewMapPropertySource("anyMapSource", make(map[string]any))
	anotherMapSource := NewMapPropertySource("anotherMapSource", make(map[string]any))
	propertySources := NewPropertySources(anyMapSource, anotherMapSource)

	// when
	count := propertySources.Len()

	// then
	assert.Equal(t, 2, count)
}

func TestPropertySources_IndexOf(t *testing.T) {
	anyMapSource := NewMapPropertySource("anyMapSource", make(map[string]any))
	anotherMapSource := NewMapPropertySource("anotherMapSource", make(map[string]any))

	testCases := []struct {
		name      string
		sources   []PropertySource
		source    PropertySource
		wantIndex int
	}{
		{
			name: "nil source",
			sources: []PropertySource{
				anyMapSource,
				anotherMapSource,
			},
			source:    nil,
			wantIndex: -1,
		},
		{
			name: "source does not exist",
			sources: []PropertySource{
				anyMapSource,
			},
			source:    anotherMapSource,
			wantIndex: -1,
		},
		{
			name: "source exists",
			sources: []PropertySource{
				anyMapSource,
				anotherMapSource,
			},
			source:    anotherMapSource,
			wantIndex: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			propertySources := NewPropertySources(tc.sources...)

			// when
			index := propertySources.IndexOf(tc.source)

			// then
			assert.Equal(t, tc.wantIndex, index)
		})
	}
}

func TestPropertySources_Slice(t *testing.T) {
	// given
	anyMapSource := NewMapPropertySource("anyMapSource", make(map[string]any))
	anotherMapSource := NewMapPropertySource("anotherMapSource", make(map[string]any))
	propertySources := NewPropertySources(anyMapSource, anotherMapSource)

	// when
	slice := propertySources.Slice()

	// then
	assert.Equal(t, []PropertySource{anyMapSource, anotherMapSource}, slice)
}

func TestNewMapPropertySource(t *testing.T) {
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
					NewMapPropertySource(tc.mapName, tc.keyValues)
				})
				return
			}

			mapPropSource := NewMapPropertySource(tc.mapName, tc.keyValues)

			// then
			require.NotNil(t, mapPropSource)
		})
	}
}

func TestMapPropertySource_Name(t *testing.T) {
	// given
	mapPropSource := NewMapPropertySource("anyMapName", map[string]any{})

	// when
	name := mapPropSource.Name()

	// then
	assert.Equal(t, "anyMapName", name)
}

func TestMapPropertySource_Origin(t *testing.T) {
	// given
	mapSource := NewMapPropertySource("anyMapName", map[string]any{})

	// when
	origin := mapSource.Origin()

	// then
	assert.Equal(t, "anyMapName", origin)
}

func TestMapPropertySource_Value(t *testing.T) {
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
			mapPropSource := NewMapPropertySource("anyMapName", tc.keyValues)

			// when
			val, exists := mapPropSource.Value(tc.propName)

			// then

			require.Equal(t, tc.wantExists, exists)

			if tc.wantValue != nil {
				assert.Equal(t, tc.wantValue, val)
			}

		})
	}
}

func TestMapPropertySource_ValueOrDefault(t *testing.T) {
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
			mapPropSource := NewMapPropertySource("anyMapName", tc.keyValues)

			// when
			val := mapPropSource.ValueOrDefault(tc.propName, tc.defaultValue)

			// then
			assert.Equal(t, tc.wantValue, val)
		})
	}
}

func TestMapPropertySource_PropertyNames(t *testing.T) {
	// given
	m := map[string]any{
		"anyKey":     "anyValue",
		"anotherKey": "anotherValue",
		"anyKeyWithSub": map[string]any{
			"anySubKey": "anySubValue",
		},
	}
	wantPropNames := []string{"anotherKey", "anyKey", "anyKeyWithSub.anySubKey"}
	mapPropSource := NewMapPropertySource("anyMapName", m)

	// when
	propNames := mapPropSource.PropertyNames()

	// then
	assert.Len(t, propNames, len(wantPropNames))

	for _, wantProp := range wantPropNames {
		assert.Contains(t, propNames, wantProp)
	}
}
