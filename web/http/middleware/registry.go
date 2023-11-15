package middleware

type Registry interface {
	Register(path string, middlewareFunction Function, options ...Option)
	Middlewares() []Middleware
}
