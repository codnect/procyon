package procyon

import (
	"codnect.io/procyon-core/container"
	"codnect.io/procyon-core/event"
	"codnect.io/procyon-core/runtime"
	"codnect.io/procyon-core/runtime/env"
	"codnect.io/reflector"
	"context"
	"time"
)

type Context struct {
	environment env.Environment
	container   container.Container
	broadcaster event.Broadcaster
	listeners   []*event.Listener

	//lifecycleProcessor *lifecycleProcessor
	values map[any]any
}

func newContext(container container.Container, broadcaster event.Broadcaster) *Context {
	return &Context{
		container:   container,
		broadcaster: broadcaster,
		listeners:   make([]*event.Listener, 0),
		values:      map[any]any{},
	}
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *Context) Done() <-chan struct{} {
	return nil
}

func (c *Context) Err() error {
	return nil
}

func (c *Context) Value(key any) any {
	return c.values[key]
}

func (c *Context) RegisterListener(listener *event.Listener) {
	c.broadcaster.RegisterListener(listener)
}

func (c *Context) Listeners() []*event.Listener {
	listeners := make([]*event.Listener, len(c.listeners))
	copy(listeners, c.listeners)
	return listeners
}

func (c *Context) PublishEvent(ctx context.Context, event event.Event) {
	c.broadcaster.BroadcastEvent(ctx, event)
}

func (c *Context) ApplicationName() string {
	return ""
}

func (c *Context) DisplayName() string {
	return ""
}

func (c *Context) StartupTime() time.Time {
	return time.Time{}
}

func (c *Context) Environment() env.Environment {
	return nil
}

func (c *Context) Container() container.Container {
	return c.container
}

func (c *Context) prepareRefresh() error {
	return nil
}

func (c *Context) prepareContainer() error {
	sharedInstances := c.container.SharedInstances()

	err := c.container.RegisterResolvable(reflector.TypeOf[container.Container](), c.container)
	if err != nil {
		return err
	}

	err = c.container.RegisterResolvable(reflector.TypeOf[runtime.Context](), c)
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

func (c *Context) Start() error {
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

func (c *Context) Stop() error {
	/*if c.lifecycleProcessor == nil {
		return nil
	}

	err := c.lifecycleProcessor.stop(c)

	if err != nil {
		return err
	}
	*/
	return nil
}

func (c *Context) registerComponentDefinitions() error {
	return nil
}

func (c *Context) initializeSharedComponents() error {
	return nil
}

func (c *Context) finalize() (err error) {
	/*c.lifecycleProcessor = defaultLifecycleProcessor(LifecycleProperties{}, c.container)

	err = c.lifecycleProcessor.start(c)

	if err != nil {
		return err
	}
	*/
	return nil
}

func (c *Context) setEnvironment(environment env.Environment) {
	c.environment = environment
}
