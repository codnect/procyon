package web

import "github.com/procyon-projects/procyon/web/mediatype"

type Response struct {
	status      HttpStatus
	contentType mediatype.MediaType
	viewName    string
	entity      any
	headers     HttpHeaders
}

func (r *Response) Status() HttpStatus {
	return r.status
}

func (r *Response) ContentType() mediatype.MediaType {
	return r.contentType
}

func (r *Response) Headers() HttpHeaders {
	return r.headers
}

func (r *Response) ViewName() string {
	return r.viewName
}

func (r *Response) Entity() any {
	return r.entity
}
