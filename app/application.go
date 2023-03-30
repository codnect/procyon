package app

import (
	"fmt"
	"github.com/procyon-projects/logy"
	"github.com/procyon-projects/procyon/app/component"
	"github.com/procyon-projects/procyon/app/env"
	"github.com/procyon-projects/procyon/app/event"
	"github.com/procyon-projects/procyon/container"
	"github.com/procyon-projects/reflector"
	"os"
	"runtime"
	"time"
)

var (
	log = logy.Get()
)

type Application interface {
	Context() Context
	Run(args ...string)
}

func New() Application {
	appContainer := container.New()
	broadcaster := event.NewBroadcaster()

	return &application{
		ctx:           newContext(appContainer, broadcaster),
		container:     appContainer,
		bannerPrinter: defaultBannerPrinter(),
	}
}

type application struct {
	ctx           *appContext
	container     container.Container
	env           env.Environment
	bannerPrinter *bannerPrinter
}

func (a *application) Context() Context {
	return a.ctx
}

func (a *application) Run(args ...string) {

	arguments, err := parseArguments(mergeArguments(args...))

	if err != nil {
		panic(fmt.Errorf("app: argument parsing failed %v", err.Error()))
	}

	startTime := time.Now()

	err = a.registerComponentDefinitions()

	if err != nil {
		panic(fmt.Errorf("app: failed to register component definitions, err: %s", err.Error()))
	}

	var listeners startupListeners
	listeners, err = a.startupListeners(arguments)

	if err != nil {
		panic(fmt.Errorf("app: failed to initialize startup listeners, err: %s", err.Error()))
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

	err = a.bannerPrinter.PrintBanner(environment, os.Stdout)
	if err != nil {
		panic(err)
	}

	err = a.prepareContext(environment, listeners)
	if err != nil {
		panic(err)
	}

	timeTaken := time.Now().Sub(startTime)
	listeners.ready(a.ctx, timeTaken)
	a.logStarted(a.ctx, timeTaken)
	listeners.started(a.ctx, timeTaken)
}

func (a *application) startupListeners(arguments *Arguments) (startupListeners, error) {
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

		listener, err := a.container.GetByNameAndArgs(a.ctx, definitionName, a, arguments)

		if err != nil {
			return nil, err
		}

		listeners = append(listeners, listener.(StartupListener))
	}

	return listeners, nil
}

func (a *application) eventCustomizers() (eventCustomizers, error) {
	customizers := make(eventCustomizers, 0)

	registry := a.container.DefinitionRegistry()
	definitionNames := registry.DefinitionNamesByType(reflector.TypeOf[env.Customizer]())

	for _, definitionName := range definitionNames {
		definition, _ := registry.Find(definitionName)

		if len(definition.Inputs()) != 0 {
			continue
		}

		customizer, err := a.container.Get(a.ctx, definitionName)

		if err != nil {
			return nil, err
		}

		customizers = append(customizers, customizer.(env.Customizer))
	}

	return customizers, nil
}

func (a *application) prepareEnvironment(arguments *Arguments, listeners startupListeners) (env.Environment, error) {
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

func (a *application) contextCustomizers() (contextCustomizers, error) {
	customizers := make(contextCustomizers, 0)

	registry := a.container.DefinitionRegistry()
	definitionNames := registry.DefinitionNamesByType(reflector.TypeOf[ContextCustomizer]())

	for _, definitionName := range definitionNames {
		definition, _ := registry.Find(definitionName)

		if len(definition.Inputs()) != 0 {
			continue
		}

		customizer, err := a.container.Get(a.ctx, definitionName)

		if err != nil {
			return nil, err
		}

		customizers = append(customizers, customizer.(ContextCustomizer))
	}

	return customizers, nil
}

func (a *application) prepareContext(environment env.Environment, listeners startupListeners) error {
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

	a.logStartup(a.ctx)
	a.logProfileInfo(environment)

	listeners.contextLoaded(a.ctx)

	err = a.ctx.Refresh()
	if err != nil {
		return err
	}

	listeners.contextRefreshed(a.ctx)
	return nil
}

func (a *application) registerComponentDefinitions() error {
	for _, registeredComponent := range component.RegisteredComponents() {
		err := a.container.DefinitionRegistry().Register(registeredComponent.Definition())
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *application) logStartup(ctx Context) {
	appName := ctx.ApplicationName()

	if appName == "" {
		appName = "application"
	}

	log.Info("Starting {} using Go {}", appName, runtime.Version()[2:])
	log.Debug("Running with Procyon {}", Version)
}

func (a *application) logProfileInfo(environment env.Environment) {
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

func (a *application) logStarted(ctx Context, timeTaken time.Duration) {
	appName := ctx.ApplicationName()

	if appName == "" {
		appName = "application"
	}

	log.Info("Started {} in {} seconds", appName, timeTaken.Seconds())
}
