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

package runtime

import (
	"context"

	"codnect.io/procyon/component"
)

// ApplicationContext represents the central runtime context of the application.
type ApplicationContext interface {
	context.Context

	// Lifecycle interface provides start/stop lifecycle methods for the application context.
	Lifecycle

	// EnvironmentCapable interface provides access to the environment.
	EnvironmentCapable

	// ContainerCapable interface provides access to the component container.
	component.ContainerCapable

	// Refresh reloads the application context contents (environment, container)
	Refresh() error
}
