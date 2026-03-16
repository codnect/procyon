// Copyright 2026 Codnect
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

import "fmt"

// nodeKind represents the type of node stored in the radix tree.
// Each kind corresponds to a different path matching strategy.
type nodeKind uint

const (
	// nodeKindStatic represents a static path segment.
	nodeKindStatic nodeKind = iota

	// nodeKindParam represents a named parameter segment (e.g. /users/{id}).
	nodeKindParam

	// nodeKindWildcard represents a single-segment wildcard (*).
	nodeKindWildcard

	// nodeKindDoubleWildcard represents a multi-segment wildcard (**).
	nodeKindDoubleWildcard

	// nodeKindPattern represents a segment containing pattern characters (* or ?),
	// such as "file-*.json".
	nodeKindPattern

	// maxMethods defines the number of supported HTTP method slots per node.
	maxMethods = 10

	// maxParams defines the maximum number of path parameters per route.
	maxParams = 16
)

// routeEntry stores metadata for a registered route pattern.
type routeEntry struct {
	endpoint *Endpoint // handler endpoint

	pattern string // original route pattern

	// parameter metadata extracted from the pattern
	paramNames [maxParams]string
	paramCount int
}

// radixNode represents a node in the radix tree used for routing.
type radixNode struct {
	prefix string   // compressed path prefix
	kind   nodeKind // node type

	// sorted first-byte index for static children
	indices  []byte
	children []*radixNode

	// special child nodes for dynamic segments
	paramChild          *radixNode
	doubleWildcardChild *radixNode
	wildcardChild       *radixNode
	patternChildren     []*radixNode

	// method-specific route entries
	routes [maxMethods]*routeEntry
}

// addChild inserts a static child node while keeping indices sorted.
// This allows fast lookup using the first byte of the segment.
func (n *radixNode) addChild(child *radixNode) {

	// first byte used for indexing
	b := child.prefix[0]

	pos := 0
	for pos < len(n.indices) && n.indices[pos] < b {
		pos++
	}

	// insert index while preserving order
	n.indices = append(n.indices, 0)
	copy(n.indices[pos+1:], n.indices[pos:])
	n.indices[pos] = b

	// insert child in same position
	n.children = append(n.children, nil)
	copy(n.children[pos+1:], n.children[pos:])
	n.children[pos] = child
}

// findChild finds a static child node using the first byte of the segment.
func (n *radixNode) findChild(b byte) (*radixNode, bool) {
	for i, c := range n.indices {
		if c == b {
			return n.children[i], true
		}
		if c > b {
			break
		}
	}
	return nil, false
}

// RequestEndpointMatcher is a router implementation based on a radix tree.
type RequestEndpointMatcher struct {
	root *radixNode
}

// NewRequestEndpointMatcher creates a new empty radix-tree router.
func NewRequestEndpointMatcher(endpointDataSource EndpointDataSource) *RequestEndpointMatcher {
	matcher := &RequestEndpointMatcher{root: &radixNode{}}

	if endpointDataSource == nil {
		return matcher
	}

	for _, endpoint := range endpointDataSource.Endpoints() {
		if err := matcher.addEndpoint(endpoint); err != nil {
			panic(fmt.Sprintf("failed to add endpoint %s %s: %v", endpoint.method, endpoint.path, err))
		}
	}

	return matcher
}

// insertStatic inserts a static path fragment into the radix tree.
// The function performs prefix compression and splits nodes when necessary.
func (t *RequestEndpointMatcher) insertStatic(n *radixNode, path string) *radixNode {

	for {
		if len(path) == 0 {
			return n
		}

		child, exists := n.findChild(path[0])
		if !exists {
			newNode := &radixNode{kind: nodeKindStatic, prefix: path}
			n.addChild(newNode)
			return newNode
		}

		// determine longest common prefix
		commonLen := 0
		minLen := len(child.prefix)
		if len(path) < minLen {
			minLen = len(path)
		}

		for commonLen < minLen && child.prefix[commonLen] == path[commonLen] {
			commonLen++
		}

		// split node if prefixes diverge
		if commonLen < len(child.prefix) {

			splitNode := &radixNode{
				kind:   nodeKindStatic,
				prefix: child.prefix[:commonLen],
			}

			child.prefix = child.prefix[commonLen:]
			splitNode.addChild(child)

			// replace original child reference
			for i, c := range n.indices {
				if c == splitNode.prefix[0] {
					n.children[i] = splitNode
					break
				}
			}

			if commonLen == len(path) {
				return splitNode
			}

			path = path[commonLen:]
			n = splitNode
			continue
		}

		path = path[commonLen:]
		n = child
	}
}

// addEndpoint registers a new endpoint into the radix tree.
func (t *RequestEndpointMatcher) addEndpoint(endpoint *Endpoint) error {
	if methodIndex(endpoint.method) < 0 {
		return fmt.Errorf("unsupported HTTP method: %s", endpoint.method)
	}

	// normalize path
	pattern := endpoint.path
	if len(pattern) == 0 || pattern[0] != '/' {
		pattern = "/" + pattern
	}
	if len(pattern) > 1 && pattern[len(pattern)-1] == '/' {
		pattern = pattern[:len(pattern)-1]
	}

	// build metadata for route
	entry, err := buildRouteEntry(pattern, endpoint)
	if err != nil {
		return err
	}

	n := t.root
	p := pattern[1:]

	pos := 0
	staticStart := 0

	// parse path segment-by-segment
	for pos <= len(p) {

		segStart := pos
		segEnd := segStart
		for segEnd < len(p) && p[segEnd] != '/' {
			segEnd++
		}

		seg := p[segStart:segEnd]

		// detect segment type
		isDoubleWildcard := seg == "**"
		isWildcard := seg == "*"
		isParam := isParamSeg(seg)
		isPattern := !isWildcard && !isDoubleWildcard && hasPatternChars(seg)

		// handle dynamic segment types
		if isDoubleWildcard || isWildcard || isParam || isPattern {

			if staticStart < segStart {
				n = t.insertStatic(n, p[staticStart:segStart])
			}

			switch {

			case isDoubleWildcard:
				if n.doubleWildcardChild == nil {
					n.doubleWildcardChild = &radixNode{
						kind:   nodeKindDoubleWildcard,
						prefix: "**",
					}
				}
				n = n.doubleWildcardChild

			case isWildcard:
				if n.wildcardChild == nil {
					n.wildcardChild = &radixNode{
						kind:   nodeKindWildcard,
						prefix: "*",
					}
				}
				n = n.wildcardChild

			case isParam:
				if n.paramChild == nil {
					n.paramChild = &radixNode{
						kind:   nodeKindParam,
						prefix: "{}",
					}
				}
				n = n.paramChild

			default:
				var child *radixNode
				for _, pc := range n.patternChildren {
					if pc.prefix == seg {
						child = pc
						break
					}
				}
				if child == nil {
					child = &radixNode{
						kind:   nodeKindPattern,
						prefix: seg,
					}
					n.patternChildren = append(n.patternChildren, child)
				}
				n = child
			}

			pos = segEnd
			if pos < len(p) && p[pos] == '/' {
				pos++
			}

			staticStart = pos
			continue
		}

		if segEnd >= len(p) {
			break
		}

		pos = segEnd + 1
	}

	// insert remaining static tail
	if staticStart < len(p) {
		n = t.insertStatic(n, p[staticStart:])
	}

	mi := methodIndex(endpoint.method)

	if n.routes[mi] != nil {
		return fmt.Errorf(
			"route already exists for path %s and method %s",
			pattern,
			endpoint.method,
		)
	}

	n.routes[mi] = entry

	return nil
}

// match recursively matches the request path against the radix tree.
func (t *RequestEndpointMatcher) match(n *radixNode, path string, ctx *Context, mi int) *radixNode {
	request := ctx.Request()

	for {

		// if path fully consumed, check for endpoint
		if len(path) == 0 {
			if n.routes[mi] != nil {
				return n
			}

			if n.doubleWildcardChild != nil &&
				n.doubleWildcardChild.routes[mi] != nil {
				return n.doubleWildcardChild
			}

			// static child whose remaining prefix is just "/" may have a ** beneath it
			for i, c := range n.indices {
				child := n.children[i]
				_ = c
				if len(child.prefix) == 1 && child.prefix[0] == '/' {
					if child.doubleWildcardChild != nil &&
						child.doubleWildcardChild.routes[mi] != nil {
						return child.doubleWildcardChild
					}
				}
			}

			return nil
		}

		// 1) Static segment match
		if child, ok := n.findChild(path[0]); ok {

			prefixLen := len(child.prefix)

			if len(path) >= prefixLen {
				match := true

				for i := 0; i < prefixLen; i++ {
					if path[i] != child.prefix[i] {
						match = false
						break
					}
				}

				if match {
					if result := t.match(child, path[prefixLen:], ctx, mi); result != nil {
						return result
					}
				}
			}

			// partial match: path is shorter than child prefix
			// if the unmatched suffix is just "/" and child has **, try it
			if len(path) < prefixLen && len(path) > 0 {
				match := true
				for i := 0; i < len(path); i++ {
					if path[i] != child.prefix[i] {
						match = false
						break
					}
				}
				if match {
					remaining := child.prefix[len(path):]
					if remaining == "/" {
						if child.doubleWildcardChild != nil &&
							child.doubleWildcardChild.routes[mi] != nil {
							return child.doubleWildcardChild
						}
					}
				}
			}
		}

		// split next segment
		segEnd := 0
		for segEnd < len(path) && path[segEnd] != '/' {
			segEnd++
		}

		segment := path[:segEnd]

		var rest string
		if segEnd < len(path) && path[segEnd] == '/' {
			rest = path[segEnd+1:]
		}

		// 2) Pattern match (e.g. file-*.json)
		if len(n.patternChildren) > 0 && len(segment) > 0 {
			for _, pc := range n.patternChildren {
				if matchPattern(pc.prefix, segment) {
					if result := t.match(pc, rest, ctx, mi); result != nil {
						return result
					}
				}
			}
		}

		// 3) Parameter segment
		if n.paramChild != nil && len(segment) > 0 {

			savedLen := request.pathValues.len()

			if request.pathValues.pushRaw(segment) {
				if result := t.match(n.paramChild, rest, ctx, mi); result != nil {
					return result
				}
			}

			request.pathValues.setLen(savedLen)
		}

		// 4) Single wildcard (*)
		if n.wildcardChild != nil && len(segment) > 0 {
			if result := t.match(n.wildcardChild, rest, ctx, mi); result != nil {
				return result
			}
		}

		// 5) Double wildcard (**)
		if n.doubleWildcardChild != nil {

			// zero segment match
			if result := t.match(n.doubleWildcardChild, path, ctx, mi); result != nil {
				return result
			}

			// consume segments progressively
			p := path

			for len(p) > 0 {

				nextSeg := 0
				for nextSeg < len(p) && p[nextSeg] != '/' {
					nextSeg++
				}

				if nextSeg < len(p) {
					p = p[nextSeg+1:]
				} else {
					p = ""
				}

				if result := t.match(n.doubleWildcardChild, p, ctx, mi); result != nil {
					return result
				}
			}
		}

		return nil
	}
}

// Match resolves the incoming request to a registered endpoint.
func (t *RequestEndpointMatcher) Match(ctx *Context) (*Endpoint, bool) {
	request := ctx.Request()
	path := request.Path()

	// normalize path
	if len(path) == 0 || path[0] != '/' {
		return nil, false
	}

	if len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	mi := methodIndex(request.Method())

	request.pathValues.reset()

	node := t.match(t.root, path[1:], ctx, mi)

	if node == nil {
		return nil, false
	}

	route := node.routes[mi]

	// assign parameter names
	limit := route.paramCount

	for i := 0; i < limit; i++ {
		request.pathValues.setName(i, route.paramNames[i])
	}

	return route.endpoint, true
}

// buildRouteEntry extracts metadata (parameters and wildcards)
// from a route pattern.
func buildRouteEntry(path string, endpoint *Endpoint) (*routeEntry, error) {
	entry := &routeEntry{
		pattern:  path,
		endpoint: endpoint,
	}

	i := 0

	for i < len(path) {

		switch {

		// parameter segment {id}
		case path[i] == '{':
			end := i + 1
			for end < len(path) && path[end] != '}' {
				end++
			}

			if end >= len(path) {
				return nil, fmt.Errorf("unclosed parameter in path: %s", path)
			}

			name := path[i+1 : end]

			if len(name) == 0 {
				return nil, fmt.Errorf("empty parameter name in path: %s", path)
			}

			if entry.paramCount >= maxParams {
				return nil, fmt.Errorf("too many parameters in path: %s", path)
			}

			entry.paramNames[entry.paramCount] = name
			entry.paramCount++

			i = end + 1

		// double wildcard **
		case i+1 < len(path) && path[i] == '*' && path[i+1] == '*':
			i += 2
			if i < len(path) && path[i] == '/' {
				i++
			}

		// single wildcard *
		case path[i] == '*':
			i++

		default:
			i++
		}
	}

	return entry, nil
}

// matchPattern performs glob-style pattern matching
// supporting '*' and '?' characters.
func matchPattern(pat, s string) bool {

	pi, si := 0, 0
	starPi, starSi := -1, -1

	for si < len(s) {

		if pi < len(pat) {

			switch pat[pi] {

			case '?':
				pi++
				si++
				continue

			case '*':
				starPi = pi
				starSi = si
				pi++
				continue

			default:
				if pat[pi] == s[si] {
					pi++
					si++
					continue
				}
			}
		}

		if starPi != -1 {
			pi = starPi + 1
			starSi++
			si = starSi
			continue
		}

		return false
	}

	for pi < len(pat) && pat[pi] == '*' {
		pi++
	}

	return pi == len(pat)
}

// isParamSeg checks whether a segment is a path parameter.
func isParamSeg(seg string) bool {
	return len(seg) >= 3 && seg[0] == '{' && seg[len(seg)-1] == '}'
}

// hasPatternChars checks whether a segment contains pattern characters.
func hasPatternChars(seg string) bool {
	for i := 0; i < len(seg); i++ {
		if seg[i] == '*' || seg[i] == '?' {
			return true
		}
	}
	return false
}

// methodIndex maps HTTP method strings to integer indices for route storage.
func methodIndex(m Method) int {
	if len(m) == 0 {
		return -1
	}

	switch m[0] {
	case 'G':
		return 0
	case 'P':
		if len(m) > 1 {
			switch m[1] {
			case 'O':
				return 1
			case 'U':
				return 2
			case 'A':
				return 4
			}
		}
	case 'D':
		return 3
	case 'H':
		return 5
	case 'O':
		return 6
	case 'C':
		return 7
	case 'T':
		return 8
	}

	return -1
}
