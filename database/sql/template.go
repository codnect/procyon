package sql

import "context"

type TransactionCallback func(ctx context.Context) (any, error)

type TransactionTemplate struct {
	manager TransactionManager
}

func NewTransactionTemplate(manager TransactionManager) *TransactionTemplate {
	return &TransactionTemplate{
		manager: manager,
	}
}

func (t *TransactionTemplate) Execute(ctx context.Context, callback TransactionCallback) (result any, err error) {
	_, err = t.manager.GetTransaction(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			err = t.manager.Rollback()
		}
	}()

	result, err = callback(ctx)

	if err != nil {
		err = t.manager.Rollback()

		if err != nil {
			return nil, err
		}
	}

	err = t.manager.Commit()

	if err != nil {
		return nil, err
	}

	return nil, nil
}
