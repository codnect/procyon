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
	"codnect.io/procyon/runtime/property"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewData(t *testing.T) {
	testCases := []struct {
		name           string
		propertySource property.Source
		wantPanic      error
	}{
		{
			name:           "nil property source",
			propertySource: nil,
			wantPanic:      errors.New("nil property source"),
		},
		{
			name:           "valid property source",
			propertySource: property.NewMapSource("anyMapName", make(map[string]any)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					NewData(tc.propertySource)
				})
				return
			}

			data := NewData(tc.propertySource)

			// then
			require.NotNil(t, data)
		})
	}
}

func TestData_PropertySource(t *testing.T) {
	// given
	mapSource := property.NewMapSource("anyMapName", make(map[string]any))
	data := NewData(mapSource)

	// when
	source := data.PropertySource()

	// then
	assert.NotNil(t, source)
	assert.Equal(t, mapSource, source)
}
