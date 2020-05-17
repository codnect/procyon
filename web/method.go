package web

import "net/http"

type HttpMethod string

const (
	HttpMethodGet   HttpMethod = http.MethodGet
	HttpMethodHead  HttpMethod = http.MethodHead
	HttpMethodPost  HttpMethod = http.MethodPost
	HttpMethodPut   HttpMethod = http.MethodPut
	HttpMethodPatch HttpMethod = http.MethodPatch
)
