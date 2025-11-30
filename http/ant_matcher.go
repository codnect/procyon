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

package http

import (
	"regexp"
	"strings"
)

type antRoute struct {
	method       string
	pattern      string
	matcher      *regexp.Regexp
	placeholders []string
	handler      Handler
	middleware   []Middleware
}

// newAntRoute builds an ant-style route with placeholders and wildcard support.
func newAntRoute(method, pattern string, handler Handler, middleware []Middleware) (*antRoute, error) {
	regex, placeholders, err := compileAntPattern(pattern)
	if err != nil {
		return nil, err
	}

	return &antRoute{
		method:       strings.ToUpper(method),
		pattern:      pattern,
		matcher:      regex,
		placeholders: placeholders,
		handler:      handler,
		middleware:   middleware,
	}, nil
}

func (r *antRoute) match(method, path string) (map[string]string, bool) {
	if r.method != method {
		return nil, false
	}

	matches := r.matcher.FindStringSubmatch(path)
	if len(matches) == 0 {
		return nil, false
	}

	params := make(map[string]string)
	for i, name := range r.placeholders {
		params[name] = matches[i+1]
	}
	return params, true
}

// compileAntPattern compiles an ant-style pattern into a regular expression.
func compileAntPattern(pattern string) (*regexp.Regexp, []string, error) {
	var builder strings.Builder
	builder.WriteString("^")

	placeholders := make([]string, 0)
	for i := 0; i < len(pattern); i++ {
		ch := pattern[i]

		switch ch {
		case '*':
			if i+1 < len(pattern) && pattern[i+1] == '*' {
				builder.WriteString("(.*)")
				i++
			} else {
				builder.WriteString("([^/]*)")
			}
		case '?':
			builder.WriteString("([^/])")
		case '{':
			end := strings.IndexByte(pattern[i:], '}')
			if end == -1 {
				return nil, nil, ErrInvalidPattern
			}

			placeholder := pattern[i+1 : i+end]
			placeholders = append(placeholders, placeholder)
			builder.WriteString("([^/]+)")
			i += end
		default:
			if strings.ContainsRune(".+()^$|[]\\", rune(ch)) {
				builder.WriteString("\\")
			}
			builder.WriteByte(ch)
		}
	}

	builder.WriteString("$")
	regex, err := regexp.Compile(builder.String())
	if err != nil {
		return nil, nil, err
	}

	return regex, placeholders, nil
}
