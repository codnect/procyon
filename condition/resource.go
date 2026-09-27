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

package condition

import (
	"context"

	"go.codnect.io/procyon/component"
	"go.codnect.io/procyon/io"
)

// Resource is a condition that matches when all specified resources exist.
type Resource struct {
	locations []string
}

// OnResource creates a condition that matches when all resources at the
// specified locations exist.
func OnResource(locations ...string) *Resource {
	return &Resource{
		locations: locations,
	}
}

// Matches reports whether all configured resources exist.
func (r *Resource) Matches(ctx context.Context, container component.Container) bool {
	resolver, err := component.ResolveType[io.ResourceResolver](ctx, container)
	if err != nil {
		return false
	}

	for _, location := range r.locations {
		var resource io.Resource
		resource, err = resolver.Resolve(ctx, location)
		if err != nil || !resource.Exists() {
			return false
		}
	}

	return true
}
