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

package io

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultResourceResolver_Resolve(t *testing.T) {

	testCases := []struct {
		name     string
		ctx      context.Context
		location string

		wantErr      error
		wantName     string
		wantLocation string
	}{
		{
			name:     "invalid url",
			location: "http://local host:8080/resources/procyon.yaml",
			wantErr:  errors.New("parse \"http://local host:8080/resources/procyon.yaml\": invalid character \" \" in host name"),
		},
		{
			name:         "path without scheme",
			location:     "resources/procyon.yaml",
			wantName:     "procyon.yaml",
			wantLocation: "resources/procyon.yaml",
		},
		{
			name:         "path with file scheme (single-colon syntax)",
			location:     "file:resources/procyon.yaml",
			wantName:     "procyon.yaml",
			wantLocation: "resources/procyon.yaml",
		},
		{
			name:         "path with file scheme (double-slash syntax)",
			location:     "file://resources/procyon.yaml",
			wantName:     "procyon.yaml",
			wantLocation: "resources/procyon.yaml",
		},
		{
			name:         "path with http scheme",
			location:     "http://localhost:8080/resources/procyon.yaml",
			wantName:     "procyon.yaml",
			wantLocation: "http://localhost:8080/resources/procyon.yaml",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			resolver := NewDefaultResourceResolver()

			// when
			resource, err := resolver.Resolve(tc.ctx, tc.location)

			// then
			if tc.wantErr != nil {
				require.Nil(t, resource)

				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			require.NotNil(t, resource)
			assert.Equal(t, tc.wantName, resource.Name())
			assert.Equal(t, tc.wantLocation, resource.Location())
		})
	}
}
