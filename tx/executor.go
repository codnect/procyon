package tx

import (
	"context"
	"reflect"
)

type ctxTransactionExecutor struct{}

var (
	ctxTransactionExecutorKey = &ctxTransactionExecutor{}
	reflTransactionExecutor   = reflect.TypeOf((*TransactionContext)(nil)).Elem()
)

type Func func(ctx context.Context, template Template) (any, error)

type TransactionExecutor interface {
	Execute(ctx context.Context, fn Func) (any, error)
	ExecutorInTransaction(ctx context.Context, fn Func, options ...Option) (any, error)
	TransactionManager() TransactionManager
}

type transactionExecutor struct {
	manager TransactionManager
	options *Options
}

func NewTransactionExecutor(manager TransactionManager, options ...Option) TransactionExecutor {
	return &transactionExecutor{
		manager: manager,
		options: NewOptions(options...),
	}
}

func (e *transactionExecutor) Execute(ctx context.Context, fn Func) (any, error) {
	db := e.manager.Connection()

	if e.options.Timeout != -1 {
		ctxWithTimeout, cancelFunc := context.WithTimeout(ctx, e.options.Timeout)
		defer cancelFunc()
		return fn(ctxWithTimeout, db)
	}

	return fn(ctx, db)
}

func (e *transactionExecutor) ExecutorInTransaction(ctx context.Context, fn Func, options ...Option) (any, error) {
	overriddenOptions := WithOptions(e.options, options...)

	if overriddenOptions.Timeout != -1 {
		var cancelFunc context.CancelFunc
		ctx, cancelFunc = context.WithTimeout(ctx, overriddenOptions.Timeout)
		defer cancelFunc()
	}

	transaction, err := e.manager.GetOrCreateTransaction(ctx,
		WithTimeout(-1),
		WithIsolation(overriddenOptions.IsolationLevel),
		WithReadOnly(overriddenOptions.ReadOnly),
		WithPropagation(overriddenOptions.Propagation),
	)

	if err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			e.manager.Rollback(ctx, transaction)
		}
	}()

	var result any
	result, err = fn(ctx, transaction.Tx())

	if err != nil {
		e.manager.Rollback(ctx, transaction)
	}

	e.manager.Commit(ctx, transaction)
	return result, err
}

func (e *transactionExecutor) TransactionManager() TransactionManager {
	return e.manager
}
