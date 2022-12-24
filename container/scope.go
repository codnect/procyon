package container

import "context"

const (
	SharedScope    = "shared"
	PrototypeScope = "prototype"
)

type Scope interface {
	GetObject(ctx context.Context, name string, supplier func(ctx context.Context) (any, error)) (any, error)
	RemoveObject(ctx context.Context, name string) (any, error)
}
