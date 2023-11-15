package web

import (
	"github.com/procyon-projects/procyon/web/http"
	"io"
	stdhttp "net/http"
)

type ServerResponse struct {
	responseWriter stdhttp.ResponseWriter
	ctx            http.Context
	writer         io.Writer
}

func (r *ServerResponse) WithWriter(writer io.Writer) http.Response {
	if writer == nil {
		panic("nil writer")
	}

	copyResponse := new(ServerResponse)
	*copyResponse = *r
	copyResponse.writer = writer
	return copyResponse
}

func (r *ServerResponse) Context() http.Context {
	return r.ctx
}

func (r *ServerResponse) AddCookie(cookie *http.Cookie) {

}

func (r *ServerResponse) ContentLength() int {
	return 0
}

func (r *ServerResponse) SetContentLength(len int) {

}

func (r *ServerResponse) CharacterEncoding() string {
	return ""
}

func (r *ServerResponse) SetCharacterEncoding(charset string) {

}

func (r *ServerResponse) ContentType() string {
	return ""
}

func (r *ServerResponse) SetContentType(contentType string) {

}

func (r *ServerResponse) AddHeader(name string, value string) {

}

func (r *ServerResponse) SetHeader(name string, value string) {

}

func (r *ServerResponse) DeleteHeader(name string) {

}

func (r *ServerResponse) Header(name string) string {
	return ""
}

func (r *ServerResponse) HeaderNames() []string {
	return nil
}

func (r *ServerResponse) Headers(name string) []string {
	return nil
}

func (r *ServerResponse) Status() http.Status {
	return 0
}

func (r *ServerResponse) SetStatus(status http.Status) {

}

func (r *ServerResponse) Writer() io.Writer {
	return nil
}

func (r *ServerResponse) Flush() {

}

func (r *ServerResponse) IsCommitted() bool {
	return false
}

func (r *ServerResponse) Reset() {

}
