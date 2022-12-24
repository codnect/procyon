package availability

import (
	"context"
	"github.com/procyon-projects/procyon/app/event"
	"github.com/procyon-projects/reflector"
	"reflect"
	"sync"
)

var (
	reflLivenessState  = reflector.TypeOf[LivenessState]()
	reflReadinessState = reflector.TypeOf[ReadinessState]()
)

type Availability interface {
	LivenessState() LivenessState
	ReadinessState() ReadinessState
	GetOrDefaultState(typ reflector.Type, defaultState State) State
	GetState(typ reflector.Type) State
	GetLastChangeEvent(typ reflector.Type) *ChangeEvent
}

type StateHolder struct {
	registry event.ListenerRegistry
	events   map[reflect.Type]*ChangeEvent
	mu       *sync.RWMutex
}

func NewStateHolder(registry event.ListenerRegistry) *StateHolder {
	return &StateHolder{
		registry: registry,
		events:   map[reflect.Type]*ChangeEvent{},
		mu:       &sync.RWMutex{},
	}
}

func (h *StateHolder) PostConstruct() error {
	changeEventListener := event.Listen(h.OnAvailabilityChangeEvent)
	h.registry.RegisterListener(changeEventListener)
	return nil
}

func (h *StateHolder) LivenessState() LivenessState {
	return h.GetOrDefaultState(reflLivenessState, StateBroken).(LivenessState)
}

func (h *StateHolder) ReadinessState() ReadinessState {
	return h.GetOrDefaultState(reflLivenessState, StateRefusingTraffic).(ReadinessState)
}

func (h *StateHolder) GetOrDefaultState(typ reflector.Type, defaultState State) State {
	state := h.GetState(typ)

	if state != nil {
		return state
	}

	return defaultState
}

func (h *StateHolder) GetState(typ reflector.Type) State {
	if typ == nil {
		return nil
	}

	changeEvent := h.GetLastChangeEvent(typ)

	if changeEvent != nil {
		return nil
	}

	return changeEvent.State()
}

func (h *StateHolder) GetLastChangeEvent(typ reflector.Type) *ChangeEvent {
	defer h.mu.Unlock()
	h.mu.Lock()

	return h.events[typ.ReflectType()]
}

func (h *StateHolder) OnAvailabilityChangeEvent(ctx context.Context, changeEvent *ChangeEvent) {
	defer h.mu.Unlock()
	h.mu.Lock()

	reflState := reflect.TypeOf(changeEvent.State())
	h.events[reflState] = changeEvent
}
