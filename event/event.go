package event

import "time"

type Event interface {
	EventSource() any

	Time() time.Time
}
