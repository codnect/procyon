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

func TestDefaultLifecycleManager_Startup(t *testing.T) {
	testCases := []struct {
		name         string
		ctx          context.Context
		preCondition func(container component.Container)

		wantErr error
	}{
		{
			name:    "nil context",
			ctx:     nil,
			wantErr: errors.New("nil context"),
		},
		{
			name: "lifecycle component resolve error",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				def, err := component.MakeDefinition(func(anyComponent AnyComponent) *AnyMockLifecycle {
					return &AnyMockLifecycle{}
				})
				require.NoError(t, err)

				err = container.RegisterDefinition(def)
				require.NoError(t, err)
			},
			wantErr: errors.New("resolve \"anyMockLifecycle\": create \"anyMockLifecycle\" (*procyon.AnyMockLifecycle): unsatisfied dependency for argument 0 (procyon.AnyComponent): resolve type procyon.AnyComponent: not found"),
		},
		{
			name: "start lifecycle components successfully",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				anyLifecycle := &AnyMockLifecycle{}
				anyLifecycle.On("Start", mock.AnythingOfType("context.backgroundCtx")).Return(nil)

				def, err := component.MakeDefinition(func() *AnyMockLifecycle {
					return anyLifecycle
				})
				require.NoError(t, err)

				err = container.RegisterDefinition(def)
				require.NoError(t, err)
			},
			wantErr: nil,
		},
		{
			name: "start lifecycle component with error",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				anyLifecycle := &AnyMockLifecycle{}
				anyLifecycle.On("Start", mock.AnythingOfType("context.backgroundCtx")).Return(errors.New("start error"))

				def, err := component.MakeDefinition(func() *AnyMockLifecycle {
					return anyLifecycle
				})
				require.NoError(t, err)

				err = container.RegisterDefinition(def)
				require.NoError(t, err)
			},
			wantErr: errors.New("start lifecycle component \"anyMockLifecycle\": start error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := component.NewStandardContainer()
			lifecycleManager := newDefaultLifecycleManager(container)

			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			// when
			err := lifecycleManager.Startup(tc.ctx)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestDefaultLifecycleManager_Shutdown(t *testing.T) {
	testCases := []struct {
		name         string
		ctx          context.Context
		preCondition func(container component.Container)
		wantErr      error
	}{
		{
			name:    "nil context",
			ctx:     nil,
			wantErr: errors.New("nil context"),
		},
		{
			name: "stop lifecycle components successfully",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				anyLifecycle := &AnyMockLifecycle{}
				anyLifecycle.On("Start", mock.AnythingOfType("context.backgroundCtx")).Return(nil)
				anyLifecycle.On("Stop", mock.AnythingOfType("*context.timerCtx")).Return(nil)

				def, err := component.MakeDefinition(func() *AnyMockLifecycle {
					return anyLifecycle
				})
				require.NoError(t, err)

				err = container.RegisterDefinition(def)
				require.NoError(t, err)
			},
			wantErr: nil,
		},
		{
			name: "stop lifecycle component with error",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				anyLifecycle := &AnyMockLifecycle{}
				anyLifecycle.On("Start", mock.AnythingOfType("context.backgroundCtx")).Return(nil)
				anyLifecycle.On("Stop", mock.AnythingOfType("*context.timerCtx")).Return(errors.New("stop error"))

				def, err := component.MakeDefinition(func() *AnyMockLifecycle {
					return anyLifecycle
				})
				require.NoError(t, err)

				err = container.RegisterDefinition(def)
				require.NoError(t, err)
			},
			wantErr: nil,
		},
		{
			name: "shutdown timeout exceeded",
			ctx:  context.Background(),
			preCondition: func(container component.Container) {
				anyLifecycle := &AnyMockLifecycle{}
				anyLifecycle.On("Start", mock.AnythingOfType("context.backgroundCtx")).Return(nil)
				anyLifecycle.On("Stop", mock.Anything).Run(func(args mock.Arguments) {
					ctx := args.Get(0).(context.Context)
					<-ctx.Done()
				}).Return(context.DeadlineExceeded)

				def, err := component.MakeDefinition(func() *AnyMockLifecycle {
					return anyLifecycle
				})
				require.NoError(t, err)

				err = container.RegisterDefinition(def)
				require.NoError(t, err)
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := component.NewStandardContainer()
			lifecycleManager := newDefaultLifecycleManager(container)
			lifecycleManager.shutdownTimeout = 500 * time.Millisecond

			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			if tc.ctx != nil {
				err := lifecycleManager.Startup(tc.ctx)
				require.NoError(t, err)
			}

			// when
			err := lifecycleManager.Shutdown(tc.ctx)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestDefaultLifecycleManager_IsRunning(t *testing.T) {
	testCases := []struct {
		name         string
		ctx          context.Context
		preCondition func(lifecycleManager runtime.LifecycleManager)
		wantResult   bool
	}{
		{
			name:       "not running",
			wantResult: false,
		},
		{
			name: "running after startup",
			preCondition: func(lifecycleManager runtime.LifecycleManager) {
				err := lifecycleManager.Startup(context.Background())
				require.NoError(t, err)
			},
			wantResult: true,
		},
		{
			name: "not running after shutdown",
			preCondition: func(lifecycleManager runtime.LifecycleManager) {
				err := lifecycleManager.Startup(context.Background())
				require.NoError(t, err)

				err = lifecycleManager.Shutdown(context.Background())
				require.NoError(t, err)
			},
			wantResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := component.NewStandardContainer()
			lifecycleManager := newDefaultLifecycleManager(container)

			if tc.preCondition != nil {
				tc.preCondition(lifecycleManager)
			}

			// when
			isRunning := lifecycleManager.IsRunning()

			// then
			assert.Equal(t, tc.wantResult, isRunning)

			// Cleanup if still running
			if isRunning {
				err := lifecycleManager.Shutdown(context.Background())
				require.NoError(t, err)
				assert.False(t, lifecycleManager.IsRunning())
			}
		})
	}
}
