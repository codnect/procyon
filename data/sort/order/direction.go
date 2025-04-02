package order

type Direction int

const (
	Ascending Direction = iota + 1
	Descending
	Default = Ascending
)
