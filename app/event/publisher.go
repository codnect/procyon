package event

import "context"

type Publisher interface {
	PublishEvent(ctx context.Context, event Event)
}
