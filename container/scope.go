package container

import "context"

const (
	SharedScope    = "shared"
	PrototypeScope = "prototype"
)

type Scope interface {
	Get(ctx context.Context, name string, supplier func(ctx context.Context) (any, error)) (any, error)
	Remove(ctx context.Context, name string) (any, error)
}
