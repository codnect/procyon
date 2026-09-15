package web

import "codnect.io/procyon/http"

type Context struct {
	*http.Context
}

type ModelContext[T any] struct {
	*http.EndpointContext[T]
}
