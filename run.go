package procyon

import (
	context "github.com/procyon-projects/procyon-context"
	core "github.com/procyon-projects/procyon-core"
)

type ApplicationRunListener interface {
	OnApplicationStarting()
	OnApplicationEnvironmentPrepared(environment core.ConfigurableEnvironment)
	OnApplicationContextPrepared(context context.ConfigurableApplicationContext)
	OnApplicationContextLoaded(context context.ConfigurableApplicationContext)
	OnApplicationStarted(context context.ConfigurableApplicationContext)
	OnApplicationRunning(context context.ConfigurableApplicationContext)
	OnApplicationFailed(context context.ConfigurableApplicationContext, err error)
}

type ApplicationRunListeners struct {
	listeners []ApplicationRunListener
}

func NewApplicationRunListeners(l []ApplicationRunListener) *ApplicationRunListeners {
	return &ApplicationRunListeners{
		listeners: l,
	}
}

func (appListeners *ApplicationRunListeners) OnApplicationStarting() {
	for _, listener := range appListeners.listeners {
		listener.OnApplicationStarting()
	}
}

func (appListeners *ApplicationRunListeners) OnApplicationEnvironmentPrepared(environment core.ConfigurableEnvironment) {
	for _, listener := range appListeners.listeners {
		listener.OnApplicationEnvironmentPrepared(environment)
	}
}

func (appListeners *ApplicationRunListeners) OnApplicationContextPrepared(context context.ConfigurableApplicationContext) {
	for _, listener := range appListeners.listeners {
		listener.OnApplicationContextPrepared(context)
	}
}

func (appListeners *ApplicationRunListeners) OnApplicationContextLoaded(context context.ConfigurableApplicationContext) {
	for _, listener := range appListeners.listeners {
		listener.OnApplicationContextLoaded(context)
	}
}

func (appListeners *ApplicationRunListeners) OnApplicationStarted(context context.ConfigurableApplicationContext) {
	for _, listener := range appListeners.listeners {
		listener.OnApplicationStarted(context)
	}
}

func (appListeners *ApplicationRunListeners) OnApplicationRunning(context context.ConfigurableApplicationContext) {
	for _, listener := range appListeners.listeners {
		listener.OnApplicationRunning(context)
	}
}

func (appListeners *ApplicationRunListeners) OnApplicationFailed(context context.ConfigurableApplicationContext, err error) {
	for _, listener := range appListeners.listeners {
		listener.OnApplicationFailed(context, err)
	}
}
