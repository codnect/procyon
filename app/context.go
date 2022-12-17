package app

import (
	"context"
	"github.com/procyon-projects/procyon/container"
	"github.com/procyon-projects/procyon/env"
	"github.com/procyon-projects/procyon/event"
	"time"
)

type Context interface {
	context.Context
	event.Publisher

	ApplicationName() string
	DisplayName() string
	StartupTime() time.Time
	Environment() env.Environment
	Container() *container.Container
	Refresh()
}

type appContext struct {
	environment env.Environment
	container   *container.Container
}

func newContext(container *container.Container) *appContext {
	return &appContext{
		container: container,
	}
}

func (c *appContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *appContext) Done() <-chan struct{} {
	return nil
}

func (c *appContext) Err() error {
	return nil
}

func (c *appContext) Value(key any) any {
	return nil
}

func (c *appContext) PublishEvent(ctx context.Context, event event.Event) {

}

func (c *appContext) ApplicationName() string {
	return ""
}

func (c *appContext) DisplayName() string {
	return ""
}

func (c *appContext) StartupTime() time.Time {
	return time.Time{}
}

func (c *appContext) Environment() env.Environment {
	return nil
}

func (c *appContext) Container() *container.Container {
	return nil
}

func (c *appContext) Refresh() {

}

func (c *appContext) setEnvironment(environment env.Environment) {
	c.environment = environment
}
