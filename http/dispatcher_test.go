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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type AnyMiddleware struct {
	err error
}

func (m *AnyMiddleware) Invoke(ctx *Context, next RequestDelegate) error {
	return m.err
}

func TestNewRequestDispatcher(t *testing.T) {
	testCases := []struct {
		name            string
		endpointMatcher EndpointMatcher
		middlewares     []Middleware
		wantPanic       error
	}{
		{
			name:            "nil endpoint matcher",
			endpointMatcher: nil,
			wantPanic:       fmt.Errorf("nil endpoint matcher"),
		},
		{
			name:            "valid endpoint matcher",
			endpointMatcher: NewRequestEndpointMatcher(nil),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					NewRequestDispatcher(tc.endpointMatcher, tc.middlewares...)
				})
				return
			}

			dispatcher := NewRequestDispatcher(tc.endpointMatcher, tc.middlewares...)
			require.NotNil(t, dispatcher, "nil dispatcher")

			// then
		})
	}
}

func TestRequestDispatcher_Dispatch(t *testing.T) {
	testCases := []struct {
		name          string
		middlewares   []Middleware
		ctx           *Context
		preCondition  func(requestDelegate *AnyMockRequestDelegate, endpointMatcher *AnyMockEndpointMatcher)
		postCondition func(t *testing.T, requestDelegate *AnyMockRequestDelegate)
		wantErr       error
	}{
		{
			name: "middleware error",
			middlewares: []Middleware{
				&AnyMiddleware{
					err: errors.New("middleware error"),
				},
			},
			ctx: &Context{},
			preCondition: func(requestDelegate *AnyMockRequestDelegate, endpointMatcher *AnyMockEndpointMatcher) {
				requestDelegate.On("Invoke", mock.AnythingOfType("*http.Context")).
					Return(nil)

				endpoint := NewEndpoint(MethodGet, "/", requestDelegate.Invoke)
				endpointMatcher.On("Match", mock.AnythingOfType("*http.Context")).
					Return(endpoint, true)
			},
			postCondition: func(t *testing.T, requestDelegate *AnyMockRequestDelegate) {
				requestDelegate.AssertNotCalled(t, "Invoke", mock.AnythingOfType("*http.Context"))
			},
			wantErr: errors.New("middleware error"),
		},
		{
			name:        "endpoint error",
			middlewares: []Middleware{},
			ctx:         &Context{},
			preCondition: func(requestDelegate *AnyMockRequestDelegate, endpointMatcher *AnyMockEndpointMatcher) {
				requestDelegate.On("Invoke", mock.AnythingOfType("*http.Context")).
					Return(errors.New("endpoint error"))

				endpoint := NewEndpoint(MethodGet, "/", requestDelegate.Invoke)
				endpointMatcher.On("Match", mock.AnythingOfType("*http.Context")).
					Return(endpoint, true)
			},
			wantErr: errors.New("endpoint error"),
		},
		{
			name:        "dispatch successfully",
			middlewares: []Middleware{},
			ctx:         &Context{},
			preCondition: func(requestDelegate *AnyMockRequestDelegate, endpointMatcher *AnyMockEndpointMatcher) {
				requestDelegate.On("Invoke", mock.AnythingOfType("*http.Context")).
					Return(nil)

				endpoint := NewEndpoint(MethodGet, "/", requestDelegate.Invoke)
				endpointMatcher.On("Match", mock.AnythingOfType("*http.Context")).
					Return(endpoint, true)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			requestDelegate := &AnyMockRequestDelegate{}
			mockEndpointMatcher := &AnyMockEndpointMatcher{}
			dispatcher := NewRequestDispatcher(mockEndpointMatcher, tc.middlewares...)

			if tc.preCondition != nil {
				tc.preCondition(requestDelegate, mockEndpointMatcher)
			}

			// when
			err := dispatcher.Dispatch(tc.ctx)

			// then
			if tc.postCondition != nil {
				tc.postCondition(t, requestDelegate)
			}

			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}
