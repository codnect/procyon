package component

import (
	"context"
	"reflect"
)

// Resolver defines methods for resolving component instances
type Resolver interface {
	// CanResolve checks if a component with the given name is resolvable.
	CanResolve(name string) bool

	// Resolve retrieves an instance of the specified type
	Resolve(ctx context.Context, typ reflect.Type) (any, error)

	// ResolveAll retrieves all instances assignable to the specified type.
	ResolveAll(ctx context.Context, typ reflect.Type) ([]any, error)

	// ResolveNamed retrieves a component instance by its name.
	ResolveNamed(ctx context.Context, name string) (any, error)

	// ResolveNamedType retrieves a component by both name and expected type.
	ResolveNamedType(ctx context.Context, typ reflect.Type, name string) (any, error)
}
