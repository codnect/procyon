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
	"codnect.io/procyon/runtime"
)

const (
	// appContextContainerKey is the key used to register the application context itself in the component container.
	appContextContainerKey = "procyonAppContext"
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
	done chan struct{}
	err  error
	mu   sync.RWMutex

	env               runtime.Environment
	containerProvider func() component.Container
	components        []*component.Component

	container        component.Container
	lifecycleManager runtime.LifecycleManager
}

// createContext creates a new application context with the given environment.
// The returned context is not started yet. You need to call Start method to start the context.
func createContext(env runtime.Environment, startupContainer component.Container) *Context {
	if env == nil {
		panic("nil environment")
	}

	return &Context{
		done: make(chan struct{}),
		mu:   sync.RWMutex{},
		env:  env,
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

	if c.err != nil {
		return c.err
	}

	if c.lifecycleManager == nil {
		return &contextError{Op: "start", Err: errors.New("context not refreshed")}
	}

	if c.lifecycleManager.IsRunning() {
		return nil
	}

	err := c.startLifecycleManager(ctx)
	if err != nil {
		return &contextError{Op: "start", Err: err}
	}

	return nil
}

// Stop stops the application context by shutting down lifecycle components. If the context is not running,
// it returns nil.
func (c *Context) Stop(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.err != nil {
		return c.err
	}

	if c.lifecycleManager == nil {
		return &contextError{Op: "stop", Err: errors.New("context not refreshed")}
	}

	if !c.lifecycleManager.IsRunning() {
		return nil
	}

	err := c.stopLifecycleManager(ctx)
	if err != nil {
		return &contextError{Op: "stop", Err: err}
	}

	return nil
}

// IsRunning returns true if the application context is currently running, false otherwise.
func (c *Context) IsRunning() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.lifecycleManager != nil && c.lifecycleManager.IsRunning()
}

// Refresh reloads the application context by recreating the container, reloading component  definitions,
// reinitializing singletons, and restarting lifecycle management.
func (c *Context) Refresh(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.err != nil {
		return &contextError{Op: "refresh", Err: c.err}
	}

	err := c.stopLifecycleManager(ctx)
	if err != nil {
		return &contextError{Op: "refresh", Err: err}
	}

	if c.container != nil {
		c.container.DestroySingletons()
		c.container = nil
	}

	c.container = c.containerProvider()

	err = c.loadComponentDefinitions(ctx)
	if err != nil {
		return &contextError{Op: "refresh", Err: err}
	}

	err = c.container.RegisterSingleton(appContextContainerKey, c)
	if err != nil {
		return &contextError{Op: "refresh", Err: err}
	}

	err = c.initializeSingletons(ctx)
	if err != nil {
		return &contextError{Op: "refresh", Err: err}
	}

	var lifecycleManager runtime.LifecycleManager
	lifecycleManager, err = component.ResolveType[runtime.LifecycleManager](ctx, c.container)
	if err != nil && !errors.Is(err, component.ErrNotFound) {
		return err
	} else if lifecycleManager != nil {
		c.lifecycleManager = lifecycleManager
	}

	if c.lifecycleManager == nil {
		c.lifecycleManager = newDefaultLifecycleManager(c.container)
	}

	err = c.container.RegisterSingleton(lifecycleManagerContainerKey, c.lifecycleManager)
	if err != nil {
		return &contextError{Op: "refresh", Err: err}
	}

	err = c.startLifecycleManager(ctx)
	if err != nil {
		return &contextError{Op: "refresh", Err: err}
	}

	return nil
}

// Close stops the application context, destroys all singleton components, releases resources, and
// marks the context as canceled.
func (c *Context) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

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
