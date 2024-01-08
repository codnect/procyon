package gdbc

import (
	"codnect.io/procyon/data/sql"
	"context"
	dsql "database/sql"
	"database/sql/driver"
	"errors"
	"time"
)

var (
	ErrConnectionDoesNotExist = errors.New("this connection has been closed")
)

type Connection struct {
	db     *dsql.DB
	closed bool
}

func (c *Connection) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if c.closed {
		return nil, ErrConnectionDoesNotExist
	}

	return c.db.ExecContext(ctx, query, args...)
}

func (c *Connection) PrepareContext(ctx context.Context, query string, args ...any) (sql.Statement, error) {
	if c.closed {
		return nil, ErrConnectionDoesNotExist
	}

	statement, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	preparedStatement := sql.NewPreparedStatement(statement)
	return preparedStatement, nil
}

func (c *Connection) QueryContext(ctx context.Context, query string, args ...any) (sql.RowSet, error) {
	if c.closed {
		return nil, ErrConnectionDoesNotExist
	}

	return c.db.QueryContext(ctx, query, args...)
}

func (c *Connection) QueryRowContext(ctx context.Context, query string, args ...any) (sql.Row, error) {
	if c.closed {
		return nil, ErrConnectionDoesNotExist
	}

	return c.db.QueryRowContext(ctx, query, args...), nil
}

func (c *Connection) PingContext(ctx context.Context) error {
	if c.closed {
		return ErrConnectionDoesNotExist
	}

	return c.db.PingContext(ctx)
}

func (c *Connection) SetConnMaxLifetime(d time.Duration) error {
	if c.closed {
		return ErrConnectionDoesNotExist
	}

	c.db.SetConnMaxLifetime(d)
	return nil
}

func (c *Connection) SetConnMaxIdleTime(d time.Duration) error {
	if c.closed {
		return ErrConnectionDoesNotExist
	}

	c.db.SetConnMaxIdleTime(d)
	return nil
}

func (c *Connection) SetMaxIdleConns(n int) error {
	if c.closed {
		return ErrConnectionDoesNotExist
	}

	c.db.SetMaxIdleConns(n)
	return nil
}

func (c *Connection) SetMaxOpenConns(n int) error {
	if c.closed {
		return ErrConnectionDoesNotExist
	}

	c.db.SetMaxOpenConns(n)
	return nil
}

func (c *Connection) BeginTransaction(ctx context.Context) (sql.TransactionStatus, error) {
	if c.closed {
		return nil, ErrConnectionDoesNotExist
	}

	return nil, nil
}

func (c *Connection) Driver() driver.Driver {
	return c.db.Driver()
}

func (c *Connection) Close() error {
	if c.closed {
		return ErrConnectionDoesNotExist
	}

	err := c.db.Close()
	if err != nil {
		return err
	}

	c.closed = true
	return nil
}

func (c *Connection) IsClosed() bool {
	return c.closed
}
