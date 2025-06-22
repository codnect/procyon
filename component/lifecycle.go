// Copyright 2025 Codnect
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

package component

import (
	"context"
)

// Initializer can be implemented by components that require custom logic during initialization.
type Initializer interface {
	// Init is called once the component has been constructed and dependencies injected.
	Init(ctx context.Context) error
}

// PreProcessor defines logic that runs before a component is initialized.
type PreProcessor interface {
	// ProcessBeforeInit is called before the Init method.
	// Allows modification or validation of the instance before initialization.
	ProcessBeforeInit(ctx context.Context, instance any) (any, error)
}

// PostProcessor defines logic that runs after a component is initialized.
type PostProcessor interface {
	// ProcessAfterInit is called after the Init method.
	// Allows additional configuration or enhancement of the instance.
	ProcessAfterInit(ctx context.Context, instance any) (any, error)
}

// Finalizer can be implemented by components that need to release resources or clean up during shutdown.
type Finalizer interface {
	// Finalize is called during application shutdown for cleanup purposes.
	Finalize() error
}

// LifecycleManager manages the registration of lifecycle hooks such as pre/post processors.
type LifecycleManager interface {
	// UsePreProcessor registers a PreProcessor to be applied before initialization.
	UsePreProcessor(processor PreProcessor) error

	// UsePostProcessor registers a PostProcessor to be applied after initialization.
	UsePostProcessor(processor PostProcessor) error
}
