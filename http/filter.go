package http

type FilterFunc func(ctx *Context, next FilterDelegate) (Result, error)

func (f FilterFunc) Filter(ctx *Context, next FilterDelegate) (Result, error) {
	return f(ctx, next)
}

type Filter interface {
	Filter(ctx *Context, next FilterDelegate) (Result, error)
}
