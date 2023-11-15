package http

type ErrorHandler interface {
	HandleError(ctx Context, err error)
}

type NotFoundError struct {
}

func (e NotFoundError) Error() string {
	return ""
}
