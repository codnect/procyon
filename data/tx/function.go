package tx

import (
	"codnect.io/procyon/runtime/types"
	"context"
	"errors"
)

type FuncOption func(attributes Attributes)

func Execute[R any](ctx context.Context, txManager Manager, callback func(txCtx context.Context) (R, error), opts ...FuncOption) (result R, err error) {
	if ctx == nil {
		return result, errors.New("nil context")
	}

	if txManager == nil {
		return result, errors.New("nil tx manager")
	}

	// apply options
	attributes := Attributes{}
	for _, opt := range opts {
		opt(attributes)
	}

	ctx, err = txManager.ContextWithTx(ctx, attributes)
	if err != nil {
		return result, err
	}

	defer func() {
		if runtimeErr := recover(); runtimeErr != nil {
			err = txManager.Rollback(ctx)
			// TODO: handle runtime error
		}
	}()

	result, err = callback(ctx)
	if err != nil {
		return result, txManager.Rollback(ctx)
	}

	return result, txManager.Commit(ctx)
}

func ExecuteWithoutResult(ctx context.Context, txManager Manager, callback func(txCtx context.Context) error, opts ...FuncOption) error {
	_, err := Execute[types.Void](ctx, txManager, func(txCtx context.Context) (types.Void, error) {
		return nil, callback(txCtx)
	}, opts...)

	return err
}
