package availability

import (
	"codnect.io/procyon/app/health"
)

type State interface {
	Status() health.Status
}
