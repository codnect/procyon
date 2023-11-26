package page

import (
	"codnect.io/procyon/data/sort"
)

type Pageable interface {
	PageNumber() int
	PageSize() int
	Offset() int
	Sort() sort.Sort
	IsPaged() bool
	IsUnpaged() bool
	First() Pageable
	HasPrevious() bool
	Next() Pageable
	PreviousOrFirst() Pageable
}

func Unpaged() Pageable {
	return unpagedRequest
}

func Of(pageNumber, pageSize int) Pageable {
	return newRequest(pageNumber, pageSize, true)
}

func OfSize(pageSize int) Pageable {
	return newRequest(0, pageSize, true)
}
