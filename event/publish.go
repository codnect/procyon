package event

import (
	context "github.com/Rollcomp/procyon-context"
	core "github.com/Rollcomp/procyon-core"
	"log"
)

type PublishRunListener struct {
}

func NewPublishRunListener() PublishRunListener {
	return PublishRunListener{}
}

func (publishRunListener PublishRunListener) Starting() {
	log.Print("Starting")
}

func (publishRunListener PublishRunListener) EnvironmentPrepared(environment core.ConfigurableEnvironment) {
	log.Print("EnvironmentPrepared")
}

func (publishRunListener PublishRunListener) ContextPrepared(context context.ConfigurableApplicationContext) {
	log.Print("ContextPrepared")
}

func (publishRunListener PublishRunListener) ContextLoaded(context context.ConfigurableApplicationContext) {
	log.Print("ContextLoaded")
}

func (publishRunListener PublishRunListener) Started(context context.ConfigurableApplicationContext) {
	log.Print("Started")
}

func (publishRunListener PublishRunListener) Running(context context.ConfigurableApplicationContext) {
	log.Print("Running")
}

func (publishRunListener PublishRunListener) Failed(context context.ConfigurableApplicationContext, err error) {
	log.Print("Failed")
}
