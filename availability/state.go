package availability

import (
	"codnect.io/procyon/health"
)

type State interface {
	Status() health.Status
}
