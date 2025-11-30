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

import "strings"

type radixNode struct {
	segment    string
	children   map[string]*radixNode
	paramChild *radixNode
	wildcard   *radixNode
	paramName  string
	handlers   map[string]routeEntry
}

type routeEntry struct {
	handler    Handler
	middleware []Middleware
}

func newRadixNode(segment string) *radixNode {
	return &radixNode{
		segment:  segment,
		children: make(map[string]*radixNode),
		handlers: make(map[string]routeEntry),
	}
}

func (n *radixNode) addRoute(method, pattern string, handler Handler, middleware []Middleware) {
	segments := splitPath(pattern)
	current := n

	for _, segment := range segments {
		switch {
		case strings.HasPrefix(segment, ":"):
			if current.paramChild == nil {
				current.paramChild = newRadixNode(segment)
				current.paramChild.paramName = strings.TrimPrefix(segment, ":")
			}
			current = current.paramChild
		case segment == "**":
			if current.wildcard == nil {
				current.wildcard = newRadixNode(segment)
			}
			current = current.wildcard
		default:
			child, ok := current.children[segment]
			if !ok {
				child = newRadixNode(segment)
				current.children[segment] = child
			}
			current = child
		}
	}

	current.handlers[strings.ToUpper(method)] = routeEntry{handler: handler, middleware: middleware}
}

func (n *radixNode) match(method, path string) (map[string]string, routeEntry, bool) {
	segments := splitPath(path)
	params := make(map[string]string)
	if entry, ok := n.matchSegments(strings.ToUpper(method), segments, params); ok {
		return params, entry, true
	}
	return nil, routeEntry{}, false
}

func (n *radixNode) matchSegments(method string, segments []string, params map[string]string) (routeEntry, bool) {
	if len(segments) == 0 {
		if entry, ok := n.handlers[method]; ok {
			return entry, true
		}
		if n.wildcard != nil {
			if entry, ok := n.wildcard.handlers[method]; ok {
				return entry, true
			}
		}
		return routeEntry{}, false
	}

	segment := segments[0]

	if child, ok := n.children[segment]; ok {
		if entry, ok := child.matchSegments(method, segments[1:], params); ok {
			return entry, true
		}
	}

	if n.paramChild != nil {
		params[n.paramChild.paramName] = segment
		if entry, ok := n.paramChild.matchSegments(method, segments[1:], params); ok {
			return entry, true
		}
		delete(params, n.paramChild.paramName)
	}

	if n.wildcard != nil {
		params["*wildcard"] = strings.Join(segments, "/")
		if entry, ok := n.wildcard.matchSegments(method, nil, params); ok {
			return entry, true
		}
		delete(params, "*wildcard")
	}

	return routeEntry{}, false
}

func splitPath(path string) []string {
	cleaned := strings.Trim(path, "/")
	if cleaned == "" {
		return nil
	}
	return strings.Split(cleaned, "/")
}
