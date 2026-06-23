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
)

type AnyResult struct {
	Status  Status
	Body    any
	Headers Header
}

func (a *AnyResult) StatusCode() Status {
	return a.Status
}

func (a *AnyResult) BodyValue() any {
	return a.Body
}

func (a *AnyResult) Header() Header {
	return a.Headers
}

func TestHandlerFunc_Handle(t *testing.T) {
	anyResult := &AnyResult{}

	testCases := []struct {
		name        string
		handlerFunc HandlerFunc
		wantResult  Result
		wantErr     error
	}{
		{
			name: "handler returns result successfully",
			handlerFunc: func(ctx *Context) (Result, error) {
				return anyResult, nil
			},
			wantResult: anyResult,
			wantErr:    nil,
		},
		{
			name: "handler returns error",
			handlerFunc: func(ctx *Context) (Result, error) {
				return nil, errors.New("handler error")
			},
			wantResult: nil,
			wantErr:    errors.New("handler error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			ctx := &Context{}

			// when
			result, err := tc.handlerFunc.Handle(ctx)

			// then
			assert.Equal(t, tc.wantResult, result)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestHandle_WithoutResult(t *testing.T) {
	testCases := []struct {
		name    string
		fn      func(*Context) error
		wantErr error
	}{
		{
			name: "returns nil error",
			fn:   func(ctx *Context) error { return nil },
		},
		{
			name:    "returns error",
			fn:      func(ctx *Context) error { return errors.New("handler error") },
			wantErr: errors.New("handler error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			handler := Handle(tc.fn)
			ctx := &Context{}

			// when
			result, err := handler.Handle(ctx)

			// then
			assert.Equal(t, tc.wantErr, err)
			assert.Nil(t, result)
		})
	}
}

func TestHandle_WithResult(t *testing.T) {
	anyResult := &AnyResult{}

	testCases := []struct {
		name       string
		fn         func(*Context) (Result, error)
		wantResult Result
		wantErr    error
	}{
		{
			name:       "returns result",
			fn:         func(ctx *Context) (Result, error) { return anyResult, nil },
			wantResult: anyResult,
		},
		{
			name:    "returns error",
			fn:      func(ctx *Context) (Result, error) { return nil, errors.New("handler error") },
			wantErr: errors.New("handler error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			handler := HandleResult(tc.fn)
			ctx := &Context{}

			// when
			result, err := handler.Handle(ctx)

			// then
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantResult, result)
		})
	}
}

func TestHandle_EndpointContextWithoutResult(t *testing.T) {
	testCases := []struct {
		name    string
		fn      func(ctx *EndpointContext[any]) error
		wantErr error
	}{
		{
			name: "returns nil error",
			fn: func(ctx *EndpointContext[any]) error {
				return nil
			},
		},
		{
			name: "returns error",
			fn: func(ctx *EndpointContext[any]) error {
				return errors.New("handler error")
			},
			wantErr: errors.New("handler error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			handler := Handle(tc.fn)
			ctx := &Context{}

			// when
			result, err := handler.Handle(ctx)

			// then
			assert.Equal(t, tc.wantErr, err)
			assert.Nil(t, result)
		})
	}
}

func TestHandle_EndpointContextWithResult(t *testing.T) {
	anyResult := &AnyResult{}

	testCases := []struct {
		name       string
		fn         func(*EndpointContext[any]) (Result, error)
		wantResult Result
		wantErr    error
	}{
		{
			name: "returns result",
			fn: func(ctx *EndpointContext[any]) (Result, error) {
				return anyResult, nil
			},
			wantResult: anyResult,
		},
		{
			name: "returns error",
			fn: func(ctx *EndpointContext[any]) (Result, error) {
				return nil, errors.New("handler error")
			},
			wantErr: errors.New("handler error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			handler := HandleResult(tc.fn)
			ctx := &Context{}

			// when
			result, err := handler.Handle(ctx)

			// then

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantResult, result)
		})
	}
}
