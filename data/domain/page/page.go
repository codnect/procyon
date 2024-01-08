package page

import (
	"codnect.io/procyon/data/domain/sort"
	"math"
)

type Paginated interface {
	TotalPages() int
	TotalElements() int
	HasNext() bool
	Number() int
	Size() int
	NumberOfElements() int
	HasPrevious() bool
	IsFirst() bool
	IsLast() bool
	NextPageable() Pageable
	PreviousPageable() Pageable
	HasContent() bool
	Content() any
	Pageable() Pageable
	Sort() sort.Sort
}

type Page[T any] interface {
	Paginated

	Items() []T
}

type page[T any] struct {
	content  []T
	pageable Pageable
	total    int
}

func newPage[T any](content []T, pageable Pageable, total int) *page[T] {
	if content == nil {
		panic("content must not be nil")
	}

	if pageable == nil {
		panic("pageable must not be nil")
	}

	if len(content) != 0 && pageable.Offset()+pageable.PageSize() > total {
		total = pageable.Offset() + len(content)
	}

	return &page[T]{
		content:  content,
		pageable: pageable,
		total:    total,
	}
}

func (p *page[T]) TotalPages() int {
	if p.Size() == 0 {
		return 1
	}

	return int(math.Ceil(float64(p.total) / float64(p.Size())))
}

func (p *page[T]) TotalElements() int {
	return p.total
}

func (p *page[T]) HasNext() bool {
	return p.Number()+1 < p.TotalPages()
}

func (p *page[T]) Number() int {
	if p.pageable.IsPaged() {
		return p.pageable.PageNumber()
	}

	return 0
}

func (p *page[T]) Size() int {
	if p.pageable.IsPaged() {
		return p.pageable.PageSize()
	}

	return len(p.content)
}

func (p *page[T]) NumberOfElements() int {
	return len(p.content)
}

func (p *page[T]) HasPrevious() bool {
	return p.Number() > 0
}

func (p *page[T]) IsFirst() bool {
	return !p.HasPrevious()
}

func (p *page[T]) IsLast() bool {
	return !p.HasNext()
}

func (p *page[T]) NextPageable() Pageable {
	if p.HasNext() {
		return p.pageable.Next()
	}

	return unpagedRequest
}

func (p *page[T]) PreviousPageable() Pageable {
	if p.HasPrevious() {
		return p.pageable.PreviousOrFirst()
	}

	return unpagedRequest
}

func (p *page[T]) HasContent() bool {
	return len(p.content) != 0
}

func (p *page[T]) Content() any {
	return p.content
}

func (p *page[T]) Items() []T {
	copyOfContent := make([]T, len(p.content))
	copy(copyOfContent, p.content)
	return copyOfContent
}

func (p *page[T]) Pageable() Pageable {
	return p.pageable
}

func (p *page[T]) Sort() sort.Sort {
	return p.pageable.Sort()
}

func Empty[T any]() Page[T] {
	return New(make([]T, 0), unpagedRequest, 0)
}

func New[T any](content []T, pageable Pageable, total int) Page[T] {
	return newPage(content, pageable, total)
}

func WithContent[T any](content []T) Page[T] {
	return New[T](content, unpagedRequest, len(content))
}
