package procyon

import (
	context "github.com/Rollcomp/procyon-context"
	core "github.com/Rollcomp/procyon-core"
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
}

func NewEventPublishRunListener(app *Application) EventPublishRunListener {
	runListener := EventPublishRunListener{
		app:         app,
		broadcaster: context.NewSimpleApplicationEventBroadcaster(),
	}
	appListeners := app.getAppListeners()
	for _, appListener := range appListeners {
		runListener.broadcaster.RegisterApplicationListener(appListener)
	}
	return runListener
}

func (listener EventPublishRunListener) Starting() {
	listener.broadcaster.BroadcastEvent(NewApplicationStarting(listener.app, nil))
}

func (listener EventPublishRunListener) EnvironmentPrepared(environment core.ConfigurableEnvironment) {
	listener.broadcaster.BroadcastEvent(NewApplicationEnvironmentPreparedEvent(listener.app, nil, environment))
}

func (listener EventPublishRunListener) ContextPrepared(ctx context.ConfigurableApplicationContext) {
	listener.broadcaster.BroadcastEvent(NewApplicationContextInitializedEvent(listener.app, nil, ctx))
}

func (listener EventPublishRunListener) ContextLoaded(ctx context.ConfigurableApplicationContext) {
	listener.broadcaster.BroadcastEvent(NewApplicationPreparedEvent(listener.app, nil, ctx))
}

func (listener EventPublishRunListener) Started(ctx context.ConfigurableApplicationContext) {
	listener.broadcaster.BroadcastEvent(NewApplicationStartedEvent(listener.app, nil, ctx))
}

func (listener EventPublishRunListener) Running(ctx context.ConfigurableApplicationContext) {
	listener.broadcaster.BroadcastEvent(NewApplicationReadyEvent(listener.app, nil, ctx))
}

func (listener EventPublishRunListener) Failed(ctx context.ConfigurableApplicationContext, err error) {
	listener.broadcaster.BroadcastEvent(NewApplicationFailedEvent(listener.app, nil, ctx, err))
}
