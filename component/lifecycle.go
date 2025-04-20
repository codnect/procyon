package component

import (
	"context"
)

// Initializer can be implemented by components that require custom logic during initialization.
type Initializer interface {
	// Init is called once the component has been constructed and dependencies injected.
	Init(ctx context.Context) error
}

// PreProcessor defines logic that runs before a component is initialized.
type PreProcessor interface {
	// ProcessBeforeInit is called before the Init method.
	// Allows modification or validation of the instance before initialization.
	ProcessBeforeInit(ctx context.Context, instance any) (any, error)
}

// PostProcessor defines logic that runs after a component is initialized.
type PostProcessor interface {
	// ProcessAfterInit is called after the Init method.
	// Allows additional configuration or enhancement of the instance.
	ProcessAfterInit(ctx context.Context, instance any) (any, error)
}

// Finalizer can be implemented by components that need to release resources or clean up during shutdown.
type Finalizer interface {
	// Finalize is called during application shutdown for cleanup purposes.
	Finalize() error
}

// LifecycleManager manages the registration of lifecycle hooks such as pre/post processors.
type LifecycleManager interface {
	// UsePreProcessor registers a PreProcessor to be applied before initialization.
	UsePreProcessor(initializer PreProcessor) error

	// UsePostProcessor registers a PostProcessor to be applied after initialization.
	UsePostProcessor(initializer PostProcessor) error
}
