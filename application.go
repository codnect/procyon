package procyon

import (
	"github.com/procyon-projects/procyon-context"
	"github.com/procyon-projects/procyon-core"
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
	Starting()
	EnvironmentPrepared(environment core.ConfigurableEnvironment)
	ContextPrepared(context context.ConfigurableApplicationContext)
	ContextLoaded(context context.ConfigurableApplicationContext)
	Started(context context.ConfigurableApplicationContext)
	Running(context context.ConfigurableApplicationContext)
	Failed(context context.ConfigurableApplicationContext, err error)
}

type ApplicationRunListeners struct {
	listeners []ApplicationRunListener
}

func NewApplicationRunListeners(l []ApplicationRunListener) *ApplicationRunListeners {
	return &ApplicationRunListeners{
		listeners: l,
	}
}

func (appListeners *ApplicationRunListeners) Starting() {
	for _, listener := range appListeners.listeners {
		listener.Starting()
	}
}

func (appListeners *ApplicationRunListeners) EnvironmentPrepared(environment core.ConfigurableEnvironment) {
	for _, listener := range appListeners.listeners {
		listener.EnvironmentPrepared(environment)
	}
}

func (appListeners *ApplicationRunListeners) ContextPrepared(context context.ConfigurableApplicationContext) {
	for _, listener := range appListeners.listeners {
		listener.ContextPrepared(context)
	}
}

func (appListeners *ApplicationRunListeners) ContextLoaded(context context.ConfigurableApplicationContext) {
	for _, listener := range appListeners.listeners {
		listener.ContextLoaded(context)
	}
}

func (appListeners *ApplicationRunListeners) Started(context context.ConfigurableApplicationContext) {
	for _, listener := range appListeners.listeners {
		listener.Started(context)
	}
}

func (appListeners *ApplicationRunListeners) Running(context context.ConfigurableApplicationContext) {
	for _, listener := range appListeners.listeners {
		listener.Running(context)
	}
}

func (appListeners *ApplicationRunListeners) Failed(context context.ConfigurableApplicationContext, err error) {
	for _, listener := range appListeners.listeners {
		listener.Failed(context, err)
	}
}
