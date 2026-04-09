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

package procyon

import (
	"context"
	"errors"
	"testing"
	"time"

	"codnect.io/procyon/component"
	"codnect.io/procyon/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type anyLifecycle struct {
	mock.Mock
}

func (a *anyLifecycle) Start(ctx context.Context) error {
	results := a.Called(ctx)
	return results.Error(0)
}

func (a *anyLifecycle) Stop(ctx context.Context) error {
	results := a.Called(ctx)
	return results.Error(0)
}

func (a *anyLifecycle) IsRunning() bool {
	results := a.Called()
	return results.Bool(0)
}

func TestContext_Deadline(t *testing.T) {
	// given
	env := NewEnvironment()
	ctx := createContext(env)

	// when
	deadline, ok := ctx.Deadline()

	// then
	assert.False(t, ok)
	assert.Zero(t, deadline)
}

func TestContext_Done(t *testing.T) {
	testCases := []struct {
		name         string
		preCondition func(ctx *Context)
		wantDone     bool
	}{
		{
			name: "context is closed",
			preCondition: func(ctx *Context) {
				err := ctx.Close(context.Background())
				assert.NoError(t, err)
			},
			wantDone: true,
		},
		{
			name: "context is stopped",
			preCondition: func(ctx *Context) {
				err := ctx.Stop(context.Background())
				assert.NoError(t, err)
			},
			wantDone: false,
		},
		{
			name:         "context is running",
			preCondition: nil,
			wantDone:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			env := NewEnvironment()
			ctx := createContext(env)
			if tc.preCondition != nil {
				tc.preCondition(ctx)
			}

			// when
			select {
			case <-ctx.Done():
				// then
				assert.True(t, true, "Done channel was closed as expected")
			case <-time.After(2 * time.Second):
				if tc.wantDone {
					assert.Fail(t, "Done channel was not closed within timeout")
				}
			}
		})
	}
}

func TestContext_Err(t *testing.T) {
	testCases := []struct {
		name         string
		preCondition func(ctx *Context)
		wantErr      error
	}{
		{
			name: "context is closed",
			preCondition: func(ctx *Context) {
				err := ctx.Close(context.Background())
				assert.NoError(t, err)
			},
			wantErr: context.Canceled,
		},
		{
			name: "stop context before invoking refresh",
			preCondition: func(ctx *Context) {
				err := ctx.Stop(context.Background())
				assert.NoError(t, err)
			},
			wantErr: errors.New("lifecycle manager is not initialized, invoke refresh method before stoping the context"),
		},
		{
			name:         "context is running",
			preCondition: nil,
			wantErr:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			env := NewEnvironment()
			ctx := createContext(env)
			if tc.preCondition != nil {
				tc.preCondition(ctx)
			}

			// when
			err := ctx.Err()

			// then
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestContext_Value(t *testing.T) {
	// given
	env := NewEnvironment()
	ctx := createContext(env)

	// when
	value := ctx.Value("anyKey")

	// then
	assert.Nil(t, value)
}

func TestContext_Start(t *testing.T) {
	// given
	lifecycle := &anyLifecycle{}
	component.Register(func() runtime.Lifecycle {
		return lifecycle
	})

	env := NewEnvironment()
	ctx := createContext(env)

	lifecycle.On("Stop", mock.AnythingOfType("context.backgroundCtx")).
		Return(nil).
		Once()
	lifecycle.On("Start", mock.AnythingOfType("context.backgroundCtx")).
		Return(nil).
		Once()

	err := ctx.Refresh(context.Background())
	if err != nil {
		assert.Fail(t, "Failed to refresh context", err)
		return
	}

	err = ctx.Stop(context.Background())
	if err != nil {
		assert.Fail(t, "Failed to stop context", err)
		return
	}

	// when
	err = ctx.Start(context.Background())

	// then
	assert.NoError(t, err)
}

func TestContext_Stop(t *testing.T) {
	// given
	env := NewEnvironment()
	ctx := createContext(env)

	// when
	err := ctx.Stop(context.Background())

	// then
	assert.NoError(t, err)
}

func TestContext_IsRunning(t *testing.T) {
	// given
	env := NewEnvironment()
	ctx := createContext(env)

	// when
	running := ctx.IsRunning()

	// then
	assert.False(t, running)
}

func TestContext_Refresh(t *testing.T) {
	// given
	env := NewEnvironment()
	ctx := createContext(env)

	// when
	err := ctx.Refresh(context.Background())

	// then
	assert.NoError(t, err)
}

func TestContext_Close(t *testing.T) {
	// given
	env := NewEnvironment()
	ctx := createContext(env)

	// when
	err := ctx.Close(context.Background())

	// then
	assert.NoError(t, err)
}

func TestContext_Environment(t *testing.T) {
	// given
	env := NewEnvironment()
	ctx := createContext(env)

	// when
	result := ctx.Environment()

	// then
	assert.Equal(t, env, result)
}

func TestContext_Container(t *testing.T) {
	// given
	env := NewEnvironment()
	ctx := createContext(env)

	// when
	container := ctx.Container()

	// then
	assert.NotNil(t, container)
}
