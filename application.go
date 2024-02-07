package procyon

import (
	"codnect.io/procyon-core/container"
	"codnect.io/procyon-core/runtime"
	"codnect.io/procyon-core/runtime/env"
	"fmt"
	"os"
	sruntime "runtime"
	"strings"
	"time"
)

type Application struct {
	ctx       *Context
	container container.Container
	env       env.Environment
	banner    Banner
}

func New() *Application {
	instanceContainer := container.New()

	return &Application{
		ctx:       newContext(instanceContainer),
		container: instanceContainer,
		banner:    newDefaultBannerPrinter(),
	}
}

func (a *Application) SetBanner(banner Banner) *Application {
	if banner != nil {
		a.banner = banner
	}

	return nil
}

func (a *Application) Run(args ...string) (err error) {
	var (
		startTime   = time.Now()
		arguments   *runtime.Arguments
		listeners   runtime.StartupListeners
		environment env.Environment
	)

	err = a.banner.PrintBanner(os.Stdout)
	if err != nil {
		return err
	}

	arguments, err = runtime.ParseArguments(args)
	if err != nil {
		return fmt.Errorf("failed to parse arguments: %v", err)
	}

	err = registerComponentDefinitions(a.container)
	if err != nil {
		return fmt.Errorf("failed to register components: %v", err)
	}

	log.Info("Starting application using Go {}", sruntime.Version()[2:])
	log.Debug("Running with Procyon {}", Version)

	listeners, err = getComponentsByType[runtime.StartupListener](a.container, a, arguments)
	if err != nil {
		return fmt.Errorf("failed to initialize startup listeners: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			listeners.Failed(a.ctx, err)
		}
	}()

	listeners.Starting(a.ctx)

	environment, err = a.prepareEnvironment(arguments, listeners)
	if err != nil {
		panic(err)
	}

	a.logProfileInfo(environment)

	err = a.prepareContext(environment, listeners, arguments)
	if err != nil {
		panic(err)
	}

	timeTakenToStartup := time.Now().Sub(startTime)
	listeners.Started(a.ctx, timeTakenToStartup)
	log.Info("Started application in {} seconds", timeTakenToStartup.Seconds())

	timeTakenToReady := time.Now().Sub(startTime)
	listeners.Ready(a.ctx, timeTakenToReady)

	// wait for context to be closed
	<-a.ctx.Done()

	return nil
}

func (a *Application) Exit() {
	a.ctx.Close()
}

func (a *Application) prepareEnvironment(arguments *runtime.Arguments, listeners runtime.StartupListeners) (env.Environment, error) {
	environment := newEnvironment()

	propertySources := environment.PropertySources()

	propertySources.AddLast(runtime.NewPropertySource(arguments))
	propertySources.AddLast(env.NewPropertySource())

	err := environment.customize(a.container)
	if err != nil {
		return nil, err
	}

	listeners.EnvironmentPrepared(a.ctx, environment)
	return environment, nil
}

func (a *Application) prepareContext(environment env.Environment, listeners runtime.StartupListeners, arguments *runtime.Arguments) error {
	a.ctx.setEnvironment(environment)

	err := a.ctx.customize()
	if err != nil {
		return err
	}

	listeners.ContextPrepared(a.ctx)

	sharedInstances := a.container.SharedInstances()
	err = sharedInstances.Register("procyonAppArguments", arguments)
	if err != nil {
		return err
	}

	listeners.ContextLoaded(a.ctx)

	err = a.ctx.start()
	if err != nil {
		return err
	}

	listeners.ContextStarted(a.ctx)
	return nil
}

func (a *Application) logProfileInfo(environment env.Environment) {
	if log.IsInfoEnabled() {
		activeProfiles := environment.ActiveProfiles()

		if len(activeProfiles) == 0 {
			defaultProfiles := environment.DefaultProfiles()
			log.Info("No active profile, using default profile(s): {}", strings.Join(defaultProfiles, ","))
		} else {
			log.Info("The application is using the following profile(s): {}", strings.Join(activeProfiles, ","))
		}
	}
}
