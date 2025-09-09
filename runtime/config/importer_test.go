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
	"testing"
)

func TestDefaultImporter_Import(t *testing.T) {
	testCases := []struct {
		name           string
		preCondition   func() ([]LocationResolver, []Loader)
		ctx            context.Context
		location       string
		profiles       []string
		wantErr        error
		wantConfigData []struct {
			name       string
			properties map[string]any
		}
	}{
		{
			name:    "nil context",
			wantErr: errors.New("nil context"),
		},
		{
			name: "no resolvable",
			ctx:  context.Background(),
			preCondition: func() ([]LocationResolver, []Loader) {
				anyLocResolver := &AnyLocationResolver{}
				anyLocResolver.On("IsResolvable", "resources/procyon.yaml").Return(false)
				return []LocationResolver{
					anyLocResolver,
				}, []Loader{}
			},
			location: "resources/procyon.yaml",
		},
		{
			name: "resolve error",
			ctx:  context.Background(),
			preCondition: func() ([]LocationResolver, []Loader) {

				anyLocResolver := &AnyLocationResolver{}
				anyLocResolver.On("IsResolvable", "resources/procyon.yaml").Return(true)

				anyLocResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yaml", []string{"dev"}).
					Return(nil, errors.New("resolve error"))

				return []LocationResolver{
					anyLocResolver,
				}, []Loader{}
			},
			location: "resources/procyon.yaml",
			profiles: []string{"dev"},
			wantErr:  errors.New("resolve error"),
		},
		{
			name: "no loaders (resolvable but not loadable)",
			ctx:  context.Background(),
			preCondition: func() ([]LocationResolver, []Loader) {

				anyLocResolver := &AnyLocationResolver{}
				anyLocResolver.On("IsResolvable", "resources/procyon.yaml").Return(true)

				resources := make([]Resource, 0)
				urlResource := newURLResource(nil, "", nil)
				resources = append(resources, urlResource)
				anyLocResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yaml", []string{"dev"}).
					Return(resources, nil)

				anyConfigLoader := &AnyConfigLoader{}
				anyConfigLoader.On("IsLoadable", urlResource).Return(false)

				return []LocationResolver{
						anyLocResolver,
					}, []Loader{
						anyConfigLoader,
					}
			},
			location: "resources/procyon.yaml",
			profiles: []string{"dev"},
			wantErr:  errors.New("no loaders"),
		},
		{
			name: "multiple loaders",
			ctx:  context.Background(),
			preCondition: func() ([]LocationResolver, []Loader) {

				anyLocResolver := &AnyLocationResolver{}
				anyLocResolver.On("IsResolvable", "resources/procyon.yaml").Return(true)

				resources := make([]Resource, 0)
				urlResource := newURLResource(nil, "", nil)
				resources = append(resources, urlResource)
				anyLocResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yaml", []string{"dev"}).
					Return(resources, nil)

				anyConfigLoader := &AnyConfigLoader{}
				anyConfigLoader.On("IsLoadable", urlResource).Return(true)

				anotherConfigLoader := &AnyConfigLoader{}
				anotherConfigLoader.On("IsLoadable", urlResource).Return(true)

				return []LocationResolver{
						anyLocResolver,
					}, []Loader{
						anyConfigLoader,
						anotherConfigLoader,
					}
			},
			location: "resources/procyon.yaml",
			profiles: []string{"dev"},
			wantErr:  errors.New("multiple loaders"),
		},
		{
			name: "load error",
			ctx:  context.Background(),
			preCondition: func() ([]LocationResolver, []Loader) {

				anyLocResolver := &AnyLocationResolver{}
				anyLocResolver.On("IsResolvable", "resources/procyon.yaml").Return(true)

				resources := make([]Resource, 0)
				urlResource := newURLResource(nil, "", nil)
				resources = append(resources, urlResource)
				anyLocResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yaml", []string{"dev"}).
					Return(resources, nil)

				anyConfigLoader := &AnyConfigLoader{}
				anyConfigLoader.On("IsLoadable", urlResource).Return(true)
				anyConfigLoader.On("Load", mock.AnythingOfType("context.backgroundCtx"), urlResource).
					Return(nil, errors.New("load error"))

				return []LocationResolver{
						anyLocResolver,
					}, []Loader{
						anyConfigLoader,
					}
			},
			location: "resources/procyon.yaml",
			profiles: []string{"dev"},
			wantErr:  errors.New("load error"),
		},
		{
			name: "import data",
			ctx:  context.Background(),
			preCondition: func() ([]LocationResolver, []Loader) {

				anyLocResolver := &AnyLocationResolver{}
				anyLocResolver.On("IsResolvable", "resources/procyon.yaml").Return(true)

				resources := make([]Resource, 0)
				urlResource := newURLResource(nil, "", nil)
				resources = append(resources, urlResource)
				anyLocResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yaml", []string{"dev"}).
					Return(resources, nil)

				anyConfigLoader := &AnyConfigLoader{}
				anyConfigLoader.On("IsLoadable", urlResource).Return(true)

				mapSource := property.NewMapSource("anyMapName", map[string]any{
					"anyKey": "anyValue",
				})
				data := NewData(mapSource)
				anyConfigLoader.On("Load", mock.AnythingOfType("context.backgroundCtx"), urlResource).
					Return(data, nil)

				return []LocationResolver{
						anyLocResolver,
					}, []Loader{
						anyConfigLoader,
					}
			},
			location: "resources/procyon.yaml",
			profiles: []string{"dev"},
			wantErr:  nil,
			wantConfigData: []struct {
				name       string
				properties map[string]any
			}{
				{
					name: "anyMapName",
					properties: map[string]any{
						"anyKey": "anyValue",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			locResolvers := make([]LocationResolver, 0)
			loaders := make([]Loader, 0)

			if tc.preCondition != nil {
				pLocResolvers, pLoaders := tc.preCondition()
				if len(pLocResolvers) != 0 {
					locResolvers = append(locResolvers, pLocResolvers...)
				}

				if len(pLoaders) != 0 {
					loaders = append(loaders, pLoaders...)
				}
			}

			importer := NewDefaultImporter(locResolvers, loaders)

			// when
			configData, err := importer.Import(tc.ctx, tc.location, tc.profiles)

			// then
			if tc.wantErr != nil {
				require.Len(t, configData, 0)

				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.Len(t, configData, len(tc.wantConfigData))
			for _, wantData := range tc.wantConfigData {
				for _, data := range configData {
					propSource := data.PropertySource()

					require.NotNil(t, propSource)

					if propSource.Name() != wantData.name {
						continue
					}

					for wantKey, wantVal := range wantData.properties {
						val, exists := propSource.Property(wantKey)
						require.True(t, exists)
						require.Equal(t, wantVal, val)
					}

					break
				}
			}
		})
	}
}
