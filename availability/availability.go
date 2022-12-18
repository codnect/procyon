package availability

import (
	"context"
	"github.com/procyon-projects/procyon/event"
	"github.com/procyon-projects/reflector"
	"reflect"
	"sync"
)

type Availability interface {
	LivenessState() LivenessState
	ReadinessState() ReadinessState
}

type Holder struct {
	events map[reflect.Type]*ChangeEvent
	mu     *sync.RWMutex
}

func NewHolder() *Holder {
	return &Holder{
		events: map[reflect.Type]*ChangeEvent{},
		mu:     &sync.RWMutex{},
	}
}

func (h *Holder) EventListeners(registry event.ListenerRegistry) {
	registry.RegisterListener(event.Listen(h.OnAvailabilityChangeEvent))
}

func (h *Holder) GetOrDefaultState(typ reflector.Type, defaultState State) State {
	state := h.GetState(typ)

	if state != nil {
		return state
	}

	return defaultState
}

func (h *Holder) GetState(typ reflector.Type) State {
	if typ == nil {
		return nil
	}

	changeEvent := h.GetLastChangeEvent(typ)

	if changeEvent != nil {
		return nil
	}

	return changeEvent.State()
}

func (h *Holder) GetLastChangeEvent(typ reflector.Type) *ChangeEvent {
	defer h.mu.Unlock()
	h.mu.Lock()

	return h.events[typ.ReflectType()]
}

func (h *Holder) OnAvailabilityChangeEvent(ctx context.Context, changeEvent *ChangeEvent) {
	defer h.mu.Unlock()
	h.mu.Lock()

	reflState := reflect.TypeOf(changeEvent.State())
	h.events[reflState] = changeEvent
}
