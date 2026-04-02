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

	"codnect.io/procyon/component"
	"codnect.io/procyon/runtime"
)

type defaultLifecycleManager struct {
	container        component.Container
	lifecycleObjects []runtime.Lifecycle

	running bool
	mu      sync.Mutex
}

func newDefaultLifecycleManager(container component.Container) *defaultLifecycleManager {
	return &defaultLifecycleManager{
		container: container,
	}
}

func (d *defaultLifecycleManager) Startup(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	lifecycleObjects, err := component.ResolveAll[runtime.Lifecycle](ctx, d.container)
	if err != nil {
		return err
	}

	d.lifecycleObjects = lifecycleObjects

	for _, lifecycle := range lifecycleObjects {
		if err = lifecycle.Start(ctx); err != nil {
			return err
		}
	}

	d.running = true
	return nil
}

func (d *defaultLifecycleManager) Shutdown(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, lifecycle := range d.lifecycleObjects {
		if err := lifecycle.Stop(ctx); err != nil {
			return err
		}
	}

	d.running = false
	return nil
}

func (d *defaultLifecycleManager) IsRunning() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.running
}
