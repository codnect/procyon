package context

import (
	"procyon/core"
	"procyon/web/server"
)

type ApplicationContext interface {
	GetApplicationName() string
	GetStartupTimeStamp() int64
}

type WebApplicationContext interface {
	ApplicationContext
}

type ConfigurableApplicationContext interface {
	SetEnvironment(environment core.ConfigurableEnvironment)
	GetEnvironment() core.ConfigurableEnvironment
}

type WebServerApplicationContext struct {
	webServer *server.WebServer
}

func NewWebServerApplicationContext() *WebServerApplicationContext {
	return &WebServerApplicationContext{}
}

func (ctx *WebServerApplicationContext) createWebServer() {
	if ctx.webServer == nil {
		factory := server.NewDefaultWebServerFactory()
		ctx.webServer = factory.GetWebServer()
	}
}

type DefaultWebServerApplicationContext struct {
}
