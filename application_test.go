package procyon

import (
	"context"
	"errors"
	stdio "io"
	"testing"

	"codnect.io/procyon/component"
	"codnect.io/procyon/io"
	"codnect.io/procyon/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type AnyUnknownComponent struct {
}

type AnyBannerPrinter struct {
	mock.Mock
}

func (p *AnyBannerPrinter) Print(env runtime.Environment, writer stdio.Writer) error {
	result := p.Called(env, writer)
	return result.Error(0)
}

type AnyCommandLineRunner struct {
	mock.Mock
}

func newAnyCommandLinerRunner(unknownComponent AnyUnknownComponent) *AnyCommandLineRunner {
	return &AnyCommandLineRunner{}
}

func (r *AnyCommandLineRunner) Run(ctx runtime.Context, args *runtime.Args) error {
	results := r.Called(ctx, args)
	return results.Error(0)
}

type AnyEnvironmentCustomizer struct {
	mock.Mock
}

func (c *AnyEnvironmentCustomizer) CustomizeEnvironment(env runtime.Environment, app runtime.Application) error {
	results := c.Called(env, app)
	return results.Error(0)
}

type AnyContextCustomizer struct {
	mock.Mock
}

func (a *AnyContextCustomizer) InitializeContext(ctx runtime.Context) error {
	results := a.Called(ctx)
	return results.Error(0)
}

type AnyServerApp struct {
	mock.Mock
}

func (a *AnyServerApp) Start(ctx context.Context) error {
	results := a.Called(ctx)
	return results.Error(0)
}

func (a *AnyServerApp) Stop(ctx context.Context) error {
	results := a.Called(ctx)
	return results.Error(0)
}

func (a *AnyServerApp) Port() int {
	results := a.Called()
	return results.Int(0)
}

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
				bannerPrinter := &AnyBannerPrinter{}
				bannerPrinter.On("Print", mock.Anything, mock.Anything).Return(errors.New("banner printer error"))

				app.SetBannerPrinter(bannerPrinter)
			},
			wantErr: errors.New("banner printer error"),
		},
		{
			name: "duplicate procyon app args in container",
			args: []string{},
			preCondition: func(app *Application) {
				container := component.NewDefaultContainer(nil)

				args, err := runtime.ParseArgs([]string{"--profile=dev"})
				require.NoError(t, err)

				err = container.RegisterSingleton("procyonAppArgs", args)
				require.NoError(t, err)

				app.startupContainer = container
			},
			wantErr: errors.New("register singleton \"procyonAppArgs\": duplicate instance"),
		},
		{
			name: "resolve environment customizer error",
			args: []string{},
			preCondition: func(app *Application) {
				app.resolveEnvCustomizers = func() ([]runtime.EnvironmentCustomizer, error) {
					return nil, errors.New("resolve environment customizer error")
				}
			},
			wantErr: errors.New("resolve environment customizer error"),
		},
		{
			name: "environment customizer error",
			args: []string{},
			preCondition: func(app *Application) {
				envCustomizer := &AnyEnvironmentCustomizer{}
				envCustomizer.On("CustomizeEnvironment", mock.Anything, mock.Anything).Return(errors.New("environment customizer error"))

				app.resolveEnvCustomizers = func() ([]runtime.EnvironmentCustomizer, error) {
					return []runtime.EnvironmentCustomizer{envCustomizer}, nil
				}
			},
			wantErr: errors.New("environment customizer error"),
		},
		{
			name: "resolve context customizer error",
			args: []string{},
			preCondition: func(app *Application) {
				app.resolveCtxInitializers = func() ([]runtime.ContextInitializer, error) {
					return nil, errors.New("resolve context customizer error")
				}
			},
			wantErr: errors.New("resolve context customizer error"),
		},
		{
			name: "context customizer error",
			args: []string{},
			preCondition: func(app *Application) {
				ctxCustomizer := &AnyContextCustomizer{}
				ctxCustomizer.On("InitializeContext", mock.Anything).Return(errors.New("context customizer error"))

				app.resolveCtxInitializers = func() ([]runtime.ContextInitializer, error) {
					return []runtime.ContextInitializer{ctxCustomizer}, nil
				}
			},
			wantErr: errors.New("context customizer error"),
		},
		{
			name: "resolve command line runner error (singleton)",
			args: []string{},
			preCondition: func(app *Application) {
				container := component.NewDefaultContainer(nil)

				definition, err := component.MakeDefinition(newAnyCommandLinerRunner)
				require.NoError(t, err)

				err = container.RegisterDefinition(definition)
				require.NoError(t, err)

				app.startupContainer = container
			},
			wantErr: errors.New("create \"anyCommandLineRunner\" (*procyon.AnyCommandLineRunner): unsatisfied dependency for argument 0 (procyon.AnyUnknownComponent): resolve type procyon.AnyUnknownComponent: not found"),
		},
		{
			name: "resolve command line runner error (prototype)",
			args: []string{},
			preCondition: func(app *Application) {
				container := component.NewDefaultContainer(nil)

				definition, err := component.MakeDefinition(newAnyCommandLinerRunner, component.WithScope(component.PrototypeScope))
				require.NoError(t, err)

				err = container.RegisterDefinition(definition)
				require.NoError(t, err)

				app.startupContainer = container
			},
			wantErr: errors.New("create \"anyCommandLineRunner\" (*procyon.AnyCommandLineRunner): unsatisfied dependency for argument 0 (procyon.AnyUnknownComponent): resolve type procyon.AnyUnknownComponent: not found"),
		},
		{
			name: "command line runner error",
			args: []string{},
			preCondition: func(app *Application) {
				container := component.NewDefaultContainer(nil)

				runner := &AnyCommandLineRunner{}
				runner.On("Run", mock.Anything, mock.Anything).Return(errors.New("command line runner error"))

				err := container.RegisterSingleton("anyCommandLineRunner", runner)
				require.NoError(t, err)

				app.startupContainer = container
			},
			wantErr: errors.New("command line runner error"),
		},
		/*
			{
				name: "server application",
				preCondition: func(app *Application) {
					container := component.NewDefaultContainer(nil)

					serverApp := &AnyServerApp{}
					serverApp.On("Start", mock.Anything).Return(nil)
					serverApp.On("Stop", mock.Anything).Return(nil)
					serverApp.On("Port").Return(8080)

					err := container.RegisterSingleton("anyServerApp", serverApp)
					require.NoError(t, err)

					app.startupContainer = container
				},
			},*/
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
