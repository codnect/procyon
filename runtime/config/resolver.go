// Copyright 2025 Codnect
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"errors"
	"fmt"
)

// PropertyResolver is an interface for resolving properties and expanding placeholders.
type PropertyResolver interface {
	// Lookup returns the value of the given property name.
	Lookup(propertyName string) (any, bool)
	// LookupOrDefault returns the value of the given property name.
	LookupOrDefault(propertyName string, defaultValue any) any
	// Expand resolves placeholders in the given text.
	Expand(text string) string
	// ExpandStrict resolves placeholders in the given text.
	// If a placeholder cannot be resolved, it returns an error.
	ExpandStrict(text string) (string, error)
}

// DefaultPropertyResolver is the default implementation of the PropertyResolver interface.
type DefaultPropertyResolver struct {
	propSources *PropertySources
}

// NewDefaultPropertyResolver creates a new DefaultPropertyResolver with the given sources.
func NewDefaultPropertyResolver(propSources *PropertySources) *DefaultPropertyResolver {
	if propSources == nil {
		panic("nil property sources")
	}

	return &DefaultPropertyResolver{
		propSources: propSources,
	}
}

// Lookup returns the value of the given property name from the sources.
func (r *DefaultPropertyResolver) Lookup(key string) (any, bool) {
	for _, cfgSource := range r.propSources.Slice() {
		if value, ok := cfgSource.Value(key); ok {
			return value, true
		}
	}

	return nil, false
}

// LookupOrDefault returns the value of the given property name from the sources.
// If the property does not exist, it returns the default value.
func (r *DefaultPropertyResolver) LookupOrDefault(name string, defaultValue any) any {
	for _, source := range r.propSources.Slice() {
		if value, ok := source.Value(name); ok {
			return value
		}
	}

	return defaultValue
}

// Expand resolves placeholders in the given text.
// If a placeholder cannot be resolved, it continues to resolve other placeholders.
func (r *DefaultPropertyResolver) Expand(s string) string {
	result, _ := r.expand(s, true)
	return result
}

// ExpandStrict resolves placeholders in the given text.
// If a placeholder cannot be resolved, it returns an error.
func (r *DefaultPropertyResolver) ExpandStrict(s string) (string, error) {
	return r.expand(s, false)
}

// expand resolves placeholders in the given text.
// If continueOnError is true, it continues to resolve placeholders even if a placeholder cannot be resolved.
func (r *DefaultPropertyResolver) expand(s string, continueOnError bool) (string, error) {
	var buf []byte

	i := 0
	for j := 0; j < len(s); j++ {
		if s[j] == '$' && j+1 < len(s) {
			if buf == nil {
				buf = make([]byte, 0, 2*len(s))
			}

			buf = append(buf, s[i:j]...)
			name, w := r.extractPlaceholder(s[j+1:])

			if name == "" && w > 0 {
				if !continueOnError {
					return "", errors.New("wrong placeholder format")
				}

				buf = append(buf, s[j:j+w+1]...)
			} else if name == "" {
				buf = append(buf, s[j])
			} else {
				value, ok := r.Lookup(name)

				if !ok {
					if !continueOnError {
						return "", fmt.Errorf("cannot resolve placeholder '${%s}'", name)
					} else {
						buf = append(buf, s[j:j+w+1]...)
					}
				} else {
					buf = append(buf, fmt.Sprint(value)...)
				}

			}

			j += w
			i = j + 1
		}
	}

	if buf == nil {
		return s, nil
	}

	return string(buf) + s[i:], nil
}

// extractPlaceholder extracts the placeholder from the given string.
func (r *DefaultPropertyResolver) extractPlaceholder(s string) (string, int) {
	switch {
	case s[0] == '{':
		for i := 1; i < len(s); i++ {
			if s[i] == '}' {
				if i == 1 {
					return "", 2
				}
				return s[1:i], i + 1
			}
		}

		return "", 1
	}

	return "", 0
}
