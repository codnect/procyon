package web

import (
	"bytes"
)

type RouterTree struct {
	methodTrees []*RouterMethodTree
}

func newRouterTree() *RouterTree {
	router := &RouterTree{}
	router.createMethodTree(MethodGet)
	router.createMethodTree(MethodPost)
	router.createMethodTree(MethodPut)
	router.createMethodTree(MethodPatch)
	router.createMethodTree(MethodDelete)
	router.createMethodTree(MethodHead)
	router.createMethodTree(MethodOptions)
	return router
}

func (tree *RouterTree) createMethodTree(method HttpMethod) {
	methodTree := new(RouterMethodTree)
	methodTree.method = method
	tree.methodTrees = append(tree.methodTrees, methodTree)
}

func (tree *RouterTree) findMethodTree(method HttpMethod) *RouterMethodTree {
	for _, methodTree := range tree.methodTrees {
		if methodTree.method == method {
			return methodTree
		}
	}
	methodTree := new(RouterMethodTree)
	methodTree.method = method
	tree.methodTrees = append(tree.methodTrees, methodTree)
	return methodTree
}

func (tree *RouterTree) addRoute(path string, method HttpMethod, handlerChain *HandlerChain) {
	methodNode := tree.findMethodTree(method)

	if methodNode.root == nil {
		methodNode.root = &RouterNode{}
	}

	methodNode.add([]byte(path), handlerChain)
	methodNode.registeredRoutes = append(methodNode.registeredRoutes, path)
}

func (tree *RouterTree) FindHandler(method HttpMethod, path string) (*HandlerChain, []string) {
	if method[0] == 'G' {
		return tree.methodTrees[0].findHandler(path)
	}

	methodNode := tree.findMethodTree(method)
	return methodNode.findHandler(path)
}

type RouterMethodTree struct {
	method           HttpMethod
	root             *RouterNode
	staticRoutes     map[string]*HandlerChain
	registeredRoutes []string
}

func (methodTree *RouterMethodTree) add(path []byte, chain *HandlerChain) {

	if bytes.IndexByte(path, ':') == -1 && bytes.IndexByte(path, '*') == -1 {
		if methodTree.staticRoutes == nil {
			methodTree.staticRoutes = make(map[string]*HandlerChain, 0)
		}
		methodTree.staticRoutes[string(path)] = chain
		return
	}

	node := methodTree.root
	index := 0
	processed := 0

	for {
	begin:

		if index == len(path) {
			if (node.nodeType == PathVariableNode || index-processed == len(node.path)) && node.handlerChain != nil {
				panic("You have already registered the same path : " + string(path))
			}
			node.handlerChain = chain
			return
		}

		char := path[index]

		if node.nodeType == PathVariableNode {

			if char == '/' {
				if char >= node.childStartIndex && char < node.childEndIndex {
					tempIndex := node.indices[char-node.childStartIndex]

					if tempIndex != 0 {
						node = node.childNodes[tempIndex]
						processed = index
						index++
						goto begin
					}
				}

				if len(node.path) == 0 {
					node.handlePathSegment(path[index:], chain)
					break
				}

				if node.pathVariableNode != nil {
					node = node.pathVariableNode
					processed = index
					goto begin
				}

				node.handlePathSegment(path[index:], chain)
				break
			}
		} else {
			if index == len(path) {
				tempIndex := index - processed
				splitNode := &RouterNode{
					path:                node.path[tempIndex:],
					length:              uint(len(node.path[tempIndex:])),
					handlerChain:        node.handlerChain,
					indices:             node.indices,
					childStartIndex:     node.childStartIndex,
					childEndIndex:       node.childEndIndex,
					childIndex:          node.childIndex,
					childNodes:          node.childNodes,
					pathVariableNode:    node.pathVariableNode,
					wildCardNode:        node.wildCardNode,
					hasPathVariableNode: node.hasPathVariableNode,
					hasWildcard:         node.hasWildcard,
					nodeType:            node.nodeType,
					childNode:           node.childNode,
				}

				node.nodeType = PathSegmentNode
				node.path = node.path[:tempIndex]
				node.length = uint(len(node.path[:tempIndex]))
				node.handlerChain = nil
				node.pathVariableNode = nil
				node.wildCardNode = nil
				node.hasWildcard = false
				node.hasPathVariableNode = false
				node.childStartIndex = 0
				node.childEndIndex = 0
				node.childIndex = 0
				node.indices = nil
				node.childNodes = nil
				node.childNode = nil

				node.handlerChain = chain
				node.addChildNode(splitNode)
				break
			}

			if index-processed == len(node.path) {

				if char >= node.childStartIndex && char < node.childEndIndex {
					tempIndex := node.indices[char-node.childStartIndex]

					if tempIndex != 0 {
						node = node.childNodes[tempIndex]
						processed = index
						index++
						goto begin
					}
				}

				if len(node.path) == 0 {
					node.handlePathSegment(path[index:], chain)
					break
				}

				if node.pathVariableNode != nil {
					node = node.pathVariableNode
					processed = index
					goto begin
				}

				node.handlePathSegment(path[index:], chain)
				break
			}

			tempIndex := index - processed
			if path[index] != node.path[index-processed] {
				splitNode := &RouterNode{
					path:                node.path[tempIndex:],
					length:              uint(len(node.path[tempIndex:])),
					handlerChain:        node.handlerChain,
					indices:             node.indices,
					childStartIndex:     node.childStartIndex,
					childEndIndex:       node.childEndIndex,
					childIndex:          node.childIndex,
					childNodes:          node.childNodes,
					pathVariableNode:    node.pathVariableNode,
					wildCardNode:        node.wildCardNode,
					hasPathVariableNode: node.hasPathVariableNode,
					hasWildcard:         node.hasWildcard,
					nodeType:            node.nodeType,
					childNode:           node.childNode,
				}

				node.nodeType = PathSegmentNode
				node.path = node.path[:tempIndex]
				node.length = uint(len(node.path[:tempIndex]))
				node.handlerChain = nil
				node.pathVariableNode = nil
				node.wildCardNode = nil
				node.hasWildcard = false
				node.hasPathVariableNode = false
				node.childStartIndex = 0
				node.childEndIndex = 0
				node.childIndex = 0
				node.indices = nil
				node.childNodes = nil
				node.childNode = nil

				if len(path[index:]) == 0 {
					node.handlerChain = chain
					node.addChildNode(splitNode)
					break
				}

				node.addChildNode(splitNode)
				node.handlePathSegment(path[index:], chain)
				break
			}
		}
		index++
	}
}

func (methodTree *RouterMethodTree) findHandler(path string) (*HandlerChain, []string) {

	if methodTree.staticRoutes != nil {
		if chain, ok := methodTree.staticRoutes[path]; ok {
			return chain, nil
		}
	}

	var handlerChain *HandlerChain

	node := methodTree.root
	pathLength := uint(len(path))
	pathVariables := make([]string, 0)

	var index uint
	var processed uint

	var lastWildcardNode *RouterNode
	var lastWildcard uint
	var existLastWildcard bool

search:
	for {

		if node == nil {
			return handlerChain, pathVariables
		}

		if index == pathLength {
			if index-processed == node.length || node.path[node.length-1] == 47 {
				handlerChain = node.handlerChain
			}
			break
		}

		if index-processed == node.length {
			if node.hasWildcard {
				lastWildcardNode = node.wildCardNode
				existLastWildcard = true
				lastWildcard = index
			}

			character := path[index]

			if character >= node.childStartIndex && character < node.childEndIndex {
				childIndex := node.indices[character-node.childStartIndex]

				if childIndex != 0 {
					node = node.childNodes[childIndex]
					processed = index
					index++
					continue search
				}
			}

			if node.hasPathVariableNode {
				node = node.pathVariableNode
				processed = index
				index++

				for {
					if index == pathLength {
						pathVariables = append(pathVariables, path[processed:index])
						return node.handlerChain, pathVariables
					}

					if path[index] == 47 {
						pathVariables = append(pathVariables, path[processed:index])
						node = node.childNode
						processed = index
						index++
						continue search
					}

					index++
				}
			}

			if node.hasWildcard {
				pathVariables = append(pathVariables, path[index:])
				handlerChain = node.wildCardNode.handlerChain
			}
			break
		}

		if path[index] != node.path[index-processed] {
			if existLastWildcard {
				pathVariables = append(pathVariables, path[lastWildcard:])
				handlerChain = lastWildcardNode.handlerChain
			}
			break
		}

		index++
	}

	return handlerChain, pathVariables
}
