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
	"fmt"
	"testing"
	"time"

	"codnect.io/procyon/component"
	"codnect.io/procyon/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestContextError_Error(t *testing.T) {
	// given
	err := contextError{
		Op:  "any operation",
		Err: errors.New("any error message"),
	}

	// when
	result := err.Error()

	// then
	assert.Equal(t, "any operation context: any error message", result)
}

func TestContextError_Unwrap(t *testing.T) {
	// given
	innerErr := errors.New("any error message")
	err := contextError{
		Op:  "any operation",
		Err: innerErr,
	}

	// when
	result := err.Unwrap()

	// then
	assert.Equal(t, innerErr, result)
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
					createContext(tc.environment, component.NewStandardContainer())
				})
				return
			}

			ctx := createContext(tc.environment, component.NewStandardContainer())

			// then
			require.NotNil(t, ctx)
		})
	}

}

func TestContext_Deadline(t *testing.T) {
	// given
	env := NewEnvironment()
	ctx := createContext(env, component.NewStandardContainer())

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
			ctx := createContext(env, component.NewStandardContainer())
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
			ctx := createContext(env, component.NewStandardContainer())
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
	ctx := createContext(env, component.NewStandardContainer())

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
			wantErr:      errors.New("start context: context not refreshed"),
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
				lifecycleManager := &anyMockLifecycleManager{}
				lifecycleManager.On("Startup", mock.AnythingOfType("context.backgroundCtx")).
					Return(errors.New("lifecycle component error"))
				lifecycleManager.On("Shutdown", mock.AnythingOfType("context.backgroundCtx")).
					Return(nil)
				lifecycleManager.On("IsRunning").
					Return(false)

				ctx.lifecycleManager = lifecycleManager
			},
			wantErr: errors.New("start context: start lifecycle manager: lifecycle component error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			env := NewEnvironment()
			ctx := createContext(env, component.NewStandardContainer())
			if tc.preCondition != nil {
				tc.preCondition(ctx)
			}

			// when
			err := ctx.Start(context.Background())

			// then
			if tc.wantErr != nil {
				assert.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
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
			wantErr:      errors.New("stop context: context not refreshed"),
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
				lifecycleManager := &anyMockLifecycleManager{}
				lifecycleManager.On("Startup", mock.AnythingOfType("context.backgroundCtx")).
					Return(nil)
				lifecycleManager.On("Shutdown", mock.AnythingOfType("context.backgroundCtx")).
					Return(errors.New("lifecycle component error"))
				lifecycleManager.On("IsRunning").
					Return(true)

				ctx.lifecycleManager = lifecycleManager
			},
			wantErr: errors.New("stop context: stop lifecycle manager: lifecycle component error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			env := NewEnvironment()
			ctx := createContext(env, component.NewStandardContainer())
			if tc.preCondition != nil {
				tc.preCondition(ctx)
			}

			// when
			err := ctx.Stop(context.Background())

			// then
			if tc.wantErr != nil {
				assert.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
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
			ctx := createContext(env, component.NewStandardContainer())
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
		name          string
		preCondition  func(ctx *Context, startupContainer component.Container)
		postCondition func(t *testing.T, ctx *Context)
		wantErr       error
	}{
		{
			name:         "non-refreshed context",
			preCondition: nil,
			wantErr:      nil,
		},
		{
			name: "already closed context",
			preCondition: func(ctx *Context, startupContainer component.Container) {
				err := ctx.Close(context.Background())
				assert.NoError(t, err)
			},
			wantErr: fmt.Errorf("refresh context: context canceled"),
		},
		{
			name: "already stopped context",
			preCondition: func(ctx *Context, startupContainer component.Container) {
				err := ctx.Refresh(context.Background())
				assert.NoError(t, err)
				err = ctx.Stop(context.Background())
				assert.NoError(t, err)
			},
			wantErr: nil,
		},
		{
			name: "lifecycle component start error",
			preCondition: func(ctx *Context, startupContainer component.Container) {
				lifecycleManager := &anyMockLifecycleManager{}
				lifecycleManager.On("Startup", mock.AnythingOfType("context.backgroundCtx")).
					Return(errors.New("lifecycle component error"))
				lifecycleManager.On("Shutdown", mock.AnythingOfType("context.backgroundCtx")).
					Return(nil)
				lifecycleManager.On("IsRunning").
					Return(false)

				ctx.lifecycleManager = lifecycleManager
			},
			wantErr: errors.New("refresh context: start lifecycle manager: lifecycle component error"),
		},
		{
			name: "lifecycle component stop error",
			preCondition: func(ctx *Context, startupContainer component.Container) {
				lifecycleManager := &anyMockLifecycleManager{}
				lifecycleManager.On("Startup", mock.AnythingOfType("context.backgroundCtx")).
					Return(nil)
				lifecycleManager.On("Shutdown", mock.AnythingOfType("context.backgroundCtx")).
					Return(errors.New("lifecycle component error"))
				lifecycleManager.On("IsRunning").
					Return(true)

				ctx.lifecycleManager = lifecycleManager
			},
			wantErr: errors.New("refresh context: stop lifecycle manager: lifecycle component error"),
		},
		{
			name: "multiple lifecycle manager",
			preCondition: func(ctx *Context, startupContainer component.Container) {
				err := startupContainer.RegisterSingleton("anyMockLifecycleManager", &anyMockLifecycleManager{})
				assert.NoError(t, err)

				err = startupContainer.RegisterSingleton("anotherLifecycleManager", &anyMockLifecycleManager{})
				assert.NoError(t, err)
			},
			wantErr: errors.New("resolve type runtime.LifecycleManager: ambiguous match"),
		},
		{
			name: "duplicate registered app context",
			preCondition: func(ctx *Context, startupContainer component.Container) {
				container := component.NewStandardContainer()

				err := container.RegisterSingleton(appContextContainerKey, &Context{})
				assert.NoError(t, err)

				ctx.containerProvider = func() component.Container {
					return container
				}
			},
			wantErr: errors.New("refresh context: register singleton \"procyonAppContext\": duplicate instance"),
		},
		{
			name: "duplicate lifecycle manager",
			preCondition: func(ctx *Context, startupContainer component.Container) {
				container := component.NewStandardContainer()

				err := container.RegisterSingleton(lifecycleManagerContainerKey, &defaultLifecycleManager{})
				assert.NoError(t, err)

				ctx.containerProvider = func() component.Container {
					return container
				}
			},
			wantErr: errors.New("refresh context: register singleton \"procyonLifecycleManager\": duplicate instance"),
		},
		{
			name: "singleton resolve error",
			preCondition: func(ctx *Context, startupContainer component.Container) {
				container := component.NewStandardContainer()

				def, err := component.MakeDefinition(func() *AnyComponent {
					panic("singleton constructor error")
					return &AnyComponent{}
				})
				require.NoError(t, err)

				err = container.RegisterDefinition(def)
				require.NoError(t, err)

				ctx.containerProvider = func() component.Container {
					return container
				}
			},
			wantErr: errors.New("refresh context: initialize singleton \"anyComponent\": resolve \"anyComponent\": invoke constructor \"anyComponent\" (*procyon.AnyComponent): constructor panic: singleton constructor error"),
		},
		{
			name: "load component error",
			preCondition: func(ctx *Context, startupContainer component.Container) {
				def, err := component.MakeDefinition(func() *AnyComponent {
					return &AnyComponent{}
				})
				require.NoError(t, err)

				comp := component.Create(def)

				ctx.components = []*component.Component{comp, comp}
			},
			wantErr: errors.New("refresh context: load component definitions: load component \"anyComponent\": register definition \"anyComponent\": duplicate definition"),
		},
		{
			name: "bootstrap types",
			preCondition: func(ctx *Context, startupContainer component.Container) {
				def, err := component.MakeDefinition(func() runtime.ContextInitializer {
					return nil
				}, component.WithName("contextInitializer"), component.WithScope(component.SingletonScope))
				require.NoError(t, err)

				comp := component.Create(def)
				ctx.components = []*component.Component{comp}
			},
			postCondition: func(t *testing.T, ctx *Context) {
				assert.False(t, ctx.Container().ContainsDefinition("contextInitializer"))
			},
			wantErr: nil,
		},
		{
			name: "refresh successfully",
			preCondition: func(ctx *Context, startupContainer component.Container) {
				container := component.NewStandardContainer()

				def, err := component.MakeDefinition(func() *AnyComponent {
					return &AnyComponent{}
				})
				require.NoError(t, err)

				err = container.RegisterDefinition(def)
				require.NoError(t, err)

				def, err = component.MakeDefinition(func(anyComponent AnyComponent) *AnyComponent {
					return &AnyComponent{}
				}, component.WithName("anotherInstance"), component.WithScope(component.PrototypeScope))

				require.NoError(t, err)

				err = container.RegisterDefinition(def)
				require.NoError(t, err)

				ctx.containerProvider = func() component.Container {
					return container
				}
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			env := NewEnvironment()
			startupContainer := component.NewStandardContainer()

			ctx := createContext(env, startupContainer)

			if tc.preCondition != nil {
				tc.preCondition(ctx, startupContainer)
			}

			// when
			err := ctx.Refresh(context.Background())

			// then
			if tc.postCondition != nil {
				tc.postCondition(t, ctx)
			}

			if tc.wantErr != nil {
				assert.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
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
			ctx := createContext(env, component.NewStandardContainer())
			if tc.preCondition != nil {
				tc.preCondition(ctx)
			}

			// when
			err := ctx.Close(context.Background())

			// then
			if tc.wantErr != nil {
				assert.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestContext_Environment(t *testing.T) {
	// given
	env := NewEnvironment()
	ctx := createContext(env, component.NewStandardContainer())

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
			wantPanic:    errors.New("nil container: context not refreshed"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			env := NewEnvironment()
			ctx := createContext(env, component.NewStandardContainer())
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
