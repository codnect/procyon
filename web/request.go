package web

import (
	"codnect.io/procyon/web/http"
	"io"
	stdhttp "net/http"
)

type ServerRequest struct {
	req    *stdhttp.Request
	ctx    *ServerContext
	reader io.Reader
}

func (r *ServerRequest) WithReader(reader io.Reader) http.Request {
	if reader == nil {
		panic("nil reader")
	}

	copyRequest := new(ServerRequest)
	*copyRequest = *r
	copyRequest.reader = reader
	return copyRequest
}

func (r *ServerRequest) Context() http.Context {
	return r.ctx
}

func (r *ServerRequest) Cookie(name string) (*http.Cookie, bool) {
	return nil, false
}

func (r *ServerRequest) Cookies() []*http.Cookie {
	return nil
}

func (r *ServerRequest) QueryParameter(name string) (string, bool) {
	return "", false
}

func (r *ServerRequest) QueryParameterNames() []string {
	return nil
}

func (r *ServerRequest) QueryParameters(name string) []string {
	return nil
}

func (r *ServerRequest) QueryString() string {
	return ""
}

func (r *ServerRequest) Header(name string) (string, bool) {
	return "", false
}

func (r *ServerRequest) HeaderNames() []string {
	return nil
}

func (r *ServerRequest) Headers(name string) []string {
	return nil
}

func (r *ServerRequest) Path() string {
	return r.req.URL.Path
}

func (r *ServerRequest) Method() http.Method {
	return http.Method(r.req.Method)
}

func (r *ServerRequest) Reader() io.Reader {
	if r.reader != nil {
		return r.reader
	}

	return r.req.Body
}

func (r *ServerRequest) Scheme() string {
	return ""
}

func (r *ServerRequest) IsSecure() bool {
	return false
}
