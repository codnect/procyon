package runtime

import "time"

// StartupEvent struct represents an event that occurs when the application starts up.
type StartupEvent struct {
	ctx  Context   // The application context.
	time time.Time // The time when the event occurred.
}

// NewStartupEvent function creates a new StartupEvent.
func NewStartupEvent(ctx Context) StartupEvent {
	return StartupEvent{
		ctx:  ctx,
		time: time.Now(),
	}
}

// Context method returns the context of the StartupEvent.
func (s StartupEvent) Context() Context {
	return s.ctx
}

// EventSource method returns the source of the event, which is the context.
func (s StartupEvent) EventSource() any {
	return s.ctx
}

// EventTime method returns the time when the event occurred.
func (s StartupEvent) EventTime() time.Time {
	return s.time
}

// ShutdownEvent struct represents an event that occurs when the application shuts down.
type ShutdownEvent struct {
	ctx  Context
	time time.Time
}

// NewShutdownEvent function creates a new ShutdownEvent.
func NewShutdownEvent(ctx Context) ShutdownEvent {
	return ShutdownEvent{
		ctx:  ctx,
		time: time.Now(),
	}
}

// Context method returns the context of the ShutdownEvent.
func (s ShutdownEvent) Context() Context {
	return s.ctx
}

// EventSource method returns the source of the event, which is the context.
func (s ShutdownEvent) EventSource() any {
	return s.ctx
}

// EventTime method returns the time when the event occurred.
func (s ShutdownEvent) EventTime() time.Time {
	return s.time
}
