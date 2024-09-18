package event

import "time"

// ApplicationEvent interface represents an application event.
type ApplicationEvent interface {
	// EventSource method returns the source of the event.
	// The source is typically the object that the event applies to.
	EventSource() any
	// EventTime method returns the time when the event occurred.
	EventTime() time.Time
}
