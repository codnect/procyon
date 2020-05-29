package procyon

import (
	context "github.com/Rollcomp/procyon-context"
	core "github.com/Rollcomp/procyon-core"
)

type BaseProcyonApplicationEvent struct {
	context.BaseApplicationEvent
	args ApplicationArguments
}

func NewBaseProcyonApplicationEvent(app *Application, args ApplicationArguments) BaseProcyonApplicationEvent {
	return BaseProcyonApplicationEvent{
		context.NewBaseApplicationEvent(app),
		args,
	}
}

func (e BaseProcyonApplicationEvent) GetProcyonApplication() *Application {
	return e.GetSource().(*Application)
}

func (e BaseProcyonApplicationEvent) GetArgs() ApplicationArguments {
	return e.args
}

type ApplicationStartingEvent struct {
	BaseProcyonApplicationEvent
}

func NewApplicationStarting(app *Application, args ApplicationArguments) ApplicationStartingEvent {
	return ApplicationStartingEvent{
		NewBaseProcyonApplicationEvent(app, args),
	}
}

type ApplicationEnvironmentPreparedEvent struct {
	BaseProcyonApplicationEvent
	environment core.ConfigurableEnvironment
}

func NewApplicationEnvironmentPreparedEvent(app *Application, args ApplicationArguments, env core.ConfigurableEnvironment) ApplicationEnvironmentPreparedEvent {
	return ApplicationEnvironmentPreparedEvent{
		NewBaseProcyonApplicationEvent(app, args),
		env,
	}
}

func (e ApplicationEnvironmentPreparedEvent) GetEnvironment() core.ConfigurableEnvironment {
	return e.environment
}

type ApplicationContextInitializedEvent struct {
	BaseProcyonApplicationEvent
	context context.ConfigurableApplicationContext
}

func NewApplicationContextInitializedEvent(app *Application, args ApplicationArguments, ctx context.ConfigurableApplicationContext) ApplicationContextInitializedEvent {
	return ApplicationContextInitializedEvent{
		NewBaseProcyonApplicationEvent(app, args),
		ctx,
	}
}

func (e ApplicationContextInitializedEvent) GetApplicationContext() context.ConfigurableApplicationContext {
	return e.context
}

type ApplicationPreparedEvent struct {
	BaseProcyonApplicationEvent
	context context.ConfigurableApplicationContext
}

func NewApplicationPreparedEvent(app *Application, args ApplicationArguments, ctx context.ConfigurableApplicationContext) ApplicationPreparedEvent {
	return ApplicationPreparedEvent{
		NewBaseProcyonApplicationEvent(app, args),
		ctx,
	}
}

func (e ApplicationPreparedEvent) GetApplicationContext() context.ConfigurableApplicationContext {
	return e.context
}

type ApplicationStartedEvent struct {
	BaseProcyonApplicationEvent
	context context.ConfigurableApplicationContext
}

func NewApplicationStartedEvent(app *Application, args ApplicationArguments, ctx context.ConfigurableApplicationContext) ApplicationStartedEvent {
	return ApplicationStartedEvent{
		NewBaseProcyonApplicationEvent(app, args),
		ctx,
	}
}

func (e ApplicationStartedEvent) GetApplicationContext() context.ConfigurableApplicationContext {
	return e.context
}

type ApplicationReadyEvent struct {
	BaseProcyonApplicationEvent
	context context.ConfigurableApplicationContext
}

func NewApplicationReadyEvent(app *Application, args ApplicationArguments, ctx context.ConfigurableApplicationContext) ApplicationReadyEvent {
	return ApplicationReadyEvent{
		NewBaseProcyonApplicationEvent(app, args),
		ctx,
	}
}

func (e ApplicationReadyEvent) GetApplicationContext() context.ConfigurableApplicationContext {
	return e.context
}

type ApplicationFailedEvent struct {
	BaseProcyonApplicationEvent
	context context.ConfigurableApplicationContext
	err     error
}

func NewApplicationFailedEvent(app *Application, args ApplicationArguments, ctx context.ConfigurableApplicationContext, err error) ApplicationFailedEvent {
	return ApplicationFailedEvent{
		NewBaseProcyonApplicationEvent(app, args),
		ctx,
		err,
	}
}

func (e ApplicationFailedEvent) GetApplicationContext() context.ConfigurableApplicationContext {
	return e.context
}

func (e ApplicationFailedEvent) GetError() error {
	return e.err
}
