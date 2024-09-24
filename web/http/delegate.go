package http

type RequestDelegate interface {
	Invoke(ctx Context)
}

type defaultRequestDelegate struct {
	ctx *defaultServerContext
}

func (d defaultRequestDelegate) Invoke(ctx Context) {
	d.ctx.Invoke(ctx)
}

/*
type contextDelegate struct {
	ctx *Context
}

func (cd *contextDelegate) Invoke(ctx *Context) {
	if cd.ctx.IsCompleted() || cd.ctx.IsAborted() || len(cd.ctx.HandlerChain.functions) <= cd.ctx.handlerIndex {
		return
	}

	next := cd.ctx.HandlerChain.functions[cd.ctx.handlerIndex]
	cd.ctx.handlerIndex++

	err := next(ctx, cd)

	if err != nil {
		cd.ctx.err = err
	}

	if cd.ctx.IsCompleted() || cd.ctx.IsAborted() {
		return
	}

	if cd.ctx.handlerIndex != len(cd.ctx.HandlerChain.functions) {
		cd.ctx.Abort()
	}
}
*/
