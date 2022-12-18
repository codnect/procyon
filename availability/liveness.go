package availability

import (
	"github.com/procyon-projects/procyon/health"
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

type LivenessStateChecker struct {
}

func (c *LivenessStateChecker) DoHealthCheck() health.Health {
	return nil
}

func (c *LivenessStateChecker) State() State {
	return nil
}
