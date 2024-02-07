package procyon

import (
	"codnect.io/procyon-core/container"
	"codnect.io/procyon-core/runtime"
	"codnect.io/procyon-core/runtime/env"
	"codnect.io/procyon-core/runtime/event"
	"codnect.io/reflector"
	"context"
	"os/signal"
	"syscall"
	"time"
)

type Context struct {
	signalCtx  context.Context
	cancelFunc context.CancelFunc

	environment env.Environment
	container   container.Container
	listeners   []*event.Listener

	//lifecycleProcessor *lifecycleProcessor
	err    error
	values map[any]any
}

func newContext(container container.Container) *Context {

	signalCtx, cancelFunc := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	/*defer stopFunc()
	<-notifyCtx.Done()
	_ = a.ctx.Stop()
	*/

	return &Context{
		signalCtx:  signalCtx,
		cancelFunc: cancelFunc,
		container:  container,
		listeners:  make([]*event.Listener, 0),
		values:     map[any]any{},
	}
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *Context) Done() <-chan struct{} {
	return c.signalCtx.Done()
}

func (c *Context) Err() error {
	return c.err
}

func (c *Context) Value(key any) any {
	return c.values[key]
}

func (c *Context) RegisterListener(listener *event.Listener) {
	//c.broadcaster.RegisterListener(listener)
}

func (c *Context) Listeners() []*event.Listener {
	listeners := make([]*event.Listener, len(c.listeners))
	copy(listeners, c.listeners)
	return listeners
}

func (c *Context) PublishEvent(ctx context.Context, event event.Event) error {
	//c.broadcaster.BroadcastEvent(ctx, event)
	return nil
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

	err = sharedInstances.Register("Environment", c.environment)
	if err != nil {
		return err
	}

	return nil
}

func (c *Context) start() error {
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

func (c *Context) Close() {
	defer c.cancelFunc()
	/*if c.lifecycleProcessor == nil {
		return nil
	}

	err := c.lifecycleProcessor.stop(c)

	if err != nil {
		return err
	}
	*/
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

func (c *Context) customize() error {
	customizers, err := getComponentsByType[runtime.ContextCustomizer](c.container)
	if err != nil {
		return err
	}

	for _, customizer := range customizers {
		err = customizer.CustomizeContext(c)

		if err != nil {
			return err
		}
	}

	return nil
}
