package router

import (
	"github.com/procyon-projects/procyon/web/http"
)

type MappingOption func(mapping *Mapping)

func WithMethod(methods ...http.Method) MappingOption {
	return func(mapping *Mapping) {
		if len(methods) != 0 {
			mapping.methods = append(mapping.methods, methods...)
		}
	}
}

func WithAccepts(accepts ...string) MappingOption {
	return func(mapping *Mapping) {
		if len(accepts) != 0 {
			mapping.accepts = append(mapping.accepts, accepts...)
		}
	}
}

func WithContentTypes(contentTypes ...string) MappingOption {
	return func(mapping *Mapping) {
		if len(contentTypes) != 0 {
			mapping.accepts = append(mapping.accepts, contentTypes...)
		}
	}
}
