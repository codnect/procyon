package procyon

import (
	context "github.com/Rollcomp/procyon-context"
	core "github.com/Rollcomp/procyon-core"
	web "github.com/Rollcomp/procyon-web"
	"log"
	"os"
)

type Application struct {
}

func NewProcyonApplication() *Application {
	return &Application{}
}

func (procyonApp *Application) Run() {
	taskWatch := core.NewTaskWatch()
	listeners := procyonApp.getAppRunListeners()
	_ = taskWatch.Start()
	listeners.Starting()
	// prepare environment
	appArguments := GetApplicationArguments(os.Args)
	environment := procyonApp.prepareEnvironment(appArguments, listeners)
	// print banner
	appBanner.PrintBanner()
	applicationContext := procyonApp.createApplicationContext()
	procyonApp.prepareContext(applicationContext, environment.(core.ConfigurableEnvironment), appArguments, listeners)
	listeners.Started(applicationContext)
	listeners.Running(applicationContext)
	procyonApp.configureContext(applicationContext)
	_ = taskWatch.Stop()
	startupLogger.LogStarted(taskWatch)
}

func (procyonApp *Application) prepareEnvironment(arguments ApplicationArguments, listeners ApplicationRunListeners) core.Environment {
	environment := procyonApp.createEnvironment()
	procyonApp.configureEnvironment(environment, arguments)
	listeners.EnvironmentPrepared(environment)
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
	arguments ApplicationArguments, listeners ApplicationRunListeners) {
	startupLogger.LogStarting()
	context.SetEnvironment(environment)
	listeners.ContextPrepared(context)
	listeners.ContextLoaded(context)
}

func (procyonApp *Application) getAppRunListeners() ApplicationRunListeners {
	listeners := core.GetComponentTypes(core.GetType((*ApplicationRunListener)(nil)))
	log.Print(listeners)
	return NewApplicationRunListeners(nil)
}

func (procyonApp *Application) configureContext(ctx context.ConfigurableApplicationContext) {
	if ctx == nil {
		panic("Context must not be null")
	}
	if configurableContextAdapter, ok := ctx.(context.ConfigurableContextAdapter); ok {
		configurableContextAdapter.Configure()
	} else {
		panic("context.ConfigurableContextAdapter methods must be implemented in your context struct")
	}
}
