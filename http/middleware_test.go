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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type AnyMockRequestDelegate struct {
	mock.Mock
}

func (a *AnyMockRequestDelegate) Invoke(ctx *Context) error {
	result := a.Called(ctx)
	return result.Error(0)
}

type AnyMockEndpointMatcher struct {
	mock.Mock
}

func (a *AnyMockEndpointMatcher) Match(ctx *Context) (*Endpoint, bool) {
	result := a.Called(ctx)
	if result.Get(0) == nil {
		return nil, false
	}

	return result.Get(0).(*Endpoint), true
}

func TestRoutingMiddleware_Invoke(t *testing.T) {
	endpoint := NewEndpoint(MethodGet, "/", func(ctx *Context) error {
		return nil
	})

	testCases := []struct {
		name           string
		nextMiddleware RequestDelegate
		preCondition   func(matcher *AnyMockEndpointMatcher)
		postCondition  func(t *testing.T, ctx *Context)
		wantErr        error
	}{
		{
			name: "no match endpoint",
			nextMiddleware: func(ctx *Context) error {
				return nil
			},
			preCondition: func(matcher *AnyMockEndpointMatcher) {
				matcher.On("Match", mock.AnythingOfType("*http.Context")).
					Return(nil, false)
			},
			postCondition: func(t *testing.T, ctx *Context) {
				assert.Nil(t, ctx.Endpoint())
			},
		},
		{
			name: "endpoint matched",
			nextMiddleware: func(ctx *Context) error {
				return nil
			},
			preCondition: func(matcher *AnyMockEndpointMatcher) {
				matcher.On("Match", mock.AnythingOfType("*http.Context")).
					Return(endpoint)
			},
			postCondition: func(t *testing.T, ctx *Context) {
				assert.NotNil(t, ctx.Endpoint())
				assert.Equal(t, endpoint, ctx.Endpoint())
			},
		},
		{
			name: "next middleware error",
			nextMiddleware: func(ctx *Context) error {
				return errors.New("middleware error")
			},
			preCondition: func(matcher *AnyMockEndpointMatcher) {
				matcher.On("Match", mock.AnythingOfType("*http.Context")).
					Return(endpoint)
			},
			postCondition: func(t *testing.T, ctx *Context) {
				assert.Nil(t, ctx.Endpoint())
				assert.Equal(t, endpoint, ctx.Endpoint())
			},
			wantErr: errors.New("middleware error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			ctx := &Context{}
			mockMatcher := &AnyMockEndpointMatcher{}

			if tc.preCondition != nil {
				tc.preCondition(mockMatcher)
			}

			middleware := newRoutingMiddleware(mockMatcher)

			// when
			err := middleware.Invoke(ctx, tc.nextMiddleware)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestEndpointMiddleware_Invoke(t *testing.T) {

	testCases := []struct {
		name          string
		preCondition  func(ctx *Context, requestDelegate *AnyMockRequestDelegate)
		postCondition func(t *testing.T, ctx *Context, requestDelegate *AnyMockRequestDelegate)
		wantErr       error
	}{
		{
			name: "no endpoint matched",
			preCondition: func(ctx *Context, requestDelegate *AnyMockRequestDelegate) {
				requestDelegate.On("Invoke", mock.AnythingOfType("*http.Context")).
					Return(nil)
			},
			postCondition: func(t *testing.T, ctx *Context, requestDelegate *AnyMockRequestDelegate) {
				requestDelegate.AssertCalled(t, "Invoke", mock.AnythingOfType("*http.Context"))
			},
		},
		{
			name: "endpoint matched",
			preCondition: func(ctx *Context, requestDelegate *AnyMockRequestDelegate) {
				requestDelegate.On("Invoke", mock.AnythingOfType("*http.Context")).
					Return(nil)

				endpoint := NewEndpoint(MethodGet, "/", requestDelegate.Invoke)
				ctx.SetEndpoint(endpoint)
			},
			postCondition: func(t *testing.T, ctx *Context, requestDelegate *AnyMockRequestDelegate) {
				requestDelegate.AssertCalled(t, "Invoke", mock.AnythingOfType("*http.Context"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			ctx := &Context{}
			requestDelegate := &AnyMockRequestDelegate{}

			if tc.preCondition != nil {
				tc.preCondition(ctx, requestDelegate)
			}

			middleware := newEndpointMiddleware()

			// when
			err := middleware.Invoke(ctx, requestDelegate.Invoke)

			// then
			if tc.postCondition != nil {
				tc.postCondition(t, ctx, requestDelegate)
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
