package gdbc

import "context"

type Operations interface {
	ExecContext(ctx context.Context, query string, args ...any) (Result, error)
	PrepareContext(ctx context.Context, query string) (*Stmt, error)
	QueryContext(ctx context.Context, query string, args ...any) (*Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *Row
}
