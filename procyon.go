package procyon

import (
	"os"
	"procyon/app"
	"procyon/context"
	"procyon/core"
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
	procyonApp.prepareContext(context, environment.(core.ConfigurableEnvironment), appArguments)
	procyonApp.appRunListeners.Started(context)
	procyonApp.appRunListeners.Running(context)
	_ = taskWatch.Stop()
	procyonApp.startupLogger.LogStarted(taskWatch)
}

func (procyonApp *Application) prepareEnvironment(arguments app.ApplicationArguments) core.Environment {
	environment := procyonApp.createEnvironment()
	procyonApp.appRunListeners.EnvironmentPrepared(environment)
	return environment
}

func (procyonApp *Application) createEnvironment() core.ConfigurableEnvironment {
	return core.NewStandardEnvironment()
}

func (procyonApp *Application) createApplicationContext() context.ConfigurableApplicationContext {
	return nil
}

func (procyonApp *Application) prepareContext(context context.ConfigurableApplicationContext,
	environment core.ConfigurableEnvironment,
	arguments app.ApplicationArguments) {
	procyonApp.startupLogger.LogStarting()
	procyonApp.appRunListeners.ContextPrepared(context)
	procyonApp.appRunListeners.ContextLoaded(context)
}
