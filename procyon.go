package procyon

import (
	"github.com/Rollcomp/procyon-context"
	"github.com/Rollcomp/procyon-core"
	web "github.com/Rollcomp/procyon-web"
	"os"
)

type Application struct {
	appRunListeners ApplicationRunListeners
	startupLogger   StartupLogger
}

func NewProcyonApplication() *Application {
	return &Application{
		startupLogger:   NewStartupLogger(),
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
	ApplicationBanner{}.PrintBanner()
	applicationContext := procyonApp.createApplicationContext()
	if environment != nil {

	}
	procyonApp.prepareContext(applicationContext, environment.(core.ConfigurableEnvironment), appArguments)
	procyonApp.appRunListeners.Started(applicationContext)
	procyonApp.appRunListeners.Running(applicationContext)
	_ = taskWatch.Stop()
	procyonApp.startupLogger.LogStarted(taskWatch)
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
	return nil
}

func (procyonApp *Application) prepareContext(context context.ConfigurableApplicationContext,
	environment core.ConfigurableEnvironment,
	arguments ApplicationArguments) {
	procyonApp.startupLogger.LogStarting()
	procyonApp.appRunListeners.ContextPrepared(context)
	procyonApp.appRunListeners.ContextLoaded(context)
}

var (
	appRunListeners = make([]ApplicationRunListener, 0)
)

func RegisterAppRunListener(appRunListener ...ApplicationRunListener) {
	appRunListeners = append(appRunListeners, appRunListener...)
}
