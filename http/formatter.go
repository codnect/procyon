package http

type InputFormatter interface {
	CanRead(ctx *Context) bool
	Read(ctx *Context) error
}

type OutputFormatter interface {
	CanWriteResult(ctx *Context, result Result) bool
	WriteResult(ctx *Context, result Result) error
}
