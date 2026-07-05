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
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"codnect.io/procyon/component"
	"codnect.io/procyon/io"
	"codnect.io/procyon/runtime"
)

const (
	// envContainerKey is the key used to register the runtime environment in the component container.
	envContainerKey = "environment"
	// lifecycleManagerContainerKey is the key used to register the lifecycle manager in the component container.
	lifecycleManagerContainerKey = "procyonLifecycleManager"
)

var (
	// bootstrapTypes is a list of types that are considered bootstrap components. These components are used during the
	// bootstrapping phase of the application and are not loaded into the main application context. This allows for
	// separation of concerns and ensures that certain components are only available during the bootstrapping process.
	bootstrapTypes = []reflect.Type{
		reflect.TypeFor[runtime.EnvironmentCustomizer](),
		reflect.TypeFor[runtime.ContextInitializer](),
	}
)

// contextError is a custom error type that wraps errors occurring during context operations (start, stop, refresh).
// It includes the operation being performed and the underlying error, providing more context for debugging and error
// handling.
type contextError struct {
	Op  string
	Err error
}

// Error method returns a string representation of the context error, including the operation and the underlying
// error message.
func (e *contextError) Error() string {
	return e.Op + " context: " + e.Err.Error()
}

// Unwrap method allows unwrapping the underlying error, enabling error chaining and compatibility with errors.Is and
// errors.As.
func (e *contextError) Unwrap() error {
	return e.Err
}

// Context struct represents the application context of Procyon.
type Context struct {
	done      chan struct{}
	err       error
	mu        sync.RWMutex
	refreshed bool

	resourceResolver  io.ResourceResolver
	env               runtime.Environment
	containerProvider func() component.Container

	components       []*component.Component
	container        component.Container
	lifecycleManager runtime.LifecycleManager
}

// createContext creates a new application context with the given environment.
// The returned context is not started yet. You need to call Start method to start the context.
func createContext(env runtime.Environment, startupContainer component.Container, resolver io.ResourceResolver) *Context {
	if env == nil {
		panic("nil environment")
	}

	if startupContainer == nil {
		panic("nil startup container")
	}

	if resolver == nil {
		panic("nil resource resolver")
	}

	return &Context{
		done:             make(chan struct{}),
		mu:               sync.RWMutex{},
		resourceResolver: resolver,
		env:              env,
		containerProvider: func() component.Container {
			container := component.NewStandardContainer()
			container.SetParentContainer(startupContainer)
			return container
		},
		components: component.List(),
	}
}

// Deadline method returns the time when work done on behalf of this context should be canceled.
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

// Done method returns a channel that's closed when work done on behalf of this context should be canceled.
func (c *Context) Done() <-chan struct{} {
	return c.done
}

// Err returns the error that caused the context to terminate. It returns nil if the context is still active.
func (c *Context) Err() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.err
}

// Value returns the value associated with the given key. Context values are currently not supported and
// this method always returns nil.
func (c *Context) Value(key any) any {
	return nil
}

// Start starts the application context. It loads component definitions, registers the context itself
// as a singleton, and initializes singleton components. After this method is called, the context is
// considered running.
func (c *Context) Start(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.doStart(ctx); err != nil {
		return &contextError{Op: "start", Err: err}
	}

	return nil
}

// doStart starts the lifecycle manager if the context has been refreshed and is not already running.
func (c *Context) doStart(ctx context.Context) error {
	if c.err != nil {
		return c.err
	}

	if !c.refreshed {
		return errors.New("context not refreshed")
	}

	if c.lifecycleManager.IsRunning() {
		return nil
	}

	if err := c.startLifecycleManager(ctx); err != nil {
		return err
	}

	return nil
}

// Stop stops the application context by shutting down lifecycle components. If the context is not running,
// it returns nil.
func (c *Context) Stop(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.doStop(ctx); err != nil {
		return &contextError{Op: "stop", Err: err}
	}

	return nil
}

// doStop stops the lifecycle manager if the context is currently running.
func (c *Context) doStop(ctx context.Context) error {
	if c.err != nil {
		return c.err
	}

	if !c.refreshed {
		return errors.New("context not refreshed")
	}

	if !c.lifecycleManager.IsRunning() {
		return nil
	}

	if err := c.stopLifecycleManager(ctx); err != nil {
		return err
	}

	return nil
}

// IsRunning returns true if the application context is currently running, false otherwise.
func (c *Context) IsRunning() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.lifecycleManager != nil && c.lifecycleManager.IsRunning()
}

// Refresh initializes the application context by preparing the container, loading component definitions,
// initializing singleton components, and starting lifecycle management.
func (c *Context) Refresh(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.doRefresh(ctx); err != nil {
		return &contextError{Op: "refresh", Err: err}
	}

	return nil
}

// doRefresh initializes the application context. It prepares the container, loads component definitions, initializes
// singleton components, resolves the lifecycle manager, and starts it.
func (c *Context) doRefresh(ctx context.Context) (err error) {
	if c.err != nil {
		return c.err
	}

	if c.refreshed {
		return errors.New("context already refreshed")
	}

	c.refreshed = true

	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = v
			default:
				err = fmt.Errorf("%v", v)
			}
		}

		if err != nil {
			log.Warn("Cancelling application context refresh attempt due to error: {}", err)
			err = errors.Join(err, c.cancelRefresh(ctx))
		}
	}()

	c.container = c.containerProvider()

	if err = c.prepareContainer(); err != nil {
		return err
	}

	if err = c.loadComponentDefinitions(ctx); err != nil {
		return err
	}

	if err = c.initializeSingletons(ctx); err != nil {
		return err
	}

	if err = c.resolveLifecycleManager(ctx); err != nil {
		return err
	}

	if err = c.container.RegisterSingleton(lifecycleManagerContainerKey, c.lifecycleManager); err != nil {
		return err
	}

	return c.startLifecycleManager(ctx)
}

// prepareContainer registers core infrastructure dependencies and singletons that must be available before
// component definitions are loaded.
func (c *Context) prepareContainer() error {

	if err := c.container.RegisterDependency(reflect.TypeFor[component.Container](), c.container); err != nil {
		return err
	}

	if err := c.container.RegisterDependency(reflect.TypeFor[runtime.Context](), c); err != nil {
		return err
	}

	if err := c.container.RegisterDependency(reflect.TypeFor[io.ResourceResolver](), c.resourceResolver); err != nil {
		return err
	}

	if err := c.container.RegisterSingleton(envContainerKey, c.env); err != nil {
		return err
	}

	return nil
}

// resolveLifecycleManager resolves a LifecycleManager from the container.
// If none is registered, a default implementation is created and used.
func (c *Context) resolveLifecycleManager(ctx context.Context) error {
	manager, err := component.ResolveType[runtime.LifecycleManager](ctx, c.container)
	if err != nil && !errors.Is(err, component.ErrNotFound) {
		return err
	}

	if manager != nil {
		c.lifecycleManager = manager
	} else if c.lifecycleManager == nil {
		c.lifecycleManager = newDefaultLifecycleManager(c.container)
	}

	return nil
}

// Close stops the application context, destroys all singleton components, releases resources, and marks
// the context as canceled.
func (c *Context) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.doClose(ctx); err != nil {
		return &contextError{Op: "close", Err: err}
	}

	return nil
}

// doClose stops lifecycle management, destroys singleton components, and marks the context as canceled.
func (c *Context) doClose(ctx context.Context) error {
	if c.err != nil {
		return c.err
	}

	err := c.stopLifecycleManager(ctx)
	if c.container != nil {
		c.container.DestroySingletons()
		c.container = nil
	}

	c.lifecycleManager = nil
	c.err = context.Canceled
	close(c.done)
	return err

}

// Environment returns the runtime environment associated with this context.
func (c *Context) Environment() runtime.Environment {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.env
}

// Container returns the component container associated with this context.
// It panics if the context has not been refreshed yet.
func (c *Context) Container() component.Container {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.container == nil {
		panic("nil container: context not refreshed")
	}

	return c.container
}

// ResourceResolver returns the resource resolver used by the application to load resources.
func (c *Context) ResourceResolver() io.ResourceResolver {
	return c.resourceResolver
}

// startLifecycleManager starts the lifecycle manager if it exists and is not already running.
func (c *Context) startLifecycleManager(ctx context.Context) error {
	if c.lifecycleManager != nil && !c.lifecycleManager.IsRunning() {
		if err := c.lifecycleManager.Startup(ctx); err != nil {
			return fmt.Errorf("start lifecycle manager: %w", err)
		}
	}

	return nil
}

// stopLifecycleManager stops the lifecycle manager if it exists and is currently running.
func (c *Context) stopLifecycleManager(ctx context.Context) error {
	if c.lifecycleManager != nil && c.lifecycleManager.IsRunning() {
		if err := c.lifecycleManager.Shutdown(ctx); err != nil {
			return fmt.Errorf("stop lifecycle manager: %w", err)
		}
	}

	return nil
}

// loadComponentDefinitions loads component definitions into the container using a ConditionalLoader.
// It retrieves the list of component definitions and loads them into the container, allowing for conditional
// loading based on the context.
func (c *Context) loadComponentDefinitions(ctx context.Context) error {
	filtered := make([]*component.Component, 0)

	for _, comp := range c.components {
		if isBootstrapType(comp.Definition().Type()) {
			continue
		}
		filtered = append(filtered, comp)
	}

	loader := component.NewConditionalLoader(c.container, filtered)
	err := loader.Load(ctx)
	if err != nil {
		return fmt.Errorf("load component definitions: %w", err)
	}

	return nil
}

// isBootstrapType reports whether the given type is considered a bootstrap component type.
func isBootstrapType(typ reflect.Type) bool {
	for _, bType := range bootstrapTypes {
		if typ.ConvertibleTo(bType) {
			return true
		}
	}

	return false
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
			return fmt.Errorf("initialize singleton %q: %w", definition.Name(), err)
		}
	}

	return nil
}

// cancelRefresh rolls back a failed refresh attempt by stopping lifecycle management, destroying initialized
// singletons, and clearing context state.
func (c *Context) cancelRefresh(ctx context.Context) error {
	if err := c.stopLifecycleManager(ctx); err != nil {
		return fmt.Errorf("cancel context refresh: %w", err)
	}

	if c.container != nil {
		c.container.DestroySingletons()
		c.container = nil
	}

	c.lifecycleManager = nil
	return nil

}
