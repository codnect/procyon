package container

import (
	"context"
	"sync"
)

// ctxObjectCreationState is an empty struct used as a key for context values.
type ctxObjectCreationState struct{}

// ctxObjectCreationStateContextKey is a key for storing and retrieving a objectCreationState from a context.
var ctxObjectCreationStateContextKey = &ctxObjectCreationState{}

// objectCreationState struct store the creation states of objects, preventing circular dependencies.
type objectCreationState struct {
	currentlyInCreation map[string]struct{}
	mu                  sync.RWMutex
}

// objectCreationStateFromContext retrieves a objectCreationState from a context
func objectCreationStateFromContext(ctx context.Context) *objectCreationState {
	return ctx.Value(ctxObjectCreationStateContextKey).(*objectCreationState)
}

// withObjectCreationState adds a objectCreationState to a context if it doesn't already exist.
// This function checks if an objectCreationState already exists in the context, if not, it creates a new one
// and stores it in the context.
func withObjectCreationState(parent context.Context) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	val := parent.Value(ctxObjectCreationStateContextKey)

	if val != nil {
		return parent
	}

	manager := &objectCreationState{
		currentlyInCreation: map[string]struct{}{},
	}

	return context.WithValue(parent, ctxObjectCreationStateContextKey, manager)
}

// putToPreparation is called before an object is created. It checks if the object is already being created and
// returns an error if it is.
func (h *objectCreationState) putToPreparation(name string) error {
	defer h.mu.Unlock()
	h.mu.Lock()

	if _, ok := h.currentlyInCreation[name]; ok {
		return ErrObjectInPreparation
	}

	h.currentlyInCreation[name] = struct{}{}
	return nil
}

// removeFromPreparation is called after an object is created. It removes the object from the map of objects currently
// being created.
func (h *objectCreationState) removeFromPreparation(name string) {
	defer h.mu.Unlock()
	h.mu.Lock()
	delete(h.currentlyInCreation, name)
}
