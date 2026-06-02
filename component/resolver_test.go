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

package component

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCanResolve_NilContainer(t *testing.T) {
	// given

	// when
	require.PanicsWithValue(t, "nil container", func() {
		CanResolve(nil, "anyInstance")
	})

	// then
}

func TestCanResolve(t *testing.T) {
	testCases := []struct {
		name         string
		container    *AnyMockContainer
		instanceName string
		preCondition func(container *AnyMockContainer)

		wantResult bool
		wantPanic  error
	}{
		{
			name:      "nil container",
			wantPanic: errors.New("nil container"),
		},
		{
			name:         "resolvable component",
			container:    &AnyMockContainer{},
			instanceName: "anyInstance",
			preCondition: func(container *AnyMockContainer) {
				container.On("CanResolve", "anyInstance").Return(true)
			},
			wantResult: true,
		},
		{
			name:         "non-resolvable component",
			container:    &AnyMockContainer{},
			instanceName: "anyInstance",
			preCondition: func(container *AnyMockContainer) {
				container.On("CanResolve", "anyInstance").Return(false)
			},
			wantResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			if tc.preCondition != nil {
				tc.preCondition(tc.container)
			}

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					CanResolve(nil, tc.instanceName)
				})
				return
			}

			result := CanResolve(tc.container, tc.instanceName)

			// then
			assert.Equal(t, tc.wantResult, result)
		})
	}
}

func TestCanResolveType_NilContainer(t *testing.T) {
	// given

	// when
	require.PanicsWithValue(t, "nil container", func() {
		CanResolveType[any](nil)
	})

	// then
}

func TestCanResolveType(t *testing.T) {
	testCases := []struct {
		name         string
		container    *AnyMockContainer
		preCondition func(container *AnyMockContainer)

		wantResult bool
	}{
		{
			name:      "resolvable component type",
			container: &AnyMockContainer{},
			preCondition: func(container *AnyMockContainer) {
				container.On("CanResolveType", reflect.TypeFor[*AnyPointerComponent]()).Return(true)
			},
			wantResult: true,
		},
		{
			name:      "non-resolvable component type",
			container: &AnyMockContainer{},
			preCondition: func(container *AnyMockContainer) {
				container.On("CanResolveType", reflect.TypeFor[*AnyPointerComponent]()).Return(false)
			},
			wantResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			if tc.preCondition != nil {
				tc.preCondition(tc.container)
			}

			// when
			result := CanResolveType[*AnyPointerComponent](tc.container)

			// then
			assert.Equal(t, tc.wantResult, result)
		})
	}
}

func TestResolve_NilContainer(t *testing.T) {
	// given

	// when
	instance, err := Resolve[any](context.Background(), nil, "anyInstance")

	// then
	assert.NotNil(t, err)
	assert.EqualError(t, err, "nil container")
	assert.Nil(t, instance)
}

func TestResolve_Error(t *testing.T) {
	// given
	anyComponentType := reflect.TypeFor[*AnyPointerComponent]()

	container := &AnyMockContainer{}
	container.On("ResolveAs", mock.AnythingOfType("context.backgroundCtx"), "anyInstance", anyComponentType).
		Return(nil, errors.New("resolve error"))

	// when
	instance, err := Resolve[*AnyPointerComponent](context.Background(), container, "anyInstance")

	// then
	assert.NotNil(t, err)
	assert.EqualError(t, err, "resolve error")
	assert.Nil(t, instance)
}

func TestResolve(t *testing.T) {
	// given
	anyComponentType := reflect.TypeFor[*AnyPointerComponent]()
	anyComponent := &AnyPointerComponent{}

	container := &AnyMockContainer{}
	container.On("ResolveAs", mock.AnythingOfType("context.backgroundCtx"), "anyInstance", anyComponentType).
		Return(anyComponent, nil)

	// when
	instance, err := Resolve[*AnyPointerComponent](context.Background(), container, "anyInstance")

	// then
	assert.Nil(t, err)
	assert.Equal(t, anyComponent, instance)
}

func TestResolveType_NilContainer(t *testing.T) {
	// given

	// when
	instance, err := ResolveType[any](context.Background(), nil)

	// then
	assert.NotNil(t, err)
	assert.EqualError(t, err, "nil container")
	assert.Nil(t, instance)
}

func TestResolveType_Error(t *testing.T) {
	// given
	anyComponentType := reflect.TypeFor[*AnyPointerComponent]()

	container := &AnyMockContainer{}
	container.On("ResolveType", mock.AnythingOfType("context.backgroundCtx"), anyComponentType).
		Return(nil, errors.New("resolve error"))

	// when
	instance, err := ResolveType[*AnyPointerComponent](context.Background(), container)

	// then
	assert.NotNil(t, err)
	assert.EqualError(t, err, "resolve error")
	assert.Nil(t, instance)
}

func TestResolveType(t *testing.T) {
	// given
	anyComponentType := reflect.TypeFor[*AnyPointerComponent]()
	anyComponent := &AnyPointerComponent{}

	container := &AnyMockContainer{}
	container.On("ResolveType", mock.AnythingOfType("context.backgroundCtx"), anyComponentType).
		Return(anyComponent, nil)

	// when
	instance, err := ResolveType[*AnyPointerComponent](context.Background(), container)

	// then
	assert.Nil(t, err)
	assert.Equal(t, anyComponent, instance)
}

func TestResolveAll_NilContainer(t *testing.T) {
	// given

	// when
	instance, err := ResolveAll[any](context.Background(), nil)

	// then
	assert.NotNil(t, err)
	assert.EqualError(t, err, "nil container")
	assert.Nil(t, instance)
}

func TestResolveAll_Error(t *testing.T) {
	// given
	anyInterfaceType := reflect.TypeFor[AnyComponent]()

	container := &AnyMockContainer{}
	container.On("ResolveAll", mock.AnythingOfType("context.backgroundCtx"), anyInterfaceType).
		Return(nil, errors.New("resolve error"))

	// when
	instances, err := ResolveAll[AnyComponent](context.Background(), container)

	// then
	assert.NotNil(t, err)
	assert.EqualError(t, err, "resolve error")
	assert.Nil(t, instances)
}

func TestResolveAll(t *testing.T) {
	// given
	anyInterfaceType := reflect.TypeFor[AnyComponent]()
	anyComponent := &AnyPointerComponent{}
	anotherComponent := &AnyDependentComponent{}

	container := &AnyMockContainer{}
	container.On("ResolveAll", mock.AnythingOfType("context.backgroundCtx"), anyInterfaceType).
		Return([]any{anyComponent, anotherComponent}, nil)

	// when
	instances, err := ResolveAll[AnyComponent](context.Background(), container)

	// then
	assert.Nil(t, err)
	assert.Equal(t, []AnyComponent{anyComponent, anotherComponent}, instances)
}
