package availability

import (
	"codnect.io/procyon/app/health"
)

type LivenessState int

const (
	StateCorrect LivenessState = iota + 1
	StateBroken
)

func (s LivenessState) Status() health.Status {
	switch s {
	case StateCorrect:
		return health.StatusUp
	case StateBroken:
		return health.StatusDown
	}

	return health.StatusUnknown
}

type LivenessStateHealthChecker struct {
	availability Availability
}

func NewLivenessStateHealthChecker(availability Availability) *LivenessStateHealthChecker {
	return &LivenessStateHealthChecker{
		availability: availability,
	}
}

func (c *LivenessStateHealthChecker) DoHealthCheck() (health.Health, error) {
	state := c.availability.LivenessState()
	return health.Of(state.Status()), nil
}
