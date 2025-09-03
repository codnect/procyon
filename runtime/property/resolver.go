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

package property

import (
	"errors"
	"fmt"
)

// Resolver interface provides methods for resolving properties.
type Resolver interface {
	// ContainsProperty checks if the given property name exists.
	ContainsProperty(name string) bool
	// Property returns the value of the given property name.
	Property(name string) (any, bool)
	// PropertyOrDefault returns the value of the given property name.
	// If the property does not exist, it returns the default value.
	PropertyOrDefault(name string, defaultValue any) any
	// ResolvePlaceholders resolves placeholders in the given text.
	// If a placeholder cannot be resolved, it continues to resolve other placeholders.
	ResolvePlaceholders(text string) string
	// ResolveRequiredPlaceholders resolves placeholders in the given text.
	// If a placeholder cannot be resolved, it returns an error.
	ResolveRequiredPlaceholders(text string) (string, error)
}

// MultiSourceResolver is an implementation of the Resolver interface.
// It resolves properties from the given sources.
type MultiSourceResolver struct {
	sources *SourceList
}

// NewMultiSourceResolver creates a new MultiSourceResolver with the given sources.
func NewMultiSourceResolver(sources *SourceList) *MultiSourceResolver {
	if sources == nil {
		panic("nil sources")
	}

	return &MultiSourceResolver{
		sources: sources,
	}
}

// ContainsProperty checks if the given property name exists in the sources.
func (r *MultiSourceResolver) ContainsProperty(name string) bool {
	for _, source := range r.sources.Slice() {
		if source.ContainsProperty(name) {
			return true
		}
	}

	return false
}

// Property returns the value of the given property name from the sources.
func (r *MultiSourceResolver) Property(name string) (any, bool) {
	for _, source := range r.sources.Slice() {
		if value, ok := source.Property(name); ok {
			return value, true
		}
	}

	return nil, false
}

// PropertyOrDefault returns the value of the given property name from the sources.
// If the property does not exist, it returns the default value.
func (r *MultiSourceResolver) PropertyOrDefault(name string, defaultValue any) any {
	for _, source := range r.sources.Slice() {
		if value, ok := source.Property(name); ok {
			return value.(string)
		}
	}

	return defaultValue
}

// ResolvePlaceholders resolves placeholders in the given text.
// If a placeholder cannot be resolved, it continues to resolve other placeholders.
func (r *MultiSourceResolver) ResolvePlaceholders(s string) string {
	result, _ := r.resolveRequiredPlaceHolders(s, true)
	return result
}

// ResolveRequiredPlaceholders resolves placeholders in the given text.
// If a placeholder cannot be resolved, it returns an error.
func (r *MultiSourceResolver) ResolveRequiredPlaceholders(s string) (string, error) {
	return r.resolveRequiredPlaceHolders(s, false)
}

// resolveRequiredPlaceHolders resolves placeholders in the given text.
// If continueOnError is true, it continues to resolve placeholders even if a placeholder cannot be resolved.
func (r *MultiSourceResolver) resolveRequiredPlaceHolders(s string, continueOnError bool) (string, error) {
	var buf []byte

	i := 0
	for j := 0; j < len(s); j++ {
		if s[j] == '$' && j+1 < len(s) {
			if buf == nil {
				buf = make([]byte, 0, 2*len(s))
			}

			buf = append(buf, s[i:j]...)
			name, w := r.getPlaceholderName(s[j+1:])

			if name == "" && w > 0 {
				if !continueOnError {
					return "", errors.New("wrong placeholder format")
				}

				buf = append(buf, s[j:j+w+1]...)
			} else if name == "" {
				buf = append(buf, s[j])
			} else {
				value, ok := r.Property(name)

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

func (r *MultiSourceResolver) getPlaceholderName(s string) (string, int) {
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
