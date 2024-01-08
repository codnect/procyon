package sql

import (
	"context"
	"database/sql/driver"
	"time"
)

type Connection interface {
	ExecContext(ctx context.Context, query string, args ...any) (Result, error)
	PrepareContext(ctx context.Context, query string, args ...any) (Statement, error)
	QueryContext(ctx context.Context, query string, args ...any) (RowSet, error)
	QueryRowContext(ctx context.Context, query string, args ...any) (Row, error)
	PingContext(ctx context.Context) error

	SetConnMaxLifetime(d time.Duration) error
	SetConnMaxIdleTime(d time.Duration) error
	SetMaxIdleConns(n int) error
	SetMaxOpenConns(n int) error

	BeginTransaction(ctx context.Context) (TransactionStatus, error)
	Driver() driver.Driver

	Close() error
	IsClosed() bool
}
