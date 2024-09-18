package container

import "context"

type Initialization interface {
	DoInit(ctx context.Context) error
}
