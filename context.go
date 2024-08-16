package procyon

import (
	"codnect.io/procyon-core/component"
	"codnect.io/procyon-core/runtime"
	"context"
	"sync"
	"time"
)

type Context struct {
	parent     context.Context
	cancelFunc context.CancelFunc

	args        *runtime.Arguments
	environment runtime.Environment
	container   component.Container

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

	loader := newDefinitionLoader(c.container)
	err := loader.load(c)
	if err != nil {
		return err
	}

	err = initializeSingletons(c, c.container)
	if err != nil {
		return err
	}

	return nil
}

func (c *Context) Stop() error {
	defer c.mu.Unlock()
	c.mu.Lock()

	c.running = false
	c.cancelFunc()
	return nil
}

func (c *Context) IsRunning() bool {
	defer c.mu.Unlock()
	c.mu.Lock()
	return c.running
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
