package procyon

import (
	"codnect.io/logy"
	"codnect.io/procyon-core/component"
	"codnect.io/procyon-core/container"
	"codnect.io/procyon-core/env"
	"codnect.io/procyon/event"
	"codnect.io/reflector"
	"context"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

var (
	log = logy.Get()
)

const (
	Version = "v0.0.1-dev"
)

type Application struct {
	ctx           *appContext
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
	arguments, err := parseArguments(mergeArguments(args...))

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
	listeners, err = a.startupListeners(arguments)

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

func (a *Application) startupListeners(arguments *Arguments) (startupListeners, error) {
	listeners := make(startupListeners, 0)

	reflApplicationType := reflector.TypeOf[Application]().ReflectType()
	reflArgumentsType := reflector.TypeOf[*Arguments]().ReflectType()

	registry := a.container.DefinitionRegistry()
	definitionNames := registry.DefinitionNamesByType(reflector.TypeOf[StartupListener]())

	for _, definitionName := range definitionNames {
		definition, _ := registry.Find(definitionName)

		if len(definition.Inputs()) != 2 {
			continue
		}

		if !reflApplicationType.ConvertibleTo(definition.Inputs()[0].Type().ReflectType()) {
			continue
		}

		if !reflArgumentsType.ConvertibleTo(definition.Inputs()[1].Type().ReflectType()) {
			continue
		}

		results, err := definition.Constructor().Invoke(a, arguments)

		if err != nil {
			return nil, err
		}

		listener := results[0].(StartupListener)
		listeners = append(listeners, listener)
	}

	return listeners, nil
}

func (a *Application) eventCustomizers() (eventCustomizers, error) {
	customizers := make(eventCustomizers, 0)

	registry := a.container.DefinitionRegistry()
	definitionNames := registry.DefinitionNamesByType(reflector.TypeOf[env.Customizer]())

	for _, definitionName := range definitionNames {
		definition, _ := registry.Find(definitionName)

		if len(definition.Inputs()) != 0 {
			continue
		}

		results, err := definition.Constructor().Invoke()

		if err != nil {
			return nil, err
		}

		customizer := results[0].(env.Customizer)
		customizers = append(customizers, customizer)
	}

	return customizers, nil
}

func (a *Application) prepareEnvironment(arguments *Arguments, listeners startupListeners) (env.Environment, error) {
	environment := env.New()

	propertySources := environment.PropertySources()

	propertySources.AddLast(newArgumentPropertySources(arguments))
	propertySources.AddLast(newSystemEnvironmentPropertySource())

	customizers, err := a.eventCustomizers()
	if err != nil {
		return nil, err
	}

	err = customizers.invoke(environment)
	if err != nil {
		return nil, err
	}

	listeners.environmentPrepared(a.ctx, environment)
	return environment, nil
}

func (a *Application) contextCustomizers() (contextCustomizers, error) {
	customizers := make(contextCustomizers, 0)

	registry := a.container.DefinitionRegistry()
	definitionNames := registry.DefinitionNamesByType(reflector.TypeOf[ContextCustomizer]())

	for _, definitionName := range definitionNames {
		definition, _ := registry.Find(definitionName)

		if len(definition.Inputs()) != 0 {
			continue
		}

		results, err := definition.Constructor().Invoke()

		if err != nil {
			return nil, err
		}

		customizer := results[0].(ContextCustomizer)
		customizers = append(customizers, customizer.(ContextCustomizer))
	}

	return customizers, nil
}

func (a *Application) prepareContext(environment env.Environment, listeners startupListeners, arguments *Arguments) error {
	a.ctx.setEnvironment(environment)

	customizers, err := a.contextCustomizers()
	if err != nil {
		return err
	}

	err = customizers.invoke(a.ctx)
	if err != nil {
		return err
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

func (a *Application) logStartup(ctx Context) {
	appName := ctx.ApplicationName()

	if appName == "" {
		appName = "application"
	}

	log.Info("Starting {} using Go {}", appName, runtime.Version()[2:])
	log.Debug("Running with Procyon {}", Version)
}

func (a *Application) logProfileInfo(environment env.Environment) {
	if log.IsInfoEnabled() {
		activeProfiles := environment.ActiveProfiles()
		if len(activeProfiles) == 0 {
			defaultProfiles := environment.DefaultProfiles()
			log.Info("No active profile, using default profile(s): {}", sliceToDelimitedString(defaultProfiles))
		} else {
			if len(activeProfiles) == 1 {
				log.Info("The following profile is active: {}", sliceToDelimitedString(activeProfiles))
			} else {
				log.Info("The following profiles are active: {}", sliceToDelimitedString(activeProfiles))
			}
		}
	}
}

func (a *Application) logStarted(ctx Context, timeTaken time.Duration) {
	appName := ctx.ApplicationName()

	if appName == "" {
		appName = "application"
	}

	log.Info("Started {} in {} seconds", appName, timeTaken.Seconds())
}

func sliceToDelimitedString(values []string) string {
	var builder strings.Builder
	for index, value := range values {
		builder.WriteByte('"')
		builder.WriteString(value)
		builder.WriteByte('"')
		if index != len(values)-1 {
			builder.WriteByte(',')
		}
	}
	return builder.String()
}
