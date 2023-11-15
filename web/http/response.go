package http

import (
	"io"
)

type Response interface {
	WithWriter(writer io.Writer) Response
	Context() Context
	AddCookie(cookie *Cookie)

	ContentLength() int
	SetContentLength(len int)

	CharacterEncoding() string
	SetCharacterEncoding(charset string)

	ContentType() string
	SetContentType(contentType string)

	AddHeader(name string, value string)
	SetHeader(name string, value string)
	DeleteHeader(name string)
	Header(name string) string
	HeaderNames() []string
	Headers(name string) []string
	Status() Status
	SetStatus(status Status)

	Writer() io.Writer
	Flush()
	IsCommitted() bool
	Reset()
}
