package availability

import (
	"github.com/procyon-projects/procyon/health"
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

type ReadinessStateChecker struct {
}

func (c *ReadinessStateChecker) DoHealthCheck() health.Health {
	return nil
}

func (c *ReadinessStateChecker) State() State {
	return nil
}
