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
	"slices"

	"go.codnect.io/procyon/component"
	"go.codnect.io/procyon/runtime"
)

// Profile is a condition that matches when at least one of the specified
// profiles is active.
type Profile struct {
	profiles []string
}

// OnProfile creates a condition that matches when at least one of the
// specified profiles is active.
func OnProfile(profiles ...string) *Profile {
	return &Profile{
		profiles: profiles,
	}
}

// Matches reports whether at least one of the configured profiles is active.
func (p *Profile) Matches(ctx context.Context, container component.Container) bool {
	env, err := component.ResolveType[runtime.Environment](ctx, container)
	if err != nil {
		return false
	}

	return slices.ContainsFunc(p.profiles, env.IsProfileActive)
}

// MissingProfile is a condition that matches when none of the specified
// profiles are active.
type MissingProfile struct {
	profiles []string
}

// OnMissingProfile creates a condition that matches when none of the
// specified profiles are active.
func OnMissingProfile(profiles ...string) *MissingProfile {
	return &MissingProfile{
		profiles: profiles,
	}
}

// Matches reports whether none of the configured profiles are active.
func (p *MissingProfile) Matches(ctx context.Context, container component.Container) bool {
	env, err := component.ResolveType[runtime.Environment](ctx, container)
	if err != nil {
		return false
	}

	return !slices.ContainsFunc(p.profiles, env.IsProfileActive)
}
