package component

import (
	"context"
	"sync"
)

// ctxCreationState is an empty type used as a context key.
type ctxCreationState struct{}

// ctxCreationStateContextKey is the key used to store creationState in a context.
var ctxCreationStateContextKey = &ctxCreationState{}

// creationState keeps track of which component names are currently being created.
// It helps detect circular dependencies.
type creationState struct {
	currentlyInCreation map[string]struct{}
	mu                  sync.RWMutex
}

// newCreationState creates and returns a new creationState instance.
func newCreationState() *creationState {
	return &creationState{
		currentlyInCreation: make(map[string]struct{}),
	}
}

// withCreationState adds a new creationState to the given context.
// If the context already has one, it returns the same context.
func withCreationState(parent context.Context) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	val := parent.Value(ctxCreationStateContextKey)
	if val != nil {
		return parent
	}

	state := newCreationState()
	return context.WithValue(parent, ctxCreationStateContextKey, state)
}

// creationStateFromContext gets the creationState from the context.
// It assumes the context contains a valid value.
func creationStateFromContext(ctx context.Context) *creationState {
	return ctx.Value(ctxCreationStateContextKey).(*creationState)
}

// putToPreparation adds the given name to the in-creation map.
// Returns an error if it's already there (to prevent circular creation).
func (s *creationState) putToPreparation(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.currentlyInCreation[name]; ok {
		return ErrInstanceInPreparation
	}

	s.currentlyInCreation[name] = struct{}{}
	return nil
}

// removeFromPreparation removes the given name from the in-creation map.
func (s *creationState) removeFromPreparation(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.currentlyInCreation, name)
}
