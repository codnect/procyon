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
