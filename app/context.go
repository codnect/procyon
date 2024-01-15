package app

import (
	"codnect.io/procyon-core/container"
	"codnect.io/procyon-core/env"
	"codnect.io/procyon/app/event"
	"codnect.io/reflector"
	"context"
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
	Container() container.Container
	Start() error
}

type appContext struct {
	environment env.Environment
	container   container.Container
	broadcaster event.Broadcaster
	listeners   []*event.Listener
	values      map[any]any
}

func newContext(container container.Container, broadcaster event.Broadcaster) *appContext {
	return &appContext{
		container:   container,
		broadcaster: broadcaster,
		listeners:   make([]*event.Listener, 0),
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

func (c *appContext) Container() container.Container {
	return c.container
}

func (c *appContext) prepareRefresh() error {
	return nil
}

func (c *appContext) prepareContainer() error {
	sharedInstances := c.container.SharedInstances()

	err := c.container.RegisterResolvable(reflector.TypeOf[container.Container](), c.container)
	if err != nil {
		return err
	}

	err = c.container.RegisterResolvable(reflector.TypeOf[Context](), c)
	if err != nil {
		return err
	}

	err = c.container.RegisterResolvable(reflector.TypeOf[event.Publisher](), c)
	if err != nil {
		return err
	}

	err = sharedInstances.Register("environment", c.environment)
	if err != nil {
		return err
	}

	return nil
}

func (c *appContext) Start() error {
	err := c.prepareRefresh()
	if err != nil {
		return err
	}

	err = c.prepareContainer()
	if err != nil {
		return err
	}

	err = c.registerComponentDefinitions()
	if err != nil {
		return err
	}

	err = c.initializeSharedComponents()
	if err != nil {
		return err
	}

	return c.finalize()
}

func (c *appContext) registerComponentDefinitions() error {
	return nil
}

func (c *appContext) initializeSharedComponents() error {
	return nil
}

func (c *appContext) finalize() error {
	return nil
}

func (c *appContext) setEnvironment(environment env.Environment) {
	c.environment = environment
}
