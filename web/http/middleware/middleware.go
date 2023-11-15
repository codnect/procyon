package middleware

type Middleware struct {
	function     Function
	pathPatterns []string
	order        int
}

func New(function Function, options ...Option) *Middleware {
	middleware := &Middleware{
		function: function,
	}

	for _, opt := range options {
		opt(middleware)
	}

	return middleware
}

func (r *Middleware) PathPatterns() []string {
	return r.pathPatterns
}

func (r *Middleware) Function() Function {
	return r.function
}

func (r *Middleware) Order() int {
	return r.order
}
