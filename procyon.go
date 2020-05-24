package procyon

import (
	context "github.com/Rollcomp/procyon-context"
	core "github.com/Rollcomp/procyon-core"
	web "github.com/Rollcomp/procyon-web"
	"os"
)

type Application struct {
	appRunListeners ApplicationRunListeners
}

func NewProcyonApplication() *Application {
	return &Application{
		appRunListeners: NewApplicationRunListeners(appRunListeners),
	}
}

func (procyonApp *Application) Run() {
	taskWatch := core.NewTaskWatch()
	_ = taskWatch.Start()
	procyonApp.appRunListeners.Starting()
	// prepare environment
	appArguments := GetApplicationArguments(os.Args)
	environment := procyonApp.prepareEnvironment(appArguments)
	// print banner
	appBanner.PrintBanner()
	applicationContext := procyonApp.createApplicationContext()
	procyonApp.prepareContext(applicationContext, environment.(core.ConfigurableEnvironment), appArguments)
	procyonApp.appRunListeners.Started(applicationContext)
	procyonApp.appRunListeners.Running(applicationContext)
	_ = taskWatch.Stop()
	startupLogger.LogStarted(taskWatch)
}

func (procyonApp *Application) prepareEnvironment(arguments ApplicationArguments) core.Environment {
	environment := procyonApp.createEnvironment()
	procyonApp.configureEnvironment(environment, arguments)
	procyonApp.appRunListeners.EnvironmentPrepared(environment)
	return environment
}

func (procyonApp *Application) createEnvironment() core.ConfigurableEnvironment {
	return web.NewStandardWebEnvironment()
}

func (procyonApp *Application) configureEnvironment(environment core.ConfigurableEnvironment, arguments ApplicationArguments) {
	propertySources := environment.GetPropertySources()
	if arguments != nil && len(arguments.GetSourceArgs()) > 0 {
		propertySources.Add(core.NewSimpleCommandLinePropertySource(arguments.GetSourceArgs()))
	}
}

func (procyonApp *Application) createApplicationContext() context.ConfigurableApplicationContext {
	return web.NewProcyonServerApplicationContext()
}

func (procyonApp *Application) prepareContext(context context.ConfigurableApplicationContext,
	environment core.ConfigurableEnvironment,
	arguments ApplicationArguments) {
	startupLogger.LogStarting()
	context.SetEnvironment(environment)
	procyonApp.appRunListeners.ContextPrepared(context)
	procyonApp.appRunListeners.ContextLoaded(context)
}

var (
	appRunListeners = make([]ApplicationRunListener, 0)
)

func RegisterAppRunListener(appRunListener ...ApplicationRunListener) {
	appRunListeners = append(appRunListeners, appRunListener...)
}
