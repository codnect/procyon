package http

import (
	"io"
)

type Request interface {
	WithReader(reader io.Reader) Request
	Context() Context

	Cookie(name string) (*Cookie, bool)
	Cookies() []*Cookie

	QueryParameter(name string) (string, bool)
	QueryParameterNames() []string
	QueryParameters(name string) []string
	QueryString() string

	Header(name string) (string, bool)
	HeaderNames() []string
	Headers(name string) []string

	Path() string
	Method() Method
	Reader() io.Reader
	Scheme() string
	IsSecure() bool
}
