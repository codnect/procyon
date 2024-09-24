package http

import (
	"strings"
)

type routeTree struct {
	children     *routeNode
	staticRoutes map[string]*Route
	routes       []string
}

func (n *routeTree) add(route *Route) {

	path := route.Pattern()

	if !strings.ContainsAny(path, ":*") {
		if n.staticRoutes == nil {
			n.staticRoutes = make(map[string]*Route, 0)
		}

		n.staticRoutes[path] = route
		return
	}

	node := n.children
	index := 0
	processed := 0

	for {
	begin:

		if index == len(path) {
			if (node.typ == variableMapping || index-processed == len(node.path)) && node.route != nil {
				panic("You have already registered the same path : " + string(path))
			}
			node.route = route
			return
		}

		char := path[index]

		if node.typ == variableMapping {

			if char == '/' {
				if char >= node.start && char < node.end {
					tempIndex := node.indices[char-node.start]

					if tempIndex != 0 {
						node = node.children[tempIndex]
						processed = index
						index++
						goto begin
					}
				}

				if len(node.path) == 0 {
					//chain.pathVariableNameMap[path[processed+1:index]] = len(chain.pathVariableNameMap)
					//chain.pathVariableIndexMap[len(chain.pathVariableIndexMap)] = path[processed+1 : index]

					node.handlePathSegment(path[index:], route)
					break
				}

				if node.variableNode != nil {
					node = node.variableNode
					processed = index
					goto begin
				}

				//chain.pathVariableNameMap[path[processed+1:index]] = len(chain.pathVariableNameMap)
				//chain.pathVariableIndexMap[len(chain.pathVariableIndexMap)] = path[processed+1 : index]

				node.handlePathSegment(path[index:], route)
				break
			}
		} else {
			if index == len(path) {
				tempIndex := index - processed
				splitNode := &routeNode{
					path:         node.path[tempIndex:],
					pathLen:      uint(len(node.path[tempIndex:])),
					route:        node.route,
					indices:      node.indices,
					start:        node.start,
					end:          node.end,
					index:        node.index,
					children:     node.children,
					variableNode: node.variableNode,
					wildcardNode: node.wildcardNode,
					hasVariable:  node.hasVariable,
					hasWildcard:  node.hasWildcard,
					typ:          node.typ,
					childNode:    node.childNode,
				}

				node.typ = segmentMapping
				node.path = node.path[:tempIndex]
				node.pathLen = uint(len(node.path[:tempIndex]))
				node.route = nil
				node.variableNode = nil
				node.wildcardNode = nil
				node.hasWildcard = false
				node.hasVariable = false
				node.start = 0
				node.end = 0
				node.index = 0
				node.indices = nil
				node.children = nil
				node.childNode = nil

				node.route = route
				node.addRoutingNode(splitNode)
				break
			}

			if index-processed == len(node.path) {

				if char >= node.start && char < node.end {
					tempIndex := node.indices[char-node.start]

					if tempIndex != 0 {
						node = node.children[tempIndex]
						processed = index
						index++
						goto begin
					}
				}

				if len(node.path) == 0 {
					node.handlePathSegment(path[index:], route)
					break
				}

				if node.variableNode != nil {
					node = node.variableNode
					processed = index
					goto begin
				}

				node.handlePathSegment(path[index:], route)
				break
			}

			tempIndex := index - processed
			if path[index] != node.path[index-processed] {
				splitNode := &routeNode{
					path:         node.path[tempIndex:],
					pathLen:      uint(len(node.path[tempIndex:])),
					route:        node.route,
					indices:      node.indices,
					start:        node.start,
					end:          node.end,
					index:        node.index,
					children:     node.children,
					variableNode: node.variableNode,
					wildcardNode: node.wildcardNode,
					hasVariable:  node.hasVariable,
					hasWildcard:  node.hasWildcard,
					typ:          node.typ,
					childNode:    node.childNode,
				}

				node.typ = segmentMapping
				node.path = node.path[:tempIndex]
				node.pathLen = uint(len(node.path[:tempIndex]))
				node.route = nil
				node.variableNode = nil
				node.wildcardNode = nil
				node.hasWildcard = false
				node.hasVariable = false
				node.start = 0
				node.end = 0
				node.index = 0
				node.indices = nil
				node.children = nil
				node.childNode = nil

				if len(path[index:]) == 0 {
					node.route = route
					node.addRoutingNode(splitNode)
					break
				}

				node.addRoutingNode(splitNode)
				node.handlePathSegment(path[index:], route)
				break
			}
		}
		index++
	}
}

func (n *routeTree) match(ctx Context) *Route {
	var (
		index     uint
		path      = ctx.Request().Path()
		processed uint

		lastWildcardMapping *routeNode
		//lastWildcard        uint
		existLastWildcard bool
		route             *Route

		node       = n.children
		pathLength = uint(len(path))
		//pathVariables      = ctx.Value(http.PathVariablesAttribute).(*http.PathVariables)
		//pathVariablesIndex int
	)

search:
	for {
		if node == nil {
			return route
		}

		if index == pathLength {
			if index-processed == node.pathLen || node.path[node.pathLen-1] == 47 {
				route = node.route
			}
			break
		}

		if index-processed == node.pathLen {
			if node.hasWildcard {
				lastWildcardMapping = node.wildcardNode
				existLastWildcard = true
				//lastWildcard = index
			}

			character := path[index]

			if character >= node.start && character < node.end {
				childIndex := node.indices[character-node.start]

				if childIndex != 0 {
					node = node.children[childIndex]
					processed = index
					index++
					continue search
				}
			}

			if node.hasVariable {
				node = node.variableNode
				processed = index
				index++

				for {
					if index == pathLength {
						//ctx.pathVariables.nameMap = node.handlerChain.pathVariableNameMap
						//pathVariableName := node.handlerChain.pathVariableIndexMap[pathVariablesIndex]

						//pathVariables.Put(pathVariableName, path[processed:index])
						return node.route
					}

					if path[index] == 47 {
						//ctx.pathVariables.nameMap = node.handlerChain.pathVariableNameMap
						//pathVariableName := node.handlerChain.pathVariableIndexMap[pathVariablesIndex]

						//pathVariables.Put(pathVariableName, path[processed:index])

						node = node.childNode
						processed = index
						index++
						continue search
					}

					index++
				}
			}

			if node.hasWildcard {
				//ctx.pathVariables.nameMap = node.handlerChain.pathVariableNameMap
				//pathVariableName := node.handlerChain.pathVariableIndexMap[pathVariablesIndex]

				//pathVariables.Put(pathVariableName, path[index:])
				route = node.wildcardNode.route
			}
			break
		}

		if path[index] != node.path[index-processed] {
			if existLastWildcard {
				//ctx.pathVariables.nameMap = node.handlerChain.pathVariableNameMap
				//pathVariableName := node.handlerChain.pathVariableIndexMap[pathVariablesIndex]

				//pathVariables.Put(pathVariableName, path[lastWildcard:])
				route = lastWildcardMapping.route
			}
			break
		}

		index++

	}

	//ctx.pathVariables.nameMap = node.handlerChain.pathVariableNameMap
	return route
}
