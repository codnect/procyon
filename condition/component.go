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

package condition

import (
	"context"
	"reflect"

	"go.codnect.io/procyon/component"
)

// Component is a condition that matches when a component with the specified
// type or name can be resolved from the container.
type Component struct {
	typ  reflect.Type
	name string
}

// OnComponent creates a condition that matches when a component of the
// specified type can be resolved from the container.
func OnComponent[T any]() *Component {
	return &Component{
		typ: reflect.TypeFor[T](),
	}
}

// OnComponentName creates a condition that matches when a component with the
// specified name can be resolved from the container.
func OnComponentName(name string) *Component {
	return &Component{
		name: name,
	}
}

// Matches reports whether the configured component can be resolved.
func (c *Component) Matches(ctx context.Context, container component.Container) bool {
	if c.typ != nil {
		_, err := container.ResolveType(ctx, c.typ)
		if err != nil {
			return false
		}

		return true
	}

	if c.name != "" {
		_, err := container.Resolve(ctx, c.name)
		if err != nil {
			return false
		}

		return true
	}

	return false
}

// MissingComponent is a condition that matches when a component with the
// specified type or name cannot be resolved from the container.
type MissingComponent struct {
	typ  reflect.Type
	name string
}

// OnMissingComponent creates a condition that matches when a component of the
// specified type cannot be resolved from the container.
func OnMissingComponent[T any]() *MissingComponent {
	return &MissingComponent{
		typ: reflect.TypeFor[T](),
	}
}

// OnMissingComponentName creates a condition that matches when a component
// with the specified name cannot be resolved from the container.
func OnMissingComponentName(name string) *MissingComponent {
	return &MissingComponent{
		name: name,
	}
}

// Matches reports whether the configured component cannot be resolved.
func (m *MissingComponent) Matches(ctx context.Context, container component.Container) bool {
	if m.typ != nil {
		_, err := container.ResolveType(ctx, m.typ)
		if err != nil {
			return true
		}

		return false
	}

	if m.name != "" {
		_, err := container.Resolve(ctx, m.name)
		if err != nil {
			return true
		}

		return false
	}

	return false
}
