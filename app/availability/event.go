package availability

import (
	"time"
)

type ChangeEvent struct {
	ctx   any
	state State
	time  time.Time
}

func NewChangeEvent(source any, state State) *ChangeEvent {
	return &ChangeEvent{
		ctx:   source,
		state: state,
		time:  time.Now(),
	}
}

func (e *ChangeEvent) EventSource() any {
	return e.ctx
}

func (e *ChangeEvent) Time() time.Time {
	return e.time
}

func (e *ChangeEvent) State() State {
	return e.state
}
