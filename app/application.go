package app

import "procyon/env"

type ApplicationArguments interface {
	ContainsOption(name string) bool
	GetOptionNames() []string
	GetOptionValues(name string) []string
	GetSourceArgs() []string
	GetNonOptionArgs() []string
}

type DefaultApplicationArguments struct {
	source env.SimpleCommandLinePropertySource
	args   []string
}

func NewDefaultApplicationArguments(args []string) ApplicationArguments {
	return &DefaultApplicationArguments{
		source: env.NewSimpleCommandLinePropertySource(args),
	}
}

func (arg DefaultApplicationArguments) ContainsOption(name string) bool {
	return arg.ContainsOption(name)
}

func (arg DefaultApplicationArguments) GetOptionNames() []string {
	return arg.source.GetPropertyNames()
}

func (arg DefaultApplicationArguments) GetOptionValues(name string) []string {
	return arg.source.GetOptionValues(name)
}

func (arg DefaultApplicationArguments) GetSourceArgs() []string {
	return arg.args
}

func (arg DefaultApplicationArguments) GetNonOptionArgs() []string {
	return arg.source.GetNonOptionArgs()
}

type ApplicationRunListener interface {
	starting()
	environmentPrepared(environment env.ConfigurableEnvironment)
	contextPrepared(context ConfigurableApplicationContext)
	contextLoaded(context ConfigurableApplicationContext)
	started(context ConfigurableApplicationContext)
	running(context ConfigurableApplicationContext)
	failed(context ConfigurableApplicationContext, err error)
}

type ApplicationRunListeners struct {
	listeners []ApplicationRunListener
}

func NewApplicationRunListeners(l []ApplicationRunListener) ApplicationRunListeners {
	return ApplicationRunListeners{
		listeners: l,
	}
}

func (appListeners ApplicationRunListeners) Starting() {
	for _, listener := range appListeners.listeners {
		listener.starting()
	}
}

func (appListeners ApplicationRunListeners) EnvironmentPrepared(environment env.ConfigurableEnvironment) {
	for _, listener := range appListeners.listeners {
		listener.environmentPrepared(environment)
	}
}

func (appListeners ApplicationRunListeners) ContextPrepared(context ConfigurableApplicationContext) {
	for _, listener := range appListeners.listeners {
		listener.contextPrepared(context)
	}
}

func (appListeners ApplicationRunListeners) ContextLoaded(context ConfigurableApplicationContext) {
	for _, listener := range appListeners.listeners {
		listener.contextLoaded(context)
	}
}

func (appListeners ApplicationRunListeners) Started(context ConfigurableApplicationContext) {
	for _, listener := range appListeners.listeners {
		listener.started(context)
	}
}

func (appListeners ApplicationRunListeners) Running(context ConfigurableApplicationContext) {
	for _, listener := range appListeners.listeners {
		listener.running(context)
	}
}

func (appListeners ApplicationRunListeners) failed(context ConfigurableApplicationContext, err error) {
	for _, listener := range appListeners.listeners {
		listener.failed(context, err)
	}
}
