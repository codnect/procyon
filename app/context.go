package app

import (
	"context"
	"procyon-test/container"
	"procyon-test/env"
	"procyon-test/event"
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
}

func newContext() *appContext {
	return &appContext{}
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

func (c *appContext) Publish(ctx context.Context, event event.Event) {

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
