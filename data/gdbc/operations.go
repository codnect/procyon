package gdbc

import (
	sql2 "codnect.io/procyon/data/sql"
	"context"
)

type Operations interface {
	Exec(ctx context.Context, query string, args ...any) (sql2.Result, error)
	Prepare(ctx context.Context, query string) (sql2.Statement, error)
	Query(ctx context.Context, query string, args ...any) (sql2.RowSet, error)
	QueryRow(ctx context.Context, query string, args ...any) sql2.Row
}
