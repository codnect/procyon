package procyon

import (
	context "github.com/procyon-projects/procyon-context"
	core "github.com/procyon-projects/procyon-core"
)

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

func (listener EventPublishRunListener) OnApplicationStarting() {
	listener.broadcaster.BroadcastEvent(nil, NewApplicationStarting(listener.app, listener.args))
}

func (listener EventPublishRunListener) OnApplicationEnvironmentPrepared(environment core.ConfigurableEnvironment) {
	listener.broadcaster.BroadcastEvent(nil, NewApplicationEnvironmentPreparedEvent(listener.app, listener.args, environment))
}

func (listener EventPublishRunListener) OnApplicationContextPrepared(ctx context.ConfigurableApplicationContext) {
	listener.broadcaster.BroadcastEvent(nil, NewApplicationContextInitializedEvent(listener.app, listener.args, ctx))
}

func (listener EventPublishRunListener) OnApplicationContextLoaded(ctx context.ConfigurableApplicationContext) {
	// when context is loaded, add application listeners registered
	appListeners := listener.app.getAppListeners()
	for _, appListener := range appListeners {
		ctx.AddApplicationListener(appListener)
	}
	// after that, broadcast an event to notify all listeners that app is prepared
	listener.broadcaster.BroadcastEvent(ctx, NewApplicationPreparedEvent(listener.app, listener.args, ctx))
}

func (listener EventPublishRunListener) OnApplicationStarted(ctx context.ConfigurableApplicationContext) {
	listener.broadcaster.BroadcastEvent(ctx, NewApplicationStartedEvent(listener.app, listener.args, ctx))
}

func (listener EventPublishRunListener) OnApplicationRunning(ctx context.ConfigurableApplicationContext) {
	listener.broadcaster.BroadcastEvent(ctx, NewApplicationReadyEvent(listener.app, listener.args, ctx))
}

func (listener EventPublishRunListener) OnApplicationFailed(ctx context.ConfigurableApplicationContext, err error) {
	listener.broadcaster.BroadcastEvent(ctx, NewApplicationFailedEvent(listener.app, listener.args, ctx, err))
}
