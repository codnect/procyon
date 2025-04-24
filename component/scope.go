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

import "context"

// SingletonScope and PrototypeScope are constants that represent the names of the singleton and prototype scopes.
const (
	SingletonScope = "singleton"
	PrototypeScope = "prototype"
)

// Scope defines methods for resolving and removing instances within a particular scope.
type Scope interface {
	// Resolve returns an instance associated with the given name.
	// If the instance does not exist, it uses the provided FactoryFunc to create one.
	Resolve(ctx context.Context, name string, fn FactoryFunc) (any, error)

	// Remove deletes the instance associated with the specified name from the scope.
	Remove(ctx context.Context, name string) error
}

// ScopeRegistry defines methods for managing scopes.
type ScopeRegistry interface {
	// RegisterScope adds a new scope with the specified name to the registry.
	RegisterScope(name string, scope Scope)

	// Scope retrieves the scope associated with the given name.
	// Returns the scope and a boolean indicating its existence.
	Scope(name string) (Scope, bool)
}
