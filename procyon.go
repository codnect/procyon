package procyon

import (
	"os"
	"procyon/app"
	"procyon/env"
	"procyon/util"
)

type Application struct {
	appRunListeners app.ApplicationRunListeners
	startupLogger   app.StartupLogger
}

func NewProcyonApplication() *Application {
	return &Application{
		startupLogger: app.NewStartupLogger(),
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
	appArguments := app.NewDefaultApplicationArguments(os.Args)
	environment := procyonApp.prepareEnvironment(appArguments)
	// print banner
	app.ProcyonBanner{}.PrintBanner()
	context := procyonApp.createApplicationContext()
	if environment != nil {

	}
	procyonApp.prepareContext(context, environment.(env.ConfigurableEnvironment), appArguments)
	procyonApp.appRunListeners.Started(context)
	procyonApp.appRunListeners.Running(context)
	_ = taskWatch.Stop()
	procyonApp.startupLogger.LogStarted(taskWatch)
}

func (procyonApp *Application) prepareEnvironment(arguments app.ApplicationArguments) env.Environment {
	environment := procyonApp.createEnvironment()
	procyonApp.appRunListeners.EnvironmentPrepared(environment)
	return environment
}

func (procyonApp *Application) createEnvironment() env.ConfigurableEnvironment {
	return env.NewStandardEnvironment()
}

func (procyonApp *Application) createApplicationContext() app.ConfigurableApplicationContext {
	return nil
}

func (procyonApp *Application) prepareContext(context app.ConfigurableApplicationContext,
	environment env.ConfigurableEnvironment,
	arguments app.ApplicationArguments) {
	procyonApp.startupLogger.LogStarting()
	procyonApp.appRunListeners.ContextPrepared(context)
	procyonApp.appRunListeners.ContextLoaded(context)
}
