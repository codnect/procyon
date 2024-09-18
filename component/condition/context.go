package condition

import (
	"codnect.io/procyon/component/container"
	"context"
	"time"
)

// Context struct wraps the standard context.Context and a container.Container.
// It is used when evaluating conditions.
type Context struct {
	ctx       context.Context
	container container.Container
}

// NewContext function creates a new Context.
// It takes a standard context and a container as parameters.
// If either of these are nil, it panics.
func NewContext(ctx context.Context, container container.Container) Context {
	if ctx == nil {
		panic("nil context")
	}

	if container == nil {
		panic("nil container")
	}

	return Context{
		ctx:       ctx,
		container: container,
	}
}

// Deadline method returns the time when work done on behalf of
// this context should be canceled.
func (c Context) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

// Done method returns a channel that's closed when work done on behalf of
// this context should be canceled.
func (c Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

// Err method returns a non-nil error value after Done is closed.
func (c Context) Err() error {
	return c.ctx.Err()
}

// Value method returns the value associated with this context for key,
// or nil if no value is associated with key.
func (c Context) Value(key any) any {
	return c.ctx.Value(key)
}

// Container method returns the container associated with this context.
func (c Context) Container() container.Container {
	return c.container
}
