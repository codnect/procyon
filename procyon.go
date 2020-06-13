package procyon

import (
	context "github.com/procyon-projects/procyon-context"
	core "github.com/procyon-projects/procyon-core"
	peas "github.com/procyon-projects/procyon-peas"
	web "github.com/procyon-projects/procyon-web"
	"os"
)

type Application struct {
	listeners []context.ApplicationListener
}

func NewProcyonApplication() *Application {
	return &Application{
		listeners: make([]context.ApplicationListener, 0),
	}
}

func (procyonApp *Application) Run() {
	taskWatch := core.NewTaskWatch()
	_ = taskWatch.Start()
	// print banner
	appBanner.PrintBanner()
	// log starting
	startupLogger.LogStarting()
	appArguments := GetApplicationArguments(os.Args)
	procyonApp.initApplicationListenerInstances()
	listeners := procyonApp.getAppRunListenerInstances(appArguments)
	listeners.Starting()
	// prepare environment
	environment := procyonApp.prepareEnvironment(appArguments, listeners)
	applicationContext := procyonApp.createApplicationContext()
	// prepare context
	procyonApp.prepareContext(applicationContext, environment.(core.ConfigurableEnvironment), appArguments, listeners)
	listeners.Started(applicationContext)
	listeners.Running(applicationContext)
	procyonApp.configureContext(applicationContext)
	_ = taskWatch.Stop()
	startupLogger.LogStarted(taskWatch)
}

func (procyonApp *Application) prepareEnvironment(arguments ApplicationArguments, listeners ApplicationRunListeners) core.Environment {
	core.Log.Debug("Started to prepare the application environment.")
	environment := procyonApp.createEnvironment()
	procyonApp.configureEnvironment(environment, arguments)
	listeners.EnvironmentPrepared(environment)
	return environment
}

func (procyonApp *Application) createEnvironment() core.ConfigurableEnvironment {
	return web.NewStandardWebEnvironment()
}

func (procyonApp *Application) configureEnvironment(environment core.ConfigurableEnvironment, arguments ApplicationArguments) {
	core.Log.Debug("Configuring the environment.")
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
	core.Log.Debug("Started to prepare the application context.")
	context.SetEnvironment(environment)
	// broadcast an event to notify that context is prepared
	listeners.ContextPrepared(context)
	// register application arguments as shared pea
	factory := context.GetPeaFactory()
	factory.RegisterSharedPea("procyonApplicationArguments", arguments)
	// broadcast an event to notify that context is loaded
	listeners.ContextLoaded(context)
}

func (procyonApp *Application) getAppRunListenerInstances(arguments ApplicationArguments) ApplicationRunListeners {
	instances := procyonApp.getInstancesWithParamTypes(core.GetType((*ApplicationRunListener)(nil)),
		[]*core.Type{core.GetType((*Application)(nil)), core.GetType((*ApplicationArguments)(nil))},
		[]interface{}{procyonApp, arguments})
	var listeners []ApplicationRunListener
	for _, instance := range instances {
		listeners = append(listeners, instance.(ApplicationRunListener))
	}
	return NewApplicationRunListeners(listeners)
}

func (procyonApp *Application) getAppListeners() []context.ApplicationListener {
	return procyonApp.listeners
}

func (procyonApp *Application) initApplicationListenerInstances() {
	instances := procyonApp.getInstances(core.GetType((*context.ApplicationListener)(nil)))
	listenerInstances := make([]context.ApplicationListener, len(instances))
	for index, instance := range instances {
		listenerInstances[index] = instance.(context.ApplicationListener)
	}
	procyonApp.listeners = listenerInstances
}

func (procyonApp *Application) getInstances(typ *core.Type) []interface{} {
	types := core.GetComponentTypes(typ)
	var instances []interface{}
	for _, t := range types {
		instance := peas.CreateInstance(t, []interface{}{})
		instances = append(instances, instance)
	}
	return instances
}

func (procyonApp *Application) getInstancesWithParamTypes(typ *core.Type, parameterTypes []*core.Type, args []interface{}) []interface{} {
	types := core.GetComponentTypesWithParam(typ, parameterTypes)
	var instances []interface{}
	for _, t := range types {
		instance := peas.CreateInstance(t, args)
		instances = append(instances, instance)
	}
	return instances
}

func (procyonApp *Application) configureContext(ctx context.ConfigurableApplicationContext) {
	core.Log.Debug("Configuring the application context.")
	if ctx == nil {
		core.Log.Panic("Context must not be null")
	}
	if configurableContextAdapter, ok := ctx.(context.ConfigurableContextAdapter); ok {
		configurableContextAdapter.Configure()
	} else {
		core.Log.Panic("context.ConfigurableContextAdapter methods must be implemented in your context struct")
	}
}
