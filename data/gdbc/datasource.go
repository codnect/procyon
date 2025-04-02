package gdbc

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

type DataSource interface {
	Conn(ctx context.Context) (Conn, error)
}

type DataSourceProperties struct {
	Url        string
	Username   string
	Password   string
	DriverName string
}

type SimpleDataSource struct {
	props DataSourceProperties
}

func NewSimpleDataSource(props DataSourceProperties) *SimpleDataSource {
	return &SimpleDataSource{
		props: props,
	}
}

func (ds *SimpleDataSource) Conn(ctx context.Context) (Conn, error) {
	var (
		db            *sql.DB
		driverName    string
		connectionUrl string
		err           error
	)

	driverName, connectionUrl, err = ds.parseDataSourceUrl()
	if err != nil {
		return nil, err
	}

	db, err = sql.Open(driverName, connectionUrl)
	if err != nil {
		return nil, err
	}

	return db.Conn(ctx)
}

func (ds *SimpleDataSource) parseDataSourceUrl() (string, string, error) {
	if len(ds.props.Url) == 0 {
		return "", "", errors.New("url must not be empty")
	}

	src := strings.Split(ds.props.Url, ":")
	if len(src) < 3 {
		return "", "", errors.New("url format is wrong")
	}

	scheme := src[0]
	if "gdbc" != scheme {
		return "", "", errors.New("scheme is not valid")
	}

	driverName := src[1]
	if len(driverName) == 0 {
		return "", "", errors.New("driver name must not be empty")
	}

	return driverName, strings.Join(append(src[:1], src[2:]...), ":"), nil
}
