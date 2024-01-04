package sql

import "context"

type TransactionManager interface {
	GetTransaction(ctx context.Context, options ...any) (TransactionStatus, error)
	Commit() error
	Rollback() error
}

type DataSourceTransactionManager struct {
	dataSource DataSource
}

func NewDataSourceTransactionManager(dataSource DataSource) *DataSourceTransactionManager {
	return &DataSourceTransactionManager{
		dataSource: dataSource,
	}
}

func (m *DataSourceTransactionManager) GetTransaction(ctx context.Context, options ...any) (TransactionStatus, error) {
	return nil, nil
}

func (m *DataSourceTransactionManager) Commit() error {
	return nil
}

func (m *DataSourceTransactionManager) Rollback() error {
	return nil
}
