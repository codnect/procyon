package gdbc

import "database/sql"

type Isolation = sql.IsolationLevel

const (
	IsolationDefault Isolation = iota
	IsolationReadUncommitted
	IsolationReadCommitted
	IsolationWriteCommitted
	IsolationRepeatableRead
	IsolationSnapshot
	IsolationSerializable
	IsolationLinearizable
)
