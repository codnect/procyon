package procyon

import (
	"procyon/app"
	"procyon/util"
)

type Application struct {
	appRunListeners app.ApplicationRunListeners
	startupLogger   app.AppStartupLogger
}

func NewProcyonApplication() *Application {
	return &Application{
		startupLogger: app.NewAppStartupLogger(),
	}
}

func (procyonApp *Application) SetApplicationRunListeners(listeners ...app.ApplicationRunListener) {
	procyonApp.appRunListeners = app.NewApplicationRunListeners(listeners)
}

func (procyonApp *Application) Run() {
	taskWatch := util.NewTaskWatch()
	_ = taskWatch.Start()
	procyonApp.appRunListeners.Starting()
	// prepare environment
	appArguments := app.NewDefaultApplicationArguments()
	environment := procyonApp.prepareEnvironment(appArguments)
	// print banner
	app.ProcyonBanner{}.PrintBanner()
	context := procyonApp.createApplicationContext()
	procyonApp.prepareContext(context, environment, appArguments)
	procyonApp.appRunListeners.Started(context)
	procyonApp.appRunListeners.Running(context)
	_ = taskWatch.Stop()
	procyonApp.startupLogger.LogStarted(taskWatch)
}

func (procyonApp *Application) prepareEnvironment(arguments app.ApplicationArguments) app.Environment {
	environment := procyonApp.createEnvironment()
	procyonApp.appRunListeners.EnvironmentPrepared(environment)
	return environment
}

func (procyonApp *Application) createEnvironment() app.ConfigurableEnvironment {
	return nil
}

func (procyonApp *Application) createApplicationContext() app.ConfigurableApplicationContext {
	return nil
}

func (procyonApp *Application) prepareContext(context app.ConfigurableApplicationContext,
	environment app.ConfigurableEnvironment,
	arguments app.ApplicationArguments) {
	procyonApp.startupLogger.LogStarting()
	procyonApp.appRunListeners.ContextPrepared(context)
	procyonApp.appRunListeners.ContextLoaded(context)
}
