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
	"fmt"
	"reflect"
	"sync"
	"time"

	"codnect.io/procyon/component"
	"codnect.io/procyon/runtime"
)

// defaultLifecycleManager is the default implementation of LifecycleManager. It manages the lifecycle of components
// that implement the runtime.Lifecycle interface. It starts all lifecycle components during application startup and
// stops them during application shutdown.
type defaultLifecycleManager struct {
	container       component.Container
	shutdownTimeout time.Duration

	lifecycleObjects map[string]runtime.Lifecycle
	running          bool
	mu               sync.RWMutex
}

// newDefaultLifecycleManager creates a new instance of defaultLifecycleManager with the provided container.
func newDefaultLifecycleManager(container component.Container) *defaultLifecycleManager {
	return &defaultLifecycleManager{
		container:        container,
		lifecycleObjects: make(map[string]runtime.Lifecycle),
	}
}

// Startup starts all lifecycle components in the container. It resolves all components that implement the
// runtime.Lifecycle interface and calls their Start method. If any component fails to start, it returns an error.
func (d *defaultLifecycleManager) Startup(ctx context.Context) error {
	if ctx == nil {
		return fmt.Errorf("nil context")
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	definitions := d.container.DefinitionsOf(reflect.TypeFor[runtime.Lifecycle]())
	for _, definition := range definitions {
		lifecycleObj, err := d.container.Resolve(ctx, definition.Name())
		if err != nil {
			return err
		}

		d.lifecycleObjects[definition.Name()] = lifecycleObj.(runtime.Lifecycle)
	}

	for objectName, lifecycle := range d.lifecycleObjects {
		if err := lifecycle.Start(ctx); err != nil {
			return fmt.Errorf("start lifecycle component %q: %w", objectName, err)
		}

		log.Debug("Started lifecycle component '{}'", objectName)
	}

	d.running = true
	return nil
}

// Shutdown stops all running lifecycle components before destruction. It calls the Stop method of each lifecycle
// component. If any component fails to stop, it logs a warning and continues stopping the remaining components.
// It uses a timeout context to ensure that the shutdown process does not hang indefinitely.
func (d *defaultLifecycleManager) Shutdown(ctx context.Context) error {
	if ctx == nil {
		return fmt.Errorf("nil context")
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	shutdownCtx, cancel := context.WithTimeout(ctx, d.shutdownTimeout)
	defer cancel()

	for name, lifecycle := range d.lifecycleObjects {
		done := make(chan error, 1)

		go func() {
			done <- lifecycle.Stop(shutdownCtx)
		}()

		select {
		case err := <-done:
			if err != nil {
				log.Warn("Failed to stopLifecycleManager component '{}'", name, err)
			}

			log.Debug("Stopped lifecycle component '{}'", name)
		case <-shutdownCtx.Done():
			d.running = false
			return nil
		}
	}

	d.running = false
	return nil
}

// IsRunning indicates whether this lifecycle manager is currently running. It returns true if the lifecycle manager
// is running, false otherwise.
func (d *defaultLifecycleManager) IsRunning() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.running
}
