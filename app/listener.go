package app

import (
	"codnect.io/procyon-core/env"
	"codnect.io/procyon/app/availability"
	"codnect.io/procyon/app/event"
	"time"
)

type StartupListener interface {
	OnStarting(ctx Context)
	OnEnvironmentPrepared(ctx Context, environment env.Environment)
	OnContextPrepared(ctx Context)
	OnContextLoaded(ctx Context)
	OnContextStarted(ctx Context)
	OnStarted(ctx Context, timeTaken time.Duration)
	OnReady(ctx Context, timeTaken time.Duration)
	OnFailed(ctx Context, err error)
}

type startupListener struct {
	app         Application
	args        *Arguments
	broadcaster event.Broadcaster
}

func newStartupListener(app Application, args *Arguments) *startupListener {
	return &startupListener{
		app:         app,
		args:        args,
		broadcaster: event.NewBroadcaster(),
	}
}

func (l *startupListener) OnStarting(ctx Context) {
	l.broadcaster.BroadcastEvent(ctx, newStartingEvent(l.app, l.args, ctx))
}

func (l *startupListener) OnEnvironmentPrepared(ctx Context, environment env.Environment) {
	l.broadcaster.BroadcastEvent(ctx, newEnvironmentPreparedEvent(l.app, l.args, ctx, environment))
}

func (l *startupListener) OnContextPrepared(ctx Context) {
	l.broadcaster.BroadcastEvent(ctx, newContextPreparedEvent(l.app, l.args, ctx))
}

func (l *startupListener) OnContextLoaded(ctx Context) {
	l.broadcaster.BroadcastEvent(ctx, newContextLoadedEvent(l.app, l.args, ctx))
	ctx.PublishEvent(ctx, availability.NewChangeEvent(ctx, availability.StateCorrect))
}

func (l *startupListener) OnContextStarted(ctx Context) {
	l.broadcaster.BroadcastEvent(ctx, newContextStartedEvent(l.app, l.args, ctx))
}

func (l *startupListener) OnStarted(ctx Context, timeTaken time.Duration) {
	l.broadcaster.BroadcastEvent(ctx, newStartedEvent(l.app, l.args, ctx, timeTaken))
	ctx.PublishEvent(ctx, availability.NewChangeEvent(ctx, availability.StateAcceptingTraffic))
}

func (l *startupListener) OnReady(ctx Context, timeTaken time.Duration) {
	l.broadcaster.BroadcastEvent(ctx, newReadyEvent(l.app, l.args, ctx, timeTaken))
}

func (l *startupListener) OnFailed(ctx Context, err error) {
	l.broadcaster.BroadcastEvent(ctx, newFailedEvent(l.app, l.args, ctx, err))
}

type startupListeners []StartupListener

func (l startupListeners) starting(ctx Context) {
	for _, listener := range l {
		listener.OnStarting(ctx)
	}
}

func (l startupListeners) environmentPrepared(ctx Context, environment env.Environment) {
	for _, listener := range l {
		listener.OnEnvironmentPrepared(ctx, environment)
	}
}

func (l startupListeners) contextPrepared(ctx Context) {
	for _, listener := range l {
		listener.OnContextPrepared(ctx)
	}
}

func (l startupListeners) contextLoaded(ctx Context) {
	for _, listener := range l {
		listener.OnContextLoaded(ctx)
	}
}

func (l startupListeners) contextStarted(ctx Context) {
	for _, listener := range l {
		listener.OnContextStarted(ctx)
	}
}

func (l startupListeners) started(ctx Context, timeTaken time.Duration) {
	for _, listener := range l {
		listener.OnStarted(ctx, timeTaken)
	}
}

func (l startupListeners) ready(ctx Context, timeTaken time.Duration) {
	for _, listener := range l {
		listener.OnReady(ctx, timeTaken)
	}
}

func (l startupListeners) failed(ctx Context, err error) {
	for _, listener := range l {
		listener.OnFailed(ctx, err)
	}
}
