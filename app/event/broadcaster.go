package event

import (
	"context"
	"reflect"
	"sync"
)

type Broadcaster interface {
	RegisterListener(listener *Listener)
	RemoveListener(listener *Listener)
	RemoveAllListeners()
	BroadcastEvent(ctx context.Context, event Event)
}

type broadcaster struct {
	mu          sync.RWMutex
	listenerMap map[reflect.Type]map[string]*Listener
}

func NewBroadcaster() Broadcaster {
	return &broadcaster{
		mu:          sync.RWMutex{},
		listenerMap: map[reflect.Type]map[string]*Listener{},
	}
}

func (b *broadcaster) RegisterListener(listener *Listener) {
	defer b.mu.Unlock()
	b.mu.Lock()

	eventType := listener.EventType()

	if _, ok := b.listenerMap[eventType]; !ok {
		b.listenerMap[eventType] = make(map[string]*Listener)
	}

	if _, exists := b.listenerMap[eventType][listener.Identifier()]; exists {
		return
	}

	b.listenerMap[eventType][listener.Identifier()] = listener
}

func (b *broadcaster) RemoveListener(listener *Listener) {
	defer b.mu.Unlock()
	b.mu.Lock()

	eventType := listener.EventType()

	if _, ok := b.listenerMap[eventType]; !ok {
		return
	}

	if _, exists := b.listenerMap[eventType][listener.Identifier()]; !exists {
		return
	}

	delete(b.listenerMap[eventType], listener.Identifier())
}

func (b *broadcaster) RemoveAllListeners() {
	defer b.mu.Unlock()
	b.mu.Lock()

	for key, _ := range b.listenerMap {
		delete(b.listenerMap, key)
	}
}

func (b *broadcaster) BroadcastEvent(ctx context.Context, event Event) {
	if event == nil {
		panic("event: event cannot be nil")
	}

	defer b.mu.Unlock()
	b.mu.Lock()

	currentEvenType := reflect.TypeOf(event)

	for eventType, listenerList := range b.listenerMap {
		if currentEvenType.ConvertibleTo(eventType) {
			for _, l := range listenerList {
				l.OnEvent(ctx, event)
			}
		}
	}
}
