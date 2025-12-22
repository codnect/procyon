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
)

// Resolver defines methods for resolving component instances.
type Resolver interface {
	// CanResolve checks if a component with the given name is resolvable.
	CanResolve(name string) bool

	// CanResolveType checks if a component of the given type is resolvable.
	CanResolveType(typ reflect.Type) bool

	// Resolve retrieves a component instance by its name.
	Resolve(ctx context.Context, name string) (any, error)

	// ResolveType retrieves an instance of the specified type.
	ResolveType(ctx context.Context, typ reflect.Type) (any, error)

	// ResolveAs retrieves a component by both name and expected type.
	ResolveAs(ctx context.Context, name string, typ reflect.Type) (any, error)

	// ResolveAll retrieves all instances assignable to the specified type.
	ResolveAll(ctx context.Context, typ reflect.Type) ([]any, error)
}

// Resolve retrieves a component instance of type T from the container by its name.
// It returns an error if the container is nil or if resolution fails.
func Resolve[T any](ctx context.Context, container Container, name string) (T, error) {
	var zeroVal T

	if container == nil {
		return zeroVal, errors.New("nil container")
	}

	instance, err := container.ResolveAs(ctx, name, reflect.TypeFor[T]())
	if err != nil {
		return zeroVal, err
	}

	return instance.(T), nil
}

// ResolveType retrieves a component instance of type T from the container by its type.
// It returns an error if the container is nil or if resolution fails.
func ResolveType[T any](ctx context.Context, container Container) (T, error) {
	var zeroVal T

	if container == nil {
		return zeroVal, errors.New("nil container")
	}

	instance, err := container.ResolveType(ctx, reflect.TypeFor[T]())
	if err != nil {
		return zeroVal, err
	}

	return instance.(T), nil
}

// ResolveAll retrieves all component instances of type T from the container.
// It returns an error if the container is nil or if resolution fails.
func ResolveAll[T any](ctx context.Context, container Container) ([]T, error) {
	if container == nil {
		return nil, errors.New("nil container")
	}

	instances, err := container.ResolveAll(ctx, reflect.TypeFor[T]())
	if err != nil {
		return nil, err
	}

	result := make([]T, len(instances))
	for index, instance := range instances {
		result[index] = instance.(T)
	}

	return result, nil
}
