package web

import (
	"fmt"
	"time"
)

type MvcContext[E, R any] struct {
	ctx Context
}

func (c *MvcContext[E, R]) Deadline() (time.Time, bool) {
	return c.ctx.Deadline()
}

func (c *MvcContext[E, R]) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *MvcContext[E, R]) Err() error {
	return c.ctx.Err()
}

func (c *MvcContext[E, R]) Value(key any) any {
	return c.ctx.Value(key)
}

func (c *MvcContext[E, R]) Path() string {
	return c.ctx.Path()
}

func (c *MvcContext[E, R]) Method() HttpMethod {
	return c.ctx.Method()
}

func (c *MvcContext[E, R]) ViewName(name string) {
	c.ctx.response.viewName = name
}

func (c *MvcContext[E, R]) Redirect(location string) {
	c.ctx.response.viewName = fmt.Sprintf("redirect:%s", location)
}
