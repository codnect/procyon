package procyon

import (
	"github.com/Rollcomp/procyon-context"
	"github.com/Rollcomp/procyon-core"
	"os"
)

type Application struct {
	appRunListeners ApplicationRunListeners
	startupLogger   StartupLogger
}

func NewProcyonApplication() *Application {
	return &Application{
		startupLogger: NewStartupLogger(),
	}
}

func (procyonApp *Application) SetApplicationRunListeners(listeners ...ApplicationRunListener) {
	procyonApp.appRunListeners = NewApplicationRunListeners(listeners)
}

func (procyonApp *Application) Run() {
	taskWatch := core.NewTaskWatch()
	_ = taskWatch.Start()
	procyonApp.appRunListeners.Starting()
	// prepare environment
	appArguments := NewDefaultApplicationArguments(os.Args)
	environment := procyonApp.prepareEnvironment(appArguments)
	// print banner
	ApplicationBanner{}.PrintBanner()
	context := procyonApp.createApplicationContext()
	if environment != nil {

	}
	procyonApp.prepareContext(context, environment.(core.ConfigurableEnvironment), appArguments)
	procyonApp.appRunListeners.Started(context)
	procyonApp.appRunListeners.Running(context)
	_ = taskWatch.Stop()
	procyonApp.startupLogger.LogStarted(taskWatch)
}

func (procyonApp *Application) prepareEnvironment(arguments ApplicationArguments) core.Environment {
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
	arguments ApplicationArguments) {
	procyonApp.startupLogger.LogStarting()
	procyonApp.appRunListeners.ContextPrepared(context)
	procyonApp.appRunListeners.ContextLoaded(context)
}
