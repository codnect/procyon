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
	"testing"

	"codnect.io/procyon/io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type AnyResourceResolver struct {
	mock.Mock
}

func (a *AnyResourceResolver) Resolve(ctx context.Context, location string) (io.Resource, error) {
	result := a.Called(ctx, location)
	resource := result.Get(0)
	if resource == nil {
		return nil, result.Error(1)
	}

	return resource.(io.Resource), result.Error(1)
}

func TestNewData(t *testing.T) {
	testCases := []struct {
		name           string
		propertySource PropertySource
		profile        string
		wantPanic      error
	}{
		{
			name:           "nil property source",
			propertySource: nil,
			profile:        DefaultProfile,
			wantPanic:      errors.New("nil property source"),
		},
		{
			name:           "valid property source",
			profile:        DefaultProfile,
			propertySource: NewMapPropertySource("anyMapName", make(map[string]any)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					NewData(tc.propertySource, tc.profile)
				})
				return
			}

			data := NewData(tc.propertySource, tc.profile)

			// then
			require.NotNil(t, data)
		})
	}
}

func TestData_PropertySource(t *testing.T) {
	// given
	mapPropSource := NewMapPropertySource("anyMapName", make(map[string]any))
	data := NewData(mapPropSource, DefaultProfile)

	// when
	source := data.PropertySource()

	// then
	assert.NotNil(t, source)
	assert.Equal(t, mapPropSource, source)
}

func TestData_Profile(t *testing.T) {
	testCases := []struct {
		name           string
		propertySource PropertySource
		profile        string
		wantProfile    string
	}{
		{
			name:           "empty profile",
			propertySource: NewMapPropertySource("anyMapName", make(map[string]any)),
			profile:        "",
			wantProfile:    DefaultProfile,
		},
		{
			name:           "any profile",
			propertySource: NewMapPropertySource("anyMapName", make(map[string]any)),
			profile:        "anyProfile",
			wantProfile:    "anyProfile",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			data := NewData(tc.propertySource, tc.profile)

			// when
			profile := data.Profile()

			// then
			assert.Equal(t, tc.wantProfile, profile)
		})
	}
}

func TestNewStandardDataLoader(t *testing.T) {
	testCases := []struct {
		name             string
		resourceResolver io.ResourceResolver
		loaders          []PropertySourceLoader
		wantPanic        error
	}{
		{
			name:             "nil resource resolver",
			resourceResolver: nil,
			wantPanic:        errors.New("nil resource resolver"),
		},
		{
			name:             "valid inputs",
			resourceResolver: io.NewDefaultResourceResolver(),
			loaders: []PropertySourceLoader{
				NewYamlPropertySourceLoader(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					NewStandardDataLoader(tc.resourceResolver, tc.loaders...)
				})
				return
			}

			dataLoader := NewStandardDataLoader(tc.resourceResolver, tc.loaders...)

			// then
			require.NotNil(t, dataLoader)
		})
	}
}

func TestStandardDataLoader_Load(t *testing.T) {
	testCases := []struct {
		name    string
		loaders []PropertySourceLoader

		preCondition func(resourceResolver *AnyResourceResolver)
		ctx          context.Context
		location     string
		profiles     []string

		wantErr        error
		wantProperties map[string]any
	}{
		{
			name:    "nil context",
			ctx:     nil,
			wantErr: errors.New("nil context"),
		},
		{
			name:     "empty location",
			ctx:      context.Background(),
			location: "",
			wantErr:  errors.New("empty or blank location"),
		},
		{
			name:     "blank location",
			ctx:      context.Background(),
			location: " ",
			wantErr:  errors.New("empty or blank location"),
		},
		{
			name: "directory: resolve error without profile",
			preCondition: func(resourceResolver *AnyResourceResolver) {
				resourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"),
					"resources/procyon.yaml").
					Return(nil, errors.New("resolve error"))
			},
			ctx: context.Background(),
			loaders: []PropertySourceLoader{
				NewYamlPropertySourceLoader(),
			},
			location: "resources/",
			wantErr:  errors.New("resolve error"),
		},
		{
			name: "directory: resolve error with profile",
			preCondition: func(resourceResolver *AnyResourceResolver) {
				resourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"),
					"resources/procyon-dev.yaml").
					Return(nil, errors.New("resolve error"))
			},
			ctx: context.Background(),
			loaders: []PropertySourceLoader{
				NewYamlPropertySourceLoader(),
			},
			location: "resources/",
			wantErr:  errors.New("resolve error"),
			profiles: []string{"dev"},
		},
		{
			name: "directory: resource does not exist",
			preCondition: func(resourceResolver *AnyResourceResolver) {
				anyResource := &AnyResource{}
				resourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"),
					"resources/procyon.yaml").
					Return(anyResource, nil)
				resourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"),
					"resources/procyon.yml").
					Return(anyResource, nil)
			},
			ctx: context.Background(),
			loaders: []PropertySourceLoader{
				NewYamlPropertySourceLoader(),
			},
			location: "resources/",
		},
		{
			name: "directory: resource does not exist with profile",
			preCondition: func(resourceResolver *AnyResourceResolver) {
				anyResource := &AnyResource{}
				resourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"),
					"resources/procyon-dev.yaml").
					Return(anyResource, nil)
				resourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"),
					"resources/procyon-dev.yml").
					Return(anyResource, nil)
			},
			ctx: context.Background(),
			loaders: []PropertySourceLoader{
				NewYamlPropertySourceLoader(),
			},
			location: "resources/",
			profiles: []string{"dev"},
		},
		{
			name: "directory: load error without profile",
			preCondition: func(resourceResolver *AnyResourceResolver) {
				anyResource := &AnyResource{
					exists: true,
					err:    errors.New("load error"),
				}
				resourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"),
					"resources/procyon.yaml").
					Return(anyResource, nil)
				resourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"),
					"resources/procyon.yml").
					Return(anyResource, nil)
			},
			ctx: context.Background(),
			loaders: []PropertySourceLoader{
				NewYamlPropertySourceLoader(),
			},
			location: "resources/",
			wantErr:  errors.New("load error"),
		},
		{
			name: "directory: load error with profile",
			preCondition: func(resourceResolver *AnyResourceResolver) {
				anyResource := &AnyResource{
					exists: true,
					err:    errors.New("load error"),
				}
				anotherResource := &AnyResource{
					exists: false,
				}
				resourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"),
					"resources/procyon-dev.yaml").
					Return(anyResource, nil)
				resourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"),
					"resources/procyon-dev.yml").
					Return(anotherResource, nil)
			},
			ctx: context.Background(),
			loaders: []PropertySourceLoader{
				NewYamlPropertySourceLoader(),
			},
			location: "resources/",
			profiles: []string{"dev"},
			wantErr:  errors.New("load error"),
		},
		{
			name: "directory: load resource",
			preCondition: func(resourceResolver *AnyResourceResolver) {
				anyResource := &AnyResource{
					exists: true,
					reader: &FakeFile{
						contents: "version: 2.1\njobs:\n  image: nginx:latest",
					},
				}
				anotherResource := &AnyResource{
					exists: false,
				}
				resourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"),
					"resources/procyon.yaml").
					Return(anyResource, nil)
				resourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"),
					"resources/procyon.yml").
					Return(anotherResource, nil)
			},
			ctx: context.Background(),
			loaders: []PropertySourceLoader{
				NewYamlPropertySourceLoader(),
			},
			location: "resources/",
			wantProperties: map[string]any{
				"version":    2.1,
				"jobs.image": "nginx:latest",
			},
		},
		{
			name: "file: resolve error without profile",
			preCondition: func(resourceResolver *AnyResourceResolver) {
				resourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"),
					"resources/procyon.yaml").
					Return(nil, errors.New("resolve error"))
			},
			ctx: context.Background(),
			loaders: []PropertySourceLoader{
				NewYamlPropertySourceLoader(),
			},
			location: "resources/procyon.yaml",
			wantErr:  errors.New("resolve error"),
		},
		{
			name: "file: resolve error with default profile",
			preCondition: func(resourceResolver *AnyResourceResolver) {
				resourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"),
					"resources/procyon.yaml").
					Return(nil, errors.New("resolve error"))
			},
			ctx: context.Background(),
			loaders: []PropertySourceLoader{
				NewYamlPropertySourceLoader(),
			},
			location: "resources/procyon.yaml",
			wantErr:  errors.New("resolve error"),
			profiles: []string{"dev"},
		},
		{
			name: "file: resource does not exist without profile",
			preCondition: func(resourceResolver *AnyResourceResolver) {
				anyResource := &AnyResource{}
				resourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"),
					"resources/procyon.yaml").
					Return(anyResource, nil)
			},
			ctx: context.Background(),
			loaders: []PropertySourceLoader{
				NewYamlPropertySourceLoader(),
			},
			location: "resources/procyon.yaml",
		},
		{
			name: "file: resource does not exist with profile",
			preCondition: func(resourceResolver *AnyResourceResolver) {
				anyResource := &AnyResource{}
				resourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"),
					"resources/procyon.yaml").
					Return(anyResource, nil)
			},
			ctx: context.Background(),
			loaders: []PropertySourceLoader{
				NewYamlPropertySourceLoader(),
			},
			location: "resources/procyon.yaml",
			profiles: []string{"dev"},
		},
		{
			name: "file: load error for without profile",
			preCondition: func(resourceResolver *AnyResourceResolver) {
				anyResource := &AnyResource{
					exists: true,
					err:    errors.New("load error"),
				}
				resourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"),
					"resources/procyon.yaml").
					Return(anyResource, nil)
			},
			ctx: context.Background(),
			loaders: []PropertySourceLoader{
				NewYamlPropertySourceLoader(),
			},
			location: "resources/procyon.yaml",
			wantErr:  errors.New("load error"),
		},
		{
			name: "file: load error for with profile",
			preCondition: func(resourceResolver *AnyResourceResolver) {
				anyResource := &AnyResource{
					exists: true,
					err:    errors.New("load error"),
				}
				resourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"),
					"resources/procyon.yaml").
					Return(anyResource, nil)
			},
			ctx: context.Background(),
			loaders: []PropertySourceLoader{
				NewYamlPropertySourceLoader(),
			},
			location: "resources/procyon.yaml",
			profiles: []string{"dev"},
			wantErr:  errors.New("load error"),
		},
		{
			name: "file: load resource",
			preCondition: func(resourceResolver *AnyResourceResolver) {
				anyResource := &AnyResource{
					exists: true,
					reader: &FakeFile{
						contents: "version: 2.1\njobs:\n  image: nginx:latest",
					},
				}
				resourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"),
					"resources/procyon.yaml").
					Return(anyResource, nil)
			},
			ctx: context.Background(),
			loaders: []PropertySourceLoader{
				NewYamlPropertySourceLoader(),
			},
			location: "resources/procyon.yaml",
			wantProperties: map[string]any{
				"version":    2.1,
				"jobs.image": "nginx:latest",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			anyResourceResolver := &AnyResourceResolver{}

			if tc.preCondition != nil {
				tc.preCondition(anyResourceResolver)
			}

			dataLoader := NewStandardDataLoader(anyResourceResolver, tc.loaders...)

			// when
			configData, err := dataLoader.Load(tc.ctx, tc.location, tc.profiles...)

			// then
			if tc.wantErr != nil {
				require.Nil(t, configData)

				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			if len(tc.wantProperties) == 0 {
				return
			}

			require.Len(t, configData, 1)

			data := configData[0]
			require.NotNil(t, data)

			for wantKey, wantValue := range tc.wantProperties {
				value, exists := data.PropertySource().Value(wantKey)
				assert.True(t, exists)
				assert.Equal(t, wantValue, value)
			}
		})
	}
}
