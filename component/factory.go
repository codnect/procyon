package component

import "context"

// FactoryFunc is a type that represents a function that provides an instance.
type FactoryFunc func(ctx context.Context) (any, error)
