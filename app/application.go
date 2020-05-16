package app

type ApplicationArguments interface {
}

type DefaultApplicationArguments struct {
}

func NewDefaultApplicationArguments() ApplicationArguments {
	return &DefaultApplicationArguments{}
}

type ApplicationRunListener interface {
	starting()
	environmentPrepared(environment ConfigurableEnvironment)
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

func (appListeners ApplicationRunListeners) EnvironmentPrepared(environment ConfigurableEnvironment) {
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
