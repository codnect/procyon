package gdbc

import "database/sql"

func NewConnection(db *sql.DB) *Connection {
	if db == nil {
		panic("nil db")
	}

	return &Connection{
		db: db,
	}
}
