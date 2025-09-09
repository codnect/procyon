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
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

type AnyPropertySourceLoader struct {
	mock.Mock
}

func (a *AnyPropertySourceLoader) Extensions() []string {
	result := a.Called()
	extensions := result.Get(0)
	if extensions == nil {
		return nil
	}

	return extensions.([]string)
}

func (a *AnyPropertySourceLoader) Load(name string, reader io.Reader) (property.Source, error) {
	result := a.Called(name, reader)
	source := result.Get(0)
	if source == nil {
		return nil, result.Error(1)
	}

	return source.(property.Source), result.Error(1)
}

type AnyConfigLoader struct {
	mock.Mock
}

func (a *AnyConfigLoader) IsLoadable(resource Resource) bool {
	result := a.Called(resource)
	return result.Bool(0)
}

func (a *AnyConfigLoader) Load(ctx context.Context, resource Resource) (*Data, error) {
	result := a.Called(ctx, resource)
	data := result.Get(0)
	if data == nil {
		return nil, result.Error(1)
	}

	return data.(*Data), nil
}

func TestDefaultLoader_IsLoadable(t *testing.T) {
	testCases := []struct {
		name         string
		preCondition func(anyDirFs *AnyDirFs) Resource
		path         string
		wantResult   bool
	}{
		{
			name: "nil resource",
			preCondition: func(anyDirFs *AnyDirFs) Resource {
				return nil
			},
			wantResult: false,
		},
		{
			name: "nil property source loader",
			preCondition: func(anyDirFs *AnyDirFs) Resource {
				return newURLResource(nil, "", nil)
			},
			wantResult: false,
		},
		{
			name: "resource does not exist",
			preCondition: func(anyDirFs *AnyDirFs) Resource {
				anyDirFs.On("Stat", "resources/procyon.yml").
					Return(nil, errors.New("no file")).Once()

				return newFileResource(anyDirFs, "procyon.yml", "resources/", "", property.NewYamlSourceLoader())
			},
			path:       "resources/",
			wantResult: false,
		},
		{
			name: "resource exists",
			preCondition: func(anyDirFs *AnyDirFs) Resource {
				fileInfo := &AnyFileInfo{}
				anyDirFs.On("Stat", "resources/procyon.yaml").
					Return(fileInfo, nil).Once()

				return newFileResource(anyDirFs, "procyon.yaml", "resources/", "", property.NewYamlSourceLoader())
			},
			path:       "resources/",
			wantResult: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			loader := NewDefaultLoader()

			anyDirFs := &AnyDirFs{
				dir: tc.path,
			}

			resource := tc.preCondition(anyDirFs)

			// when
			result := loader.IsLoadable(resource)

			// then
			assert.Equal(t, tc.wantResult, result)
		})
	}
}

func TestDefaultLoader_Load(t *testing.T) {
	testCases := []struct {
		name         string
		preCondition func(anyDirFs *AnyDirFs) Resource
		ctx          context.Context
		path         string
		wantErr      error
		wantConfig   map[string]any
	}{
		{
			name: "nil context",
			preCondition: func(anyDirFs *AnyDirFs) Resource {
				return newURLResource(nil, "", nil)
			},
			wantErr: errors.New("nil context"),
		},
		{
			name: "nil resource",
			preCondition: func(anyDirFs *AnyDirFs) Resource {
				return nil
			},
			ctx:     context.Background(),
			wantErr: errors.New("nil resource"),
		},
		{
			name: "nil property source loader",
			preCondition: func(anyDirFs *AnyDirFs) Resource {
				return newURLResource(nil, "", nil)
			},
			ctx:     context.Background(),
			wantErr: errors.New("nil resource loader"),
		},
		{
			name: "no resource found",
			preCondition: func(anyDirFs *AnyDirFs) Resource {
				anyDirFs.On("Stat", "resources/procyon.yml").
					Return(nil, errors.New("no file")).Once()

				return newFileResource(anyDirFs, "procyon.yml", "resources/procyon.yaml", "", property.NewYamlSourceLoader())
			},
			ctx:     context.Background(),
			path:    "resources/",
			wantErr: errors.New("no resource found"),
			wantConfig: map[string]any{
				"anyKey": "anyValue",
			},
		},
		{
			name: "reader error",
			preCondition: func(anyDirFs *AnyDirFs) Resource {
				fileInfo := &AnyFileInfo{}
				anyDirFs.On("Stat", "resources/procyon.yaml").
					Return(fileInfo, nil).Once()

				anyDirFs.On("Open", "resources/procyon.yaml").
					Return(nil, errors.New("reader error")).Once()

				return newFileResource(anyDirFs, "procyon.yaml", "resources/procyon.yaml", "", property.NewYamlSourceLoader())
			},
			ctx:     context.Background(),
			path:    "resources/",
			wantErr: errors.New("reader error"),
		},
		{
			name: "load error",
			preCondition: func(anyDirFs *AnyDirFs) Resource {
				fileInfo := &AnyFileInfo{}
				anyDirFs.On("Stat", "resources/procyon.yaml").
					Return(fileInfo, nil).Once()

				fakeFile := &FakeFile{}
				anyDirFs.On("Open", "resources/procyon.yaml").
					Return(fakeFile, nil).Once()

				anyLoader := &AnyPropertySourceLoader{}
				anyLoader.On("Load", "resources/procyon.yaml", fakeFile).
					Return(nil, errors.New("load error")).Once()

				return newFileResource(anyDirFs, "procyon.yaml", "resources/procyon.yaml", "", anyLoader)
			},
			ctx:     context.Background(),
			path:    "resources/",
			wantErr: errors.New("load error"),
		},
		{
			name: "load config data",
			preCondition: func(anyDirFs *AnyDirFs) Resource {
				fileInfo := &AnyFileInfo{}
				anyDirFs.On("Stat", "resources/procyon.yaml").
					Return(fileInfo, nil).Once()

				fakeFile := &FakeFile{}
				anyDirFs.On("Open", "resources/procyon.yaml").
					Return(fakeFile, nil).Once()

				anyLoader := &AnyPropertySourceLoader{}
				anyMapSource := property.NewMapSource("anyMapName", map[string]any{
					"anyKey": "anyValue",
				})
				anyLoader.On("Load", "resources/procyon.yaml", fakeFile).
					Return(anyMapSource, nil).Once()

				return newFileResource(anyDirFs, "procyon.yaml", "resources/procyon.yaml", "", anyLoader)
			},
			ctx:     context.Background(),
			path:    "resources/",
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			loader := NewDefaultLoader()

			anyDirFs := &AnyDirFs{
				dir: tc.path,
			}

			resource := tc.preCondition(anyDirFs)

			// when
			data, err := loader.Load(tc.ctx, resource)

			// then
			if tc.wantErr != nil {
				require.Nil(t, data)

				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NotNil(t, data)
			assert.NotNil(t, data.PropertySource())

			for wantKey, wantVal := range tc.wantConfig {
				val, exists := data.PropertySource().Property(wantKey)
				assert.True(t, exists)
				assert.Equal(t, wantVal, val)
			}
		})
	}
}
