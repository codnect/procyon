package sql

import (
	"database/sql"
	"database/sql/driver"
)

func Register(name string, driver driver.Driver) {
	sql.Register(name, driver)
}
