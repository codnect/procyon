package gdbc

import "context"

type Conn interface {
	BeginTx(ctx context.Context, opts *TxOptions) (*Tx, error)

	PingContext(ctx context.Context) error
	ExecContext(ctx context.Context, query string, args ...any) (Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *Row
	PrepareContext(ctx context.Context, query string) (*Stmt, error)

	Close() error
}
