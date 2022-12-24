package event

import (
	"context"
	"fmt"
	"github.com/procyon-projects/reflector"
	"reflect"
)

type ListenerFunc[E Event] interface {
	~func(ctx context.Context, event E)
}

type Listener struct {
	identifier string
	eventType  reflect.Type
	wrapper    func(ctx context.Context, event Event)
}

func (l *Listener) Identifier() string {
	return l.identifier
}

func (l *Listener) EventType() reflect.Type {
	return l.eventType
}

func (l *Listener) OnEvent(ctx context.Context, event Event) {
	l.wrapper(ctx, event)
}

func Listen[E Event, F ListenerFunc[E]](handler F) *Listener {
	reflFunc := reflector.TypeOfAny(handler)
	eventType := reflector.TypeOf[E]().ReflectType()

	return &Listener{
		identifier: fmt.Sprintf("%s.%s", reflFunc.PackagePath(), reflFunc.Name()),
		eventType:  eventType,
		wrapper: func(ctx context.Context, event Event) {
			if event == nil && eventType.Kind() == reflect.Pointer {
				handler(ctx, reflect.Zero(eventType).Interface().(E))
			} else {
				handler(ctx, event.(E))
			}
		},
	}
}

type ListenerRegistry interface {
	RegisterListener(listener *Listener)
	Listeners() []*Listener
}
