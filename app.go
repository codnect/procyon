package procyon

type Application struct {
}

func New() *Application {
	return &Application{}
}

func (app *Application) Run(args ...string) error {
	return nil
}
