package container

import (
	"context"
	"errors"
	"strings"
	"sync"
)

// SingletonScope and PrototypeScope are constants that represent the names of the singleton and prototype scopes.
const (
	SingletonScope = "singleton"
	PrototypeScope = "prototype"
)

// ScopeRegistry is an interface that defines methods for registering, finding, and listing the names of scopes.
type ScopeRegistry interface {
	// Register adds a new scope to the registry.
	Register(name string, scope Scope) error
	// Find retrieves a scope from the registry by its name.
	Find(name string) (Scope, error)
	// Names returns the names of all scopes in the registry.
	Names() []string
}

// Scope is an interface that defines methods for getting and removing objects from a scope.
type Scope interface {
	// GetObject retrieves an object from the scope by its name.
	// If the object does not exist, it is created using the provided provider function.
	GetObject(ctx context.Context, name string, provider ObjectProviderFunc) (any, error)
	// RemoveObject removes an object from the scope by its name.
	RemoveObject(ctx context.Context, name string) (any, error)
}

// simpleScopeRegistry struct represents a registry for managing scopes.
type simpleScopeRegistry struct {
	// scopes is a map that stores the registered scopes.
	// The key is the name of the scope and the value is the Scope instance.
	scopes map[string]Scope
	mu     sync.RWMutex
}

// newSimpleScopeRegistry creates a new simple scope registry.
func newSimpleScopeRegistry() *simpleScopeRegistry {
	return &simpleScopeRegistry{
		scopes: make(map[string]Scope),
		mu:     sync.RWMutex{},
	}
}

// Register adds a new scope to the registry.
func (r *simpleScopeRegistry) Register(name string, scope Scope) error {
	if strings.TrimSpace(name) == "" {
		return ErrInvalidScopeName
	}

	if scope == nil {
		return errors.New("nil scope")
	}

	if SingletonScope != name && PrototypeScope != name {
		defer r.mu.Unlock()
		r.mu.Lock()

		r.scopes[name] = scope
		return nil
	}

	return ErrScopeReplacementNotAllowed

}

// Names returns the names of all scopes in the registry.
func (r *simpleScopeRegistry) Names() []string {
	defer r.mu.Unlock()
	r.mu.Lock()

	scopeNames := make([]string, 0)
	for scopeName := range r.scopes {
		scopeNames = append(scopeNames, scopeName)
	}

	return scopeNames
}

// Find retrieves a scope from the registry by its name.
func (r *simpleScopeRegistry) Find(name string) (Scope, error) {
	defer r.mu.Unlock()
	r.mu.Lock()

	if scope, ok := r.scopes[name]; ok {
		return scope, nil
	}

	return nil, ErrScopeNotFound
}
