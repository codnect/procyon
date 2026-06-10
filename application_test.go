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
	"errors"
	"os"
	"syscall"
	"testing"
	"time"

	"codnect.io/procyon/component"
	"codnect.io/procyon/io"
	"codnect.io/procyon/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	// given

	// when
	app := New()

	// then
	assert.NotNil(t, app)
}

func TestApplication_SetBannerPrinter(t *testing.T) {

	testCases := []struct {
		name      string
		printer   runtime.BannerPrinter
		wantPanic error
	}{
		{
			name:      "nil printer",
			printer:   nil,
			wantPanic: errors.New("nil printer"),
		},
		{
			name:      "valid printer",
			printer:   NewBannerPrinter(),
			wantPanic: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			app := New()

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					app.SetBannerPrinter(tc.printer)
				})
				return
			}

			app.SetBannerPrinter(tc.printer)

			// then
			assert.Equal(t, tc.printer, app.bannerPrinter)
		})
	}
}

func TestApplication_ResourceResolver(t *testing.T) {
	// given
	app := New()
	customResourceResolver := io.NewDefaultResourceResolver()
	app.resourceResolver = customResourceResolver

	// when
	resolver := app.ResourceResolver()

	// then
	assert.NotNil(t, resolver)
	assert.Equal(t, customResourceResolver, resolver)
}

func TestApplication_Run(t *testing.T) {
	testCases := []struct {
		name         string
		args         []string
		preCondition func(app *Application)
		wantErr      error
	}{
		{
			name:    "wrong argument format",
			args:    []string{"--invalid"},
			wantErr: errors.New("wrong argument format '--invalid'"),
		},
		{
			name:    "valid arguments",
			args:    []string{"--profile=dev"},
			wantErr: nil,
		},
		{
			name: "banner printer error",
			args: []string{},
			preCondition: func(app *Application) {
				bannerPrinter := &AnyMockBannerPrinter{}
				bannerPrinter.On("Print", mock.Anything, mock.Anything).Return(errors.New("banner printer error"))

				app.SetBannerPrinter(bannerPrinter)
			},
			wantErr: errors.New("banner printer error"),
		},
		{
			name: "duplicate procyon app args in container",
			args: []string{},
			preCondition: func(app *Application) {
				container := component.NewStandardContainer()

				args, err := runtime.ParseArgs([]string{"--profile=dev"})
				require.NoError(t, err)

				err = container.RegisterSingleton("procyonAppArgs", args)
				require.NoError(t, err)

				app.startupContainer = container
			},
			wantErr: errors.New("register singleton \"procyonAppArgs\": duplicate instance"),
		},
		{
			name: "environment customizer load error",
			args: []string{},
			preCondition: func(app *Application) {
				def, err := component.MakeDefinition(func() runtime.EnvironmentCustomizer {
					return nil
				})
				require.NoError(t, err)

				comp := component.Create(def)
				app.envCustomizers = []*component.Component{comp}

			},
			wantErr: errors.New("load component \"environmentCustomizer\": not found"),
		},
		{
			name: "environment customization error",
			args: []string{},
			preCondition: func(app *Application) {
				envCustomizer := &AnyMockEnvironmentCustomizer{}
				envCustomizer.On("CustomizeEnvironment", mock.AnythingOfType("*procyon.Environment"),
					mock.AnythingOfType("*procyon.Application")).
					Return(errors.New("environment customization error"))

				def, err := component.MakeDefinition(func() runtime.EnvironmentCustomizer {
					return envCustomizer
				})
				require.NoError(t, err)

				comp := component.Create(def)
				app.envCustomizers = []*component.Component{comp}
				app.envCustomizerLoadFunc = func(name string) (runtime.EnvironmentCustomizer, error) {
					return envCustomizer, nil
				}
			},
			wantErr: errors.New("environment customization error"),
		},
		{
			name: "context initializer load error",
			args: []string{},
			preCondition: func(app *Application) {
				def, err := component.MakeDefinition(func() runtime.ContextInitializer {
					return nil
				})
				require.NoError(t, err)

				comp := component.Create(def)
				app.ctxInitializers = []*component.Component{comp}
			},
			wantErr: errors.New("load component \"contextInitializer\": not found"),
		},
		{
			name: "context initialization error",
			args: []string{},
			preCondition: func(app *Application) {
				contextInitializer := &AnyMockContextInitializer{}
				contextInitializer.On("InitializeContext", mock.AnythingOfType("*procyon.Context")).Return(errors.New("context initialization error"))

				def, err := component.MakeDefinition(func() runtime.ContextInitializer {
					return contextInitializer
				})
				require.NoError(t, err)

				comp := component.Create(def)
				app.ctxInitializers = []*component.Component{comp}
				app.ctxInitializerLoadFunc = func(name string) (runtime.ContextInitializer, error) {
					return contextInitializer, nil
				}
			},
			wantErr: errors.New("context initialization error"),
		},
		{
			name: "refresh context error",
			args: []string{},
			preCondition: func(app *Application) {
				container := component.NewStandardContainer()

				lifecycleManager := &anyMockLifecycleManager{}
				lifecycleManager.On("IsRunning").Return(false)
				lifecycleManager.On("Startup", mock.AnythingOfType("context.backgroundCtx")).Return(errors.New("lifecycle manager startup error"))

				err := container.RegisterSingleton("lifecycleManager", lifecycleManager)
				require.NoError(t, err)

				app.startupContainer = container
			},
			wantErr: errors.New("refresh context: start lifecycle manager: lifecycle manager startup error"),
		},
		{
			name: "resolve command line runner error (singleton)",
			args: []string{},
			preCondition: func(app *Application) {
				container := component.NewStandardContainer()

				definition, err := component.MakeDefinition(newAnyMockCommandLinerRunner)
				require.NoError(t, err)

				err = container.RegisterDefinition(definition)
				require.NoError(t, err)

				app.startupContainer = container
			},
			wantErr: errors.New("resolve \"anyMockCommandLineRunner\": create \"anyMockCommandLineRunner\" (*procyon.AnyMockCommandLineRunner): unsatisfied dependency for argument 0 (procyon.AnyComponent): resolve type procyon.AnyComponent: not found"),
		},
		{
			name: "resolve command line runner error (prototype)",
			args: []string{},
			preCondition: func(app *Application) {
				container := component.NewStandardContainer()

				definition, err := component.MakeDefinition(newAnyMockCommandLinerRunner, component.WithScope(component.PrototypeScope))
				require.NoError(t, err)

				err = container.RegisterDefinition(definition)
				require.NoError(t, err)

				app.startupContainer = container
			},
			wantErr: errors.New("resolve \"anyMockCommandLineRunner\": create \"anyMockCommandLineRunner\" (*procyon.AnyMockCommandLineRunner): unsatisfied dependency for argument 0 (procyon.AnyComponent): resolve type procyon.AnyComponent: not found"),
		},
		{
			name: "command line runner error",
			args: []string{},
			preCondition: func(app *Application) {
				container := component.NewStandardContainer()

				runner := &AnyMockCommandLineRunner{}
				runner.On("Run", mock.Anything, mock.Anything).Return(errors.New("command line runner error"))

				err := container.RegisterSingleton("anyCommandLineRunner", runner)
				require.NoError(t, err)

				app.startupContainer = container
			},
			wantErr: errors.New("command line runner error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			app := New()

			if tc.preCondition != nil {
				tc.preCondition(app)
			}

			// when
			err := app.Run(tc.args...)

			// then
			if tc.wantErr != nil {
				assert.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestApplication_Run_ServerApplication(t *testing.T) {
	testCases := []struct {
		name     string
		shutdown func(app *Application)
	}{
		{
			name: "shutdown via SIGINT",
			shutdown: func(_ *Application) {
				p, err := os.FindProcess(os.Getpid())
				if err == nil {
					_ = p.Signal(syscall.SIGINT)
				}
			},
		},
		{
			name: "shutdown via SIGTERM",
			shutdown: func(_ *Application) {
				p, err := os.FindProcess(os.Getpid())
				if err == nil {
					_ = p.Signal(syscall.SIGTERM)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			app := New()

			container := component.NewStandardContainer()
			serverApp := &AnyMockServerApp{}
			err := container.RegisterSingleton("anyServerApp", serverApp)
			require.NoError(t, err)

			app.startupContainer = container

			errCh := make(chan error, 1)

			// when
			go func() {
				errCh <- app.Run()
			}()

			time.Sleep(100 * time.Millisecond)
			tc.shutdown(app)

			// then
			select {
			case err = <-errCh:
				assert.NoError(t, err)
			case <-time.After(3 * time.Second):
				t.Fatal("Run did not return in time")
			}
		})
	}
}
