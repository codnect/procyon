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
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

type FakeReader struct {
	n   int
	err error
}

func (f *FakeReader) Read(p []byte) (n int, err error) {
	return f.n, f.err
}

func TestYamlSourceLoader_Extensions(t *testing.T) {
	// given
	loader := NewYamlSourceLoader()

	// when
	extensions := loader.Extensions()

	// then
	assert.Equal(t, []string{"yaml", "yml"}, extensions)
}

func TestYamlSourceLoader_Load(t *testing.T) {

	testCases := []struct {
		name           string
		sourceName     string
		reader         io.Reader
		wantErr        error
		wantProperties map[string]string
	}{
		{
			name:       "empty source name",
			sourceName: "",
			wantErr:    errors.New("empty source name"),
			reader:     bytes.NewReader([]byte("key: value")),
		},
		{
			name:       "empty source name",
			sourceName: "anySourceName",
			wantErr:    errors.New("nil reader"),
			reader:     nil,
		},
		{
			name:       "invalid yaml",
			sourceName: "anySourceName",
			reader:     bytes.NewReader([]byte("version 2.1\njobs:\n  image: 'nginx:latest'")),
			wantErr:    errors.New("yaml: line 2: mapping values are not allowed in this context"),
		},
		{
			name:       "reader error",
			sourceName: "anySourceName",
			reader: &FakeReader{
				err: errors.New("reader error"),
			},
			wantErr: errors.New("reader error"),
		},
		{
			name:       "valid yaml",
			sourceName: "anySourceName",
			reader:     bytes.NewReader([]byte("version: 2.1\njobs:\n  image: nginx:latest")),
			wantProperties: map[string]string{
				"version":    "2.1",
				"jobs.image": "nginx:latest",
			},
		},
		{
			name:       "valid yaml with array",
			sourceName: "anySourceName",
			reader:     bytes.NewReader([]byte("version: 2.1\njobs:\n  build:\n    docker:\n      image: cimg/base:2023.03\n    steps:\n      - checkout\n      - echo \"this is the build job\"")),
			wantProperties: map[string]string{
				"version":                 "2.1",
				"jobs.build.docker.image": "cimg/base:2023.03",
				"jobs.build.steps.0":      "checkout",
				"jobs.build.steps.1":      "echo \"this is the build job\"",
			},
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
				require.Nil(t, source)

				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			require.NotNil(t, source)

			for wantKey, wantValue := range tc.wantProperties {
				value, exists := source.Property(wantKey)
				assert.True(t, exists)
				assert.Equal(t, wantValue, value)
			}
		})
	}
}
