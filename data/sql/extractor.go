package sql

func ExtractRow[T any](row Row, callback func(row Row) T) (t T) {
	return callback(row)
}

func ExtractRowSet[T any](rowSet RowSet, callback func(rowSet RowSet) []T) []T {
	return callback(rowSet)
}
