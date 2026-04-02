// Copyright 2026 Codnect
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package procyon

import (
	"context"
	"sync"
	"time"

	"codnect.io/procyon/component"
	"codnect.io/procyon/runtime"
)

const (
	// appContextContainerKey is the key used to register the application context itself in the component container.
	appContextContainerKey = "procyonAppContext"
	// lifecycleManagerContainerKey is the key used to register the lifecycle manager in the component container.
	lifecycleManagerContainerKey = "procyonLifecycleManager"
)

// Context struct represents the application context of Procyon.
type Context struct {
	running bool
	err     error
	mu      sync.RWMutex

	done      chan struct{}
	env       runtime.Environment
	container component.Container

	lifecycleManager runtime.LifecycleManager
}

// NewContext creates a new application context with the given environment.
// The returned context is not started yet. You need to call Start method to start the context.
func newContext(env runtime.Environment) *Context {
	return &Context{
		env:       env,
		container: component.NewDefaultContainer(),
	}
}

// Deadline method returns the time when work done on behalf of
// this context should be canceled.
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

// Done method returns a channel that's closed when work done on behalf of
// this context should be canceled.
func (c *Context) Done() <-chan struct{} {
	return c.done
}

func (c *Context) Err() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.err
}

func (c *Context) Value(key any) any {
	return nil
}

// Start starts the application context. It loads component definitions, registers the context itself
// as a singleton, and initializes singleton components. After this method is called, the context is
// considered running.
func (c *Context) Start(ctx context.Context) error {

	c.running = true
	return nil
}

// Stop stops the application context. It cancels the parent context and marks the context as not running.
func (c *Context) Stop(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.running = false
	c.err = context.Canceled
	close(c.done)
	return nil
}

// IsRunning returns true if the application context is currently running, false otherwise.
func (c *Context) IsRunning() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.running
}

func (c *Context) Refresh(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.loadComponentDefinitions(ctx)
	if err != nil {
		return err
	}

	err = c.container.RegisterSingleton(appContextContainerKey, c)
	if err != nil {
		return err
	}

	err = c.initializeSingletons(ctx)
	if err != nil {
		return err
	}

	c.lifecycleManager = newDefaultLifecycleManager(c.container)
	err = c.container.RegisterSingleton(lifecycleManagerContainerKey, c.lifecycleManager)
	if err != nil {
		return err
	}

	return c.lifecycleManager.Startup(ctx)
}

func (c *Context) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lifecycleManager != nil && c.lifecycleManager.IsRunning() {
		if err := c.lifecycleManager.Shutdown(ctx); err != nil {
			return err
		}
	}

	c.running = false
	return nil
}

func (c *Context) Environment() runtime.Environment {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.env
}

// Container returns the component container associated with this application context.
func (c *Context) Container() component.Container {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.container
}

// loadComponentDefinitions loads component definitions into the container using a ConditionalLoader.
// It retrieves the list of component definitions and loads them into the container, allowing for conditional
// loading based on the context.
func (c *Context) loadComponentDefinitions(ctx context.Context) error {
	loader := component.NewConditionalLoader(c.container, component.List())
	return loader.Load(ctx)
}

// initializeSingletons initializes all singleton components defined in the container. It iterates through
// the component definitions, checks for singleton definitions, and resolves them to ensure they are initialized
// and ready for use.
func (c *Context) initializeSingletons(ctx context.Context) error {
	for _, definition := range c.container.Definitions() {
		if !definition.IsSingleton() {
			continue
		}

		_, err := c.container.Resolve(ctx, definition.Name())

		if err != nil {
			return err
		}
	}

	return nil
}
