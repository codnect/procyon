// Copyright 2026 Codnect
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

package http

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type AnyEndpointMapper struct{}

func (m *AnyEndpointMapper) MapEndpoints(endpoints Endpoints) {
	endpoints.MapGet("/test", &AnyHandler{})
}

func TestNewEndpoint(t *testing.T) {
	testCases := []struct {
		name     string
		method   Method
		path     string
		delegate RequestDelegate

		wantPanic error
	}{
		{
			name:      "nil request delegate",
			method:    MethodGet,
			path:      "/",
			delegate:  nil,
			wantPanic: errors.New("nil request delegate"),
		},
		{
			name:   "valid endpoint",
			method: MethodGet,
			path:   "/",
			delegate: func(ctx *Context) error {
				return nil
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					NewEndpoint(tc.method, tc.path, tc.delegate)
				})
				return
			}

			endpoint := NewEndpoint(tc.method, tc.path, tc.delegate)

			// then
			require.NotNil(t, endpoint)
		})
	}
}

func TestEndpoint_Path(t *testing.T) {
	// given
	endpoint := NewEndpoint(MethodGet, "/test", func(ctx *Context) error {
		return nil
	})

	// when
	path := endpoint.Path()

	// then
	assert.Equal(t, "/test", path)
}

func TestEndpoint_Method(t *testing.T) {
	// given
	endpoint := NewEndpoint(MethodGet, "/test", func(ctx *Context) error {
		return nil
	})

	// when
	method := endpoint.Method()

	// then
	assert.Equal(t, MethodGet, method)
}

func TestEndpoint_RequestDelegate(t *testing.T) {
	// given
	endpoint := NewEndpoint(MethodGet, "/test", func(ctx *Context) error {
		return nil
	})

	// when
	requestDelegate := endpoint.RequestDelegate()

	// then
	assert.NotNil(t, requestDelegate)
}

func TestDefaultEndpointRegistry_Endpoints(t *testing.T) {
	// given
	endpoint := NewEndpoint(MethodGet, "/test", func(ctx *Context) error {
		return nil
	})

	registry := NewDefaultEndpointRegistry()
	err := registry.Register(endpoint)
	assert.NoError(t, err)

	// when
	endpoints := registry.Endpoints()

	// then
	assert.Len(t, endpoints, 1)
	assert.Equal(t, endpoint, endpoints[0])
}

func TestNewEndpointGroup(t *testing.T) {
	testCases := []struct {
		name       string
		prefix     string
		wantPrefix string
	}{
		{
			name:       "empty prefix",
			prefix:     "",
			wantPrefix: "/",
		},
		{
			name:       "slash prefix",
			prefix:     "/",
			wantPrefix: "/",
		},
		{
			name:       "no slash prefix",
			prefix:     "/api/v1/test",
			wantPrefix: "/api/v1/test",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			endpointGroup := newEndpointGroup(tc.prefix)

			// then
			assert.Equal(t, tc.wantPrefix, endpointGroup.prefix)
		})
	}
}

func TestEndpointGroup_MapAny(t *testing.T) {
	// given
	anyHandler := HandlerFunc(func(ctx *Context) (Result, error) {
		return nil, nil
	})

	endpointGroup := newEndpointGroup("/prefix")

	// when
	endpointBuilder := endpointGroup.MapAny("/test", anyHandler)

	// then
	assert.NotNil(t, endpointBuilder)
	assert.Len(t, endpointGroup.routes, 1)

	route := endpointGroup.routes[0]
	assert.Equal(t, "/prefix/test", route.path)
	assert.Len(t, route.methods, 0)
	assert.NotNil(t, route.handler)
}

func TestEndpointGroup_MapMethods(t *testing.T) {
	// given
	anyHandler := HandlerFunc(func(ctx *Context) (Result, error) {
		return nil, nil
	})
	methods := []Method{MethodGet, MethodPatch}

	endpointGroup := newEndpointGroup("/prefix")

	// when
	endpointBuilder := endpointGroup.MapMethods("/test", methods, anyHandler)

	// then
	assert.NotNil(t, endpointBuilder)
	assert.Len(t, endpointGroup.routes, 1)

	route := endpointGroup.routes[0]
	assert.Equal(t, "/prefix/test", route.path)
	assert.ElementsMatch(t, methods, route.methods)
	assert.NotNil(t, route.handler)
}

func TestEndpointGroup_MapGet(t *testing.T) {
	// given
	anyHandler := HandlerFunc(func(ctx *Context) (Result, error) {
		return nil, nil
	})
	methods := []Method{MethodGet}

	endpointGroup := newEndpointGroup("/prefix")

	// when
	endpointBuilder := endpointGroup.MapGet("/test", anyHandler)

	// then
	assert.NotNil(t, endpointBuilder)
	assert.Len(t, endpointGroup.routes, 1)

	route := endpointGroup.routes[0]
	assert.Equal(t, "/prefix/test", route.path)
	assert.ElementsMatch(t, methods, route.methods)
	assert.NotNil(t, route.handler)
}

func TestEndpointGroup_MapPost(t *testing.T) {
	// given
	anyHandler := HandlerFunc(func(ctx *Context) (Result, error) {
		return nil, nil
	})
	methods := []Method{MethodPost}

	endpointGroup := newEndpointGroup("/prefix")

	// when
	endpointBuilder := endpointGroup.MapPost("/test", anyHandler)

	// then
	assert.NotNil(t, endpointBuilder)
	assert.Len(t, endpointGroup.routes, 1)

	route := endpointGroup.routes[0]
	assert.Equal(t, "/prefix/test", route.path)
	assert.ElementsMatch(t, methods, route.methods)
	assert.NotNil(t, route.handler)
}

func TestEndpointGroup_MapPut(t *testing.T) {
	// given
	anyHandler := HandlerFunc(func(ctx *Context) (Result, error) {
		return nil, nil
	})
	methods := []Method{MethodPut}

	endpointGroup := newEndpointGroup("/prefix")

	// when
	endpointBuilder := endpointGroup.MapPut("/test", anyHandler)

	// then
	assert.NotNil(t, endpointBuilder)
	assert.Len(t, endpointGroup.routes, 1)

	route := endpointGroup.routes[0]
	assert.Equal(t, "/prefix/test", route.path)
	assert.ElementsMatch(t, methods, route.methods)
	assert.NotNil(t, route.handler)
}

func TestEndpointGroup_MapDelete(t *testing.T) {
	// given
	anyHandler := HandlerFunc(func(ctx *Context) (Result, error) {
		return nil, nil
	})
	methods := []Method{MethodDelete}

	endpointGroup := newEndpointGroup("/prefix")

	// when
	endpointBuilder := endpointGroup.MapDelete("/test", anyHandler)

	// then
	assert.NotNil(t, endpointBuilder)
	assert.Len(t, endpointGroup.routes, 1)

	route := endpointGroup.routes[0]
	assert.Equal(t, "/prefix/test", route.path)
	assert.ElementsMatch(t, methods, route.methods)
	assert.NotNil(t, route.handler)
}

func TestEndpointGroup_MapPatch(t *testing.T) {
	// given
	anyHandler := HandlerFunc(func(ctx *Context) (Result, error) {
		return nil, nil
	})
	methods := []Method{MethodPatch}

	endpointGroup := newEndpointGroup("/prefix")

	// when
	endpointBuilder := endpointGroup.MapPatch("/test", anyHandler)

	// then
	assert.NotNil(t, endpointBuilder)
	assert.Len(t, endpointGroup.routes, 1)

	route := endpointGroup.routes[0]
	assert.Equal(t, "/prefix/test", route.path)
	assert.ElementsMatch(t, methods, route.methods)
	assert.NotNil(t, route.handler)
}

func TestEndpointGroup_MapGroup(t *testing.T) {
	// given
	endpointGroup := newEndpointGroup("/prefix")

	// when
	group := endpointGroup.MapGroup("/test")

	// then
	assert.NotNil(t, group)
	assert.Equal(t, "/prefix/test", group.prefix)
	assert.Len(t, group.routes, 0)
}

func TestNewEndpointMappingProcessor(t *testing.T) {
	testCases := []struct {
		name              string
		endpointRegistrar EndpointRegistrar
		executors         ResultExecutorRegistry
		wantPanic         error
	}{
		{
			name:              "nil endpoint registrar",
			endpointRegistrar: nil,
			executors:         newDefaultResultExecutorRegistry(),
			wantPanic:         errors.New("nil endpoint registrar"),
		},
		{
			name:              "nil result executor registry",
			endpointRegistrar: &DefaultEndpointRegistry{},
			executors:         nil,
			wantPanic:         errors.New("nil result executor registry"),
		},
		{
			name:              "valid dependencies",
			endpointRegistrar: &DefaultEndpointRegistry{},
			executors:         newDefaultResultExecutorRegistry(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					newEndpointMappingProcessor(
						tc.endpointRegistrar,
						tc.executors,
					)
				})
				return
			}

			processor := newEndpointMappingProcessor(
				tc.endpointRegistrar,
				tc.executors,
			)

			// then
			require.NotNil(t, processor)
			assert.Equal(t, tc.endpointRegistrar, processor.endpointRegistrar)
			assert.Equal(t, tc.executors, processor.executors)
		})
	}
}

func TestEndpointMappingProcessor_ProcessAfterInit(t *testing.T) {
	testCases := []struct {
		name         string
		instance     any
		preCondition func(registrar *AnyEndpointRegistrar)
		wantErr      error
	}{
		{
			name:     "not endpoint mapper",
			instance: &struct{}{},
		},
		{
			name:     "endpoint mapper",
			instance: &AnyEndpointMapper{},
			preCondition: func(registrar *AnyEndpointRegistrar) {
				registrar.
					On("Register", mock.AnythingOfType("*http.Endpoint")).
					Return(nil)
			},
		},
		{
			name:     "endpoint registration error",
			instance: &AnyEndpointMapper{},
			preCondition: func(registrar *AnyEndpointRegistrar) {
				registrar.
					On("Register", mock.AnythingOfType("*http.Endpoint")).
					Return(errors.New("register error"))
			},
			wantErr: errors.New(
				`map endpoints for "endpoint registration error": register error`,
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			registrar := &AnyEndpointRegistrar{}

			if tc.preCondition != nil {
				tc.preCondition(registrar)
			}

			processor := newEndpointMappingProcessor(
				registrar,
				newDefaultResultExecutorRegistry(),
			)

			// when
			result, err := processor.ProcessAfterInit(
				context.Background(),
				tc.name,
				tc.instance,
			)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				assert.Nil(t, result)
				return
			}

			require.NoError(t, err)
			assert.Same(t, tc.instance, result)

			registrar.AssertExpectations(t)
		})
	}
}

func TestJoinPaths(t *testing.T) {
	testCases := []struct {
		name     string
		paths    []string
		wantPath string
	}{
		{
			name:     "prefix without leading slash",
			paths:    []string{"api", "users"},
			wantPath: "/api/users",
		},
		{
			name:     "last element has trailing slash",
			paths:    []string{"/api", "/users/"},
			wantPath: "/api/users/",
		},
		{
			name:     "last element has no trailing slash",
			paths:    []string{"/api", "/users"},
			wantPath: "/api/users",
		},
		{
			name:     "double slashes are cleaned",
			paths:    []string{"/api/", "/users"},
			wantPath: "/api/users",
		},
		{
			name:     "single element",
			paths:    []string{"/api"},
			wantPath: "/api",
		},
		{
			name:     "empty second element",
			paths:    []string{"/api", ""},
			wantPath: "/api",
		},
		{
			name:     "elements with extra slashes are cleaned",
			paths:    []string{"/api/", "/users/", "/profile"},
			wantPath: "/api/users/profile",
		},
		{
			name:     "empty middle element is ignored",
			paths:    []string{"/api", "", "/users"},
			wantPath: "/api/users",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := joinPaths(tc.paths...)
			assert.Equal(t, tc.wantPath, result)
		})
	}
}
