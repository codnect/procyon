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
	return []context.ApplicationEvent{}
}

func (listener BootstrapListener) OnApplicationEvent(event context.ApplicationEvent) {

}

type EventPublishRunListener struct {
	app         *Application
	broadcaster context.ApplicationEventBroadcaster
}

func NewEventPublishRunListener(app *Application) EventPublishRunListener {
	return EventPublishRunListener{
		app: app,
	}
}

func (listener EventPublishRunListener) Starting() {
	listener.broadcaster.BroadcastEvent(nil)
}

func (listener EventPublishRunListener) EnvironmentPrepared(environment core.ConfigurableEnvironment) {
	listener.broadcaster.BroadcastEvent(nil)
}

func (listener EventPublishRunListener) ContextPrepared(ctx context.ConfigurableApplicationContext) {
	listener.broadcaster.BroadcastEvent(nil)
}

func (listener EventPublishRunListener) ContextLoaded(ctx context.ConfigurableApplicationContext) {
	listener.broadcaster.BroadcastEvent(nil)
}

func (listener EventPublishRunListener) Started(ctx context.ConfigurableApplicationContext) {
	listener.broadcaster.BroadcastEvent(nil)
}

func (listener EventPublishRunListener) Running(ctx context.ConfigurableApplicationContext) {
	listener.broadcaster.BroadcastEvent(nil)
}

func (listener EventPublishRunListener) Failed(context context.ConfigurableApplicationContext, err error) {
	listener.broadcaster.BroadcastEvent(nil)
}
