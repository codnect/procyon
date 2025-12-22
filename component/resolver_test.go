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
)

func TestResolveNilContainer(t *testing.T) {
	// given

	// when
	instance, err := Resolve[any](context.Background(), nil, "anyInstance")

	// then
	assert.NotNil(t, err)
	assert.EqualError(t, err, "nil container")
	assert.Nil(t, instance)
}

func TestResolveError(t *testing.T) {
	// given
	anyComponentType := reflect.TypeFor[*AnyComponent]()

	container := &AnyContainer{}
	container.On("ResolveAs", mock.AnythingOfType("context.backgroundCtx"), "anyInstance", anyComponentType).
		Return(nil, errors.New("resolve error"))

	// when
	instance, err := Resolve[*AnyComponent](context.Background(), container, "anyInstance")

	// then
	assert.NotNil(t, err)
	assert.EqualError(t, err, "resolve error")
	assert.Nil(t, instance)
}

func TestResolve(t *testing.T) {
	// given
	anyComponentType := reflect.TypeFor[*AnyComponent]()
	anyComponent := &AnyComponent{}

	container := &AnyContainer{}
	container.On("ResolveAs", mock.AnythingOfType("context.backgroundCtx"), "anyInstance", anyComponentType).
		Return(anyComponent, nil)

	// when
	instance, err := Resolve[*AnyComponent](context.Background(), container, "anyInstance")

	// then
	assert.Nil(t, err)
	assert.Equal(t, anyComponent, instance)
}

func TestResolveTypeNilContainer(t *testing.T) {
	// given

	// when
	instance, err := ResolveType[any](context.Background(), nil)

	// then
	assert.NotNil(t, err)
	assert.EqualError(t, err, "nil container")
	assert.Nil(t, instance)
}

func TestResolveTypeError(t *testing.T) {
	// given
	anyComponentType := reflect.TypeFor[*AnyComponent]()

	container := &AnyContainer{}
	container.On("ResolveType", mock.AnythingOfType("context.backgroundCtx"), anyComponentType).
		Return(nil, errors.New("resolve error"))

	// when
	instance, err := ResolveType[*AnyComponent](context.Background(), container)

	// then
	assert.NotNil(t, err)
	assert.EqualError(t, err, "resolve error")
	assert.Nil(t, instance)
}

func TestResolveType(t *testing.T) {
	// given
	anyComponentType := reflect.TypeFor[*AnyComponent]()
	anyComponent := &AnyComponent{}

	container := &AnyContainer{}
	container.On("ResolveType", mock.AnythingOfType("context.backgroundCtx"), anyComponentType).
		Return(anyComponent, nil)

	// when
	instance, err := ResolveType[*AnyComponent](context.Background(), container)

	// then
	assert.Nil(t, err)
	assert.Equal(t, anyComponent, instance)
}

func TestResolveAllNilContainer(t *testing.T) {
	// given

	// when
	instance, err := ResolveAll[any](context.Background(), nil)

	// then
	assert.NotNil(t, err)
	assert.EqualError(t, err, "nil container")
	assert.Nil(t, instance)
}

func TestResolveAllError(t *testing.T) {
	// given
	anyInterfaceType := reflect.TypeFor[AnyInterface]()

	container := &AnyContainer{}
	container.On("ResolveAll", mock.AnythingOfType("context.backgroundCtx"), anyInterfaceType).
		Return(nil, errors.New("resolve error"))

	// when
	instances, err := ResolveAll[AnyInterface](context.Background(), container)

	// then
	assert.NotNil(t, err)
	assert.EqualError(t, err, "resolve error")
	assert.Nil(t, instances)
}

func TestResolveAll(t *testing.T) {
	// given
	anyInterfaceType := reflect.TypeFor[AnyInterface]()
	anyComponent := &AnyComponent{}
	anotherComponent := &AnotherComponent{}

	container := &AnyContainer{}
	container.On("ResolveAll", mock.AnythingOfType("context.backgroundCtx"), anyInterfaceType).
		Return([]any{anyComponent, anotherComponent}, nil)

	// when
	instances, err := ResolveAll[AnyInterface](context.Background(), container)

	// then
	assert.Nil(t, err)
	assert.Equal(t, []AnyInterface{anyComponent, anotherComponent}, instances)
}
