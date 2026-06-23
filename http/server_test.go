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

package http

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	testCases := []struct {
		name       string
		dispatcher Dispatcher
		wantPanic  error
	}{
		{
			name:       "nil dispatcher",
			dispatcher: nil,
			wantPanic:  errors.New("nil dispatcher"),
		},
		{
			name:       "valid dispatcher",
			dispatcher: &RequestDispatcher{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					NewServer(ServerProperties{}, tc.dispatcher)
				})
				return
			}

			server := NewServer(ServerProperties{}, tc.dispatcher)

			// then
			require.NotNil(t, server)
		})
	}
}

func TestServer_Port(t *testing.T) {
	testCases := []struct {
		name     string
		props    ServerProperties
		wantPort int
	}{
		{
			name:     "with port",
			props:    ServerProperties{Port: 9090},
			wantPort: 9090,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			server := NewServer(tc.props, &RequestDispatcher{})

			// when
			port := server.Port()

			// then
			assert.Equal(t, tc.wantPort, port)
		})
	}
}
