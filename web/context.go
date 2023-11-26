package web

import (
	"codnect.io/procyon/web/http"
	"context"
	stdhttp "net/http"
	"time"
)

type ServerContext struct {
	parent   *ServerContext
	context  context.Context
	request  ServerRequest
	response ServerResponse

	HandlerChain     http.HandlerChain
	nextHandlerIndex int

	err       error
	completed bool
	aborted   bool

	delegate      ServerContextDelegate
	pathVariables http.PathVariables
}

func newServerContext() *ServerContext {
	return &ServerContext{
		pathVariables: http.PathVariables{},
	}
}

func (c *ServerContext) WithValue(key, val any) http.Context {
	copyContext := new(ServerContext)
	*copyContext = *c

	ctx := c.context
	if ctx == nil {
		ctx = context.Background()
	}

	copyContext.context = context.WithValue(ctx, key, val)
	return copyContext
}

func (c *ServerContext) With(request http.Request, response http.Response) http.Context {
	if request == nil {
		panic("nil request")
	}

	if response == nil {
		panic("nil response")
	}

	copyContext := new(ServerContext)
	*copyContext = *c
	copyContext.request = *(request.(*ServerRequest))
	copyContext.response = *(response.(*ServerResponse))

	if c.parent == nil {
		copyContext.parent = c
	}

	return copyContext
}

func (c *ServerContext) WithRequest(request http.Request) http.Context {
	if request == nil {
		panic("nil request")
	}

	copyContext := new(ServerContext)
	*copyContext = *c
	copyContext.request = *(request.(*ServerRequest))

	if c.parent == nil {
		copyContext.parent = c
	}

	return copyContext
}

func (c *ServerContext) WithResponse(response http.Response) http.Context {
	if response == nil {
		panic("nil response")
	}

	copyContext := new(ServerContext)
	*copyContext = *c
	copyContext.response = *(response.(*ServerResponse))

	if c.parent == nil {
		copyContext.parent = c
	}

	return copyContext
}

func (c *ServerContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *ServerContext) Done() <-chan struct{} {
	return nil
}

func (c *ServerContext) setErr(err error) {
	if c.parent != nil {
		c.parent.setErr(err)
	} else {
		c.err = err
	}
}

func (c *ServerContext) Err() error {
	if c.parent != nil {
		return c.parent.Err()
	}

	return c.err
}

func (c *ServerContext) Value(key any) any {
	if key == http.PathVariablesAttribute {
		return &c.pathVariables
	}

	if c.context == nil {
		return nil
	}

	return c.context.Value(key)
}

func (c *ServerContext) Parent() *ServerContext {
	return c.parent
}

func (c *ServerContext) complete() {
	if c.parent != nil {
		c.parent.complete()
	} else {
		c.completed = true
	}
}

func (c *ServerContext) IsCompleted() bool {
	if c.parent != nil {
		return c.parent.IsCompleted()
	}

	return c.completed
}

func (c *ServerContext) Abort() {
	if c.parent != nil {
		c.parent.Abort()
	} else {
		c.aborted = true
	}
}

func (c *ServerContext) IsAborted() bool {
	if c.parent != nil {
		return c.parent.IsAborted()
	}

	return c.aborted
}

func (c *ServerContext) Request() http.Request {
	return &c.request
}

func (c *ServerContext) Response() http.Response {
	return &c.response
}

func (c *ServerContext) Reset(req *stdhttp.Request, writer stdhttp.ResponseWriter) {
	/*if !c.IsCompleted() {
		return
	}*/

	c.request.req = req
	c.response.writer = writer
	c.delegate.ctx = c

	c.parent = nil
	c.err = nil
	c.context = nil
	c.completed = false
	c.aborted = false

	c.nextHandlerIndex = 0
	//c.pathVariables.currentIndex = 0
}

func (c *ServerContext) nextHandler() int {
	if c.parent != nil {
		return c.parent.nextHandler()
	}

	return c.nextHandlerIndex
}

func (c *ServerContext) setNextHandler(nextHandler int) {
	if c.parent != nil {
		c.parent.setNextHandler(nextHandler)
	} else {
		c.nextHandlerIndex = nextHandler
	}
}

func (c *ServerContext) Invoke(ctx http.Context) {
	if len(c.HandlerChain) == 0 {
		return
	}

	nextHandler := c.nextHandler()

	if c.IsCompleted() || c.IsAborted() || len(c.HandlerChain) <= nextHandler {
		return
	}

	next := c.HandlerChain[nextHandler]
	nextHandler++
	c.setNextHandler(nextHandler)

	err := next(ctx, c.delegate)

	if err != nil {
		c.setErr(err)
	}

	if c.IsCompleted() || c.IsAborted() {
		return
	}

	if nextHandler != len(c.HandlerChain) {
		c.Abort()
	} else {
		c.complete()
	}
}
