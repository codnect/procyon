package procyon

import (
	"github.com/procyon-projects/procyon/app"
	"github.com/procyon-projects/procyon/container"
	"github.com/procyon-projects/procyon/web"
)

type AppType int

const (
	None AppType = -1
	Web  AppType = 0
)

type AppBuilder interface {
	Type(appType AppType) AppBuilder
	Run(args ...string)
}

type appBuilder struct {
	appType AppType
}

func newAppBuilder() *appBuilder {
	return &appBuilder{
		appType: Web,
	}
}

func (b *appBuilder) Type(appType AppType) AppBuilder {
	b.appType = appType
	return b
}

func (b *appBuilder) Run(args ...string) {
	if b.appType == Web {
		container.Register(web.NewContextCustomizer)
	}

	app.New().Run(args...)
}

func New() AppBuilder {
	return newAppBuilder()
}
