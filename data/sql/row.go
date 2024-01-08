package sql

type Row interface {
	Scan(dest ...any) error
}

type RowSet interface {
	Row

	Next() bool
}
