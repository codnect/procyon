package http

type routingNodeType byte

const (
	segmentMapping  routingNodeType = 0
	variableMapping routingNodeType = 1
	wildcardMapping routingNodeType = 2
)

type routeNode struct {
	path    string
	pathLen uint
	typ     routingNodeType
	route   *Route

	wildcardNode *routeNode
	hasWildcard  bool

	variableNode *routeNode
	hasVariable  bool

	index   byte
	start   byte
	end     byte
	indices []byte

	childNode *routeNode
	children  []*routeNode
}

func (n *routeNode) addRoutingNode(child *routeNode) {
	character := child.path[0]

	if len(n.children) == 0 {
		n.children = append(n.children, nil)
		n.start = character
		n.indices = []byte{0}
	} else {
		var diff byte
		if character < n.start {
			diff = n.start - character
		} else if character >= n.end {
			diff = character - n.end + 1
		}

		if character < n.start {
			newIndices := make([]byte, diff+byte(len(n.indices)))
			copy(newIndices[diff:], n.indices)
			n.start = character
			n.indices = newIndices
		} else if character >= n.end {
			newIndices := make([]byte, diff+byte(len(n.indices)))
			copy(newIndices, n.indices)
			n.indices = newIndices
		}
	}

	n.end = n.start + byte(len(n.indices))
	index := n.indices[character-n.start]

	if index == 0 {
		n.indices[character-n.start] = byte(len(n.children))
		n.children = append(n.children, child)
	} else {
		n.children[index] = child
	}

	tempIndex := 47 - int(n.start)
	if tempIndex >= 0 && len(n.indices) > tempIndex {
		n.index = n.indices[tempIndex]
		n.childNode = n.children[n.index]
	}
}

func (n *routeNode) handlePathSegment(path string, route *Route) {

search:
	for {

		if len(path) == 0 {
			n.route = route
			break
		}

		pathVariableIndex := -1
		tempIndex := 0
		for tempIndex < len(path) {
			if path[tempIndex] == ':' || path[tempIndex] == '*' {
				pathVariableIndex = tempIndex
				break
			}
			tempIndex++
		}

		if pathVariableIndex == -1 {
			if len(n.path) == 0 {
				n.path = path
				n.pathLen = uint(len(path))
				n.route = route
				break
			}

			child := &routeNode{
				path:    path,
				pathLen: uint(len(path)),
				route:   route,
			}

			n.addRoutingNode(child)
			break
		}

		if pathVariableIndex == 0 {
			tempIndex := 1
			for tempIndex < len(path) {
				if path[tempIndex] == '/' {
					break
				}
				tempIndex++
			}

			pathVariableName := path[1:tempIndex]
			child := &routeNode{
				path:    "*",
				pathLen: 1,
			}

			if path[0] == ':' {
				child.typ = variableMapping

				if len(pathVariableName) == 0 {
					panic("Path variable cannot be empty " + string(path))
				}

				//chain.pathVariableNameMap[pathVariableName] = len(chain.pathVariableNameMap)
				//chain.pathVariableIndexMap[len(chain.pathVariableIndexMap)] = pathVariableName

				n.variableNode = child
				n.hasVariable = true
				n = child
				path = path[tempIndex:]
				continue search
			}

			child.typ = wildcardMapping
			child.route = route
			n.wildcardNode = child
			n.hasWildcard = true
			break
		}

		if len(n.path) == 0 {
			n.path = path[:pathVariableIndex]
			n.pathLen = uint(len(path[:pathVariableIndex]))
			path = path[pathVariableIndex:]
			continue search
		}

		child := &routeNode{
			path:    path[:pathVariableIndex],
			pathLen: uint(len(path[:pathVariableIndex])),
		}

		if child.path[0] == '/' {
			child.route = n.route
		}

		n.addRoutingNode(child)
		n = child
		path = path[pathVariableIndex:]
	}
}
