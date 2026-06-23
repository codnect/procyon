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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateContext(t *testing.T) {
	testCases := []struct {
		name         string
		nativeReq    *http.Request
		nativeWriter http.ResponseWriter
		wantPanic    error
	}{
		{
			name:         "nil http request",
			nativeReq:    nil,
			nativeWriter: nil,
			wantPanic:    errors.New("nil http request"),
		},
		{
			name:         "nil writer",
			nativeReq:    httptest.NewRequest("GET", "/", nil),
			nativeWriter: nil,
			wantPanic:    errors.New("nil response writer"),
		},
		{
			name:         "valid context",
			nativeReq:    httptest.NewRequest("GET", "/", nil),
			nativeWriter: httptest.NewRecorder(),
			wantPanic:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					CreateContext(tc.nativeReq, tc.nativeWriter)
				})
				return
			}

			ctx := CreateContext(tc.nativeReq, tc.nativeWriter)

			// then
			require.NotNil(t, ctx, "nil http context")
		})
	}
}

func TestContext_Deadline(t *testing.T) {
	// given
	deadline := time.Now()
	ctx, _ := context.WithDeadline(context.Background(), deadline)

	nativeRequest := httptest.NewRequestWithContext(ctx, "GET", "/", nil)
	nativeWriter := httptest.NewRecorder()

	httpCtx := CreateContext(nativeRequest, nativeWriter)

	// when
	result, ok := httpCtx.Deadline()

	// then
	assert.True(t, ok)
	assert.Equal(t, deadline, result)
}

func TestContext_Done(t *testing.T) {
	// given
	ctx, cancelFn := context.WithCancel(context.Background())

	// cancel context
	cancelFn()

	nativeRequest := httptest.NewRequestWithContext(ctx, "GET", "/", nil)
	nativeWriter := httptest.NewRecorder()

	httpCtx := CreateContext(nativeRequest, nativeWriter)

	// when
	<-httpCtx.Done()

	// then
	assert.Equal(t, context.Canceled, ctx.Err())
}

func TestContext_Err(t *testing.T) {
	customErr := errors.New("custom error")

	testCases := []struct {
		name            string
		contextProvider func() context.Context
		preCondition    func(ctx *Context)
		wantErrors      []error
	}{
		{
			name: "canceled context",
			contextProvider: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErrors: []error{context.Canceled},
		},
		{
			name: "custom error",
			contextProvider: func() context.Context {
				return context.Background()
			},
			preCondition: func(ctx *Context) {
				ctx.err = customErr
			},
			wantErrors: []error{customErr},
		},
		{
			name: "multiple errors",
			contextProvider: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			preCondition: func(ctx *Context) {
				ctx.err = customErr
			},
			wantErrors: []error{context.Canceled, customErr},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			ctx := tc.contextProvider()

			nativeRequest := httptest.NewRequestWithContext(ctx, "GET", "/", nil)
			nativeWriter := httptest.NewRecorder()

			httpCtx := CreateContext(nativeRequest, nativeWriter)

			if tc.preCondition != nil {
				tc.preCondition(httpCtx)
			}

			// when
			err := httpCtx.Err()

			// then
			for _, wantError := range tc.wantErrors {
				assert.ErrorIs(t, err, wantError)
			}
		})
	}
}

func TestContext_Value(t *testing.T) {
	testCases := []struct {
		name            string
		contextProvider func() context.Context
		preCondition    func(ctx *Context)
		key             string
		wantValue       string
	}{
		{
			name: "value from parent context",
			contextProvider: func() context.Context {
				return context.WithValue(context.Background(), "anyKey", "anyValue")
			},
			preCondition: func(ctx *Context) {},
			key:          "anyKey",
			wantValue:    "anyValue",
		},
		{
			name: "value from current context",
			contextProvider: func() context.Context {
				return context.Background()
			},
			preCondition: func(ctx *Context) {
				ctx.SetValue("anotherKey", "anotherValue")
			},
			key:       "anotherKey",
			wantValue: "anotherValue",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			ctx := tc.contextProvider()
			nativeReq := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			recorder := httptest.NewRecorder()

			httpCtx := CreateContext(nativeReq, recorder)

			if tc.preCondition != nil {
				tc.preCondition(httpCtx)
			}

			// when
			value := httpCtx.Value(tc.key)

			// then
			assert.Equal(t, tc.wantValue, value)
		})
	}
}

func TestContext_SetValue(t *testing.T) {
	// given
	nativeReq := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	httpCtx := CreateContext(nativeReq, recorder)

	// when
	httpCtx.SetValue("anyKey", "anyValue")

	// then
	assert.Contains(t, httpCtx.values, "anyKey")

	value := httpCtx.values["anyKey"]
	assert.Contains(t, value, "anyValue")
}

func TestContext_Request(t *testing.T) {
	// given
	nativeReq := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	httpCtx := CreateContext(nativeReq, recorder)

	// when
	request := httpCtx.Request()

	// then
	assert.NotNil(t, request)
	assert.Equal(t, nativeReq, request.nativeReq)
}

func TestContext_Response(t *testing.T) {
	// given
	nativeReq := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	httpCtx := CreateContext(nativeReq, recorder)

	// when
	response := httpCtx.Response()

	// then
	assert.NotNil(t, response)
	assert.Equal(t, recorder, response.writer)
}

func TestContext_Endpoint(t *testing.T) {
	// given
	reqEndpoint := NewEndpoint(MethodGet, "/test", func(ctx *Context) error {
		return nil
	})

	nativeReq := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	httpCtx := CreateContext(nativeReq, recorder)
	httpCtx.SetEndpoint(reqEndpoint)

	// when
	endpoint := httpCtx.Endpoint()

	// then
	assert.Equal(t, reqEndpoint, endpoint)
}

func TestContext_SetEndpoint(t *testing.T) {
	// given
	endpoint := NewEndpoint(MethodGet, "/test", func(ctx *Context) error {
		return nil
	})

	nativeReq := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	httpCtx := CreateContext(nativeReq, recorder)

	// when
	httpCtx.SetEndpoint(endpoint)

	// then
	assert.Equal(t, endpoint, httpCtx.Endpoint())
}

func TestContext_reset(t *testing.T) {
	// given
	nativeReq := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	httpCtx := CreateContext(nativeReq, recorder)

	httpCtx.SetValue("anyKey", "anyValue")
	httpCtx.Response().SetHeader("anyHeaderKey", "anyHeaderValue")

	newNativeReq := httptest.NewRequest(http.MethodGet, "/", nil)
	newRecorder := httptest.NewRecorder()

	// when
	httpCtx.reset(newNativeReq, newRecorder)

	// then
	assert.Nil(t, httpCtx.err)
	assert.Empty(t, httpCtx.values)
	assert.Same(t, newNativeReq, httpCtx.req.nativeReq)
	assert.Nil(t, httpCtx.req.cookiesCache)
	assert.Nil(t, httpCtx.req.queryCache)
	assert.Same(t, newRecorder, httpCtx.res.writer)
	assert.Equal(t, StatusOK, httpCtx.res.status)
	assert.False(t, httpCtx.res.writtenHeaders)
	assert.False(t, httpCtx.res.writerUsed)
	assert.Empty(t, httpCtx.res.headers)
}

func TestEndpointContext_Deadline(t *testing.T) {
	// given
	deadline := time.Now()
	ctx, _ := context.WithDeadline(context.Background(), deadline)

	nativeRequest := httptest.NewRequestWithContext(ctx, "GET", "/", nil)
	nativeWriter := httptest.NewRecorder()

	endpointCtx := &EndpointContext[any]{
		ctx: CreateContext(nativeRequest, nativeWriter),
	}

	// when
	result, ok := endpointCtx.Deadline()

	// then
	assert.True(t, ok)
	assert.Equal(t, deadline, result)
}

func TestEndpointContext_Done(t *testing.T) {
	// given
	ctx, cancelFn := context.WithCancel(context.Background())

	// cancel context
	cancelFn()

	nativeRequest := httptest.NewRequestWithContext(ctx, "GET", "/", nil)
	nativeWriter := httptest.NewRecorder()

	endpointCtx := &EndpointContext[any]{
		ctx: CreateContext(nativeRequest, nativeWriter),
	}

	// when
	<-endpointCtx.Done()

	// then
	assert.Equal(t, context.Canceled, ctx.Err())
}

func TestEndpointContext_Err(t *testing.T) {
	customErr := errors.New("custom error")

	testCases := []struct {
		name            string
		contextProvider func() context.Context
		preCondition    func(ctx *Context)
		wantErrors      []error
	}{
		{
			name: "canceled context",
			contextProvider: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErrors: []error{context.Canceled},
		},
		{
			name: "custom error",
			contextProvider: func() context.Context {
				return context.Background()
			},
			preCondition: func(ctx *Context) {
				ctx.err = customErr
			},
			wantErrors: []error{customErr},
		},
		{
			name: "multiple errors",
			contextProvider: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			preCondition: func(ctx *Context) {
				ctx.err = customErr
			},
			wantErrors: []error{context.Canceled, customErr},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			ctx := tc.contextProvider()

			nativeRequest := httptest.NewRequestWithContext(ctx, "GET", "/", nil)
			nativeWriter := httptest.NewRecorder()

			endpointCtx := &EndpointContext[any]{
				ctx: CreateContext(nativeRequest, nativeWriter),
			}

			if tc.preCondition != nil {
				tc.preCondition(endpointCtx.ctx)
			}

			// when
			err := endpointCtx.Err()

			// then
			for _, wantError := range tc.wantErrors {
				assert.ErrorIs(t, err, wantError)
			}
		})
	}
}

func TestEndpointContext_Value(t *testing.T) {
	testCases := []struct {
		name            string
		contextProvider func() context.Context
		preCondition    func(ctx *Context)
		key             string
		wantValue       string
	}{
		{
			name: "value from parent context",
			contextProvider: func() context.Context {
				return context.WithValue(context.Background(), "anyKey", "anyValue")
			},
			preCondition: func(ctx *Context) {},
			key:          "anyKey",
			wantValue:    "anyValue",
		},
		{
			name: "value from current context",
			contextProvider: func() context.Context {
				return context.Background()
			},
			preCondition: func(ctx *Context) {
				ctx.SetValue("anotherKey", "anotherValue")
			},
			key:       "anotherKey",
			wantValue: "anotherValue",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			ctx := tc.contextProvider()
			nativeReq := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			recorder := httptest.NewRecorder()

			endpointCtx := &EndpointContext[any]{
				ctx: CreateContext(nativeReq, recorder),
			}

			if tc.preCondition != nil {
				tc.preCondition(endpointCtx.ctx)
			}

			// when
			value := endpointCtx.Value(tc.key)

			// then
			assert.Equal(t, tc.wantValue, value)
		})
	}
}

func TestEndpointContext_SetValue(t *testing.T) {
	// given
	nativeReq := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	endpointCtx := &EndpointContext[any]{
		ctx: CreateContext(nativeReq, recorder),
	}

	// when
	endpointCtx.SetValue("anyKey", "anyValue")

	// then
	assert.Contains(t, endpointCtx.ctx.values, "anyKey")

	value := endpointCtx.ctx.values["anyKey"]
	assert.Contains(t, value, "anyValue")
}

func TestEndpointContext_Request(t *testing.T) {
	// given
	nativeReq := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	endpointCtx := &EndpointContext[any]{
		ctx: CreateContext(nativeReq, recorder),
	}

	// when
	request := endpointCtx.Request()

	// then
	assert.NotNil(t, request)
	assert.Equal(t, nativeReq, request.nativeReq)
}

func TestEndpointContext_Response(t *testing.T) {
	// given
	nativeReq := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	endpointCtx := &EndpointContext[any]{
		ctx: CreateContext(nativeReq, recorder),
	}

	// when
	response := endpointCtx.Response()

	// then
	assert.NotNil(t, response)
	assert.Equal(t, recorder, response.writer)
}

func TestEndpointContext_Input(t *testing.T) {
	// given
	anyValue := "anyValue"

	nativeReq := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	endpointCtx := &EndpointContext[any]{
		ctx:   CreateContext(nativeReq, recorder),
		input: anyValue,
	}

	// when
	input := endpointCtx.Input()

	// then
	assert.Equal(t, anyValue, input)
}

func TestEndpointContext_NativeContext(t *testing.T) {
	// given
	nativeReq := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()
	httpCtx := CreateContext(nativeReq, recorder)

	endpointCtx := &EndpointContext[any]{
		ctx: httpCtx,
	}

	// when
	ctx := endpointCtx.NativeContext()

	// then
	assert.Equal(t, httpCtx, ctx)
}

func TestEndpointContext_setContext(t *testing.T) {
	// given
	nativeReq := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()
	httpCtx := CreateContext(nativeReq, recorder)

	endpointCtx := &EndpointContext[any]{}

	// when
	endpointCtx.setContext(httpCtx)

	// then
	assert.Equal(t, httpCtx, endpointCtx.NativeContext())
}
