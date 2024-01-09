package condition

import (
	"codnect.io/procyon/core/container"
	"codnect.io/procyon/core/env"
)

type Context interface {
	DefinitionRegistry() container.DefinitionRegistry
	Container() container.Container
	Environment() env.Environment
}

type context struct {
	container   container.Container
	registry    container.DefinitionRegistry
	environment env.Environment
}

func newContext(container container.Container, environment env.Environment) Context {
	return &context{
		container:   container,
		registry:    container.DefinitionRegistry(),
		environment: environment,
	}
}

func (c *context) DefinitionRegistry() container.DefinitionRegistry {
	return c.registry
}

func (c *context) Container() container.Container {
	return c.container
}

func (c *context) Environment() env.Environment {
	return c.environment
}
