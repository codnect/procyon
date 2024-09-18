package container

import "context"

// ObjectProcessor interface provides methods for processing objects before and after initialization.
// These methods allow for custom logic to be executed before and after an object is initialized.
type ObjectProcessor interface {
	// ProcessBeforeInit is called before an object is initialized.
	// It takes a context and the object as input and returns the processed object and an error if the processing fails.
	ProcessBeforeInit(ctx context.Context, object any) (any, error)
	// ProcessAfterInit is called after an object is initialized.
	// It takes a context and the object as input and returns the processed object and an error if the processing fails.
	ProcessAfterInit(ctx context.Context, object any) (any, error)
}
