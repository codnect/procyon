package procyon

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

func newApplicationRunListeners(l []ApplicationRunListener) ApplicationRunListeners {
	return ApplicationRunListeners{
		listeners: l,
	}
}

func (appListeners ApplicationRunListeners) starting() {
	for _, listener := range appListeners.listeners {
		listener.starting()
	}
}

func (appListeners ApplicationRunListeners) environmentPrepared(environment ConfigurableEnvironment) {
	for _, listener := range appListeners.listeners {
		listener.environmentPrepared(environment)
	}
}

func (appListeners ApplicationRunListeners) contextPrepared(context ConfigurableApplicationContext) {
	for _, listener := range appListeners.listeners {
		listener.contextPrepared(context)
	}
}

func (appListeners ApplicationRunListeners) contextLoaded(context ConfigurableApplicationContext) {
	for _, listener := range appListeners.listeners {
		listener.contextLoaded(context)
	}
}

func (appListeners ApplicationRunListeners) started(context ConfigurableApplicationContext) {
	for _, listener := range appListeners.listeners {
		listener.started(context)
	}
}

func (appListeners ApplicationRunListeners) running(context ConfigurableApplicationContext) {
	for _, listener := range appListeners.listeners {
		listener.running(context)
	}
}

func (appListeners ApplicationRunListeners) failed(context ConfigurableApplicationContext, err error) {
	for _, listener := range appListeners.listeners {
		listener.failed(context, err)
	}
}
