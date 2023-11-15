package router

import (
	"github.com/procyon-projects/procyon/web/http"
)

type Mapping struct {
	pattern      string
	methods      []http.Method
	accepts      []string
	contentTypes []string
	attributes   map[string]any
}

func NewMapping(pattern string, opts ...MappingOption) *Mapping {
	mapping := &Mapping{
		pattern: pattern,
	}

	for _, option := range opts {
		option(mapping)
	}

	return mapping
}

func (r *Mapping) Pattern() string {
	return r.pattern
}

func (r *Mapping) Methods() []http.Method {
	methods := make([]http.Method, 0)
	methods = append(methods, r.methods...)
	return methods
}

func (r *Mapping) Accepts() []string {
	accepts := make([]string, 0)
	accepts = append(accepts, r.accepts...)
	return accepts
}

func (r *Mapping) ContentTypes() []string {
	contentTypes := make([]string, 0)
	contentTypes = append(contentTypes, r.contentTypes...)
	return contentTypes
}

func (r *Mapping) Attributes() map[string]any {
	attributes := make(map[string]any)

	for k, v := range r.attributes {
		attributes[k] = v
	}

	return attributes
}
