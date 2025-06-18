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

// SingletonRegistry defines methods for managing singleton instances within the component system.
type SingletonRegistry interface {
	// RegisterSingleton registers a singleton instance with the given name.
	// Returns an error if a singleton with the same name already exists.
	RegisterSingleton(name string, instance any) error

	// ContainsSingleton checks whether a singleton with the specified name exists.
	ContainsSingleton(name string) bool

	// Singleton retrieves the singleton instance associated with the given name.
	// Returns the instance and a boolean indicating its existence.
	Singleton(name string) (any, bool)

	// RemoveSingleton removes the singleton instance associated with the specified name.
	RemoveSingleton(name string) error
}
