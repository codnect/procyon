package config

import (
	"codnect.io/procyon/runtime/prop"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name           string
		propertySource prop.Source
		wantPanic      error
	}{
		{
			name:      "nil property source",
			wantPanic: errors.New("nil source"),
		},
		{
			name:           "valid property source",
			propertySource: prop.NewMapSource("anyMapSource", map[string]any{}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					New(tc.propertySource)
				})
				return
			}

			data := New(tc.propertySource)

			// then
			require.NotNil(t, data)
		})
	}
}

func TestData_PropertySource(t *testing.T) {
	// given
	mapSource := prop.NewMapSource("anyMapSource", map[string]any{})

	// when
	data := New(mapSource)

	// then
	assert.Equal(t, mapSource, data.PropertySource())
}
