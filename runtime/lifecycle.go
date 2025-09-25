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

import "context"

// Lifecycle interface defines methods for managing the start/stop lifecycle of a component.
type Lifecycle interface {
	// Start starts this component.
	Start(ctx context.Context) error
	// Stop stops this component.
	Stop(ctx context.Context) error
	// IsRunning indicates whether this component is currently running.
	IsRunning() bool
}
