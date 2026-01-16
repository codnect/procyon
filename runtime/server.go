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

package runtime

import "context"

// Server interface represents a server with start and stop capabilities.
type Server interface {
	// Start method starts the server.
	Start(ctx context.Context) error
	// Stop method stops the server.
	Stop(ctx context.Context) error
	// Port method returns the port the server is running on.
	Port() int
}
