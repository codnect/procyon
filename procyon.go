package procyon

import (
	"errors"
	"github.com/google/uuid"
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

func (procyonApp *Application) createApplicationAndContextId() (uuid.UUID, uuid.UUID) {
	var err error
	var applicationId uuid.UUID
	applicationId, err = uuid.NewUUID()
	if err != nil {
		panic("Could not application id")
	}
	var contextId uuid.UUID
	contextId, err = uuid.NewUUID()
	if err != nil {
		panic("Could not context id")
	}
	return applicationId, contextId
}

func (procyonApp *Application) Run() {
	// create the application id
	applicationId, mainContextId := procyonApp.createApplicationAndContextId()

	// startup logger
	logger := core.NewSimpleLogger(applicationId.String(), mainContextId.String())
	startupLogger := NewStartupLogger(logger)

	// it is executed during panic
	defer func() {
		if r := recover(); r != nil {
			logger.Panic(r)
		}
	}()

	// start a new task to watch time which will pass
	taskWatch := core.NewTaskWatch()
	_ = taskWatch.Start()

	// print application banner
	appBanner.PrintBanner()

	// log starting
	startupLogger.LogStarting(applicationId.String(), mainContextId.String())
	appArguments := GetApplicationArguments(os.Args)

	// application listener
	err := procyonApp.initApplicationListenerInstances()
	if err != nil {
		logger.Fatal(err)
	}

	// app run listeners
	var listeners *ApplicationRunListeners
	listeners, err = procyonApp.getAppRunListenerInstances(appArguments)
	if err != nil {
		logger.Fatal(err)
	}

	// broadcast an event to inform the application is starting
	listeners.Starting()

	// prepare environment
	var environment core.Environment
	environment, err = procyonApp.prepareEnvironment(appArguments, listeners)
	if err != nil {
		logger.Fatal(err)
	}

	// create application context
	var applicationContext context.ConfigurableApplicationContext
	applicationContext, err = procyonApp.createApplicationContext(applicationId, mainContextId)
	if err != nil {
		logger.Fatal(err)
	}

	// prepare context
	err = procyonApp.prepareContext(applicationContext,
		environment.(core.ConfigurableEnvironment),
		appArguments,
		listeners,
	)
	if err != nil {
		logger.Fatal(err)
	}

	listeners.Started(applicationContext)
	listeners.Running(applicationContext)

	// configure context
	err = procyonApp.configureContext(applicationContext)
	if err != nil {
		logger.Fatal(err)
	}
	_ = taskWatch.Stop()
	startupLogger.LogStarted(taskWatch)
}

func (procyonApp *Application) prepareEnvironment(arguments ApplicationArguments, listeners *ApplicationRunListeners) (core.Environment, error) {
	environment := procyonApp.createEnvironment()
	err := procyonApp.configureEnvironment(environment, arguments)
	if err != nil {
		return nil, err
	}
	listeners.EnvironmentPrepared(environment)
	return environment, nil
}

func (procyonApp *Application) createEnvironment() core.ConfigurableEnvironment {
	return web.NewStandardWebEnvironment()
}

func (procyonApp *Application) configureEnvironment(environment core.ConfigurableEnvironment, arguments ApplicationArguments) error {
	propertySources := environment.GetPropertySources()
	if arguments != nil && len(arguments.GetSourceArgs()) > 0 {
		propertySources.Add(core.NewSimpleCommandLinePropertySource(arguments.GetSourceArgs()))
	}
	return nil
}

func (procyonApp *Application) createApplicationContext(appId uuid.UUID, contextId uuid.UUID) (context.ConfigurableApplicationContext, error) {
	return web.NewProcyonServerApplicationContext(appId, contextId), nil
}

func (procyonApp *Application) prepareContext(context context.ConfigurableApplicationContext,
	environment core.ConfigurableEnvironment,
	arguments ApplicationArguments,
	listeners *ApplicationRunListeners) error {
	// set environment
	context.SetEnvironment(environment)
	factory := context.GetPeaFactory()
	// broadcast an event to notify that context is prepared
	listeners.ContextPrepared(context)
	// register application arguments as shared pea
	err := factory.RegisterSharedPea("procyonApplicationArguments", arguments)
	if err != nil {
		return err
	}
	// broadcast an event to notify that context is loaded
	listeners.ContextLoaded(context)
	return nil
}

func (procyonApp *Application) getAppRunListenerInstances(arguments ApplicationArguments) (*ApplicationRunListeners, error) {
	instances, err := procyonApp.getInstancesWithParamTypes(core.GetType((*ApplicationRunListener)(nil)),
		[]*core.Type{core.GetType((*Application)(nil)), core.GetType((*ApplicationArguments)(nil))},
		[]interface{}{procyonApp, arguments})
	if err != nil {
		return nil, err
	}
	var listeners []ApplicationRunListener
	for _, instance := range instances {
		listeners = append(listeners, instance.(ApplicationRunListener))
	}
	return NewApplicationRunListeners(listeners), nil
}

func (procyonApp *Application) getAppListeners() []context.ApplicationListener {
	return procyonApp.listeners
}

func (procyonApp *Application) initApplicationListenerInstances() error {
	instances, err := procyonApp.getInstances(core.GetType((*context.ApplicationListener)(nil)))
	if err != nil {
		return err
	}
	listenerInstances := make([]context.ApplicationListener, len(instances))
	for index, instance := range instances {
		listenerInstances[index] = instance.(context.ApplicationListener)
	}
	procyonApp.listeners = listenerInstances
	return nil
}

func (procyonApp *Application) getInstances(typ *core.Type) (result []interface{}, err error) {
	var types []*core.Type
	types, err = core.GetComponentTypes(typ)
	if err != nil {
		return
	}
	var instances []interface{}
	for _, t := range types {
		var instance interface{}
		instance, err = peas.CreateInstance(t, []interface{}{})
		if err != nil {
			return
		}
		result = append(result, instance)
	}
	return instances, nil
}

func (procyonApp *Application) getInstancesWithParamTypes(typ *core.Type, parameterTypes []*core.Type, args []interface{}) (result []interface{}, err error) {
	var types []*core.Type
	types, err = core.GetComponentTypesWithParam(typ, parameterTypes)
	if err != nil {
		return
	}
	var instances []interface{}
	for _, t := range types {
		var instance interface{}
		instance, err = peas.CreateInstance(t, args)
		if err != nil {
			return
		}
		instances = append(instances, instance)
	}
	return instances, nil
}

func (procyonApp *Application) configureContext(ctx context.ConfigurableApplicationContext) error {
	if ctx == nil {
		return errors.New("context must not be null")
	}
	if configurableContextAdapter, ok := ctx.(context.ConfigurableContextAdapter); ok {
		configurableContextAdapter.Configure()
		return nil
	}
	return errors.New("context.ConfigurableContextAdapter methods must be implemented in your context struct")
}
