package runtime

import (
	"codnect.io/procyon/runtime/property"
	"context"
	"time"
)

// Lifecycle interface provides the methods for start/stop application lifecycle control.
type Lifecycle interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	IsRunning() bool
}

// LifecycleProperties struct represents the properties of application lifecycle.
type LifecycleProperties struct {
	property.Properties `prefix:"procyon.lifecycle"`

	ShutdownTimeout time.Duration `prop:"shutdown-timeout" default:"30000"`
}

// NewLifecycleProperties function creates a new LifecycleProperties.
func NewLifecycleProperties() *LifecycleProperties {
	return &LifecycleProperties{}
}
