package http

import "reflect"

func Typed[E, R any](fn func(ctx *RequestContext[E]) (ResponseEntity[R], error)) Handler {
	v := reflect.TypeFor[E]()
	r := reflect.TypeFor[R]()
	if v != nil {

	}

	if r != nil {

	}
	return nil
}

type RequestContext[E any] struct {
}

func (c *RequestContext[E]) Body() (*E, error) {
	return nil, nil
}

type ResponseEntity[T any] struct {
	Body   T
	Status Status
}

func EntityOk[T any](val T) ResponseEntity[T] {
	return ResponseEntity[T]{
		Body:   val,
		Status: StatusOK,
	}
}

func EntityCreated[T any](location string, val T) ResponseEntity[T] {
	return ResponseEntity[T]{
		Body:   val,
		Status: StatusCreated,
	}
}

func EntityNoContent[T any]() ResponseEntity[T] {
	return ResponseEntity[T]{
		Status: StatusNoContent,
	}
}

func EntityNotFound[T any]() ResponseEntity[T] {
	return ResponseEntity[T]{
		Status: StatusNotFound,
	}
}
