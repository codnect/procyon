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

// Container provides a unified interface for managing components,
// including definition registration, instance resolving, custom scopes,
// lifecycle management, and manual bindings.
type Container interface {
	// DefinitionRegistry provides access to component definitions and their metadata.
	DefinitionRegistry

	// SingletonRegistry manages singleton instances of components.
	SingletonRegistry

	// Resolver resolves component instances by type or name.
	Resolver

	// ResolvableRegistry registers non-component dependencies for injection.
	ResolvableRegistry

	// ScopeRegistry manages custom scopes.
	ScopeRegistry

	// LifecycleManager manages lifecycle hooks.
	LifecycleManager
}
