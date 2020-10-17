package procyon

import (
	"errors"
	"fmt"
	"github.com/codnect/goo"
	"github.com/google/uuid"
	context "github.com/procyon-projects/procyon-context"
	core "github.com/procyon-projects/procyon-core"
	peas "github.com/procyon-projects/procyon-peas"
	web "github.com/procyon-projects/procyon-web"
	"os"
)

type Application struct {
	listeners           []context.ApplicationListener
	contextInitializers []context.ApplicationContextInitializer
}

func NewProcyonApplication() *Application {
	return &Application{
		listeners:           make([]context.ApplicationListener, 0),
		contextInitializers: make([]context.ApplicationContextInitializer, 0),
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
	logger := context.NewSimpleLogger(mainContextId.String())
	startupLogger := NewStartupLogger(logger)

	// it is executed during panic
	defer func() {
		if r := recover(); r != nil {
			logger.P(mainContextId.String(), r)
		}
	}()
	// start a new task to watch time which will pass
	taskWatch := core.NewTaskWatch()
	_ = taskWatch.Start()

	// print application banner
	appBanner.PrintBanner()

	// log starting
	startupLogger.LogStarting(applicationId.String(), mainContextId.String())

	// get the application arguments
	appArguments := GetApplicationArguments(os.Args)

	argumentComponentScan := appArguments.GetOptionValues("procyon.component.scan")
	scanComponents := true
	if argumentComponentScan != nil && len(argumentComponentScan) == 1 && argumentComponentScan[0] == "false" {
		scanComponents = false
	}

	if scanComponents {
		// scan components
		err := procyonApp.scanComponents(mainContextId.String(), logger)
		if err != nil {
			logger.F(mainContextId.String(), err)
		}
	}

	// application listener
	err := procyonApp.initApplicationListenerInstances()
	if err != nil {
		logger.F(mainContextId.String(), err)
	}

	// application context initializers
	err = procyonApp.initApplicationContextInitializers()
	if err != nil {
		logger.F(mainContextId.String(), err)
	}

	// app run listeners
	var listeners *ApplicationRunListeners
	listeners, err = procyonApp.getAppRunListenerInstances(logger, appArguments)
	if err != nil {
		logger.F(mainContextId.String(), err)
	}

	// broadcast an event to inform the application is starting
	listeners.Starting()

	// prepare environment
	var environment core.Environment
	environment, err = procyonApp.prepareEnvironment(appArguments, listeners)
	if err != nil {
		logger.F(mainContextId.String(), err)
	}

	// create application context
	var applicationContext context.ConfigurableApplicationContext
	applicationContext, err = procyonApp.createApplicationContext(applicationId, mainContextId)
	if err != nil {
		logger.Fatal(applicationContext, err)
	}

	// prepare context
	err = procyonApp.prepareContext(applicationContext,
		environment.(core.ConfigurableEnvironment),
		appArguments,
		listeners,
		logger,
	)
	if err != nil {
		logger.Fatal(applicationContext, err)
	}

	listeners.Started(applicationContext)
	listeners.Running(applicationContext)

	// configure context
	err = procyonApp.configureContext(applicationContext)
	if err != nil {
		logger.Fatal(applicationContext, err)
	}
	_ = taskWatch.Stop()
	startupLogger.LogStarted(mainContextId.String(), taskWatch)
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

func (procyonApp *Application) scanComponents(contextId string, logger context.Logger) error {
	logger.I(contextId, "Scanning components...")
	componentScanner := newComponentScanner()
	componentCount, err := componentScanner.scan(contextId, logger)
	if err != nil {
		return err
	}
	logger.I(contextId, fmt.Sprintf("Found (%d) components.", componentCount))
	return nil
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
	listeners *ApplicationRunListeners,
	logger context.Logger) error {
	// set environment
	context.SetEnvironment(environment)
	// set logger
	context.SetLogger(logger)
	factory := context.GetPeaFactory()

	// apply context initializers
	for _, contextInitializer := range procyonApp.getAppContextInitializers() {
		contextInitializer.InitializeContext(context)
	}

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

func (procyonApp *Application) getAppRunListenerInstances(logger context.Logger, arguments ApplicationArguments) (*ApplicationRunListeners, error) {
	instances, err := procyonApp.getInstancesWithParamTypes(goo.GetType((*ApplicationRunListener)(nil)),
		[]goo.Type{goo.GetType((*context.Logger)(nil)), goo.GetType((*Application)(nil)), goo.GetType((*ApplicationArguments)(nil))},
		[]interface{}{logger, procyonApp, arguments})
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

func (procyonApp *Application) getAppContextInitializers() []context.ApplicationContextInitializer {
	return procyonApp.contextInitializers
}

func (procyonApp *Application) initApplicationListenerInstances() error {
	instances, err := procyonApp.getInstances(goo.GetType((*context.ApplicationListener)(nil)))
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

func (procyonApp *Application) initApplicationContextInitializers() error {
	instances, err := procyonApp.getInstances(goo.GetType((*context.ApplicationContextInitializer)(nil)))
	if err != nil {
		return err
	}
	initializerInstances := make([]context.ApplicationContextInitializer, len(instances))
	for index, instance := range instances {
		initializerInstances[index] = instance.(context.ApplicationContextInitializer)
	}
	procyonApp.contextInitializers = initializerInstances
	return nil
}

func (procyonApp *Application) getInstances(typ goo.Type) (result []interface{}, err error) {
	var types []goo.Type
	types, err = core.GetComponentTypes(typ)
	if err != nil {
		return
	}
	for _, t := range types {
		var instance interface{}
		instance, err = peas.CreateInstance(t, []interface{}{})
		if err != nil {
			return
		}
		if instance != nil {
			result = append(result, instance)
		} else {
			err = errors.New("Instance cannot be created by using the method " + t.GetName())
		}
	}
	return
}

func (procyonApp *Application) getInstancesWithParamTypes(typ goo.Type, parameterTypes []goo.Type, args []interface{}) (result []interface{}, err error) {
	var types []goo.Type
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
