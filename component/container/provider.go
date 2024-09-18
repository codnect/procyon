package container

import "context"

// ObjectProviderFunc is a type that represents a function that provides an object.
// This function takes a context as input and returns an object and an error.
// The object is the result of the function and the error is returned if the function fails.
type ObjectProviderFunc func(ctx context.Context) (any, error)
