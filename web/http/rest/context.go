package rest

import (
	"github.com/procyon-projects/procyon/web/http"
	"time"
)

type Context[T, E any] struct {
	ctx            http.Context
	responseEntity ResponseEntity

	headersBuilder HeadersBuilder
	bodyBuilder    BodyBuilder[E]
}

func (c *Context[T, E]) Deadline() (time.Time, bool) {
	return c.ctx.Deadline()
}

func (c *Context[T, E]) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *Context[T, E]) Err() error {
	return c.ctx.Err()
}

func (c *Context[T, E]) Value(key any) any {
	return c.ctx.Value(key)
}

func (c *Context[T, E]) Path() string {
	//return c.ctx.Path()
	return ""
}

func (c *Context[T, E]) Method() string {
	// return c.ctx.Method()
	return ""
}

func (c *Context[T, E]) Get() (e E) {
	return e
}

func (c *Context[T, E]) Ok() BodyBuilder[E] {
	//c.ctx.SetStatus(http.StatusOK)
	return c.bodyBuilder
}

func (c *Context[T, E]) NotFound() HeadersBuilder {
	//c.ctx.SetStatus(http.StatusNotFound)
	return c.headersBuilder
}

func (c *Context[T, E]) NoContent() HeadersBuilder {
	//c.ctx.SetStatus(http.StatusNoContent)
	return c.headersBuilder
}

func (c *Context[T, E]) InternalServerError() BodyBuilder[E] {
	//c.ctx.SetStatus(http.StatusInternalServerError)
	return c.bodyBuilder
}

func (c *Context[T, E]) Created(location string) BodyBuilder[E] {
	//c.ctx.SetStatus(http.StatusCreated)
	//c.headersBuilder.Header("Location", location)
	return c.bodyBuilder
}

func (c *Context[T, E]) BadRequest() BodyBuilder[E] {
	//c.ctx.SetStatus(http.StatusBadRequest)
	return c.bodyBuilder
}

func (c *Context[T, E]) Accepted() BodyBuilder[E] {
	//c.ctx.SetStatus(http.StatusAccepted)
	return c.bodyBuilder
}

func (c *Context[T, E]) Status(status int) BodyBuilder[E] {
	//c.ctx.SetStatus(status)
	return c.bodyBuilder
}
