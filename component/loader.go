package component

import "context"

// Loader defines the interface for loading components into a container.
// It evaluates conditions and registers components if applicable.
type Loader interface {
	// LoadComponents processes the given components and loads them based on the current context.
	// Returns an error if any step of the loading process fails.
	LoadComponents(ctx context.Context, components []*Component) error
}
