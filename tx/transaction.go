package tx

import "database/sql"

type Transaction interface {
	Tx() *sql.Tx
	IsNew() bool
	IsReadOnly() bool
	MarkAsRollbackOnly()
	IsRollbackOnly() bool
	MarkAsCompleted()
	IsCompleted() bool
}

type transaction struct {
	tx             *sql.Tx
	newTransaction bool
	readOnly       bool
	rollbackOnly   bool
	completed      bool
}

func NewTransaction(tx *sql.Tx, newTransaction bool, readOnly bool) Transaction {
	return &transaction{
		tx,
		false,
		false,
		false,
		false,
	}
}

func (t *transaction) Tx() *sql.Tx {
	return t.tx
}

func (t *transaction) IsNew() bool {
	return t.newTransaction
}

func (t *transaction) IsReadOnly() bool {
	return t.readOnly
}

func (t *transaction) MarkAsRollbackOnly() {
	t.rollbackOnly = true
}

func (t *transaction) IsRollbackOnly() bool {
	return t.rollbackOnly
}

func (t *transaction) MarkAsCompleted() {
	t.completed = true
}

func (t *transaction) IsCompleted() bool {
	return t.completed
}
