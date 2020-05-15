package procyon

type Application struct {
	appRunListeners ApplicationRunListeners
}

func NewProcyonApplication() *Application {
	return &Application{}
}

func (app *Application) SetApplicationRunListeners(listeners ...ApplicationRunListener) {
	app.appRunListeners = newApplicationRunListeners(listeners)
}

func (app *Application) Run() {
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

	app.appRunListeners.contextPrepared(context)
	app.appRunListeners.contextLoaded(context)
}
