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
	"reflect"
)

// Resolver defines methods for resolving component instances
type Resolver interface {
	// CanResolve checks if a component with the given name is resolvable.
	CanResolve(ctx context.Context, name string) bool

	// Resolve retrieves an instance of the specified type
	Resolve(ctx context.Context, typ reflect.Type) (any, error)

	// ResolveAll retrieves all instances assignable to the specified type.
	ResolveAll(ctx context.Context, typ reflect.Type) ([]any, error)

	// ResolveNamed retrieves a component instance by its name.
	ResolveNamed(ctx context.Context, name string) (any, error)

	// ResolveNamedType retrieves a component by both name and expected type.
	ResolveNamedType(ctx context.Context, typ reflect.Type, name string) (any, error)
}
