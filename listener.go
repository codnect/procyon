package procyon

import (
	"codnect.io/procyon-core/runtime"
	"codnect.io/procyon-core/runtime/env"
	"time"
)

type startupListeners []runtime.StartupListener

func (l startupListeners) starting(ctx runtime.Context) {
	for _, listener := range l {
		listener.OnStarting(ctx)
	}
}

func (l startupListeners) environmentPrepared(ctx runtime.Context, environment env.Environment) {
	for _, listener := range l {
		listener.OnEnvironmentPrepared(ctx, environment)
	}
}

func (l startupListeners) contextPrepared(ctx runtime.Context) {
	for _, listener := range l {
		listener.OnContextPrepared(ctx)
	}
}

func (l startupListeners) contextLoaded(ctx runtime.Context) {
	for _, listener := range l {
		listener.OnContextLoaded(ctx)
	}
}

func (l startupListeners) contextStarted(ctx runtime.Context) {
	for _, listener := range l {
		listener.OnContextStarted(ctx)
	}
}

func (l startupListeners) started(ctx runtime.Context, timeTaken time.Duration) {
	for _, listener := range l {
		listener.OnStarted(ctx, timeTaken)
	}
}

func (l startupListeners) ready(ctx runtime.Context, timeTaken time.Duration) {
	for _, listener := range l {
		listener.OnReady(ctx, timeTaken)
	}
}

func (l startupListeners) failed(ctx runtime.Context, err error) {
	for _, listener := range l {
		listener.OnFailed(ctx, err)
	}
}
