package sql

import (
	"context"
	"reflect"
)

type ctxTransactionContext struct{}

var (
	ctxTransactionContextKey = &ctxTransactionContext{}
	rTransactionContextType  = reflect.TypeOf((*TransactionContext)(nil))
)

type TransactionContext struct {
}

func FromContext[T *TransactionContext](ctx context.Context) T {
	var value T
	rType := reflect.TypeOf((*T)(nil)).Elem()

	if rType.ConvertibleTo(rTransactionContextType) {
		value = ctx.Value(ctxTransactionContextKey)
	}

	return value
}
