package event

import "context"

// TypedPublisher is an interface that represents a typed event publisher.
type TypedPublisher[E ApplicationEvent] interface {
	// PublishEvent method publishes a specific type of event.
	PublishEvent(ctx context.Context, event E) error
}

// AsyncTypedPublisher is an interface that represents an asynchronous typed event publisher.
type AsyncTypedPublisher[E ApplicationEvent] interface {
	// PublishEventAsync method publishes a specific type of event asynchronously.
	PublishEventAsync(ctx context.Context, event E) error
}

// Publisher is an interface that represents an event publisher.
type Publisher interface {
	// PublishEvent method publishes an event.
	PublishEvent(ctx context.Context, event ApplicationEvent) error
	// PublishEventAsync method publishes an event asynchronously.
	PublishEventAsync(ctx context.Context, event ApplicationEvent) error
}
