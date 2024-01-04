package sql

type TransactionStatus interface {
	IsNewTransaction() bool
	SetRollbackOnly()
	IsRollbackOnly() bool
	IsCompleted() bool
}

type DefaultTransactionStatus struct {
}

func NewDefaultTransactionStatus() *DefaultTransactionStatus {
	return &DefaultTransactionStatus{}
}
