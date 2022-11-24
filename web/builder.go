package web

import "procyon-test/web/mediatype"

type HeadersBuilder struct {
	response *Response
}

func (b HeadersBuilder) Header(headerName string, headerValues ...string) HeadersBuilder {
	for _, val := range headerValues {
		b.response.headers[headerName] = append(b.response.headers[headerName], val)
	}

	return b
}

func (b HeadersBuilder) Headers(headers HttpHeaders) HeadersBuilder {
	for name, value := range headers {
		b.Header(name, value...)
	}

	return b

}

type BodyBuilder[R any] struct {
	response *Response
}

func (b BodyBuilder[R]) Header(headerName string, headerValues ...string) BodyBuilder[R] {
	for _, val := range headerValues {
		b.response.headers[headerName] = append(b.response.headers[headerName], val)
	}

	return b
}

func (b BodyBuilder[R]) Headers(headers HttpHeaders) BodyBuilder[R] {
	for name, value := range headers {
		b.Header(name, value...)
	}

	return b
}

func (b BodyBuilder[R]) Body(body R) BodyBuilder[R] {
	b.response.entity = body

	return b
}

func (b BodyBuilder[R]) ContentType(contentType mediatype.MediaType) BodyBuilder[R] {
	b.response.contentType = contentType

	return b
}
