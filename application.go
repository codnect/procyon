package procyon

import (
	"codnect.io/procyon-core/component"
	"codnect.io/procyon-core/container"
	"codnect.io/procyon-core/event"
	"codnect.io/procyon-core/runtime"
	"codnect.io/procyon-core/runtime/env"
	"context"
	"os"
	"os/signal"
	sruntime "runtime"
	"strings"
	"syscall"
	"time"
)

type Application struct {
	ctx           *Context
	container     container.Container
	env           env.Environment
	bannerPrinter *bannerPrinter
}

func New() *Application {
	appContainer := container.New()
	broadcaster := event.NewBroadcaster()

	return &Application{
		ctx:           newContext(appContainer, broadcaster),
		container:     appContainer,
		bannerPrinter: defaultBannerPrinter(),
	}
}

func (a *Application) Run(args ...string) {
	startTime := time.Now()

	a.bannerPrinter.PrintBanner(os.Stdout)
	arguments, err := runtime.ParseArguments(args)

	if err != nil {
		log.Error("Argument parsing failed", err)
		os.Exit(1)
	}

	err = a.registerComponentDefinitions()

	if err != nil {
		log.Error("Failed to register component definitions", err)
		os.Exit(1)
	}

	a.logStartup(a.ctx)

	var listeners startupListeners
	listeners, err = getComponentsByType[runtime.StartupListener](a.container, a, arguments)

	if err != nil {
		log.Error("Failed to initialize startup listeners", err)
		os.Exit(1)
	}

	defer func() {
		if r := recover(); r != nil {
			listeners.failed(a.ctx, err)
		}
	}()

	listeners.starting(a.ctx)

	var environment env.Environment
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

	listeners.started(a.ctx, timeTakenToStartup)
	a.logStarted(a.ctx, timeTakenToStartup)

	timeTakenToReady := time.Now().Sub(startTime)
	listeners.ready(a.ctx, timeTakenToReady)

	notifyCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-notifyCtx.Done()
	_ = a.ctx.Stop()
}

func (a *Application) prepareEnvironment(arguments *runtime.Arguments, listeners startupListeners) (env.Environment, error) {
	environment := env.New()

	propertySources := environment.PropertySources()

	propertySources.AddLast(runtime.NewPropertySource(arguments))
	propertySources.AddLast(env.NewPropertySource())

	customizers, err := getComponentsByType[env.Customizer](a.container)
	if err != nil {
		return nil, err
	}

	for _, customizer := range customizers {
		err = customizer.CustomizeEnvironment(environment)

		if err != nil {
			return nil, err
		}
	}

	listeners.environmentPrepared(a.ctx, environment)
	return environment, nil
}

func (a *Application) prepareContext(environment env.Environment, listeners startupListeners, arguments *runtime.Arguments) error {
	a.ctx.setEnvironment(environment)

	customizers, err := getComponentsByType[runtime.ContextCustomizer](a.container)
	if err != nil {
		return err
	}

	for _, customizer := range customizers {
		err = customizer.CustomizeContext(a.ctx)

		if err != nil {
			return err
		}
	}

	listeners.contextPrepared(a.ctx)

	sharedInstances := a.container.SharedInstances()
	err = sharedInstances.Register("procyonApplicationArguments", arguments)
	if err != nil {
		return err
	}

	listeners.contextLoaded(a.ctx)

	err = a.ctx.Start()
	if err != nil {
		return err
	}

	listeners.contextStarted(a.ctx)
	return nil
}

func (a *Application) registerComponentDefinitions() error {
	for _, registeredComponent := range component.RegisteredComponents() {
		err := a.container.DefinitionRegistry().Register(registeredComponent.Definition())
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *Application) logStartup(ctx *Context) {
	appName := ctx.ApplicationName()

	if appName == "" {
		appName = "application"
	}

	log.Info("Starting {} using Go {}", appName, sruntime.Version()[2:])
	log.Debug("Running with Procyon {}", Version)
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

func (a *Application) logStarted(ctx *Context, timeTaken time.Duration) {
	appName := ctx.ApplicationName()

	if appName == "" {
		appName = "application"
	}

	log.Info("Started {} in {} seconds", appName, timeTaken.Seconds())
}
