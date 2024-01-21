package procyon

import (
	"codnect.io/procyon-core/container"
	"codnect.io/procyon-core/env/property"
	"codnect.io/reflector"
	"context"
	"time"
)

type Lifecycle interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	IsRunning() bool
}

type LifecycleProperties struct {
	property.Properties `prefix:"procyon.lifecycle"`

	ShutdownTimeout time.Duration `prop:"shutdown.timeout" default:"30000"`
}

type lifecycleProcessor struct {
	properties LifecycleProperties
	container  container.Container
}

func defaultLifecycleProcessor(properties LifecycleProperties, container container.Container) *lifecycleProcessor {
	return &lifecycleProcessor{
		properties: properties,
		container:  container,
	}
}

func (p *lifecycleProcessor) start(ctx context.Context) error {
	err := p.startLifecycleComponents(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (p *lifecycleProcessor) stop(ctx context.Context) error {
	err := p.stopLifecycleComponents(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (p *lifecycleProcessor) startLifecycleComponents(ctx context.Context) (err error) {
	sharedInstances := p.container.SharedInstances()
	lifecycleInstances := sharedInstances.FindAllByType(reflector.TypeOf[Lifecycle]())

	for _, instance := range lifecycleInstances {
		lifecycle := instance.(Lifecycle)

		err = lifecycle.Start(ctx)

		if err != nil {
			return
		}
	}

	return
}

func (p *lifecycleProcessor) stopLifecycleComponents(ctx context.Context) (err error) {
	sharedInstances := p.container.SharedInstances()
	lifecycleInstances := sharedInstances.FindAllByType(reflector.TypeOf[Lifecycle]())

	for _, instance := range lifecycleInstances {
		lifecycle := instance.(Lifecycle)

		err = lifecycle.Stop(ctx)

		if err != nil {
			return
		}
	}

	return
}
