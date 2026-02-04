package http

import (
	"errors"
	"sync"
)

const (
	nodeKindStatic uint8 = iota
	nodeKindParam
	nodeKindCatchAll
	maxMethods  = 10
	maxParams   = 16
	maxSegments = 32
)

var (
	ErrDuplicateRoute = errors.New("duplicate route")
	ErrEmptyParamName = errors.New("empty parameter name")
	ErrUnclosedParam  = errors.New("unclosed parameter brace")
	ErrTooManyParams  = errors.New("too many parameters")
)

// Candidate - routing metadata (Endpoint'ten ayrı)
type Candidate struct {
	Endpoint    *Endpoint
	Pattern     string
	ParamNames  [maxParams]string // fixed array - no allocation
	ParamCount  int
	HasCatchAll bool
	CatchAllIdx int // paramNames'deki index
}

// Param - match sonucu
type Param struct {
	Key   string
	Value string
}

// Params - fixed size param array
type Params struct {
	values [maxParams]Param
	count  int
}

func (ps *Params) Get(key string) string {
	for i := 0; i < ps.count; i++ {
		if ps.values[i].Key == key {
			return ps.values[i].Value
		}
	}
	return ""
}

func (ps *Params) add(key, value string) {
	if ps.count < maxParams {
		ps.values[ps.count] = Param{Key: key, Value: value}
		ps.count++
	}
}

func (ps *Params) reset() {
	ps.count = 0
}

func (ps *Params) Count() int {
	return ps.count
}

func (ps *Params) Index(i int) Param {
	return ps.values[i]
}

// MatchResult
type MatchResult struct {
	Endpoint *Endpoint
	Params   Params
	Found    bool
}

// routeEntry - endpoint + candidate
type routeEntry struct {
	endpoint  *Endpoint
	candidate *Candidate
}

// radixNode - compressed radix tree node
type radixNode struct {
	prefix     string // compressed path segment(s)
	kind       uint8
	indices    []byte       // first byte of each child's prefix (sorted)
	children   []*radixNode // static children
	paramChild *radixNode   // {param} child
	catchChild *radixNode   // ** child
	routes     [maxMethods]*routeEntry
}

func (n *radixNode) findChild(b byte) *radixNode {
	for i, c := range n.indices {
		if c == b {
			return n.children[i]
		}
		if c > b {
			break
		}
	}
	return nil
}

func (n *radixNode) addChild(child *radixNode) {
	b := child.prefix[0]

	// Find insert position (keep sorted)
	pos := 0
	for pos < len(n.indices) && n.indices[pos] < b {
		pos++
	}

	// Insert
	n.indices = append(n.indices, 0)
	copy(n.indices[pos+1:], n.indices[pos:])
	n.indices[pos] = b

	n.children = append(n.children, nil)
	copy(n.children[pos+1:], n.children[pos:])
	n.children[pos] = child
}

// matchContext - reusable match state (pooled)
type matchContext struct {
	params      Params
	paramValues [maxParams]string // captured values during traversal
	paramCount  int
	catchAll    string
}

func (mc *matchContext) reset() {
	mc.params.reset()
	mc.paramCount = 0
	mc.catchAll = ""
}

func (mc *matchContext) pushParam(value string) {
	if mc.paramCount < maxParams {
		mc.paramValues[mc.paramCount] = value
		mc.paramCount++
	}
}

func (mc *matchContext) popParam() {
	if mc.paramCount > 0 {
		mc.paramCount--
	}
}

var matchCtxPool = sync.Pool{
	New: func() interface{} {
		return &matchContext{}
	},
}

// RadixTree
type RadixTree struct {
	root *radixNode
}

func NewRadixTree() *RadixTree {
	return &RadixTree{
		root: &radixNode{kind: nodeKindStatic},
	}
}

// AddRoute adds a route to the tree
func (t *RadixTree) AddRoute(path string, method Method, handler interface{}) error {
	if len(path) == 0 || path[0] != '/' {
		path = "/" + path
	}

	// Remove trailing slash (except root)
	if len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	ep := &Endpoint{}
	candidate, err := buildCandidate(path, ep)
	if err != nil {
		return err
	}

	n := t.root
	remaining := path[1:] // skip leading /

	for len(remaining) > 0 {
		// Check for catch-all
		if len(remaining) >= 2 && remaining[0] == '*' && remaining[1] == '*' {
			if n.catchChild == nil {
				n.catchChild = &radixNode{kind: nodeKindCatchAll, prefix: "**"}
			}
			n = n.catchChild
			break
		}

		// Check for param
		if remaining[0] == '{' {
			// Find end of param
			end := 1
			for end < len(remaining) && remaining[end] != '}' {
				end++
			}
			if end < len(remaining) {
				if n.paramChild == nil {
					n.paramChild = &radixNode{kind: nodeKindParam, prefix: "{}"}
				}
				n = n.paramChild
				remaining = remaining[end+1:]
				// Skip slash after param
				if len(remaining) > 0 && remaining[0] == '/' {
					remaining = remaining[1:]
				}
				continue
			}
		}

		// Static segment - find where it ends
		segEnd := 0
		for segEnd < len(remaining) {
			if remaining[segEnd] == '{' || (remaining[segEnd] == '*' && segEnd+1 < len(remaining) && remaining[segEnd+1] == '*') {
				break
			}
			segEnd++
		}

		// Remove trailing slash from static segment
		seg := remaining[:segEnd]
		if len(seg) > 0 && seg[len(seg)-1] == '/' {
			seg = seg[:len(seg)-1]
		}
		remaining = remaining[segEnd:]

		if len(seg) > 0 {
			n = t.insertStatic(n, seg)
		}

		// Skip slash
		if len(remaining) > 0 && remaining[0] == '/' {
			remaining = remaining[1:]
		}
	}

	mi := methodIndex(method)
	if n.routes[mi] != nil {
		return ErrDuplicateRoute
	}
	n.routes[mi] = &routeEntry{endpoint: ep, candidate: candidate}
	return nil
}

func (t *RadixTree) insertStatic(n *radixNode, path string) *radixNode {
	for {
		if len(path) == 0 {
			return n
		}

		child := n.findChild(path[0])
		if child == nil {
			// No matching child, create new node
			newNode := &radixNode{kind: nodeKindStatic, prefix: path}
			n.addChild(newNode)
			return newNode
		}

		// Find common prefix length
		commonLen := 0
		minLen := len(child.prefix)
		if len(path) < minLen {
			minLen = len(path)
		}
		for commonLen < minLen && child.prefix[commonLen] == path[commonLen] {
			commonLen++
		}

		// Need to split?
		if commonLen < len(child.prefix) {
			// Split existing node
			splitNode := &radixNode{
				kind:   nodeKindStatic,
				prefix: child.prefix[:commonLen],
			}
			child.prefix = child.prefix[commonLen:]
			splitNode.addChild(child)

			// Replace child in parent
			for i, c := range n.indices {
				if c == splitNode.prefix[0] {
					n.children[i] = splitNode
					break
				}
			}

			if commonLen == len(path) {
				return splitNode
			}

			// Continue with new segment
			path = path[commonLen:]
			n = splitNode
			continue
		}

		// Full prefix match, continue down
		path = path[commonLen:]
		n = child
	}
}

func buildCandidate(path string, ep *Endpoint) (*Candidate, error) {
	c := &Candidate{
		Endpoint:    ep,
		Pattern:     path,
		CatchAllIdx: -1,
	}

	i := 0
	for i < len(path) {
		if path[i] == '{' {
			// Parameter
			end := i + 1
			for end < len(path) && path[end] != '}' {
				end++
			}
			if end >= len(path) {
				return nil, ErrUnclosedParam
			}
			name := path[i+1 : end]
			if len(name) == 0 {
				return nil, ErrEmptyParamName
			}
			if c.ParamCount >= maxParams {
				return nil, ErrTooManyParams
			}
			c.ParamNames[c.ParamCount] = name
			c.ParamCount++
			i = end + 1
		} else if i+1 < len(path) && path[i] == '*' && path[i+1] == '*' {
			// Catch-all
			c.HasCatchAll = true
			c.CatchAllIdx = c.ParamCount
			c.ParamNames[c.ParamCount] = "**"
			c.ParamCount++
			break
		} else {
			i++
		}
	}

	return c, nil
}

// Match finds a matching route - zero allocation in hot path
func (t *RadixTree) Match(path string, method Method) MatchResult {
	if len(path) == 0 || path[0] != '/' {
		return MatchResult{Found: false}
	}

	// Get pooled context
	ctx := matchCtxPool.Get().(*matchContext)
	ctx.reset()

	// Skip leading slash
	remaining := path[1:]

	node := t.match(t.root, remaining, ctx)

	result := MatchResult{Found: false}

	if node != nil {
		mi := methodIndex(method)
		route := node.routes[mi]
		if route != nil {
			result.Found = true
			result.Endpoint = route.endpoint

			// Build params from candidate + captured values
			for i := 0; i < route.candidate.ParamCount; i++ {
				name := route.candidate.ParamNames[i]
				if name == "**" {
					result.Params.add(name, ctx.catchAll)
				} else if i < ctx.paramCount {
					result.Params.add(name, ctx.paramValues[i])
				}
			}
		}
	}

	// Return context to pool
	matchCtxPool.Put(ctx)

	return result
}

func (t *RadixTree) match(n *radixNode, path string, ctx *matchContext) *radixNode {
	for {
		// Base case: path consumed
		if len(path) == 0 {
			// Check if this node has any routes
			for _, r := range n.routes {
				if r != nil {
					return n
				}
			}
			return nil
		}

		// Try static children first (highest priority)
		if child := n.findChild(path[0]); child != nil {
			// Check if prefix matches
			if len(path) >= len(child.prefix) {
				match := true
				for i := 0; i < len(child.prefix); i++ {
					if path[i] != child.prefix[i] {
						match = false
						break
					}
				}
				if match {
					newPath := path[len(child.prefix):]
					// Skip slash after static segment
					if len(newPath) > 0 && newPath[0] == '/' {
						newPath = newPath[1:]
					}
					if result := t.match(child, newPath, ctx); result != nil {
						return result
					}
				}
			}
		}

		// Extract current segment for param/catchall matching
		segEnd := 0
		for segEnd < len(path) && path[segEnd] != '/' {
			segEnd++
		}
		segment := path[:segEnd]
		rest := path[segEnd:]
		if len(rest) > 0 && rest[0] == '/' {
			rest = rest[1:]
		}

		// Try param child
		if n.paramChild != nil && len(segment) > 0 {
			ctx.pushParam(segment)
			if result := t.match(n.paramChild, rest, ctx); result != nil {
				return result
			}
			ctx.popParam() // backtrack
		}

		// Try catch-all (lowest priority, matches everything)
		if n.catchChild != nil {
			ctx.catchAll = path
			for _, r := range n.catchChild.routes {
				if r != nil {
					return n.catchChild
				}
			}
		}

		return nil
	}
}

func methodIndex(m Method) int {
	if len(m) == 0 {
		return 9
	}
	switch m[0] {
	case 'G':
		return 0 // GET
	case 'P':
		if len(m) > 1 {
			switch m[1] {
			case 'O':
				return 1 // POST
			case 'U':
				return 2 // PUT
			case 'A':
				return 4 // PATCH
			}
		}
	case 'D':
		return 3 // DELETE
	case 'H':
		return 5 // HEAD
	case 'O':
		return 6 // OPTIONS
	case 'C':
		return 7 // CONNECT
	case 'T':
		return 8 // TRACE
	}
	return 9
}
