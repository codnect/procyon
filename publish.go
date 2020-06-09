package procyon

import (
	context "github.com/procyon-projects/procyon-context"
	core "github.com/procyon-projects/procyon-core"
)

type BootstrapListener struct {
}

func NewBootstrapListener() BootstrapListener {
	return BootstrapListener{}
}

func (listener BootstrapListener) SubscribeEvents() []context.ApplicationEvent {
	return []context.ApplicationEvent{
		(*ApplicationStartingEvent)(nil),
	}
}

func (listener BootstrapListener) OnApplicationEvent(event context.ApplicationEvent) {
	//source := event.GetSource()
	//timestamp := event.GetTimestamp()
}

type EventPublishRunListener struct {
	app         *Application
	broadcaster context.ApplicationEventBroadcaster
	args        ApplicationArguments
}

func NewEventPublishRunListener(app *Application, arguments ApplicationArguments) EventPublishRunListener {
	runListener := EventPublishRunListener{
		app:         app,
		broadcaster: context.NewSimpleApplicationEventBroadcaster(),
		args:        arguments,
	}
	appListeners := app.getAppListeners()
	for _, appListener := range appListeners {
		runListener.broadcaster.RegisterApplicationListener(appListener)
	}
	return runListener
}

func (listener EventPublishRunListener) Starting() {
	listener.broadcaster.BroadcastEvent(NewApplicationStarting(listener.app, listener.args))
}

func (listener EventPublishRunListener) EnvironmentPrepared(environment core.ConfigurableEnvironment) {
	listener.broadcaster.BroadcastEvent(NewApplicationEnvironmentPreparedEvent(listener.app, listener.args, environment))
}

func (listener EventPublishRunListener) ContextPrepared(ctx context.ConfigurableApplicationContext) {
	listener.broadcaster.BroadcastEvent(NewApplicationContextInitializedEvent(listener.app, listener.args, ctx))
}

func (listener EventPublishRunListener) ContextLoaded(ctx context.ConfigurableApplicationContext) {
	// when context is loaded, add application listeners registered
	appListeners := listener.app.getAppListeners()
	for _, appListener := range appListeners {
		ctx.AddApplicationListener(appListener)
	}
	// after that, broadcast an event to notify all listeners that app is prepared
	listener.broadcaster.BroadcastEvent(NewApplicationPreparedEvent(listener.app, listener.args, ctx))
}

func (listener EventPublishRunListener) Started(ctx context.ConfigurableApplicationContext) {
	listener.broadcaster.BroadcastEvent(NewApplicationStartedEvent(listener.app, listener.args, ctx))
}

func (listener EventPublishRunListener) Running(ctx context.ConfigurableApplicationContext) {
	listener.broadcaster.BroadcastEvent(NewApplicationReadyEvent(listener.app, listener.args, ctx))
}

func (listener EventPublishRunListener) Failed(ctx context.ConfigurableApplicationContext, err error) {
	listener.broadcaster.BroadcastEvent(NewApplicationFailedEvent(listener.app, listener.args, ctx, err))
}
