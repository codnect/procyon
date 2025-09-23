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
	"context"
	"errors"
	stdio "io"
	"io/fs"
	"testing"

	"codnect.io/procyon/io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type FakeFile struct {
	contents string
	offset   int
	readErr  error
	fileInfo fs.FileInfo
}

func (f *FakeFile) Reset() *FakeFile {
	f.offset = 0
	return f
}

func (f *FakeFile) Stat() (fs.FileInfo, error) {
	return f.fileInfo, nil
}

func (f *FakeFile) Read(p []byte) (int, error) {
	if f.readErr != nil {
		return 0, f.readErr
	}
	if f.offset >= len(f.contents) {
		return 0, stdio.EOF
	}
	n := copy(p, f.contents[f.offset:])
	f.offset += n
	return n, nil
}

func (f *FakeFile) Close() error {
	return nil
}

type AnyResource struct {
	name     string
	location string
	exists   bool
	err      error
	reader   stdio.ReadCloser
}

func (a *AnyResource) Name() string {
	return a.name
}

func (a *AnyResource) Location() string {
	return a.location
}

func (a *AnyResource) Exists() bool {
	return a.exists
}

func (a *AnyResource) Reader() (stdio.ReadCloser, error) {
	if a.err != nil {
		return nil, a.err
	}
	return a.reader, nil
}

func TestYamlPropertySourceLoader_Extensions(t *testing.T) {
	// given
	loader := NewYamlPropertySourceLoader()

	// when
	extensions := loader.Extensions()

	// then
	assert.Equal(t, []string{"yaml", "yml"}, extensions)
}

func TestYamlPropertySourceLoader_Load(t *testing.T) {

	testCases := []struct {
		name           string
		ctx            context.Context
		sourceName     string
		resource       io.Resource
		wantErr        error
		wantProperties map[string]string
	}{
		{
			name:       "nil context",
			ctx:        nil,
			sourceName: "anySourceName",
			wantErr:    errors.New("nil context"),
			resource:   nil,
		},
		{
			name:       "empty source name",
			ctx:        context.Background(),
			sourceName: "",
			wantErr:    errors.New("empty source name"),
			resource:   io.NewFileResource("resources/procyon.yaml"),
		},
		{
			name:       "nil resource",
			ctx:        context.Background(),
			sourceName: "anySourceName",
			wantErr:    errors.New("nil resource"),
			resource:   nil,
		},
		{
			name:       "invalid yaml",
			ctx:        context.Background(),
			sourceName: "anySourceName",
			resource: &AnyResource{
				reader: &FakeFile{
					contents: "version 2.1\njobs:\n  image: 'nginx:latest'",
				},
			},
			wantErr: errors.New("yaml: line 2: mapping values are not allowed in this context"),
		},
		{
			name:       "reader error",
			ctx:        context.Background(),
			sourceName: "anySourceName",
			resource: &AnyResource{
				err: errors.New("reader error"),
			},
			wantErr: errors.New("reader error"),
		},
		{
			name:       "read error",
			ctx:        context.Background(),
			sourceName: "anySourceName",
			resource: &AnyResource{
				reader: &FakeFile{
					readErr: errors.New("read error"),
				},
			},
			wantErr: errors.New("read error"),
		},
		{
			name:       "valid yaml",
			ctx:        context.Background(),
			sourceName: "anySourceName",
			resource: &AnyResource{
				reader: &FakeFile{
					contents: "version: 2.1\njobs:\n  image: nginx:latest",
				},
			},
			wantProperties: map[string]string{
				"version":    "2.1",
				"jobs.image": "nginx:latest",
			},
		},
		{
			name:       "valid yaml with array",
			ctx:        context.Background(),
			sourceName: "anySourceName",
			resource: &AnyResource{
				reader: &FakeFile{
					contents: "version: 2.1\njobs:\n  build:\n    docker:\n      image: cimg/base:2023.03\n    steps:\n      - checkout\n      - echo \"this is the build job\"",
				},
			},
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
			loader := NewYamlPropertySourceLoader()

			// when

			propSource, err := loader.Load(tc.ctx, tc.sourceName, tc.resource)

			// then
			if tc.wantErr != nil {
				require.Nil(t, propSource)

				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			require.NotNil(t, propSource)

			for wantKey, wantValue := range tc.wantProperties {
				value, exists := propSource.Value(wantKey)
				assert.True(t, exists)
				assert.Equal(t, wantValue, value)
			}
		})
	}
}
