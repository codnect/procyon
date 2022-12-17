package app

import (
	"fmt"
	"github.com/procyon-projects/procyon/container"
	"github.com/procyon-projects/procyon/env"
	"github.com/procyon-projects/procyon/event"
	"github.com/procyon-projects/reflector"
	"log"
	"time"
)

type Application interface {
	Context() Context
	Run(args ...string)
}

func New() Application {
	appContainer := container.New()
	broadcaster := event.NewBroadcaster()

	return &application{
		ctx:       newContext(appContainer, broadcaster),
		container: appContainer,
	}
}

type application struct {
	ctx       *appContext
	container *container.Container
}

func (a *application) Context() Context {
	return a.ctx
}

func (a *application) Run(args ...string) {
	log.SetFlags(0)

	arguments, err := parseArguments(mergeArguments(args...))

	if err != nil {
		panic(fmt.Errorf("app: argument parsing failed %v", err.Error()))
	}

	startTime := time.Now()

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
	a.prepareEnvironment(arguments, listeners)

	listeners.ready(a.ctx, startTime.Sub(time.Now()))

	listeners.started(a.ctx, startTime.Sub(time.Now()))
}

func (a *application) startupListeners(arguments *Arguments) (startupListeners, error) {
	listeners := make(startupListeners, 0)

	reflApplicationType := reflector.TypeOf[Application]().ReflectType()
	reflArgumentsType := reflector.TypeOf[*Arguments]().ReflectType()

	registry := a.container.DefinitionRegistry()
	definitionNames := registry.DefinitionNamesByType(container.TypeOf[StartupListener]())

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

func (a *application) prepareEnvironment(arguments *Arguments, listeners startupListeners) env.Environment {
	environment := env.New()
	propertySources := environment.PropertySources()

	propertySources.AddFirst(newArgumentPropertySources(arguments))

	listeners.environmentPrepared(a.ctx, environment)
	return nil
}

func (a *application) prepareContext(environment env.Environment, listeners startupListeners) {
	a.ctx.setEnvironment(environment)
	listeners.contextPrepared(a.ctx)

	listeners.contextLoaded(a.ctx)
}
