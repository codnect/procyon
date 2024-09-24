package web

import (
	"codnect.io/procyon/runtime/property"
	"context"
)

// Server interface represents a server that can be started and stopped.
type Server interface {
	// Start method starts the server.
	Start(ctx context.Context) error
	// Stop method stops the server.
	Stop(ctx context.Context) error
	// Port method returns the port the server is running on.
	Port() int
}

// ServerProperties struct represents the properties of a server.
type ServerProperties struct {
	property.Properties `prefix:"procyon.server"` // The prefix for server properties.

	ContextPath string `prop:"context-path"`        // The context path of the server.
	Port        int    `prop:"port" default:"8080"` // The port the server is running on.
}

// NewServerProperties function creates a new ServerProperties.
func NewServerProperties() *ServerProperties {
	return &ServerProperties{}
}
