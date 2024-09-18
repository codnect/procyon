package runtime

import (
	"codnect.io/procyon/component/container"
	"codnect.io/procyon/runtime/event"
	"context"
)

// Context interface represents the application context.
// It extends the standard context.Context and event.Publisher interfaces.
// It provides methods for starting and stopping the application context,
// checking if the application context is running, adding event listeners,
// and accessing the environment and container.
type Context interface {
	context.Context
	event.Publisher

	// Start method starts the application context.
	Start() error
	// Stop method stops the application context.
	Stop() error
	// IsRunning method checks if the application context is running.
	IsRunning() bool
	// AddEventListeners method adds event listeners to the application
	AddEventListeners(listeners ...event.Listener) error
	// Environment method returns the environment of the application context.
	Environment() Environment
	// Container method returns the container of the application context.
	Container() container.Container
}

// ContextConfigurer interface provides a method for configuring the application context.
type ContextConfigurer interface {
	// ConfigureContext method configures the application context.
	ConfigureContext(ctx Context) error
}
