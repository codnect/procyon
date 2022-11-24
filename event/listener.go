package event

import (
	"context"
	"github.com/procyon-projects/reflector"
)

type ListenerFunc[E Event] interface {
	~func(ctx context.Context, event E)
}

type Listener interface {
	EventType() reflector.Type
}

func Listen[E Event, F ListenerFunc[E]](f F) Listener {
	return nil
}

type ListenerRegistry interface {
	Register(l Listener)
	Listeners() []Listener
}

type RegisterListeners interface {
	EventListeners(r ListenerRegistry)
}
