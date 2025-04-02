package sort

import (
	order2 "codnect.io/procyon/data/sort/order"
)

var (
	unsorted = newSort()
)

type Sort interface {
	IsEmpty() bool
	IsSorted() bool
	IsUnsorted() bool
	Orders() []order2.Order
	And(sort Sort) Sort
}

type sort struct {
	orders []order2.Order
}

func newSort(orders ...order2.Order) Sort {
	return &sort{
		orders: orders,
	}
}

func (s *sort) IsEmpty() bool {
	return len(s.orders) == 0
}

func (s *sort) IsSorted() bool {
	return !s.IsEmpty()
}

func (s *sort) IsUnsorted() bool {
	return s.IsEmpty()
}

func (s *sort) Orders() []order2.Order {
	return s.orders
}

func (s *sort) And(sort Sort) Sort {
	if sort == nil {
		panic("sort must not be nil")
	}

	orders := make([]order2.Order, len(s.orders))
	copy(orders, s.orders)

	if len(sort.Orders()) != 0 {
		orders = append(orders, sort.Orders()...)
	}

	return newSort(orders...)
}

func Unsorted() Sort {
	return unsorted
}

func By(direction order2.Direction, properties ...string) Sort {
	if len(properties) == 0 {
		return unsorted
	}

	orders := make([]order2.Order, len(properties))

	for index, property := range properties {
		orders[index] = order2.By(property, order2.WithDirection(direction))
	}

	return newSort(orders...)
}

func ByOrder(orders ...order2.Order) Sort {
	if len(orders) == 0 {
		return unsorted
	}

	return newSort(orders...)
}

func ByProperties(properties ...string) Sort {
	return By(order2.Default, properties...)
}
