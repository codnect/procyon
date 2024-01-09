package gdbc

import (
	"codnect.io/procyon/data/sql"
	"context"
)

type Operations interface {
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
	Prepare(ctx context.Context, query string) (sql.Statement, error)
	Query(ctx context.Context, query string, args ...any) (sql.RowSet, error)
	QueryRow(ctx context.Context, query string, args ...any) sql.Row
}
