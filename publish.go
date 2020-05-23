package procyon

import (
	context "github.com/Rollcomp/procyon-context"
	core "github.com/Rollcomp/procyon-core"
	"log"
)

type EventPublishRunListener struct {
}

func NewEventPublishRunListener() EventPublishRunListener {
	return EventPublishRunListener{}
}

func (publishRunListener EventPublishRunListener) Starting() {
	log.Print("Starting")
}

func (publishRunListener EventPublishRunListener) EnvironmentPrepared(environment core.ConfigurableEnvironment) {
	log.Print("EnvironmentPrepared")
}

func (publishRunListener EventPublishRunListener) ContextPrepared(context context.ConfigurableApplicationContext) {
	log.Print("ContextPrepared")
}

func (publishRunListener EventPublishRunListener) ContextLoaded(context context.ConfigurableApplicationContext) {
	log.Print("ContextLoaded")
}

func (publishRunListener EventPublishRunListener) Started(context context.ConfigurableApplicationContext) {
	log.Print("Started")
}

func (publishRunListener EventPublishRunListener) Running(context context.ConfigurableApplicationContext) {
	log.Print("Running")
}

func (publishRunListener EventPublishRunListener) Failed(context context.ConfigurableApplicationContext, err error) {
	log.Print("Failed")
}
