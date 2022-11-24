package app

import (
	"fmt"
	"github.com/procyon-projects/procyon/env"
	"log"
	"time"
)

type Application interface {
	Run(args ...string)
}

func New() Application {
	return &application{}
}

type application struct {
}

func (a *application) Run(args ...string) {
	log.SetFlags(0)

	arguments, err := parseArguments(mergeArguments(args...))

	if err != nil {
		panic(fmt.Errorf("app: argument parsing failed %v", err.Error()))
	}

	startTime := time.Now()

	var listeners startupListeners
	ctx := newContext()

	defer func() {
		if r := recover(); r != nil {
			listeners.failed(ctx, err)
		}
	}()

	listeners.starting(ctx)
	a.prepareEnvironment(ctx, arguments, listeners)

	listeners.ready(ctx, startTime.Sub(time.Now()))

	listeners.started(ctx, startTime.Sub(time.Now()))
}

func (a *application) prepareEnvironment(ctx *appContext, arguments *Arguments, listeners startupListeners) env.Environment {
	environment := env.New()
	propertySources := environment.PropertySources()

	propertySources.AddFirst(newArgumentPropertySources(arguments))

	listeners.environmentPrepared(ctx, environment)
	return nil
}

func (a *application) prepareContext(ctx *appContext, environment env.Environment, listeners startupListeners) {
	ctx.setEnvironment(environment)
	listeners.contextPrepared(ctx)

	listeners.contextLoaded(ctx)
}
