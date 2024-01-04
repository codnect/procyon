package sql

type DataSource interface {
	GetConnection() (Connection, error)
}
