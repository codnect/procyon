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
	"go/version"
	"runtime"

	"go.codnect.io/procyon/component"
)

// VersionRange defines how a version is compared against the current
// Go version.
type VersionRange uint8

const (
	// EqualOrNewer matches when the current Go version is equal to or newer
	// than the configured version.
	EqualOrNewer VersionRange = iota

	// OlderThan matches when the current Go version is older than the
	// configured version.
	OlderThan
)

// GoVersion is a condition that matches against the current Go version.
type GoVersion struct {
	version      string
	versionRange VersionRange
}

// OnGoVersion creates a condition that matches when the current Go version
// is equal to or newer than the specified version by default.
func OnGoVersion(v string) *GoVersion {
	v = normalizeGoVersion(v)
	if !version.IsValid(v) {
		panic("invalid Go version: " + v)
	}

	return &GoVersion{
		version:      v,
		versionRange: EqualOrNewer,
	}
}

// Range sets the version range used to compare the configured Go version
// against the current Go version.
func (g *GoVersion) Range(r VersionRange) *GoVersion {
	g.versionRange = r
	return g
}

// Matches reports whether the current Go version satisfies the configured
// version range.
func (g *GoVersion) Matches(ctx context.Context, container component.Container) bool {
	current := runtime.Version()

	if !version.IsValid(current) || !version.IsValid(g.version) {
		return false
	}

	switch g.versionRange {
	case EqualOrNewer:
		return version.Compare(current, g.version) >= 0
	case OlderThan:
		return version.Compare(current, g.version) < 0
	default:
		return false
	}
}

// normalizeGoVersion ensures the version has the "go" prefix required by
// the go/version package.
func normalizeGoVersion(v string) string {
	if len(v) >= 2 && v[:2] == "go" {
		return v
	}

	return "go" + v
}
