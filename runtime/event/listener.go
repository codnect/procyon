package event

import (
	"context"
	"errors"
)

// Listener is an interface that represents an event listener.
type Listener interface {
	// OnEvent method handles an event.
	OnEvent(ctx context.Context, event ApplicationEvent) error
	// SupportsEvent method checks if an event is supported.
	SupportsEvent(event ApplicationEvent) bool
}

// Listen function creates a new Listener with the provided event handler.
func Listen[E ApplicationEvent](handler func(ctx context.Context, event E) error) Listener {
	return listenerAdapter[E]{
		handler: handler,
	}
}

// listenerAdapter struct is an adapter that implements the Listener interface.
type listenerAdapter[E any] struct {
	handler func(ctx context.Context, event E) error
}

// OnEvent method handles an event.
// It returns an error if the event is not supported.
func (a listenerAdapter[E]) OnEvent(ctx context.Context, event ApplicationEvent) error {
	if !a.SupportsEvent(event) {
		return errors.New("")
	}

	return a.handler(ctx, event.(E))
}

// SupportsEvent method checks if an event is supported.
// It returns true if the event is supported, false otherwise.
func (a listenerAdapter[E]) SupportsEvent(event ApplicationEvent) bool {
	if _, ok := event.(E); ok {
		return true
	}

	return false
}
