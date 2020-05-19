package procyon

import (
	"github.com/Rollcomp/procyon-context"
	"github.com/Rollcomp/procyon-core"
)

type ApplicationArguments interface {
	ContainsOption(name string) bool
	GetOptionNames() []string
	GetOptionValues(name string) []string
	GetSourceArgs() []string
	GetNonOptionArgs() []string
}

type DefaultApplicationArguments struct {
	source core.SimpleCommandLinePropertySource
	args   []string
}

func GetApplicationArguments(args []string) ApplicationArguments {
	return &DefaultApplicationArguments{
		args:   args,
		source: core.NewSimpleCommandLinePropertySource(args),
	}
}

func (arg DefaultApplicationArguments) ContainsOption(name string) bool {
	return arg.source.ContainsOption(name)
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
	environmentPrepared(environment core.ConfigurableEnvironment)
	contextPrepared(context context.ConfigurableApplicationContext)
	contextLoaded(context context.ConfigurableApplicationContext)
	started(context context.ConfigurableApplicationContext)
	running(context context.ConfigurableApplicationContext)
	failed(context context.ConfigurableApplicationContext, err error)
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

func (appListeners ApplicationRunListeners) EnvironmentPrepared(environment core.ConfigurableEnvironment) {
	for _, listener := range appListeners.listeners {
		listener.environmentPrepared(environment)
	}
}

func (appListeners ApplicationRunListeners) ContextPrepared(context context.ConfigurableApplicationContext) {
	for _, listener := range appListeners.listeners {
		listener.contextPrepared(context)
	}
}

func (appListeners ApplicationRunListeners) ContextLoaded(context context.ConfigurableApplicationContext) {
	for _, listener := range appListeners.listeners {
		listener.contextLoaded(context)
	}
}

func (appListeners ApplicationRunListeners) Started(context context.ConfigurableApplicationContext) {
	for _, listener := range appListeners.listeners {
		listener.started(context)
	}
}

func (appListeners ApplicationRunListeners) Running(context context.ConfigurableApplicationContext) {
	for _, listener := range appListeners.listeners {
		listener.running(context)
	}
}

func (appListeners ApplicationRunListeners) failed(context context.ConfigurableApplicationContext, err error) {
	for _, listener := range appListeners.listeners {
		listener.failed(context, err)
	}
}
