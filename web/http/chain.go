package http

type HandlerChain []HandlerFunction

/*
func NewHandlerChain(functions []HandlerFunction) *HandlerChain {
	return &HandlerChain{
		functions: functions,
		size:      len(functions),
	}
}

func (c *HandlerChain) Invoke(ctx Context) {
	if ctx == nil || len(c.functions) == 0 {
		return
	}

	nextHandler := ctx.nextHandler()

	if ctx.IsCompleted() || ctx.IsAborted() || c.size <= nextHandler {
		return
	}

	next := c.functions[nextHandler]
	nextHandler++
	ctx.setNextHandler(nextHandler)

	err := next(ctx, c)

	if err != nil {
		ctx.setErr(err)
	}

	if ctx.IsCompleted() || ctx.IsAborted() {
		return
	}

	if nextHandler != c.size {
		ctx.Abort()
	} else {
		ctx.complete()
	}
}*/

/*
func (c *HandlerChain) Invoke(ctx *Context) {
	if ctx == nil || len(c.functions) == 0 {
		return
	}

	if ctx.HandlerChain == nil {
		ctx.HandlerChain = c
	}

	nextHandler := ctx.nextHandler()
	next := c.functions[nextHandler]
	ctx.handlerIndex = 1
	err := next(ctx, &ctx.delegate)

	if err != nil {
		ctx.err = err
	}

	ctx.completed = true
	return
}
*/
