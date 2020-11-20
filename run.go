package procyon

import (
	context "github.com/procyon-projects/procyon-context"
	core "github.com/procyon-projects/procyon-core"
)

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
