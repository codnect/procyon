package middleware

type Option func(middleware *Middleware)

func WithOrder(order int) Option {
	return func(middleware *Middleware) {
		middleware.order = order
	}
}
