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
	"github.com/stretchr/testify/require"
)

type anyLifecycleManager struct {
	mock.Mock
}

func (a *anyLifecycleManager) Startup(ctx context.Context) error {
	result := a.Called(ctx)
	return result.Error(0)
}

func (a *anyLifecycleManager) Shutdown(ctx context.Context) error {
	result := a.Called(ctx)
	return result.Error(0)
}

func (a *anyLifecycleManager) IsRunning() bool {
	result := a.Called()
	return result.Bool(0)
}

func TestCreateContext(t *testing.T) {
	testCases := []struct {
		name        string
		environment runtime.Environment
		wantPanic   error
	}{
		{
			name:        "nil environment",
			environment: nil,
			wantPanic:   errors.New("nil environment"),
		},
		{
			name:        "valid environment",
			environment: NewEnvironment(),
			wantPanic:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					createContext(tc.environment)
				})
				return
			}

			ctx := createContext(tc.environment)

			// then
			require.NotNil(t, ctx)
		})
	}

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
		wantClosed   bool
	}{
		{
			name:         "non-refreshed context",
			preCondition: nil,
			wantClosed:   false,
		},
		{
			name: "stopped context",
			preCondition: func(ctx *Context) {
				err := ctx.Refresh(context.Background())
				assert.NoError(t, err)
				err = ctx.Stop(context.Background())
				assert.NoError(t, err)
			},
			wantClosed: false,
		},
		{
			name: "already closed context",
			preCondition: func(ctx *Context) {
				err := ctx.Close(context.Background())
				assert.NoError(t, err)
			},
			wantClosed: true,
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
			doneCh := ctx.Done()

			// then
			if tc.wantClosed {
				assert.Eventually(t, func() bool {
					select {
					case <-doneCh:
						return true
					default:
						return false
					}
				}, time.Second, 10*time.Millisecond)
			} else {
				select {
				case <-doneCh:
					t.Fatal("done channel should not be closed")
				default:
					// do nothing
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
			name:         "non-refreshed context",
			preCondition: nil,
			wantErr:      nil,
		},
		{
			name: "stopped context",
			preCondition: func(ctx *Context) {
				err := ctx.Refresh(context.Background())
				assert.NoError(t, err)
				err = ctx.Stop(context.Background())
				assert.NoError(t, err)
			},
			wantErr: nil,
		},
		{
			name: "already closed context",
			preCondition: func(ctx *Context) {
				err := ctx.Close(context.Background())
				assert.NoError(t, err)
			},
			wantErr: context.Canceled,
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
	testCases := []struct {
		name         string
		preCondition func(ctx *Context)
		wantErr      error
	}{
		{
			name:         "non-refreshed context",
			preCondition: nil,
			wantErr:      errors.New("lifecycle manager is not initialized, invoke refresh method before starting the context"),
		},
		{
			name: "refreshed context",
			preCondition: func(ctx *Context) {
				err := ctx.Refresh(context.Background())
				assert.NoError(t, err)
			},
			wantErr: nil,
		},
		{
			name: "already closed context",
			preCondition: func(ctx *Context) {
				err := ctx.Close(context.Background())
				assert.NoError(t, err)
			},
			wantErr: context.Canceled,
		},
		{
			name: "already stopped context",
			preCondition: func(ctx *Context) {
				err := ctx.Refresh(context.Background())
				assert.NoError(t, err)
				err = ctx.Stop(context.Background())
				assert.NoError(t, err)
			},
			wantErr: nil,
		},
		{
			name: "already started context",
			preCondition: func(ctx *Context) {
				err := ctx.Refresh(context.Background())
				assert.NoError(t, err)
				err = ctx.Start(context.Background())
				assert.NoError(t, err)
			},
			wantErr: nil,
		},
		{
			name: "lifecycle component error",
			preCondition: func(ctx *Context) {
				lifecycleManager := &anyLifecycleManager{}
				lifecycleManager.On("Startup", mock.AnythingOfType("context.backgroundCtx")).
					Return(errors.New("lifecycle component error"))
				lifecycleManager.On("Shutdown", mock.AnythingOfType("context.backgroundCtx")).
					Return(nil)
				lifecycleManager.On("IsRunning").
					Return(false)

				ctx.lifecycleManager = lifecycleManager
			},
			wantErr: errors.New("lifecycle component error"),
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
			err := ctx.Start(context.Background())

			// then
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestContext_Stop(t *testing.T) {
	testCases := []struct {
		name         string
		preCondition func(ctx *Context)
		wantErr      error
	}{
		{
			name:         "non-refreshed context",
			preCondition: nil,
			wantErr:      errors.New("lifecycle manager is not initialized, invoke refresh method before stopping the context"),
		},
		{
			name: "refreshed context",
			preCondition: func(ctx *Context) {
				err := ctx.Refresh(context.Background())
				assert.NoError(t, err)
			},
			wantErr: nil,
		},
		{
			name: "already closed context",
			preCondition: func(ctx *Context) {
				err := ctx.Close(context.Background())
				assert.NoError(t, err)
			},
			wantErr: context.Canceled,
		},
		{
			name: "already stopped context",
			preCondition: func(ctx *Context) {
				err := ctx.Refresh(context.Background())
				assert.NoError(t, err)
				err = ctx.Stop(context.Background())
				assert.NoError(t, err)
			},
			wantErr: nil,
		},
		{
			name: "lifecycle component error",
			preCondition: func(ctx *Context) {
				lifecycleManager := &anyLifecycleManager{}
				lifecycleManager.On("Startup", mock.AnythingOfType("context.backgroundCtx")).
					Return(nil)
				lifecycleManager.On("Shutdown", mock.AnythingOfType("context.backgroundCtx")).
					Return(errors.New("lifecycle component error"))
				lifecycleManager.On("IsRunning").
					Return(true)

				ctx.lifecycleManager = lifecycleManager
			},
			wantErr: errors.New("lifecycle component error"),
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
			err := ctx.Stop(context.Background())

			// then
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestContext_IsRunning(t *testing.T) {
	testCases := []struct {
		name          string
		preCondition  func(ctx *Context)
		wantIsRunning bool
	}{
		{
			name:          "non-refreshed context",
			preCondition:  nil,
			wantIsRunning: false,
		},
		{
			name: "refreshed context",
			preCondition: func(ctx *Context) {
				err := ctx.Refresh(context.Background())
				assert.NoError(t, err)
			},
			wantIsRunning: true,
		},
		{
			name: "stopped context",
			preCondition: func(ctx *Context) {
				err := ctx.Refresh(context.Background())
				assert.NoError(t, err)
				err = ctx.Stop(context.Background())
				assert.NoError(t, err)
			},
			wantIsRunning: false,
		},
		{
			name: "closed context",
			preCondition: func(ctx *Context) {
				err := ctx.Close(context.Background())
				assert.NoError(t, err)
			},
			wantIsRunning: false,
		},
		{
			name: "started context",
			preCondition: func(ctx *Context) {
				err := ctx.Refresh(context.Background())
				assert.NoError(t, err)
				err = ctx.Start(context.Background())
				assert.NoError(t, err)
			},
			wantIsRunning: true,
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
			isRunning := ctx.IsRunning()

			// then
			assert.Equal(t, tc.wantIsRunning, isRunning)
		})
	}
}

func TestContext_Refresh(t *testing.T) {
	testCases := []struct {
		name         string
		preCondition func(ctx *Context)
		wantErr      error
	}{
		{
			name:         "non-refreshed context",
			preCondition: nil,
			wantErr:      nil,
		},
		{
			name: "already closed context",
			preCondition: func(ctx *Context) {
				err := ctx.Close(context.Background())
				assert.NoError(t, err)
			},
			wantErr: context.Canceled,
		},
		{
			name: "already stopped context",
			preCondition: func(ctx *Context) {
				err := ctx.Refresh(context.Background())
				assert.NoError(t, err)
				err = ctx.Stop(context.Background())
				assert.NoError(t, err)
			},
			wantErr: nil,
		},
		{
			name: "lifecycle component start error",
			preCondition: func(ctx *Context) {
				lifecycleManager := &anyLifecycleManager{}
				lifecycleManager.On("Startup", mock.AnythingOfType("context.backgroundCtx")).
					Return(errors.New("lifecycle component error"))
				lifecycleManager.On("Shutdown", mock.AnythingOfType("context.backgroundCtx")).
					Return(nil)
				lifecycleManager.On("IsRunning").
					Return(false)

				ctx.lifecycleManager = lifecycleManager
			},
			wantErr: errors.New("lifecycle component error"),
		},
		{
			name: "lifecycle component stop error",
			preCondition: func(ctx *Context) {
				lifecycleManager := &anyLifecycleManager{}
				lifecycleManager.On("Startup", mock.AnythingOfType("context.backgroundCtx")).
					Return(nil)
				lifecycleManager.On("Shutdown", mock.AnythingOfType("context.backgroundCtx")).
					Return(errors.New("lifecycle component error"))
				lifecycleManager.On("IsRunning").
					Return(true)

				ctx.lifecycleManager = lifecycleManager
			},
			wantErr: errors.New("lifecycle component error"),
		},
		{
			name: "multiple lifecycle manager",
			preCondition: func(ctx *Context) {
				ctx.containerProvider = func() component.Container {
					container := component.NewDefaultContainer()
					err := container.RegisterSingleton("anyLifecycleManager", &anyLifecycleManager{})
					assert.NoError(t, err)
					err = container.RegisterSingleton("anotherLifecycleManager", &anyLifecycleManager{})
					assert.NoError(t, err)
					return container
				}
			},
			wantErr: errors.New("multiple singletons found"),
		},
		{
			name: "already registered app context",
			preCondition: func(ctx *Context) {
				ctx.containerProvider = func() component.Container {
					container := component.NewDefaultContainer()
					err := container.RegisterSingleton(appContextContainerKey, &Context{})
					assert.NoError(t, err)
					return container
				}
			},
			wantErr: errors.New("instance already exists"),
		},
		{
			name: "already registered lifecycle manager",
			preCondition: func(ctx *Context) {
				ctx.containerProvider = func() component.Container {
					container := component.NewDefaultContainer()
					err := container.RegisterSingleton(lifecycleManagerContainerKey, &defaultLifecycleManager{})
					assert.NoError(t, err)
					return container
				}
			},
			wantErr: errors.New("instance already exists"),
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
			err := ctx.Refresh(context.Background())

			// then
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestContext_Close(t *testing.T) {
	testCases := []struct {
		name         string
		preCondition func(ctx *Context)
		wantErr      error
	}{
		{
			name:         "non-refreshed context",
			preCondition: func(ctx *Context) {},
			wantErr:      nil,
		},
		{
			name: "refreshed context",
			preCondition: func(ctx *Context) {
				err := ctx.Refresh(context.Background())
				assert.NoError(t, err)
			},
			wantErr: nil,
		},
		{
			name: "already closed context",
			preCondition: func(ctx *Context) {
				err := ctx.Close(context.Background())
				assert.NoError(t, err)
			},
			wantErr: context.Canceled,
		},
		{
			name: "already stopped context",
			preCondition: func(ctx *Context) {
				err := ctx.Refresh(context.Background())
				assert.NoError(t, err)
				err = ctx.Stop(context.Background())
				assert.NoError(t, err)
			},
			wantErr: nil,
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
			err := ctx.Close(context.Background())

			// then
			assert.Equal(t, tc.wantErr, err)
		})
	}
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
	testCases := []struct {
		name         string
		preCondition func(ctx *Context)
		wantPanic    error
	}{
		{
			name: "refresh context before accessing container",
			preCondition: func(ctx *Context) {
				err := ctx.Refresh(context.Background())
				assert.NoError(t, err)
			},
		},
		{
			name:         "access container without refreshing context",
			preCondition: nil,
			wantPanic:    errors.New("container is not initialized, invoke refresh method before accessing the container"),
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
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					ctx.Container()
				})
				return
			}

			container := ctx.Container()

			// then
			assert.NotNil(t, container)
		})
	}
}
