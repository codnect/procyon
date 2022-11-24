package event

type Broadcaster interface {
	AddListener(listener Listener)
	RemoveListener(listener Listener)
	RemoveAllListeners()
	BroadcastEvent(event Event)
}
