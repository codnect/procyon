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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestEndpointDataSource_Endpoints(t *testing.T) {
	// given
	endpoint := NewEndpoint(MethodGet, "/test", func(ctx *Context) error {
		return nil
	})

	reqEndpointDataSource := NewEndpointDataSource(endpoint)

	// when
	endpoints := reqEndpointDataSource.Endpoints()

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
