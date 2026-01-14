package http

func FromQuery[T any](ctx *Context, name string) (T, error) {
	var zero T
	return zero, nil
}

func FromPath[T any](ctx *Context, path string) (T, error) {
	var zero T
	return zero, nil
}

func FromBody[T any](ctx *Context) (T, error) {
	var zero T
	return zero, nil
}

func FromForm[T any](ctx *Context) (T, error) {
	var zero T
	return zero, nil
}

func FromHeader[T any](ctx *Context) (T, error) {
	var zero T
	return zero, nil
}

func FromContext[T any](ctx *Context, name string) (T, error) {
	var zero T
	return zero, nil
}

func ReadFromJson[T any](ctx *Context) (T, error) {
	var zero T
	return zero, nil
}

func ReadForm[T any](ctx *Context) (T, error) {
	var zero T
	return zero, nil
}

func ReadAsString(ctx *Context) (string, error) {
	return "", nil
}

func ReadAsBytes(ctx *Context) ([]byte, error) {
	return nil, nil
}
