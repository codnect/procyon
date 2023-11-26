package page

import (
	"codnect.io/procyon/data/sort"
)

var (
	unpagedRequest = newRequest(0, 0, false)
)

type request struct {
	page    int
	size    int
	sort    sort.Sort
	isPaged bool
}

func newRequest(pageNumber, pageSize int, isPaged bool) *request {
	return &request{
		page:    pageNumber,
		size:    pageSize,
		isPaged: isPaged,
	}
}

func (r *request) PageNumber() int {
	if !r.isPaged {
		panic("unsupported operation")
	}

	return r.page
}

func (r *request) PageSize() int {
	if !r.isPaged {
		panic("unsupported operation")
	}

	return r.size
}

func (r *request) Offset() int {
	if !r.isPaged {
		panic("unsupported operation")
	}

	return r.page * r.size
}

func (r *request) Sort() sort.Sort {
	/* todo fix sort method */
	return r.sort
}

func (r *request) IsPaged() bool {
	return r.isPaged
}

func (r *request) IsUnpaged() bool {
	return !r.isPaged
}

func (r *request) HasPrevious() bool {
	if !r.isPaged {
		return false
	}

	return r.page > 0
}

func (r *request) PreviousOrFirst() Pageable {
	if !r.isPaged {
		return r
	}

	if r.HasPrevious() {
		return r.Previous()
	}

	return r.First()
}

func (r *request) Previous() Pageable {
	if r.PageNumber() == 0 {
		return r
	}

	return newRequest(r.PageNumber()-1, r.PageSize(), r.IsPaged())
}

func (r *request) First() Pageable {
	if !r.isPaged {
		return r
	}

	return newRequest(0, r.PageSize(), r.IsPaged())
}

func (r *request) Next() Pageable {
	if !r.isPaged {
		return r
	}

	return newRequest(r.PageNumber()+1, r.PageSize(), r.IsPaged())
}
