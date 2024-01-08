package gdbc

import (
	"codnect.io/procyon/data/sql"
	dsql "database/sql"
	"errors"
	"strings"
)

type DataSource struct {
	props DataSourceProperties
}

type DataSourceProperties struct {
	Url      string
	Username string
	Password string
}

func NewDataSource(props DataSourceProperties) *DataSource {
	return &DataSource{
		props: props,
	}
}

func (ds *DataSource) GetConnection() (sql.Connection, error) {
	var (
		db            *dsql.DB
		driverName    string
		connectionUrl string
		err           error
	)

	driverName, connectionUrl, err = ds.parseDataSourceUrl()
	if err != nil {
		return nil, err
	}

	db, err = dsql.Open(driverName, connectionUrl)

	if err != nil {
		return nil, err
	}

	return NewConnection(db), nil
}

func (ds *DataSource) parseDataSourceUrl() (string, string, error) {
	if len(ds.props.Url) == 0 {
		return "", "", errors.New("url must not be empty")
	}

	src := strings.Split(ds.props.Url, ":")
	if len(src) < 3 {
		return "", "", errors.New("url format is wrong : " + ds.props.Url)
	}

	scheme := src[0]
	if "gdbc" != scheme {
		return "", "", errors.New("driver name must not be empty")
	}

	driverName := src[1]
	if len(driverName) == 0 {
		return "", "", errors.New("driver name must not be empty")
	}

	return driverName, strings.Join(append(src[:1], src[2:]...), ":"), nil
}
