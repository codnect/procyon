package procyon

import (
	"log"
)

type Application struct {
	appRunListeners ApplicationRunListeners
	startupLogger   AppStartupLogger
}

func NewProcyonApplication() *Application {
	return &Application{
		startupLogger: NewAppStartupLogger(),
	}
}

func (app *Application) SetApplicationRunListeners(listeners ...ApplicationRunListener) {
	app.appRunListeners = newApplicationRunListeners(listeners)
}

func (app *Application) Run() {
	taskWatch := NewTaskWatch()
	_ = taskWatch.Start()
	app.appRunListeners.starting()
	// prepare environment
	appArguments := NewDefaultApplicationArguments()
	environment := app.prepareEnvironment(appArguments)
	// print banner
	ProcyonBanner{}.printBanner()
	context := app.createApplicationContext()
	app.prepareContext(context, environment, appArguments)
	app.appRunListeners.started(context)
	app.appRunListeners.running(context)
	_ = taskWatch.Stop()
	app.startupLogger.logStarted(taskWatch)
}

func (app *Application) prepareEnvironment(arguments ApplicationArguments) Environment {
	environment := app.createEnvironment()
	app.appRunListeners.environmentPrepared(environment)
	return environment
}

func (app *Application) createEnvironment() ConfigurableEnvironment {
	return nil
}

func (app *Application) createApplicationContext() ConfigurableApplicationContext {
	return nil
}

func (app *Application) prepareContext(context ConfigurableApplicationContext,
	environment ConfigurableEnvironment,
	arguments ApplicationArguments) {
	app.startupLogger.logStarting()
	app.appRunListeners.contextPrepared(context)
	app.appRunListeners.contextLoaded(context)
}

type AppStartupLogger struct {
}

func NewAppStartupLogger() AppStartupLogger {
	return AppStartupLogger{}
}

func (logger AppStartupLogger) logStarting() {
	log.Println("Starting...")
	log.Println("Running with Procyon, Procyon " + Version)
}

func (logger AppStartupLogger) logStarted(watch *TaskWatch) {
	lastTime := float32(watch.totalTimeNanoSeconds) / 1e9
	log.Printf("Started in %.2f second(s)\n", lastTime)
}
