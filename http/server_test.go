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
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewServerProperties(t *testing.T) {
	// given

	// when
	props := newServerProperties()

	// then
	require.NotNil(t, props)
}

func TestServerProperties_Prefix(t *testing.T) {
	// given
	props := newServerProperties()

	// when
	prefix := props.Prefix()

	// then
	assert.Equal(t, "server", prefix)

}

func TestNewServerAdapter(t *testing.T) {
	testCases := []struct {
		name       string
		dispatcher Dispatcher
		wantPanic  error
	}{
		{
			name:      "nil dispatcher",
			wantPanic: errors.New("nil dispatcher"),
		},
		{
			name:       "valid dispatcher",
			dispatcher: &AnyDispatcher{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					newServerAdapter(tc.dispatcher)
				})
				return
			}
			adapter := newServerAdapter(tc.dispatcher)
			// then
			require.NotNil(t, adapter)
		})
	}
}

func TestServerAdapter_ServeHTTP(t *testing.T) {
	testCases := []struct {
		name         string
		preCondition func(dispatcher *AnyDispatcher)
	}{
		{
			name: "successfully dispatch",
			preCondition: func(dispatcher *AnyDispatcher) {
				dispatcher.
					On("Dispatch", mock.AnythingOfType("*http.Context")).
					Return(nil)
			},
		},
		{
			name: "dispatch error",
			preCondition: func(dispatcher *AnyDispatcher) {
				dispatcher.
					On("Dispatch", mock.AnythingOfType("*http.Context")).
					Return(errors.New("dispatch error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			dispatcher := &AnyDispatcher{}

			if tc.preCondition != nil {
				tc.preCondition(dispatcher)
			}

			adapter := newServerAdapter(dispatcher)
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			res := httptest.NewRecorder()

			// when
			adapter.ServeHTTP(res, req)

			// then
			dispatcher.AssertExpectations(t)
		})
	}
}

func TestNewDefaultServer(t *testing.T) {
	testCases := []struct {
		name       string
		properties *ServerProperties
		dispatcher Dispatcher
		wantPanic  error
	}{
		{
			name:       "nil server properties",
			properties: nil,
			wantPanic:  errors.New("nil server properties"),
		},
		{
			name:       "nil dispatcher",
			properties: &ServerProperties{},
			dispatcher: nil,
			wantPanic:  errors.New("nil dispatcher"),
		},
		{
			name:       "valid properties and dispatcher",
			properties: &ServerProperties{},
			dispatcher: &DefaultDispatcher{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					NewDefaultServer(tc.properties, tc.dispatcher)
				})
				return
			}

			server := NewDefaultServer(tc.properties, tc.dispatcher)

			// then
			require.NotNil(t, server)
		})
	}
}

func TestDefaultServer_Start(t *testing.T) {
	// given
	server := NewDefaultServer(
		&ServerProperties{
			Port: 0,
		},
		&DefaultDispatcher{},
	)

	// when
	err := server.Start(context.Background())

	// then
	require.NoError(t, err)
	require.NotNil(t, server.httpServer)
	assert.NotZero(t, server.boundPort)

	// cleanup
	require.NoError(t, server.Shutdown(context.Background()))
}

func TestDefaultServer_Shutdown(t *testing.T) {
	testCases := []struct {
		name         string
		preCondition func(t *testing.T, server *DefaultServer)
	}{
		{
			name: "server not started",
		},
		{
			name: "server started",
			preCondition: func(t *testing.T, server *DefaultServer) {
				err := server.Start(context.Background())
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			server := NewDefaultServer(
				&ServerProperties{Port: 0},
				&DefaultDispatcher{},
			)

			if tc.preCondition != nil {
				tc.preCondition(t, server)
			}

			// when
			err := server.Shutdown(context.Background())

			// then
			require.NoError(t, err)
		})
	}
}

func TestDefaultServer_Port(t *testing.T) {
	testCases := []struct {
		name     string
		props    *ServerProperties
		wantPort int
	}{
		{
			name:     "with port",
			props:    &ServerProperties{Port: 9090},
			wantPort: 9090,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			server := NewDefaultServer(tc.props, &DefaultDispatcher{})

			// when
			port := server.Port()

			// then
			assert.Equal(t, tc.wantPort, port)
		})
	}
}
