package procyon

import (
	"errors"
	"fmt"
	"github.com/procyon-projects/goo"
	context "github.com/procyon-projects/procyon-context"
	core "github.com/procyon-projects/procyon-core"
	web "github.com/procyon-projects/procyon-web"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
)

var bannerText = []string{"",
	"   ___",
	"  / _ \\  _ __   ___    ___  _   _   ___   _ __",
	" / /_)/ | '__| / _ \\  / __|| | | | / _ \\ | '_ \\",
	"/ ___/  | |   | (_) || (__ | |_| || (_) || | | |",
	"\\/     |_|     \\___/  \\___| \\__, | \\___/ |_| |_|",
	"                            |___/",
}

type application interface {
	getLogger() context.Logger
	getLoggingProperties(arguments ApplicationArguments) *context.LoggingProperties
	configureLogger(logger context.Logger, loggingProperties *context.LoggingProperties)
	getTaskWatch() *core.TaskWatch
	getApplicationId() context.ApplicationId
	getContextId() context.ContextId
	generateApplicationAndContextId()
	getApplicationArguments() ApplicationArguments
	printBanner()
	logStarting()
	scanComponents(arguments ApplicationArguments) error
	prepareEnvironment(arguments ApplicationArguments, listeners *ApplicationRunListeners) (core.Environment, error)
	prepareContext(environment core.ConfigurableEnvironment, arguments ApplicationArguments, listeners *ApplicationRunListeners) (context.ConfigurableApplicationContext, error)
	getApplicationRunListenerInstances(arguments ApplicationArguments) (*ApplicationRunListeners, error)
	getApplicationListeners() []context.ApplicationListener
	getApplicationContextInitializers() []context.ApplicationContextInitializer
	initApplicationListenerInstances() error
	initApplicationContextInitializers() error
	logStarted()
	invokeApplicationRunners(ctx context.ApplicationContext, arguments ApplicationArguments)
	finish()
}

type environmentProvider interface {
	getNewEnvironment() core.ConfigurableEnvironment
}

type contextProvider interface {
	getNewContext(applicationId context.ApplicationId, contextId context.ContextId) context.ConfigurableApplicationContext
}

type ProcyonApplication struct {
	application
}

func NewProcyonApplication() *ProcyonApplication {
	baseApplication := newBaseApplication()
	app := &ProcyonApplication{
		baseApplication,
	}
	baseApplication.procyonApplication = app
	return app
}

func (procyonApp *ProcyonApplication) Run() {
	taskWatch := procyonApp.getTaskWatch()
	taskWatch.Start()

	// get the application arguments
	arguments := procyonApp.getApplicationArguments()

	logger := procyonApp.getLogger()
	loggerProperties := procyonApp.getLoggingProperties(arguments)
	procyonApp.configureLogger(logger, loggerProperties)

	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case error:
				err := r.(error)
				errorString := err.Error()
				logger.Fatal(procyonApp.getContextId(), errorString+"\n"+string(debug.Stack()))
			case string:
				errorString := r.(string)
				logger.Fatal(procyonApp.getContextId(), errorString+"\n"+string(debug.Stack()))
			default:
				logger.Error(procyonApp.getContextId(), r)
				logger.Fatal(procyonApp.getContextId(), string(debug.Stack()))
			}
		}
	}()

	procyonApp.printBanner()

	// log starting
	procyonApp.logStarting()

	// scan components
	err := procyonApp.scanComponents(arguments)
	if err != nil {
		logger.Fatal(procyonApp.getContextId(), err)
	}

	// application listener
	err = procyonApp.initApplicationListenerInstances()
	if err != nil {
		logger.Fatal(procyonApp.getContextId(), err)
	}

	// application context initializers
	err = procyonApp.initApplicationContextInitializers()
	if err != nil {
		logger.Fatal(procyonApp.getContextId(), err)
	}

	// app run listeners
	var listeners *ApplicationRunListeners
	listeners, err = procyonApp.getApplicationRunListenerInstances(arguments)
	if err != nil {
		logger.Fatal(procyonApp.getContextId(), err)
	}

	// broadcast an event to inform the application is starting
	listeners.OnApplicationStarting()

	// prepare environment
	var environment core.Environment
	environment, err = procyonApp.prepareEnvironment(arguments, listeners)
	if err != nil {
		logger.Fatal(procyonApp.getContextId(), err)
	}

	// prepare context
	var applicationContext context.ConfigurableApplicationContext
	applicationContext, err = procyonApp.prepareContext(environment.(core.ConfigurableEnvironment),
		arguments,
		listeners,
	)
	if err != nil {
		logger.Fatal(procyonApp.getContextId(), err)
	}
	taskWatch.Stop()
	procyonApp.logStarted()

	listeners.OnApplicationStarted(applicationContext)
	procyonApp.invokeApplicationRunners(applicationContext, arguments)
	listeners.OnApplicationRunning(applicationContext)

	procyonApp.finish()
}

type baseApplication struct {
	procyonApplication  *ProcyonApplication
	applicationId       context.ApplicationId
	contextId           context.ContextId
	logger              context.Logger
	customLogger        context.Logger
	taskWatch           *core.TaskWatch
	listeners           []context.ApplicationListener
	contextInitializers []context.ApplicationContextInitializer
	contextProvider     contextProvider
	environmentProvider environmentProvider
}

func newBaseApplication() *baseApplication {
	baseApplication := &baseApplication{
		listeners:           make([]context.ApplicationListener, 0),
		contextInitializers: make([]context.ApplicationContextInitializer, 0),
		taskWatch:           core.NewTaskWatch(),
		logger:              context.NewSimpleLogger(),
		contextProvider:     newDefaultContextProvider(),
		environmentProvider: newDefaultEnvironmentProvider(),
	}
	baseApplication.generateApplicationAndContextId()

	return baseApplication
}

func (application *baseApplication) getLogger() context.Logger {
	if application.customLogger != nil {
		return application.customLogger
	}
	return application.logger
}

func (application *baseApplication) getTaskWatch() *core.TaskWatch {
	return application.taskWatch
}

func (application *baseApplication) getApplicationId() context.ApplicationId {
	return application.applicationId
}

func (application *baseApplication) getContextId() context.ContextId {
	return application.contextId
}

func (application *baseApplication) printBanner() {
	for _, line := range bannerText {
		fmt.Println(line)
	}
}

func (application *baseApplication) getApplicationArguments() ApplicationArguments {
	return getApplicationArguments(os.Args)
}

func (application *baseApplication) generateApplicationAndContextId() {
	var applicationId [36]byte
	core.GenerateUUID(applicationId[:])
	var contextId [36]byte
	core.GenerateUUID(contextId[:])

	application.applicationId = context.ApplicationId(applicationId[:])
	application.contextId = context.ContextId(contextId[:])
}

func (application *baseApplication) prepareEnvironment(arguments ApplicationArguments, listeners *ApplicationRunListeners) (core.Environment, error) {
	environment := application.environmentProvider.getNewEnvironment()

	propertySources := environment.GetPropertySources()
	if arguments != nil && len(arguments.GetSourceArgs()) > 0 {
		propertySources.Add(core.NewSimpleCommandLinePropertySource(arguments.GetSourceArgs()))
	}

	listeners.OnApplicationEnvironmentPrepared(environment)
	return environment, nil
}

func (application *baseApplication) scanComponents(arguments ApplicationArguments) error {
	if arguments == nil {
		return nil
	}
	argumentComponentScan := arguments.GetOptionValues("procyon.component.scan")
	if argumentComponentScan != nil && len(argumentComponentScan) == 1 && argumentComponentScan[0] == "false" {
		return nil
	}

	application.logger.Info(application.contextId, "Scanning components...")
	componentScanner := newComponentScanner()
	componentCount, err := componentScanner.scan(application.contextId, application.logger)
	if err != nil {
		return err
	}

	application.logger.Info(application.contextId, fmt.Sprintf("Found (%d) components.", componentCount))
	return nil
}

func (application *baseApplication) prepareContext(environment core.ConfigurableEnvironment,
	arguments ApplicationArguments,
	listeners *ApplicationRunListeners) (context.ConfigurableApplicationContext, error) {

	applicationContext := application.contextProvider.getNewContext(application.applicationId, application.contextId)

	if applicationContext == nil {
		return nil, errors.New("context could not be created")
	}

	// set environment
	applicationContext.SetEnvironment(environment)
	// set logger
	applicationContext.SetLogger(application.logger)
	factory := applicationContext.GetPeaFactory()

	// apply context initializers
	for _, contextInitializer := range application.getApplicationContextInitializers() {
		contextInitializer.InitializeContext(applicationContext)
	}
	factory.ExcludeType(goo.GetType((*ApplicationRunListener)(nil)))

	// broadcast an event to notify that context is prepared
	listeners.OnApplicationContextPrepared(applicationContext)

	// register application arguments as shared pea
	err := factory.RegisterSharedPea("procyonApplicationArguments", arguments)
	if err != nil {
		return nil, err
	}
	// broadcast an event to notify that context is loaded
	listeners.OnApplicationContextLoaded(applicationContext)

	if configurableContextAdapter, ok := applicationContext.(context.ConfigurableContextAdapter); ok {
		configurableContextAdapter.Configure()
		return applicationContext, nil
	}
	return nil, errors.New("context.ConfigurableContextAdapter methods must be implemented in your context struct")
}

func (application *baseApplication) getApplicationRunListenerInstances(arguments ApplicationArguments) (*ApplicationRunListeners, error) {
	instances, err := getInstancesWithParamTypes(goo.GetType((*ApplicationRunListener)(nil)),
		[]goo.Type{goo.GetType((*ProcyonApplication)(nil)), goo.GetType((*ApplicationArguments)(nil))},
		[]interface{}{application.procyonApplication, arguments})
	if err != nil {
		return nil, err
	}
	var listeners []ApplicationRunListener
	for _, instance := range instances {
		listeners = append(listeners, instance.(ApplicationRunListener))
	}
	return NewApplicationRunListeners(listeners), nil
}

func (application *baseApplication) getApplicationListeners() []context.ApplicationListener {
	return application.listeners
}

func (application *baseApplication) getApplicationContextInitializers() []context.ApplicationContextInitializer {
	return application.contextInitializers
}

func (application *baseApplication) initApplicationListenerInstances() error {
	instances, err := getInstances(goo.GetType((*context.ApplicationListener)(nil)))
	if err != nil {
		return err
	}
	listenerInstances := make([]context.ApplicationListener, len(instances))
	for index, instance := range instances {
		listenerInstances[index] = instance.(context.ApplicationListener)
	}
	application.listeners = listenerInstances
	return nil
}

func (application *baseApplication) initApplicationContextInitializers() error {
	instances, err := getInstances(goo.GetType((*context.ApplicationContextInitializer)(nil)))
	if err != nil {
		return err
	}
	initializerInstances := make([]context.ApplicationContextInitializer, len(instances))
	for index, instance := range instances {
		initializerInstances[index] = instance.(context.ApplicationContextInitializer)
	}
	application.contextInitializers = initializerInstances
	return nil
}

func (application *baseApplication) invokeApplicationRunners(ctx context.ApplicationContext, arguments ApplicationArguments) {
	applicationRunners := ctx.GetSharedPeasByType(goo.GetType((*ApplicationRunner)(nil)))
	for _, applicationRunner := range applicationRunners {
		applicationRunner.(ApplicationRunner).OnApplicationRun(ctx, arguments)
	}
}

func (application *baseApplication) logStarting() {
	application.logger.Info(application.contextId, "Starting...")
	application.logger.Infof(application.contextId, "Application Id : %s", application.applicationId)
	application.logger.Infof(application.contextId, "Application Context Id : %s", application.contextId)
	application.logger.Info(application.contextId, "Running with Procyon, Procyon "+Version)
}

func (application *baseApplication) logStarted() {
	lastTime := float32(application.taskWatch.GetTotalTime()) / 1e9
	formattedText := fmt.Sprintf("Started in %.2f second(s)", lastTime)
	application.logger.Info(application.contextId, formattedText)
}

func (application *baseApplication) finish() {
	exitSignalChannel := make(chan os.Signal, 1)
	signal.Notify(exitSignalChannel, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignalChannel
}

func (application *baseApplication) getCustomLogger() {
	customLoggers, err := getInstances(goo.GetType((*context.Logger)(nil)))
	if err != nil {
		panic(err)
	}

	if customLoggers != nil {
		if len(customLoggers) != 1 {
			panic("Custom logger cannot be distinguished because there are more than one")
		}

		if len(customLoggers) != 0 {
			application.customLogger = customLoggers[0].(context.Logger)
		}
	}
}

func (application *baseApplication) configureLogger(logger context.Logger, loggingProperties *context.LoggingProperties) {
	if logger == nil {
		return
	}

	if configurableLogger, ok := logger.(context.LoggingConfiguration); ok {
		configurableLogger.ApplyLoggingProperties(*loggingProperties)
	}
}

func (application *baseApplication) getLoggingProperties(arguments ApplicationArguments) *context.LoggingProperties {
	if arguments == nil {
		return nil
	}

	properties := &context.LoggingProperties{}
	loggingLevel := arguments.GetOptionValues("logging.level")
	if len(loggingLevel) != 0 {
		properties.Level = loggingLevel[0]
	} else {
		properties.Level = "TRACE"
	}

	loggingFile := arguments.GetOptionValues("logging.file.name")
	if len(loggingFile) != 0 {
		properties.FileName = loggingFile[0]
	}

	loggingPath := arguments.GetOptionValues("logging.file.path")
	if len(loggingFile) != 0 {
		properties.FilePath = loggingPath[0]
	}

	return properties
}

type defaultEnvironmentProvider struct {
}

func newDefaultEnvironmentProvider() defaultEnvironmentProvider {
	return defaultEnvironmentProvider{}
}

func (provider defaultEnvironmentProvider) getNewEnvironment() core.ConfigurableEnvironment {
	return web.NewStandardWebEnvironment()
}

type defaultContextProvider struct {
}

func newDefaultContextProvider() defaultContextProvider {
	return defaultContextProvider{}
}

func (provider defaultContextProvider) getNewContext(applicationId context.ApplicationId, contextId context.ContextId) context.ConfigurableApplicationContext {
	return web.NewProcyonServerApplicationContext(applicationId, contextId)
}
