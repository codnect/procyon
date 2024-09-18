package container

import (
	"errors"
)

var (
	// ErrNoFilterProvided is an error that occurs when no filter is provided but at least one is required.
	ErrNoFilterProvided = errors.New("no filter provided, at least one filter is required")
	// ErrObjectNotFound is an error that occurs when an object is not found.
	ErrObjectNotFound = errors.New("object not found")
	// ErrDefinitionNotFound is an error that occurs when a definition is not found.
	ErrDefinitionNotFound = errors.New("definition not found")
	// ErrObjectAlreadyExists is an error that occurs when an object already exists.
	ErrObjectAlreadyExists = errors.New("object already exists")
	// ErrDefinitionAlreadyExists is an error that occurs when a definition already exists.
	ErrDefinitionAlreadyExists = errors.New("definition already exists")
	// ErrMultipleObjectsFound is an error that occurs when multiple objects are found but only one was expected.
	ErrMultipleObjectsFound = errors.New("multiple objects found, expected only one")
	// ErrMultipleDefinitionsFound is an error that occurs when multiple definitions are found but only one was expected.
	ErrMultipleDefinitionsFound = errors.New("multiple definitions found, expected only one")
	// ErrObjectInPreparation is an error that occurs when an object is in preparation.
	ErrObjectInPreparation = errors.New("object is in preparation, maybe it has got circular dependency cycle")
	// ErrInvalidScopeName is an error that occurs when an invalid scope name is provided.
	ErrInvalidScopeName = errors.New("invalid scope name")
	// ErrScopeReplacementNotAllowed is an error that occurs when a scope replacement is attempted for singleton and prototype scopes.
	ErrScopeReplacementNotAllowed = errors.New("scope replacement is not allowed for singleton and prototype scopes")
	// ErrScopeNotFound is an error that occurs when a scope is not found.
	ErrScopeNotFound = errors.New("scope not found")
)
