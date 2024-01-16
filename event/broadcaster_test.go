package event

import (
	"codnect.io/reflector"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testEvent struct {
}

func (c *testEvent) EventSource() any {
	return nil
}

func (c *testEvent) Time() time.Time {
	return time.Time{}
}

func TestBroadcaster_RegisterListener(t *testing.T) {
	broadcaster := NewBroadcaster().(*broadcaster)

	listener := Listen(func(ctx context.Context, event *testEvent) {})
	broadcaster.RegisterListener(listener)
	anotherListener := Listen(func(ctx context.Context, event *testEvent) {})
	broadcaster.RegisterListener(anotherListener)

	assert.Len(t, broadcaster.listenerMap, 1)
	assert.Len(t, broadcaster.listenerMap[reflector.TypeOf[*testEvent]().ReflectType()], 2)
	assert.Equal(t, map[string]*Listener{
		listener.Identifier():        listener,
		anotherListener.Identifier(): anotherListener,
	}, broadcaster.listenerMap[reflector.TypeOf[*testEvent]().ReflectType()])
}

func TestBroadcaster_RemoveListener(t *testing.T) {
	broadcaster := NewBroadcaster().(*broadcaster)

	listener := Listen(func(ctx context.Context, event *testEvent) {})
	broadcaster.RegisterListener(listener)
	anotherListener := Listen(func(ctx context.Context, event *testEvent) {})
	broadcaster.RegisterListener(anotherListener)

	broadcaster.RemoveListener(listener)

	assert.Len(t, broadcaster.listenerMap, 1)
	assert.Len(t, broadcaster.listenerMap[reflector.TypeOf[*testEvent]().ReflectType()], 1)
	assert.Equal(t, map[string]*Listener{
		anotherListener.Identifier(): anotherListener,
	}, broadcaster.listenerMap[reflector.TypeOf[*testEvent]().ReflectType()])
}

func TestBroadcaster_BroadcastEvent(t *testing.T) {
	broadcaster := NewBroadcaster()
	var (
		actualEvent   = &testEvent{}
		numberOfCalls = 0
	)
	broadcaster.RegisterListener(Listen(func(ctx context.Context, event *testEvent) {
		assert.Equal(t, actualEvent, event)
		numberOfCalls++
	}))
	broadcaster.RegisterListener(Listen(func(ctx context.Context, event *testEvent) {
		assert.Equal(t, actualEvent, event)
		numberOfCalls++
	}))

	broadcaster.BroadcastEvent(context.Background(), actualEvent)
	assert.Equal(t, 2, numberOfCalls)
}

func TestBroadcaster_RemoveAllListeners(t *testing.T) {
	broadcaster := NewBroadcaster().(*broadcaster)

	listener := Listen(func(ctx context.Context, event *testEvent) {})
	broadcaster.RegisterListener(listener)
	anotherListener := Listen(func(ctx context.Context, event *testEvent) {})
	broadcaster.RegisterListener(anotherListener)

	broadcaster.RemoveAllListeners()

	assert.Len(t, broadcaster.listenerMap, 0)
}
