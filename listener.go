package procyon

import (
	"codnect.io/procyon-core/runtime"
	"codnect.io/procyon-core/runtime/env"
	"time"
)

type StartupListener struct {
	app  Application
	args *runtime.Arguments
}

func newStartupListener(app Application, args *runtime.Arguments) *StartupListener {
	return &StartupListener{
		app:  app,
		args: args,
	}
}

func (l *StartupListener) OnStarting(ctx runtime.Context) {
	//l.broadcaster.BroadcastEvent(ctx, newStartingEvent(l.app, l.args, ctx))
}

func (l *StartupListener) OnEnvironmentPrepared(ctx runtime.Context, environment env.Environment) {
	//l.broadcaster.BroadcastEvent(ctx, newEnvironmentPreparedEvent(l.app, l.args, ctx, environment))
}

func (l *StartupListener) OnContextPrepared(ctx runtime.Context) {
	//l.broadcaster.BroadcastEvent(ctx, newContextPreparedEvent(l.app, l.args, ctx))
}

func (l *StartupListener) OnContextLoaded(ctx runtime.Context) {
	//l.broadcaster.BroadcastEvent(ctx, newContextLoadedEvent(l.app, l.args, ctx))
	//ctx.PublishEvent(ctx, availability.NewChangeEvent(ctx, availability.StateCorrect))
}

func (l *StartupListener) OnContextStarted(ctx runtime.Context) {
	//l.broadcaster.BroadcastEvent(ctx, newContextStartedEvent(l.app, l.args, ctx))
}

func (l *StartupListener) OnStarted(ctx runtime.Context, timeTaken time.Duration) {
	//l.broadcaster.BroadcastEvent(ctx, newStartedEvent(l.app, l.args, ctx, timeTaken))
	//ctx.PublishEvent(ctx, availability.NewChangeEvent(ctx, availability.StateAcceptingTraffic))
}

func (l *StartupListener) OnReady(ctx runtime.Context, timeTaken time.Duration) {
	//l.broadcaster.BroadcastEvent(ctx, newReadyEvent(l.app, l.args, ctx, timeTaken))
}

func (l *StartupListener) OnFailed(ctx runtime.Context, err error) {
	//l.broadcaster.BroadcastEvent(ctx, newFailedEvent(l.app, l.args, ctx, err))
}
