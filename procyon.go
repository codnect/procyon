package procyon

import (
	"codnect.io/procyon/app"
	"codnect.io/procyon/component"
	"codnect.io/procyon/web"
)

type ApplicationType int

const (
	None ApplicationType = -1
	Web  ApplicationType = 0
)

type Application interface {
	ApplicationType(typ ApplicationType) Application
	Run(args ...string)
}

type application struct {
	typ ApplicationType
}

func newApplication() *application {
	return &application{
		typ: Web,
	}
}

func (a *application) ApplicationType(typ ApplicationType) Application {
	a.typ = typ
	return a
}

func (a *application) Run(args ...string) {
	if a.typ == Web {
		component.Register(web.NewContextCustomizer)
	}

	app.New().Run(args...)
}

func New() Application {
	return newApplication()
}
