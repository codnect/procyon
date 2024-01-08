package gdbc

import (
	"codnect.io/procyon/data/sql"
	"context"
)

type Template struct {
	dataSource sql.DataSource
}

func NewTemplate(dataSource sql.DataSource) *Template {
	return &Template{
		dataSource: dataSource,
	}
}

func (t *Template) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return nil, nil
}

func (t *Template) Prepare(ctx context.Context, query string) (sql.Statement, error) {
	return nil, nil
}

func (t *Template) Query(ctx context.Context, query string, args ...any) (sql.RowSet, error) {
	return nil, nil
}

func (t *Template) QueryRow(ctx context.Context, query string, args ...any) sql.Row {
	return nil
}
