package sql

import (
	"context"
	"database/sql"
)

type Statement interface {
	ExecContext(ctx context.Context, args ...any) (Result, error)
	QueryContext(ctx context.Context, args ...any) (RowSet, error)
	QueryRowContext(ctx context.Context, args ...any) (Row, error)
	Close() error
}

type PreparedStatement struct {
	stmt *sql.Stmt
}

func NewPreparedStatement(stmt *sql.Stmt) *PreparedStatement {
	return &PreparedStatement{
		stmt: stmt,
	}
}

func (p *PreparedStatement) ExecContext(ctx context.Context, args ...any) (Result, error) {
	return p.stmt.ExecContext(ctx, args...)
}

func (p *PreparedStatement) QueryContext(ctx context.Context, args ...any) (RowSet, error) {
	return p.stmt.QueryContext(ctx, args...)
}

func (p *PreparedStatement) QueryRowContext(ctx context.Context, args ...any) (Row, error) {
	return p.stmt.QueryRowContext(ctx, args...), nil
}

func (p *PreparedStatement) Close() error {
	return p.stmt.Close()
}
