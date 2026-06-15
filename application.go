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
	"os"
	"os/signal"
	goruntime "runtime"
	"syscall"
	"time"

	"codnect.io/procyon/component"
	"codnect.io/procyon/io"
	"codnect.io/procyon/runtime"
	"codnect.io/procyon/runtime/config"
)

const (
	// procyonArgsContainerKey is the key used to register the procyon args in the component container.
	procyonArgsContainerKey = "procyonAppArgs"
)

func init() {
	// runtime/config
	component.Register(config.NewYamlPropertySourceLoader)
	// main
	component.Register(newConfigEnvCustomizer)
}

// Application is the main entry point of the Procyon framework. It is responsible for initializing the application
// context, loading resources, and running the application.
type Application struct {
	bannerPrinter    runtime.BannerPrinter
	resourceResolver io.ResourceResolver

	startupContainer component.Container
	runtimeCtx       runtime.Context
	env              runtime.Environment

	envCustomizers  []*component.Component
	ctxInitializers []*component.Component

	envCustomizerLoadFunc  func(name string) (runtime.EnvironmentCustomizer, error)
	ctxInitializerLoadFunc func(name string) (runtime.ContextInitializer, error)
}

// New creates a new instance of the application with default banner printer and resource resolver.
func New() *Application {
	return &Application{
		bannerPrinter:    NewBannerPrinter(),
		resourceResolver: io.NewDefaultResourceResolver(),
		startupContainer: component.NewStandardContainer(),
		envCustomizers:   component.ListOf[runtime.EnvironmentCustomizer](),
		ctxInitializers:  component.ListOf[runtime.ContextInitializer](),
		envCustomizerLoadFunc: func(name string) (runtime.EnvironmentCustomizer, error) {
			return component.Load[runtime.EnvironmentCustomizer](name)
		},
		ctxInitializerLoadFunc: func(name string) (runtime.ContextInitializer, error) {
			return component.Load[runtime.ContextInitializer](name)
		},
	}
}

// SetBannerPrinter sets the banner printer to be used by the application to print the banner at startup.
func (a *Application) SetBannerPrinter(printer runtime.BannerPrinter) {
	if printer == nil {
		panic("nil printer")
	}

	a.bannerPrinter = printer
}

// ResourceResolver returns the resource resolver used by the application to load resources.
func (a *Application) ResourceResolver() io.ResourceResolver {
	return a.resourceResolver
}

// Run starts the application with the given command-line arguments. It initializes the environment, prepares
// the application context, and invokes any command-line runners defined in the application context.
func (a *Application) Run(args ...string) (err error) {
	startTime := time.Now()

	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = v
			default:
				err = fmt.Errorf("%v", v)
			}
		}

		err = a.close(err)
	}()

	var rArgs *runtime.Args
	rArgs, err = runtime.ParseArgs(args)
	if err != nil {
		return
	}

	a.env, err = a.prepareEnvironment(rArgs)
	if err != nil {
		return
	}

	err = a.bannerPrinter.Print(a.env, os.Stdout)
	if err != nil {
		return
	}

	signalCtx, stopSignals := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopSignals()

	a.runtimeCtx, err = a.prepareRuntimeContext(rArgs)
	if err != nil {
		return
	}

	startupCtx := context.Background()
	err = a.runtimeCtx.Refresh(startupCtx)

	if err != nil {
		return
	}

	timeTakenToStartup := time.Now().Sub(startTime)
	log.Info("Started application in {} seconds", timeTakenToStartup.Seconds())

	err = a.invokeCmdLineRunners(rArgs)
	if err != nil {
		return
	}

	if a.isServerApplication() {
		<-signalCtx.Done()
	}

	return
}

// close handles application shutdown by recovering from panics, logging run failures,
// and closing the runtime context if it is still running.
func (a *Application) close(err error) error {
	if err != nil {
		log.Error("Application run failed", err)
	}

	if a.runtimeCtx != nil && a.runtimeCtx.IsRunning() {
		closeErr := a.runtimeCtx.Close(context.Background())
		if closeErr != nil {
			log.Error("Application context close failed {}", closeErr)
			err = errors.Join(err, closeErr)
		}
	}

	return err
}

// invokeCmdLineRunners retrieves all CommandLineRunner components from the application context and executes
// them with the provided command-line arguments.
func (a *Application) invokeCmdLineRunners(args *runtime.Args) error {
	runners, err := component.ResolveAll[runtime.CommandLineRunner](a.runtimeCtx, a.runtimeCtx.Container())
	if err != nil {
		return err
	}

	for _, runner := range runners {
		err = runner.Run(a.runtimeCtx, args)
		if err != nil {
			return err
		}
	}

	return nil
}

// prepareEnvironment initializes the application environment by creating a new environment instance, adding
// property sources for command-line arguments and environment variables, and allowing customizers to modify
// the environment.
func (a *Application) prepareEnvironment(args *runtime.Args) (runtime.Environment, error) {
	env := NewEnvironment()

	propertySources := env.PropertySources()
	propertySources.PushFront(runtime.NewArgsPropertySource(args))
	propertySources.PushBack(runtime.NewEnvPropertySource())

	err := a.customizeEnv(env)
	if err != nil {
		return nil, err
	}

	return env, nil
}

// customizeEnv retrieves all EnvironmentCustomizer components and invokes their CustomizeEnvironment method to allow
// them to modify the environment before it is used by the application.
func (a *Application) customizeEnv(env runtime.Environment) error {
	customizers, err := a.loadEnvCustomizers()
	if err != nil {
		return err
	}

	for _, customizer := range customizers {
		err = customizer.CustomizeEnvironment(env, a)
		if err != nil {
			return err
		}
	}

	return nil
}

// prepareRuntimeContext creates the application context, allows customizers to modify it, and registers
// the command-line arguments in the context's container.
func (a *Application) prepareRuntimeContext(args *runtime.Args) (runtime.Context, error) {
	runtimeCtx := createContext(a.env, a.startupContainer)

	err := a.initializeRuntimeContext(runtimeCtx)
	if err != nil {
		return nil, err
	}

	log.Info("Starting application using Go {} ({}/{})", goruntime.Version()[2:], goruntime.GOOS, goruntime.GOARCH)
	log.Info("Running with Procyon {}", Version)

	err = a.startupContainer.RegisterSingleton(procyonArgsContainerKey, args)
	if err != nil {
		return nil, err
	}

	return runtimeCtx, nil
}

// initializeRuntimeContext retrieves all ContextInitializer components and invokes their InitializeContext method to allow
// them to modify the application context before it is refreshed and used by the application.
func (a *Application) initializeRuntimeContext(runtimeCtx *Context) error {
	customizers, err := a.loadCtxInitializers()
	if err != nil {
		return err
	}

	for _, customizer := range customizers {
		err = customizer.InitializeContext(runtimeCtx)
		if err != nil {
			return err
		}
	}

	return nil
}

// isServerApplication checks if the application is a server application by checking if there is a Server component
// registered in the application context.
func (a *Application) isServerApplication() bool {
	return component.CanResolveType[runtime.Server](a.runtimeCtx.Container())
}

// loadEnvCustomizers loads all EnvironmentCustomizer components from the application and returns them as a slice.
func (a *Application) loadEnvCustomizers() ([]runtime.EnvironmentCustomizer, error) {
	var customizers []runtime.EnvironmentCustomizer
	for _, comp := range a.envCustomizers {
		c, err := a.envCustomizerLoadFunc(comp.Definition().Name())
		if err != nil {
			return nil, err
		}

		customizers = append(customizers, c)
	}
	return customizers, nil
}

// loadCtxInitializers loads all ContextInitializer components from the application and returns them as a slice.
func (a *Application) loadCtxInitializers() ([]runtime.ContextInitializer, error) {
	var initializers []runtime.ContextInitializer
	for _, comp := range a.ctxInitializers {
		c, err := a.ctxInitializerLoadFunc(comp.Definition().Name())
		if err != nil {
			return nil, err
		}

		initializers = append(initializers, c)
	}

	return initializers, nil
}
