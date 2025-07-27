package property

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

type ErrorReader struct{}

func (e *ErrorReader) Read(p []byte) (int, error) {
	return 0, errors.New("AnyError")
}

func TestYamlSourceLoader_FileExtensions(t *testing.T) {
	// given
	loader := NewYamlSourceLoader()

	// when
	extensions := loader.FileExtensions()

	// then
	assert.Len(t, extensions, 2)
	assert.Contains(t, extensions, "yml")
	assert.Contains(t, extensions, "yaml")
}

func TestYamlSourceLoader_LoadSource(t *testing.T) {
	testCases := []struct {
		name string

		sourceName string
		reader     io.Reader

		wantErr        error
		wantSourceName string
		wantProps      map[string]any
	}{
		{
			name:    "empty source name",
			wantErr: errors.New("empty source name"),
		},
		{
			name:       "nil reader",
			sourceName: "yaml source",
			wantErr:    errors.New("nil reader"),
		},
		{
			name:       "reader returns error",
			sourceName: "yaml source",
			reader:     &ErrorReader{},
			wantErr:    errors.New("AnyError"),
		},
		{
			name:       "invalid yaml",
			sourceName: "yaml source",
			reader:     bytes.NewReader([]byte("key value")),
			wantErr:    errors.New("yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `key value` into map[string]interface {}"),
		},
		{
			name:           "valid yaml",
			sourceName:     "yaml source",
			reader:         bytes.NewReader([]byte("key: value")),
			wantSourceName: "yaml source",
			wantProps:      map[string]any{"key": "value"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			loader := NewYamlSourceLoader()

			// when
			source, err := loader.Load(tc.sourceName, tc.reader)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, source)
			assert.Equal(t, tc.wantSourceName, source.Name())

			for propKey, propVal := range tc.wantProps {
				result, ok := source.Property(propKey)
				require.True(t, ok)
				require.Equal(t, propVal, result)
			}
		})
	}
}
