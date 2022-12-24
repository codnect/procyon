package availability

import (
	"github.com/procyon-projects/procyon/app/health"
)

type State interface {
	Status() health.Status
}
