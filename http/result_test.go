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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type AnyResult struct {
	StatusCode Status
	Body       any
	Headers    Header
}

func (a *AnyResult) Status() Status {
	return a.StatusCode
}

func (a *AnyResult) Value() any {
	return a.Body
}

func (a *AnyResult) Header() Header {
	return a.Headers
}

type AnyResultExecutor struct {
	mock.Mock
}

func (e *AnyResultExecutor) CanExecute(result Result) bool {
	r := e.Called(result)
	return r.Bool(0)
}

func (e *AnyResultExecutor) Execute(ctx *Context, result Result) error {
	r := e.Called(ctx, result)
	return r.Error(0)
}

func TestTypedResult_StatusCode(t *testing.T) {
	testCases := []struct {
		name       string
		status     Status
		wantStatus Status
	}{
		{
			name:       "No Status",
			wantStatus: StatusOK,
		},
		{
			name:       "OK",
			status:     StatusOK,
			wantStatus: StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			typedResult := TypedResult[any]{
				Status: tc.status,
			}

			// code
			assert.Equal(t, tc.wantStatus, typedResult.StatusCode())
		})
	}
}

func TestTypedResult_BodyValue(t *testing.T) {
	// given
	body := any("anyBody")

	// when
	typedResult := TypedResult[any]{
		Body: body,
	}

	// then
	assert.Equal(t, body, typedResult.BodyValue())
}

func TestTypedResult_Header(t *testing.T) {
	// given
	headers := Header{}

	// when
	typedResult := TypedResult[any]{
		Headers: headers,
	}

	// then
	assert.Equal(t, headers, typedResult.Header())
}

func TestResultExecutorRegistry_Register(t *testing.T) {
	resultExecutor := &AnyResultExecutor{}

	testCases := []struct {
		name          string
		executor      ResultExecutor
		wantExecutors []ResultExecutor
		wantErr       string
	}{
		{
			name:          "nil result executor",
			executor:      nil,
			wantExecutors: []ResultExecutor{},
			wantErr:       "nil result executor",
		},
		{
			name:          "valid result executor",
			executor:      resultExecutor,
			wantExecutors: []ResultExecutor{resultExecutor},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			registry := NewResultExecutorRegistry()

			// when
			err := registry.Register(tc.executor)

			// then
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
				require.Equal(t, tc.wantExecutors, registry.executors)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.wantExecutors, registry.executors)
		})
	}
}

func TestResultExecutorRegistry_Resolve(t *testing.T) {
	resultExecutor := &AnyResultExecutor{}
	anyResult := &AnyResult{}

	testCases := []struct {
		name         string
		executors    []ResultExecutor
		result       Result
		wantExecutor ResultExecutor
		wantFound    bool
	}{
		{
			name:         "no registered executors",
			executors:    nil,
			result:       anyResult,
			wantExecutor: nil,
			wantFound:    false,
		},
		{
			name: "matching result executor",
			executors: []ResultExecutor{
				resultExecutor,
			},
			result:       anyResult,
			wantExecutor: resultExecutor,
			wantFound:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			registry := NewResultExecutorRegistry()
			for _, executor := range tc.executors {
				require.NoError(t, registry.Register(executor))
			}

			// when
			executor, ok := registry.Resolve(tc.result)

			// then
			require.Equal(t, tc.wantFound, ok)
			require.Equal(t, tc.wantExecutor, executor)
		})
	}
}
