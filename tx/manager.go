package tx

import (
	"context"
	"database/sql"
	"reflect"
)

type ctxTransactionManager struct{}

var (
	ctxTransactionManagerKey = &ctxTransactionManager{}
	reflTransactionManager   = reflect.TypeOf((*TransactionManager)(nil)).Elem()
)

type TransactionManager interface {
	CreateContext() TransactionContext
	Connection() *sql.DB
	GetOrCreateTransaction(ctx context.Context, options ...Option) (Transaction, error)
	Commit(ctx context.Context, tx Transaction) error
	Rollback(ctx context.Context, tx Transaction) error
}

type transactionManager struct {
	db      *sql.DB
	options *Options
}

func NewTransactionManager(db *sql.DB, options ...Option) TransactionManager {
	return &transactionManager{
		db:      db,
		options: NewOptions(options...),
	}
}

func (m *transactionManager) CreateContext() TransactionContext {
	return nil
}

func (m *transactionManager) Connection() *sql.DB {
	return m.db
}

func (m *transactionManager) GetOrCreateTransaction(ctx context.Context, options ...Option) (Transaction, error) {
	return nil, nil
}

func (m *transactionManager) Commit(ctx context.Context, tx Transaction) error {
	return nil
}

func (m *transactionManager) Rollback(ctx context.Context, tx Transaction) error {
	return nil
}
