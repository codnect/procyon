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
	"os"
	"os/signal"
	goruntime "runtime"
	"syscall"
	"time"

	"codnect.io/procyon/component"
	"codnect.io/procyon/io"

	"codnect.io/procyon/runtime"
)

// Application is the main entry point of the Procyon framework. It is responsible for initializing the application
// context, loading resources, and running the application.
type Application struct {
	bannerPrinter    runtime.BannerPrinter
	resourceResolver io.ResourceResolver

	runtimeCtx runtime.Context
	env        runtime.Environment
}

// New creates a new instance of the application with default banner printer and resource resolver.
func New() *Application {
	return &Application{
		bannerPrinter:    NewBannerPrinter(),
		resourceResolver: io.NewDefaultResourceResolver(),
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
func (a *Application) Run(args ...string) error {
	startTime := time.Now()

	rArgs, err := runtime.ParseArgs(args)
	if err != nil {
		return err
	}

	a.env, err = prepareEnvironment(rArgs, a)

	err = a.bannerPrinter.Print(a.env, os.Stdout)
	if err != nil {
		return err
	}

	signalCtx, stopSignals := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopSignals()

	a.runtimeCtx, err = prepareRuntimeContext(a.env, rArgs)
	if err != nil {
		return err
	}

	startupCtx := context.Background()
	err = a.runtimeCtx.Refresh(startupCtx)

	if err != nil {
		return err
	}

	timeTakenToStartup := time.Now().Sub(startTime)
	log.Info("Started application in {} seconds", timeTakenToStartup.Seconds())

	err = invokeCmdLineRunners(a.runtimeCtx, a.runtimeCtx.Container(), rArgs)
	if err != nil {
		return err
	}

	if isServerApplication() {
		<-signalCtx.Done()
	}

	if a.runtimeCtx.IsRunning() {
		closeCtx := context.Background()
		return a.runtimeCtx.Close(closeCtx)
	}

	return nil
}

// invokeCmdLineRunners retrieves all CommandLineRunner components from the application context and executes
// them with the provided command-line arguments.
func invokeCmdLineRunners(ctx runtime.Context, container component.Container, args *runtime.Args) error {
	runners, err := component.ResolveAll[runtime.CommandLineRunner](ctx, container)
	if err != nil {
		return err
	}

	for _, runner := range runners {
		err = runner.Run(ctx, args)
		if err != nil {
			return err
		}
	}

	return nil
}

// prepareEnvironment initializes the application environment by creating a new environment instance, adding
// property sources for command-line arguments and environment variables, and allowing customizers to modify
// the environment.
func prepareEnvironment(args *runtime.Args, app runtime.Application) (runtime.Environment, error) {
	env := NewEnvironment()

	propertySources := env.PropertySources()
	propertySources.PushFront(runtime.NewArgsPropertySource(args))
	propertySources.PushBack(runtime.NewEnvPropertySource())

	err := customizeEnv(env, app)
	if err != nil {
		return nil, err
	}

	return env, nil
}

// customizeEnv retrieves all EnvironmentCustomizer components and invokes their CustomizeEnvironment method to allow
// them to modify the environment before it is used by the application.
func customizeEnv(env runtime.Environment, app runtime.Application) error {
	components := component.ListOf[runtime.EnvironmentCustomizer]()

	for _, comp := range components {
		customizer, err := component.Load[runtime.EnvironmentCustomizer](comp.Definition().Name())
		if err != nil {
			return err
		}

		err = customizer.CustomizeEnvironment(env, app)

		if err != nil {
			return nil
		}
	}

	return nil
}

// prepareRuntimeContext creates the application context, allows customizers to modify it, and registers
// the command-line arguments in the context's container.
func prepareRuntimeContext(env runtime.Environment, args *runtime.Args) (runtime.Context, error) {
	runtimeCtx := createContext(env)
	err := customizeRuntimeContext(runtimeCtx)
	if err != nil {
		return nil, err
	}

	log.Info("Starting application using Go {} ({}/{})", goruntime.Version()[2:], goruntime.GOOS, goruntime.GOARCH)
	log.Info("Running with Procyon {}", Version)

	container := runtimeCtx.Container()
	err = container.RegisterSingleton("procyonAppArgs", args)
	if err != nil {
		return nil, err
	}

	return runtimeCtx, nil
}

// customizeRuntimeContext retrieves all ContextCustomizer components and invokes their CustomizeContext method to allow
// them to modify the application context before it is refreshed and used by the application.
func customizeRuntimeContext(runtimeCtx *Context) error {
	components := component.ListOf[runtime.ContextCustomizer]()

	for _, comp := range components {
		customizer, err := component.Load[runtime.ContextCustomizer](comp.Definition().Name())
		if err != nil {
			return err
		}

		err = customizer.CustomizeContext(runtimeCtx)

		if err != nil {
			return nil
		}
	}

	return nil
}

// isServerApplication checks whether the application is a server application by looking for components of
// type runtime.Server in the component list.
func isServerApplication() bool {
	components := component.ListOf[runtime.Server]()
	return len(components) > 0
}
