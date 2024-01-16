package availability

import (
	"codnect.io/procyon/health"
)

type ReadinessState int

const (
	StateAcceptingTraffic ReadinessState = iota + 1
	StateRefusingTraffic
)

func (s ReadinessState) Status() health.Status {
	switch s {
	case StateAcceptingTraffic:
		return health.StatusUp
	case StateRefusingTraffic:
		return health.StatusDown
	}

	return health.StatusUnknown
}

type ReadinessStateHealthChecker struct {
	availability Availability
}

func NewReadinessStateHealthChecker(availability Availability) *ReadinessStateHealthChecker {
	return &ReadinessStateHealthChecker{
		availability: availability,
	}
}

func (c *ReadinessStateHealthChecker) DoHealthCheck() (health.Health, error) {
	state := c.availability.ReadinessState()
	return health.Of(state.Status()), nil
}
