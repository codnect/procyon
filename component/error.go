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
