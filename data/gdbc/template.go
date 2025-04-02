package gdbc

import "context"

type Template struct {
	dataSource DataSource
}

func NewTemplate(dataSource DataSource) *Template {
	if dataSource == nil {
		panic("nil datasource")
	}

	return &Template{
		dataSource: dataSource,
	}
}

func (t *Template) obtainConnOperations(ctx context.Context) (Operations, error) {
	// TODO: complete implementation
	//tx, exists := txFromContext(ctx, t.dataSource)
	return t.dataSource.Conn(ctx)
}

func (t *Template) Exec(ctx context.Context, query string, args ...any) (Result, error) {
	conn, err := t.obtainConnOperations(ctx)
	if err != nil {
		return nil, err
	}

	return conn.ExecContext(ctx, query, args...)
}

func (t *Template) Prepare(ctx context.Context, query string) (*Stmt, error) {
	conn, err := t.obtainConnOperations(ctx)
	if err != nil {
		return nil, err
	}

	return conn.PrepareContext(ctx, query)
}

func (t *Template) Query(ctx context.Context, query string, args ...any) (*Rows, error) {
	conn, err := t.obtainConnOperations(ctx)
	if err != nil {
		return nil, err
	}

	return conn.QueryContext(ctx, query, args...)
}

func (t *Template) QueryRow(ctx context.Context, query string, args ...any) (*Row, error) {
	conn, err := t.obtainConnOperations(ctx)
	if err != nil {
		return nil, err
	}

	row := conn.QueryRowContext(ctx, query, args...)
	return row, row.Err()
}
