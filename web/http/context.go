package http

import (
	"context"
)

type Context interface {
	context.Context

	WithValue(key, val any) Context
	With(request Request, response Response) Context
	WithRequest(request Request) Context
	WithResponse(response Response) Context

	IsCompleted() bool
	Abort()
	IsAborted() bool
	Request() Request
	Response() Response
}
