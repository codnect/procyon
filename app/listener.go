package app

import (
	"procyon-test/env"
	"procyon-test/event"
	"time"
)

type StartupListener interface {
	OnStarting(ctx Context)
	OnEnvironmentPrepared(ctx Context, environment env.Environment)
	OnContextPrepared(ctx Context)
	OnContextLoaded(ctx Context)
	OnStarted(ctx Context, timeTaken time.Duration)
	OnReady(ctx Context, timeTaken time.Duration)
	OnFailed(ctx Context, err error)
}

type startupListener struct {
	app         Application
	args        []string
	broadcaster event.Broadcaster
}

func (l *startupListener) OnStarting(ctx Context) {
	l.broadcaster.BroadcastEvent(newStartingEvent(l.app, l.args, ctx))
}

func (l *startupListener) OnEnvironmentPrepared(ctx Context, environment env.Environment) {
	l.broadcaster.BroadcastEvent(newEnvironmentPreparedEvent(l.app, l.args, ctx, environment))
}

func (l *startupListener) OnContextPrepared(ctx Context) {
	l.broadcaster.BroadcastEvent(newContextPreparedEvent(l.app, l.args, ctx))
}

func (l *startupListener) OnContextLoaded(ctx Context) {
	l.broadcaster.BroadcastEvent(newContextLoadedEvent(l.app, l.args, ctx))
}

func (l *startupListener) OnStarted(ctx Context, timeTaken time.Duration) {
	l.broadcaster.BroadcastEvent(newStartedEvent(l.app, l.args, ctx, timeTaken))
}

func (l *startupListener) OnReady(ctx Context, timeTaken time.Duration) {
	l.broadcaster.BroadcastEvent(newReadyEvent(l.app, l.args, ctx, timeTaken))
}

func (l *startupListener) OnFailed(ctx Context, err error) {
	l.broadcaster.BroadcastEvent(newFailedEvent(l.app, l.args, ctx, err))
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
