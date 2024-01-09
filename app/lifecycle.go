package app

import (
	"codnect.io/procyon/core/env/property"
	"context"
	"time"
)

type Lifecycle interface {
	Start() error
	Stop(ctx context.Context) error
	IsRunning() bool
}

type LifecycleProcessor interface {
	Start() error
	Stop() error
	IsRunning() bool
}

type LifecycleProperties struct {
	property.Properties `prefix:"procyon.lifecycle"`

	ShutdownTimeout time.Duration `prop:"shutdown.timeout" default:"30000"`
}

type DefaultLifecycleProcessor struct {
	shutdownTimeout time.Duration
	running         bool
}

func NewDefaultLifecycleProcessor() *DefaultLifecycleProcessor {
	return &DefaultLifecycleProcessor{}
}

func (p *DefaultLifecycleProcessor) Start() error {
	p.running = true
	return nil
}

func (p *DefaultLifecycleProcessor) Stop() error {
	p.running = false

	lifecycleInstances := make([]Lifecycle, 0)

	for _, instance := range lifecycleInstances {
		instance.Stop(nil)
	}

	return nil
}

func (p *DefaultLifecycleProcessor) IsRunning() bool {
	return p.running
}
