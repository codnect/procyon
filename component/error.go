package component

import "errors"

var (
	// ErrDefinitionAlreadyExists is an error that occurs when a definition already exists.
	ErrDefinitionAlreadyExists = errors.New("definition already exists")
	// ErrDefinitionNotFound is an error that occurs when a definition is not found.
	ErrDefinitionNotFound = errors.New("definition not found")
)
