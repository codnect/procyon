package web

import (
	"time"
)

type RestContext[E, R any] struct {
	ctx            Context
	headersBuilder HeadersBuilder
	bodyBuilder    BodyBuilder[R]
	request        E
}

func (c *RestContext[E, R]) Deadline() (time.Time, bool) {
	return c.ctx.Deadline()
}

func (c *RestContext[E, R]) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *RestContext[E, R]) Err() error {
	return c.ctx.Err()
}

func (c *RestContext[E, R]) Value(key any) any {
	return c.ctx.Value(key)
}

func (c *RestContext[E, R]) Path() string {
	return c.ctx.Path()
}

func (c *RestContext[E, R]) Method() HttpMethod {
	return c.ctx.Method()
}

func (c *RestContext[E, R]) Get() E {
	return c.request
}

func (c *RestContext[E, R]) Ok() BodyBuilder[R] {
	c.ctx.SetStatus(StatusOK)
	return c.bodyBuilder
}

func (c *RestContext[E, R]) NotFound() HeadersBuilder {
	c.ctx.SetStatus(StatusNotFound)
	return c.headersBuilder
}

func (c *RestContext[E, R]) NoContent() HeadersBuilder {
	c.ctx.SetStatus(StatusNoContent)
	return c.headersBuilder
}

func (c *RestContext[E, R]) InternalServerError() BodyBuilder[R] {
	c.ctx.SetStatus(StatusInternalServerError)
	return c.bodyBuilder
}

func (c *RestContext[E, R]) Created(location string) BodyBuilder[R] {
	c.ctx.SetStatus(StatusCreated)
	c.headersBuilder.Header("Location", location)
	return c.bodyBuilder
}

func (c *RestContext[E, R]) BadRequest() BodyBuilder[R] {
	c.ctx.SetStatus(StatusBadRequest)
	return c.bodyBuilder
}

func (c *RestContext[E, R]) Accepted() BodyBuilder[R] {
	c.ctx.SetStatus(StatusAccepted)
	return c.bodyBuilder
}

func (c *RestContext[E, R]) Status(status HttpStatus) BodyBuilder[R] {
	c.ctx.SetStatus(status)
	return c.bodyBuilder
}
