package availability

import "github.com/procyon-projects/procyon/health"

type State interface {
	Status() health.Status
}

type StateChecker interface {
	health.Checker

	State() State
}
