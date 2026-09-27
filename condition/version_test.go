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
	"errors"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.codnect.io/procyon/component"
)

func TestOnGoVersion(t *testing.T) {
	testCases := []struct {
		name    string
		version string

		wantVersion string
		wantPanic   error
	}{
		{
			name:        "normalizes version",
			version:     "1.27",
			wantVersion: "go1.27",
		},
		{
			name:        "preserves normalized version",
			version:     "go1.27",
			wantVersion: "go1.27",
		},
		{
			name:        "normalizes patch version",
			version:     "1.27.1",
			wantVersion: "go1.27.1",
		},
		{
			name:      "panics on empty version",
			version:   "",
			wantPanic: errors.New("invalid Go version: go"),
		},
		{
			name:      "panics on invalid version",
			version:   "invalid",
			wantPanic: errors.New("invalid Go version: goinvalid"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					OnGoVersion(tc.version)
				})
				return
			}

			// when
			condition := OnGoVersion(tc.version)

			// then
			require.NotNil(t, condition)
			assert.Equal(t, tc.wantVersion, condition.version)
			assert.Equal(t, EqualOrNewer, condition.versionRange)
		})
	}
}

func TestGoVersion_Range(t *testing.T) {
	testCases := []struct {
		name         string
		versionRange VersionRange
	}{
		{
			name:         "sets equal or newer range",
			versionRange: EqualOrNewer,
		},
		{
			name:         "sets older than range",
			versionRange: OlderThan,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			condition := OnGoVersion("1.27")

			// when
			result := condition.Range(tc.versionRange)

			// then
			assert.Same(t, condition, result)
			assert.Equal(t, tc.versionRange, condition.versionRange)
		})
	}
}

func TestGoVersion_Matches(t *testing.T) {
	current := runtime.Version()

	testCases := []struct {
		name       string
		condition  *GoVersion
		wantResult bool
	}{
		{
			name:       "matches equal or newer version",
			condition:  OnGoVersion(current),
			wantResult: true,
		},
		{
			name:       "matches equal version",
			condition:  OnGoVersion(current).Range(EqualOrNewer),
			wantResult: true,
		},
		{
			name:       "does not match equal version as older",
			condition:  OnGoVersion(current).Range(OlderThan),
			wantResult: false,
		},
		{
			name: "does not match invalid version",
			condition: &GoVersion{
				version:      "invalid",
				versionRange: EqualOrNewer,
			},
			wantResult: false,
		},
		{
			name: "does not match invalid version range",
			condition: &GoVersion{
				version:      current,
				versionRange: VersionRange(255),
			},
			wantResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			ctx := context.Background()
			container := component.NewStandardContainer()

			// when
			result := tc.condition.Matches(ctx, container)

			// then
			assert.Equal(t, tc.wantResult, result)
		})
	}
}
