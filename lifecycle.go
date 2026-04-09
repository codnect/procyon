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
	"reflect"
	"sync"
	"time"

	"codnect.io/procyon/component"
	"codnect.io/procyon/runtime"
)

type defaultLifecycleManager struct {
	container       component.Container
	shutdownTimeout time.Duration

	lifecycleObjects map[string]runtime.Lifecycle
	running          bool
	mu               sync.RWMutex
}

func newDefaultLifecycleManager(container component.Container) *defaultLifecycleManager {
	return &defaultLifecycleManager{
		container:        container,
		lifecycleObjects: make(map[string]runtime.Lifecycle),
	}
}

func (d *defaultLifecycleManager) Startup(ctx context.Context) error {
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

	for _, lifecycle := range d.lifecycleObjects {
		if err := lifecycle.Start(ctx); err != nil {
			return err
		}
	}

	d.running = true
	return nil
}

func (d *defaultLifecycleManager) Shutdown(ctx context.Context) error {
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
				log.Warn("Failed to stop component {}", name, err)
			}
		case <-shutdownCtx.Done():
			d.running = false
			return nil
		}
	}

	d.running = false
	return nil
}

func (d *defaultLifecycleManager) IsRunning() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.running
}
