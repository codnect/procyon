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
	event.ListenerRegistry

	ApplicationName() string
	DisplayName() string
	StartupTime() time.Time
	Environment() env.Environment
	Container() *container.Container
	Refresh() error
}

type appContext struct {
	environment env.Environment
	container   *container.Container
	broadcaster event.Broadcaster
	listeners   []*event.Listener
	customizers *contextCustomizers
	values      map[any]any
}

func newContext(container *container.Container, broadcaster event.Broadcaster) *appContext {
	return &appContext{
		container:   container,
		broadcaster: broadcaster,
		listeners:   make([]*event.Listener, 0),
		customizers: newContextCustomizers(make([]ContextCustomizer, 0)),
		values:      map[any]any{},
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
	return c.values[key]
}

func (c *appContext) RegisterListener(listener *event.Listener) {
	c.broadcaster.RegisterListener(listener)
}

func (c *appContext) Listeners() []*event.Listener {
	listeners := make([]*event.Listener, len(c.listeners))
	copy(listeners, c.listeners)
	return listeners
}

func (c *appContext) PublishEvent(ctx context.Context, event event.Event) {
	c.broadcaster.BroadcastEvent(ctx, event)
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
	return c.container
}

func (c *appContext) prepareRefresh() error {
	return nil
}

func (c *appContext) prepareContainer() error {
	sharedInstances := c.container.SharedInstances()

	err := c.container.RegisterResolvable(container.TypeOf[*container.Container](), c.container)
	if err != nil {
		return err
	}

	err = c.container.RegisterResolvable(container.TypeOf[Context](), c)
	if err != nil {
		return err
	}

	err = c.container.RegisterResolvable(container.TypeOf[event.Publisher](), c)
	if err != nil {
		return err
	}

	err = sharedInstances.Add("environment", c.environment)
	if err != nil {
		return err
	}

	return nil
}

func (c *appContext) Refresh() error {
	err := c.prepareRefresh()
	if err != nil {
		return err
	}

	err = c.prepareContainer()
	if err != nil {
		return err
	}

	err = c.customizers.invokeCustomizers(c)
	if err != nil {
		return err
	}

	return nil
}

func (c *appContext) setEnvironment(environment env.Environment) {
	c.environment = environment
}
