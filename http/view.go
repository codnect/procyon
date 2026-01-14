package http

type Template interface {
	Render(ctx *Context, model map[string]any) error
}

type TemplateResolver interface {
	Resolve(ctx *Context, view string) (Template, error)
}
