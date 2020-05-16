package app

import (
	"procyon/env"
	"procyon/server"
)

type ApplicationContext interface {
	GetApplicationName() string
	GetStartupTimeStamp() int64
}

type WebApplicationContext interface {
	ApplicationContext
}

type ConfigurableApplicationContext interface {
	SetEnvironment(environment env.ConfigurableEnvironment)
	GetEnvironment() env.ConfigurableEnvironment
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
