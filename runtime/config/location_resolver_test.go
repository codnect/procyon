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
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type AnyLocationResolver struct {
	mock.Mock
}

func (a *AnyLocationResolver) IsResolvable(location string) bool {
	result := a.Called(location)
	return result.Bool(0)
}

func (a *AnyLocationResolver) Resolve(ctx context.Context, location string, profiles ...string) ([]Resource, error) {
	result := a.Called(ctx, location, profiles)
	resources := result.Get(0)
	if resources == nil {
		return nil, result.Error(1)
	}

	return resources.([]Resource), nil
}

func TestDefaultLocationResolver_IsResolvable(t *testing.T) {
	// given
	resolver := NewDefaultLocationResolver()

	// when
	result := resolver.IsResolvable("anyLocation")

	// then
	assert.True(t, result)
}

func TestDefaultLocationResolver_Resolve(t *testing.T) {

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.RequestURI, ".yml") {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})

	testServer := httptest.NewServer(handler)
	defer testServer.Close()

	testCases := []struct {
		name         string
		preCondition func(resolver *DefaultLocationResolver)
		loaders      []property.SourceLoader
		ctx          context.Context
		location     string
		profiles     []string

		wantErr       error
		wantResources []string
	}{
		{
			name:     "nil context",
			ctx:      nil,
			location: "anyLocation",
			profiles: []string{"anyProfile"},
			wantErr:  errors.New("nil context"),
		},
		{
			name:     "empty location",
			ctx:      context.Background(),
			location: "",
			profiles: []string{"anyProfile"},
			wantErr:  errors.New("empty or blank location"),
		},
		{
			name:     "blank location",
			ctx:      context.Background(),
			location: " ",
			profiles: []string{"anyProfile"},
			wantErr:  errors.New("empty or blank location"),
		},
		{
			name:     "wrong url directory format",
			ctx:      context.Background(),
			location: "://missing-scheme/config/",
			loaders:  []property.SourceLoader{property.NewYamlSourceLoader()},
			profiles: []string{""},
			wantErr:  errors.New("parse \"://missing-scheme/config/procyon.yaml\": missing protocol scheme"),
		},
		{
			name:     "wrong url resource format",
			ctx:      context.Background(),
			location: "://missing-scheme/config.yaml",
			loaders:  []property.SourceLoader{property.NewYamlSourceLoader()},
			profiles: []string{""},
			wantErr:  errors.New("parse \"://missing-scheme/config.yaml\": missing protocol scheme"),
		},
		{
			name: "single-file resources with a file scheme, using default profiles",
			preCondition: func(resolver *DefaultLocationResolver) {
				resolver.fsProvider = func(path string) fs.FS {
					mockAnyFs := &AnyDirFs{
						dir: path,
					}
					fileInfo := &AnyFileInfo{}

					mockAnyFs.On("Stat", "resources/redis-config.yaml").Return(fileInfo, nil).Once()
					return mockAnyFs
				}
			},
			ctx:           context.Background(),
			loaders:       []property.SourceLoader{property.NewYamlSourceLoader()},
			location:      "file://resources/redis-config.yaml",
			profiles:      []string{},
			wantResources: []string{"resources/redis-config.yaml"},
		},
		{
			name: "single-file resources without a file scheme, using default profiles",
			preCondition: func(resolver *DefaultLocationResolver) {
				resolver.fsProvider = func(path string) fs.FS {
					mockAnyFs := &AnyDirFs{
						dir: path,
					}
					fileInfo := &AnyFileInfo{}

					mockAnyFs.On("Stat", "resources/db-config.yaml").Return(fileInfo, nil).Once()
					return mockAnyFs
				}
			},
			ctx:           context.Background(),
			loaders:       []property.SourceLoader{property.NewYamlSourceLoader()},
			location:      "resources/db-config.yaml",
			profiles:      []string{},
			wantResources: []string{"resources/db-config.yaml"},
		},
		{
			name: "multi-file resources using default profiles",
			preCondition: func(resolver *DefaultLocationResolver) {
				resolver.fsProvider = func(path string) fs.FS {
					mockAnyFs := &AnyDirFs{
						dir: path,
					}

					fileInfo := &AnyFileInfo{}
					mockAnyFs.On("Stat", "resources/redis-config.yaml").Return(fileInfo, nil).Once()
					mockAnyFs.On("Stat", "resources/db-config.yaml").Return(fileInfo, nil).Once()
					return mockAnyFs
				}
			},
			ctx:           context.Background(),
			loaders:       []property.SourceLoader{property.NewYamlSourceLoader()},
			location:      "file://resources/redis-config.yaml;resources/db-config.yaml",
			profiles:      []string{},
			wantResources: []string{"resources/redis-config.yaml", "resources/db-config.yaml"},
		},
		{
			name: "file-directory resources using default profiles",
			preCondition: func(resolver *DefaultLocationResolver) {
				resolver.fsProvider = func(path string) fs.FS {
					mockAnyFs := &AnyDirFs{
						dir: path,
					}
					fileInfo := &AnyFileInfo{}

					mockAnyFs.On("Stat", "resources/procyon.yml").Return(nil, errors.New("no file")).Once()
					mockAnyFs.On("Stat", "resources/procyon.yaml").Return(fileInfo, nil).Once()
					return mockAnyFs
				}
			},
			ctx:           context.Background(),
			loaders:       []property.SourceLoader{property.NewYamlSourceLoader()},
			location:      "resources/",
			profiles:      []string{},
			wantResources: []string{"resources/procyon.yaml"},
		},
		{
			name: "file-directory resources with custom profiles",
			preCondition: func(resolver *DefaultLocationResolver) {
				resolver.fsProvider = func(path string) fs.FS {
					mockAnyFs := &AnyDirFs{
						dir: path,
					}
					fileInfo := &AnyFileInfo{}

					mockAnyFs.On("Stat", "resources/procyon-dev.yml").Return(nil, errors.New("no file")).Once()
					mockAnyFs.On("Stat", "resources/procyon-test.yml").Return(nil, errors.New("no file")).Once()
					mockAnyFs.On("Stat", "resources/procyon-dev.yaml").Return(fileInfo, nil).Once()
					mockAnyFs.On("Stat", "resources/procyon-test.yaml").Return(fileInfo, nil).Once()
					return mockAnyFs
				}
			},
			ctx:           context.Background(),
			loaders:       []property.SourceLoader{property.NewYamlSourceLoader()},
			location:      "resources/",
			profiles:      []string{"dev", "test"},
			wantResources: []string{"resources/procyon-dev.yaml", "resources/procyon-test.yaml"},
		},
		{
			name:     "single-url resources using default profiles",
			ctx:      context.Background(),
			loaders:  []property.SourceLoader{property.NewYamlSourceLoader()},
			location: fmt.Sprintf("%s/redis-config.yaml", testServer.URL),
			profiles: []string{},
			wantResources: []string{
				fmt.Sprintf("%s/redis-config.yaml", testServer.URL),
			},
		},
		{
			name:     "multi-url resources using default profiles",
			ctx:      context.Background(),
			loaders:  []property.SourceLoader{property.NewYamlSourceLoader()},
			location: fmt.Sprintf("%s/redis-config.yaml;%s/db-config.yaml", testServer.URL, testServer.URL),
			profiles: []string{},
			wantResources: []string{
				fmt.Sprintf("%s/redis-config.yaml", testServer.URL),
				fmt.Sprintf("%s/db-config.yaml", testServer.URL),
			},
		},
		{
			name:     "url-directory resources with custom profiles",
			ctx:      context.Background(),
			loaders:  []property.SourceLoader{property.NewYamlSourceLoader()},
			location: fmt.Sprintf("%s/resources/", testServer.URL),
			profiles: []string{"dev", "test"},
			wantResources: []string{
				fmt.Sprintf("%s/resources/procyon-dev.yaml", testServer.URL),
				fmt.Sprintf("%s/resources/procyon-test.yaml", testServer.URL),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			resolver := NewDefaultLocationResolver(tc.loaders...)

			if tc.preCondition != nil {
				tc.preCondition(resolver)
			}

			// when
			resources, err := resolver.Resolve(tc.ctx, tc.location, tc.profiles...)

			// then
			if tc.wantErr != nil {
				require.Len(t, resources, 0)

				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			require.Len(t, resources, len(tc.wantResources))
			for _, wantResource := range tc.wantResources {
				foundResource := false
				for _, resource := range resources {
					if resource.Location() == wantResource {
						foundResource = true
						break
					}
				}

				if !foundResource {
					t.Errorf("no resource found: %s", wantResource)
				}
			}
		})
	}
}
