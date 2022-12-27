package app

import (
	"github.com/procyon-projects/procyon/app/env"
	"time"
)

type StartingEvent struct {
	app  Application
	ctx  Context
	args *Arguments
	time time.Time
}

func newStartingEvent(app Application, args *Arguments, ctx Context) *StartingEvent {
	return &StartingEvent{
		app:  app,
		ctx:  ctx,
		args: args,
		time: time.Now(),
	}
}

func (e *StartingEvent) EventSource() any {
	return e.app
}

func (e *StartingEvent) Time() time.Time {
	return e.time
}

func (e *StartingEvent) Args() *Arguments {
	return e.args
}

func (e *StartingEvent) Application() Application {
	return e.app
}

func (e *StartingEvent) Context() Context {
	return e.ctx
}

type EnvironmentPreparedEvent struct {
	app         Application
	ctx         Context
	args        *Arguments
	environment env.Environment
	time        time.Time
}

func newEnvironmentPreparedEvent(app Application, args *Arguments, ctx Context, environment env.Environment) *EnvironmentPreparedEvent {
	return &EnvironmentPreparedEvent{
		app:         app,
		ctx:         ctx,
		args:        args,
		environment: environment,
		time:        time.Now(),
	}
}

func (e *EnvironmentPreparedEvent) EventSource() any {
	return e.app
}

func (e *EnvironmentPreparedEvent) Time() time.Time {
	return e.time
}

func (e *EnvironmentPreparedEvent) Args() *Arguments {
	return e.args
}

func (e *EnvironmentPreparedEvent) Application() Application {
	return e.app
}

func (e *EnvironmentPreparedEvent) Context() Context {
	return e.ctx
}

func (e *EnvironmentPreparedEvent) Environment() env.Environment {
	return e.environment
}

type ContextPreparedEvent struct {
	app  Application
	ctx  Context
	args *Arguments
	time time.Time
}

func newContextPreparedEvent(app Application, args *Arguments, ctx Context) *ContextPreparedEvent {
	return &ContextPreparedEvent{
		app:  app,
		ctx:  ctx,
		args: args,
		time: time.Now(),
	}
}

func (e *ContextPreparedEvent) EventSource() any {
	return e.app
}

func (e *ContextPreparedEvent) Time() time.Time {
	return e.time
}

func (e *ContextPreparedEvent) Args() *Arguments {
	return e.args
}

func (e *ContextPreparedEvent) Application() Application {
	return e.app
}

func (e *ContextPreparedEvent) Context() Context {
	return e.ctx
}

type ContextLoadedEvent struct {
	app  Application
	ctx  Context
	args *Arguments
	time time.Time
}

func newContextLoadedEvent(app Application, args *Arguments, ctx Context) *ContextLoadedEvent {
	return &ContextLoadedEvent{
		app:  app,
		ctx:  ctx,
		args: args,
		time: time.Now(),
	}
}

func (e *ContextLoadedEvent) EventSource() any {
	return e.app
}

func (e *ContextLoadedEvent) Time() time.Time {
	return e.time
}

func (e *ContextLoadedEvent) Args() *Arguments {
	return e.args
}

func (e *ContextLoadedEvent) Application() Application {
	return e.app
}

func (e *ContextLoadedEvent) Context() Context {
	return e.ctx
}

type ContextRefreshedEvent struct {
	app  Application
	ctx  Context
	args *Arguments
	time time.Time
}

func newContextRefreshedEvent(app Application, args *Arguments, ctx Context) *ContextRefreshedEvent {
	return &ContextRefreshedEvent{
		app:  app,
		ctx:  ctx,
		args: args,
		time: time.Now(),
	}
}

func (e *ContextRefreshedEvent) EventSource() any {
	return e.app
}

func (e *ContextRefreshedEvent) Time() time.Time {
	return e.time
}

func (e *ContextRefreshedEvent) Args() *Arguments {
	return e.args
}

func (e *ContextRefreshedEvent) Application() Application {
	return e.app
}

func (e *ContextRefreshedEvent) Context() Context {
	return e.ctx
}

type StartedEvent struct {
	app       Application
	ctx       Context
	args      *Arguments
	time      time.Time
	timeTaken time.Duration
}

func newStartedEvent(app Application, args *Arguments, ctx Context, timeTaken time.Duration) *StartedEvent {
	return &StartedEvent{
		app:       app,
		ctx:       ctx,
		args:      args,
		time:      time.Now(),
		timeTaken: timeTaken,
	}
}

func (e *StartedEvent) EventSource() any {
	return e.app
}

func (e *StartedEvent) Time() time.Time {
	return e.time
}

func (e *StartedEvent) Args() *Arguments {
	return e.args
}

func (e *StartedEvent) Application() Application {
	return e.app
}

func (e *StartedEvent) Context() Context {
	return e.ctx
}

func (e *StartedEvent) TimeTaken() time.Duration {
	return e.timeTaken
}

type ReadyEvent struct {
	app       Application
	ctx       Context
	args      *Arguments
	time      time.Time
	timeTaken time.Duration
}

func newReadyEvent(app Application, args *Arguments, ctx Context, timeTaken time.Duration) *ReadyEvent {
	return &ReadyEvent{
		app:       app,
		ctx:       ctx,
		args:      args,
		time:      time.Now(),
		timeTaken: timeTaken,
	}
}

func (e *ReadyEvent) EventSource() any {
	return e.app
}

func (e *ReadyEvent) Time() time.Time {
	return e.time
}

func (e *ReadyEvent) Args() *Arguments {
	return e.args
}

func (e *ReadyEvent) Application() Application {
	return e.app
}

func (e *ReadyEvent) Context() Context {
	return e.ctx
}

func (e *ReadyEvent) TimeTaken() time.Duration {
	return e.timeTaken
}

type FailedEvent struct {
	app  Application
	ctx  Context
	args *Arguments
	time time.Time
	err  error
}

func newFailedEvent(app Application, args *Arguments, ctx Context, err error) *FailedEvent {
	return &FailedEvent{
		app:  app,
		ctx:  ctx,
		args: args,
		time: time.Now(),
		err:  err,
	}
}

func (e *FailedEvent) EventSource() any {
	return e.app
}

func (e *FailedEvent) Time() time.Time {
	return e.time
}

func (e *FailedEvent) Args() *Arguments {
	return e.args
}

func (e *FailedEvent) Application() Application {
	return e.app
}

func (e *FailedEvent) Context() Context {
	return e.ctx
}

func (e *FailedEvent) Err() error {
	return e.err
}
