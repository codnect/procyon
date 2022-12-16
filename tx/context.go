package tx

import (
	"context"
	"fmt"
	"reflect"
)

type ctxTransactionContext struct{}

var (
	ctxTransactionContextKey = &ctxTransactionContext{}
	reflTransactionContext   = reflect.TypeOf((*TransactionContext)(nil)).Elem()
)

type TransactionContext interface {
	Transaction() Transaction
	Parent() TransactionContext
	Resources() any
}

type transactionContext struct {
}

func NewTransactionContext() TransactionContext {
	return &transactionContext{}
}

func (t *transactionContext) Transaction() Transaction {
	return nil
}

func (t *transactionContext) Parent() TransactionContext {
	return nil
}

func (t *transactionContext) Resources() any {
	return nil
}

func FromContext[T any](ctx context.Context) T {
	var value T
	reflType := reflect.TypeOf((*T)(nil)).Elem()

	switch {
	case reflType.ConvertibleTo(reflTransactionContext):
		value, _ = ctx.Value(ctxTransactionContextKey).(T)
	case reflType.ConvertibleTo(reflTransactionExecutor):
		value, _ = ctx.Value(ctxTransactionExecutorKey).(T)
	case reflType.ConvertibleTo(reflTransactionManager):
		value, _ = ctx.Value(ctxTransactionManagerKey).(T)
	default:
		panic(fmt.Sprintf("tx: type %s is not supported", reflType.Name()))
	}

	return value
}

func ToContext[T any](parent context.Context, value T) context.Context {
	switch any(value).(type) {
	case TransactionContext:
		return context.WithValue(parent, ctxTransactionContextKey, value)
	case TransactionExecutor:
		return context.WithValue(parent, ctxTransactionExecutorKey, value)
	case TransactionManager:
		return context.WithValue(parent, ctxTransactionManagerKey, value)
	}

	reflType := reflect.TypeOf((*T)(nil)).Elem()
	panic(fmt.Sprintf("tx: type %s is not supported", reflType.Name()))
}
