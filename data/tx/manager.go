package tx

import "context"

type Manager interface {
	ContextWithTx(ctx context.Context, attributes Attributes) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
