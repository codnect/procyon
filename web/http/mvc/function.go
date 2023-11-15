package mvc

type Function[T, E any] func(ctx *Context[T, E]) error
