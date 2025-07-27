package config

import (
	"codnect.io/procyon/runtime/prop"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type AnyLoader struct {
	data     *Data
	loadable bool
	err      error
}

func (m *AnyLoader) IsLoadable(resource Resource) bool {
	return m.loadable
}

func (m *AnyLoader) LoadConfig(ctx context.Context, resource Resource) (*Data, error) {
	return m.data, m.err
}

func TestNewFileLoader(t *testing.T) {
	// given

	// when
	fileLoader := NewFileLoader()

	// then
	assert.NotNil(t, fileLoader)
}

func TestFileLoader_IsLoadable(t *testing.T) {
	testCases := []struct {
		name     string
		resource Resource

		wantResult bool
	}{
		{
			name:       "nil resource",
			resource:   nil,
			wantResult: false,
		},
		{
			name:       "not loadable",
			resource:   &AnyResource{},
			wantResult: false,
		},
		{
			name:       "loadable",
			resource:   newFileResource("anyPath", &FakeFile{}, prop.NewYamlSourceLoader()),
			wantResult: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			loader := NewFileLoader()

			// when
			result := loader.IsLoadable(tc.resource)

			// then
			assert.Equal(t, tc.wantResult, result)
		})
	}
}

func TestFileLoader_LoadConfig(t *testing.T) {
	testCases := []struct {
		name     string
		ctx      context.Context
		resource Resource

		wantErr        error
		wantProperties map[string]any
	}{
		{
			name:     "nil context",
			ctx:      nil,
			resource: nil,
			wantErr:  errors.New("nil context"),
		},
		{
			name:     "nil resource",
			ctx:      context.Background(),
			resource: nil,
			wantErr:  errors.New("nil resource"),
		},
		{
			name:     "unknown resource",
			ctx:      context.Background(),
			resource: &AnyResource{},
			wantErr:  errors.New("unknown resource"),
		},
		{
			name: "invalid file resource",
			ctx:  context.Background(),
			resource: newFileResource("anyPath", &FakeFile{
				Content: []byte("key value"),
			}, prop.NewYamlSourceLoader()),
			wantErr: errors.New("yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `key value` into map[string]interface {}"),
		},
		{
			name: "valid file resource",
			ctx:  context.Background(),
			resource: newFileResource("anyPath", &FakeFile{
				Content: []byte("key: value"),
			}, prop.NewYamlSourceLoader()),
			wantProperties: map[string]any{
				"key": "value",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			loader := NewFileLoader()

			// when
			data, err := loader.LoadConfig(tc.ctx, tc.resource)

			// then

			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, data)

			propertySource := data.PropertySource()
			assert.NotNil(t, propertySource)

			for wantKey, wantValue := range tc.wantProperties {
				value, ok := propertySource.Property(wantKey)
				assert.True(t, ok)
				assert.Equal(t, wantValue, value)
			}
		})
	}
}
