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
	"codnect.io/procyon/io"
)

// Application defines the interface for an application.
type Application interface {
	// SetBannerPrinter sets the banner printer to use for displaying the application banner.
	SetBannerPrinter(printer BannerPrinter)
	// ResourceResolver returns the resource resolver used by the application.
	ResourceResolver() io.ResourceResolver
	// Run starts the application with the given command-line arguments.
	Run(args ...string) error
}

// Context represents the central runtime context of the application.
type Context interface {
	context.Context

	// Lifecycle interface provides start/stop lifecycle methods for the application context.
	Lifecycle

	// EnvironmentCapable interface provides access to the environment.
	EnvironmentCapable

	// ContainerCapable interface provides access to the component container.
	component.ContainerCapable

	// Refresh reloads the application context contents (environment, container)
	Refresh(ctx context.Context) error

	// Close closes the application context and releases all resources.
	Close(ctx context.Context) error
}

// ContextCustomizer is an interface for customizing the application context.
type ContextCustomizer interface {
	// CustomizeContext customizes the given application context.
	CustomizeContext(ctx Context) error
}
