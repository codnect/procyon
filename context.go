package procyon

import (
	"codnect.io/procyon-core/component"
	"codnect.io/procyon-core/component/filter"
	"codnect.io/procyon-core/runtime"
	"codnect.io/procyon-core/runtime/event"
	"context"
	"sync"
	"time"
)

type Context struct {
	parent     context.Context
	cancelFunc context.CancelFunc

	args             *runtime.Arguments
	environment      runtime.Environment
	container        component.Container
	eventMulticaster event.Multicaster

	running bool
	mu      sync.RWMutex
}

func createContext(args *runtime.Arguments) *Context {
	parent, cancel := context.WithCancel(context.Background())

	return &Context{
		parent:      parent,
		cancelFunc:  cancel,
		args:        args,
		environment: prepareEnvironment(args),
		container:   prepareContainer(args),
	}
}

func prepareContainer(args *runtime.Arguments) component.Container {
	container := component.NewObjectContainer()
	_ = container.Singletons().Register("procyonApplicationArgs", args)
	return container
}

func prepareEnvironment(args *runtime.Arguments) runtime.Environment {
	environment := runtime.NewDefaultEnvironment()

	propertySources := environment.PropertySources()
	propertySources.AddLast(runtime.NewArgumentsSource(args))
	propertySources.AddLast(runtime.NewEnvironmentSource())

	return environment
}

func (c *Context) Start() error {
	defer c.mu.Unlock()
	c.mu.Lock()

	err := c.loadComponentDefinitions()
	if err != nil {
		return err
	}

	err = c.container.Singletons().Register("procyonApplicationContext", c)
	if err != nil {
		return err
	}

	err = c.initializeEventMulticaster()
	if err != nil {
		return err
	}

	err = c.initializeSingletons()
	if err != nil {
		return err
	}

	err = c.startLifecycleObjects()
	if err != nil {
		return err
	}

	c.running = true
	return c.PublishEvent(c, runtime.NewStartupEvent(c))
}

func (c *Context) Stop() error {
	defer c.mu.Unlock()
	c.mu.Lock()

	err := c.PublishEvent(c, runtime.NewShutdownEvent(c))
	if err != nil {
		// log the error
	}

	/*
		shutdown gracefully

		timeoutCtx, timeoutCancelFunc := context.WithTimeout(context.Background(), time.Second*3)
		defer timeoutCancelFunc()

		go func() {
			err = c.stopLifecycleObjects(timeoutCtx)
		}()

		<-timeoutCtx.Done()
	*/
	err = c.stopLifecycleObjects()

	if err != nil {
		// log the error
	}

	c.running = false
	c.cancelFunc()
	return nil
}

func (c *Context) IsRunning() bool {
	defer c.mu.Unlock()
	c.mu.Lock()
	return c.running
}

func (c *Context) AddEventListeners(listeners ...event.Listener) error {
	for _, listener := range listeners {
		err := c.eventMulticaster.AddEventListener(listener)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Context) Deadline() (time.Time, bool) {
	return c.parent.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.parent.Done()
}

func (c *Context) Err() error {
	return c.parent.Err()
}

func (c *Context) Value(key any) any {
	return c.parent.Value(key)
}

func (c *Context) Environment() runtime.Environment {
	return c.environment
}

func (c *Context) Container() component.Container {
	return c.container
}

func (c *Context) loadComponentDefinitions() error {
	loader := component.NewDefinitionLoader(c.container)
	return loader.LoadDefinitions(c, component.List())
}

func (c *Context) initializeEventMulticaster() error {
	multicaster, err := c.container.GetObject(c, filter.ByTypeOf[event.Multicaster]())

	if err != nil {
		return err
	}

	c.eventMulticaster = multicaster.(event.Multicaster)
	return nil
}

func (c *Context) initializeSingletons() error {

	for _, definition := range c.container.Definitions().List() {
		_, err := c.container.GetObject(c, filter.ByName(definition.Name()))

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Context) PublishEvent(ctx context.Context, event event.ApplicationEvent) error {
	return c.eventMulticaster.MulticastEvent(ctx, event)
}

func (c *Context) PublishEventAsync(ctx context.Context, event event.ApplicationEvent) error {
	return c.eventMulticaster.MulticastEventAsync(ctx, event)
}

func (c *Context) startLifecycleObjects() error {
	lifecycleObjects := c.container.ListObjects(c, filter.ByTypeOf[runtime.Lifecycle]())
	for _, object := range lifecycleObjects {
		lifecycle := object.(runtime.Lifecycle)

		err := lifecycle.Start(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Context) stopLifecycleObjects() error {
	lifecycleObjects := c.container.ListObjects(c, filter.ByTypeOf[runtime.Lifecycle]())
	for _, object := range lifecycleObjects {
		lifecycle := object.(runtime.Lifecycle)

		err := lifecycle.Stop(c)
		if err != nil {
			return err
		}
	}

	return nil
}
