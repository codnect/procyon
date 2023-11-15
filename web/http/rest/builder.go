package rest

import "net/http"

type HeadersBuilder struct {
	responseEntity *ResponseEntity
}

func (b HeadersBuilder) Header(headerName string, headerValues ...string) HeadersBuilder {
	for _, val := range headerValues {
		b.responseEntity.headers[headerName] = append(b.responseEntity.headers[headerName], val)
	}

	return b
}

func (b HeadersBuilder) Headers(headers http.Header) HeadersBuilder {
	for name, value := range headers {
		b.Header(name, value...)
	}

	return b

}

type BodyBuilder[R any] struct {
	responseEntity *ResponseEntity
}

func (b BodyBuilder[R]) Header(headerName string, headerValues ...string) BodyBuilder[R] {
	for _, val := range headerValues {
		b.responseEntity.headers[headerName] = append(b.responseEntity.headers[headerName], val)
	}

	return b
}

func (b BodyBuilder[R]) Headers(headers http.Header) BodyBuilder[R] {
	for name, value := range headers {
		b.Header(name, value...)
	}

	return b
}

func (b BodyBuilder[R]) Body(body R) BodyBuilder[R] {
	b.responseEntity.body = body

	return b
}

func (b BodyBuilder[R]) ContentType(contentType string) BodyBuilder[R] {
	//b.responseEntity.contentType = contentType

	return b
}
