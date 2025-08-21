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

import "errors"

var (
	// ErrDefinitionAlreadyExists is an error that occurs when a definition already exists.
	ErrDefinitionAlreadyExists = errors.New("definition already exists")
	// ErrDefinitionNotFound is an error that occurs when a definition is not found.
	ErrDefinitionNotFound = errors.New("definition not found")
	// ErrInvalidScopeName is an error that occurs when an invalid scope name is provided.
	ErrInvalidScopeName = errors.New("invalid scope name")
	// ErrInstanceAlreadyExists is an error that occurs when a singleton already exists.
	ErrInstanceAlreadyExists = errors.New("instance already exists")
	// ErrInstanceNotFound is an error that occurs when a singleton is not found.
	ErrInstanceNotFound = errors.New("instance not found")
	// ErrInstanceInPreparation is an error that occurs when an instance is in preparation.
	ErrInstanceInPreparation = errors.New("instance is in preparation, maybe it has got circular dependency cycle")
	// ErrScopeNotFound is an error that occurs when a scope is not found.
	ErrScopeNotFound = errors.New("scope not found")
	// ErrScopeReplacementNotAllowed is an error that occurs when a scope replacement is attempted for singleton and prototype scopes.
	ErrScopeReplacementNotAllowed = errors.New("scope replacement is not allowed for singleton and prototype scopes")
)
