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

package http

// SecurityContext represents authorization and authentication data associated with the request.
type SecurityContext interface {
	// Principal returns the authenticated principal or nil when unauthenticated.
	Principal() any
	// IsAuthenticated returns true when the principal is authenticated.
	IsAuthenticated() bool
	// HasAuthority checks for a granted authority/role.
	HasAuthority(string) bool
	// Attributes exposes additional security attributes.
	Attributes() map[string]any
}

// SimpleSecurityContext provides a minimal SecurityContext implementation.
type SimpleSecurityContext struct {
	principal     any
	authenticated bool
	authorities   []string
	attributes    map[string]any
}

// NewUnauthenticatedSecurityContext creates an unauthenticated security context.
func NewUnauthenticatedSecurityContext() *SimpleSecurityContext {
	return &SimpleSecurityContext{
		principal:     nil,
		authenticated: false,
		attributes:    make(map[string]any),
	}
}

// NewSecurityContext creates an authenticated security context with the given principal and authorities.
func NewSecurityContext(principal any, authorities ...string) *SimpleSecurityContext {
	return &SimpleSecurityContext{
		principal:     principal,
		authenticated: true,
		authorities:   authorities,
		attributes:    make(map[string]any),
	}
}

// Principal returns the associated principal.
func (c *SimpleSecurityContext) Principal() any {
	return c.principal
}

// IsAuthenticated returns true if the principal is authenticated.
func (c *SimpleSecurityContext) IsAuthenticated() bool {
	return c.authenticated
}

// HasAuthority checks if the requested authority is present.
func (c *SimpleSecurityContext) HasAuthority(authority string) bool {
	for _, granted := range c.authorities {
		if granted == authority {
			return true
		}
	}
	return false
}

// Attributes returns additional security attributes.
func (c *SimpleSecurityContext) Attributes() map[string]any {
	if c.attributes == nil {
		c.attributes = make(map[string]any)
	}
	return c.attributes
}
