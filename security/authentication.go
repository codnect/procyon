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

package security

import "context"

// Authentication represents the result of a successful authentication.
type Authentication interface {
	// Principal returns the authenticated principal.
	Principal() Principal
	// Authorities returns the authority granted to the authenticated principal.
	Authorities() []Authority
}

// Authenticator authenticates credentials and produces an authentication result.
type Authenticator interface {
	// Authenticate authenticates the given credential.
	Authenticate(ctx context.Context, credential Credential) (Authentication, error)
}
