package app

import (
	"github.com/procyon-projects/procyon/health"
	"time"
)

type AvailabilityState interface {
	Status() health.Status
}

type LivenessState int

const (
	Correct LivenessState = iota + 1
	Broken
)

func (s LivenessState) Status() health.Status {
	switch s {
	case Correct:
		return health.StatusUp
	case Broken:
		return health.StatusDown
	}

	return health.StatusUnknown
}

type ReadinessState int

const (
	AcceptingTraffic ReadinessState = iota + 1
	RefusingTraffic
)

func (s ReadinessState) Status() health.Status {
	switch s {
	case AcceptingTraffic:
		return health.StatusUp
	case RefusingTraffic:
		return health.StatusDown
	}

	return health.StatusUnknown
}

type Availability interface {
	LivenessState() LivenessState
	ReadinessState() ReadinessState
}

type AvailabilityChangeEvent struct {
	ctx   Context
	state AvailabilityState
	time  time.Time
}

func (e *AvailabilityChangeEvent) EventSource() any {
	return e.ctx
}

func (e *AvailabilityChangeEvent) Time() time.Time {
	return e.time
}

func (e *AvailabilityChangeEvent) State() AvailabilityState {
	return e.state
}

type AvailabilityStateChecker interface {
	health.Checker

	State() AvailabilityState
}

type LivenessStateChecker struct {
}

func (c *LivenessStateChecker) DoHealthCheck() health.Health {
	return nil
}

func (c *LivenessStateChecker) State() AvailabilityState {
	return nil
}

type ReadinessStateChecker struct {
}

func (c *ReadinessStateChecker) DoHealthCheck() health.Health {
	return nil
}

func (c *ReadinessStateChecker) State() AvailabilityState {
	return nil
}
