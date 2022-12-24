package web

import "net/http"

type HttpMethod string

const (
	MethodGet     HttpMethod = http.MethodGet
	MethodPost    HttpMethod = http.MethodPost
	MethodPut     HttpMethod = http.MethodPut
	MethodDelete  HttpMethod = http.MethodDelete
	MethodPatch   HttpMethod = http.MethodPatch
	MethodOptions HttpMethod = http.MethodOptions
	MethodHead    HttpMethod = http.MethodHead
)
