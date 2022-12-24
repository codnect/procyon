package condition

import (
	"github.com/procyon-projects/procyon/app/env"
	"github.com/procyon-projects/procyon/container"
)

type Context interface {
	DefinitionRegistry() container.DefinitionRegistry
	Container() *container.Container
	Environment() env.Environment
}
