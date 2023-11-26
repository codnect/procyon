package router

import (
	"codnect.io/procyon/web/http"
	"strings"
)

type mappingTree struct {
	root             *mappingNode
	staticRoutes     map[string]any
	registeredRoutes []string
}

func (t *mappingTree) addMapping(mapping *Mapping, chain *http.HandlerChain) {

	path := mapping.pattern

	if !strings.ContainsAny(path, ":*") {
		if t.staticRoutes == nil {
			t.staticRoutes = make(map[string]any, 0)
		}

		t.staticRoutes[path] = chain
		return
	}

	node := t.root
	index := 0
	processed := 0

	for {
	begin:

		if index == len(path) {
			if (node.mappingType == pathVariableMapping || index-processed == len(node.path)) && node.handlerChain != nil {
				panic("You have already registered the same path : " + string(path))
			}
			node.handlerChain = chain
			return
		}

		char := path[index]

		if node.mappingType == pathVariableMapping {

			if char == '/' {
				if char >= node.childStartIndex && char < node.childEndIndex {
					tempIndex := node.indices[char-node.childStartIndex]

					if tempIndex != 0 {
						node = node.childMappings[tempIndex]
						processed = index
						index++
						goto begin
					}
				}

				if len(node.path) == 0 {
					//chain.pathVariableNameMap[path[processed+1:index]] = len(chain.pathVariableNameMap)
					//chain.pathVariableIndexMap[len(chain.pathVariableIndexMap)] = path[processed+1 : index]

					node.handlePathSegment(path[index:], chain)
					break
				}

				if node.pathVariableMapping != nil {
					node = node.pathVariableMapping
					processed = index
					goto begin
				}

				//chain.pathVariableNameMap[path[processed+1:index]] = len(chain.pathVariableNameMap)
				//chain.pathVariableIndexMap[len(chain.pathVariableIndexMap)] = path[processed+1 : index]

				node.handlePathSegment(path[index:], chain)
				break
			}
		} else {
			if index == len(path) {
				tempIndex := index - processed
				splitNode := &mappingNode{
					path:                node.path[tempIndex:],
					length:              uint(len(node.path[tempIndex:])),
					handlerChain:        node.handlerChain,
					indices:             node.indices,
					childStartIndex:     node.childStartIndex,
					childEndIndex:       node.childEndIndex,
					childIndex:          node.childIndex,
					childMappings:       node.childMappings,
					pathVariableMapping: node.pathVariableMapping,
					wildCardMapping:     node.wildCardMapping,
					hasPathVariableNode: node.hasPathVariableNode,
					hasWildcard:         node.hasWildcard,
					mappingType:         node.mappingType,
					childMapping:        node.childMapping,
				}

				node.mappingType = pathSegmentMapping
				node.path = node.path[:tempIndex]
				node.length = uint(len(node.path[:tempIndex]))
				node.handlerChain = nil
				node.pathVariableMapping = nil
				node.wildCardMapping = nil
				node.hasWildcard = false
				node.hasPathVariableNode = false
				node.childStartIndex = 0
				node.childEndIndex = 0
				node.childIndex = 0
				node.indices = nil
				node.childMappings = nil
				node.childMapping = nil

				node.handlerChain = chain
				node.addChildNode(splitNode)
				break
			}

			if index-processed == len(node.path) {

				if char >= node.childStartIndex && char < node.childEndIndex {
					tempIndex := node.indices[char-node.childStartIndex]

					if tempIndex != 0 {
						node = node.childMappings[tempIndex]
						processed = index
						index++
						goto begin
					}
				}

				if len(node.path) == 0 {
					node.handlePathSegment(path[index:], chain)
					break
				}

				if node.pathVariableMapping != nil {
					node = node.pathVariableMapping
					processed = index
					goto begin
				}

				node.handlePathSegment(path[index:], chain)
				break
			}

			tempIndex := index - processed
			if path[index] != node.path[index-processed] {
				splitNode := &mappingNode{
					path:                node.path[tempIndex:],
					length:              uint(len(node.path[tempIndex:])),
					handlerChain:        node.handlerChain,
					indices:             node.indices,
					childStartIndex:     node.childStartIndex,
					childEndIndex:       node.childEndIndex,
					childIndex:          node.childIndex,
					childMappings:       node.childMappings,
					pathVariableMapping: node.pathVariableMapping,
					wildCardMapping:     node.wildCardMapping,
					hasPathVariableNode: node.hasPathVariableNode,
					hasWildcard:         node.hasWildcard,
					mappingType:         node.mappingType,
					childMapping:        node.childMapping,
				}

				node.mappingType = pathSegmentMapping
				node.path = node.path[:tempIndex]
				node.length = uint(len(node.path[:tempIndex]))
				node.handlerChain = nil
				node.pathVariableMapping = nil
				node.wildCardMapping = nil
				node.hasWildcard = false
				node.hasPathVariableNode = false
				node.childStartIndex = 0
				node.childEndIndex = 0
				node.childIndex = 0
				node.indices = nil
				node.childMappings = nil
				node.childMapping = nil

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

func (t *mappingTree) getHandlerChain(ctx http.Context) *http.HandlerChain {
	var (
		index     uint
		path      = ctx.Request().Path()
		processed uint

		lastWildcardMapping *mappingNode
		//lastWildcard        uint
		existLastWildcard bool
		handlerChain      *http.HandlerChain

		node       = t.root
		pathLength = uint(len(path))
		//pathVariables      = ctx.Value(http.PathVariablesAttribute).(*http.PathVariables)
		//pathVariablesIndex int
	)

search:
	for {
		if node == nil {
			return handlerChain
		}

		if index == pathLength {
			if index-processed == node.length || node.path[node.length-1] == 47 {
				handlerChain = node.handlerChain
			}
			break
		}

		if index-processed == node.length {
			if node.hasWildcard {
				lastWildcardMapping = node.wildCardMapping
				existLastWildcard = true
				//lastWildcard = index
			}

			character := path[index]

			if character >= node.childStartIndex && character < node.childEndIndex {
				childIndex := node.indices[character-node.childStartIndex]

				if childIndex != 0 {
					node = node.childMappings[childIndex]
					processed = index
					index++
					continue search
				}
			}

			if node.hasPathVariableNode {
				node = node.pathVariableMapping
				processed = index
				index++

				for {
					if index == pathLength {
						//ctx.pathVariables.nameMap = node.handlerChain.pathVariableNameMap
						//pathVariableName := node.handlerChain.pathVariableIndexMap[pathVariablesIndex]

						//pathVariables.Put(pathVariableName, path[processed:index])
						return node.handlerChain
					}

					if path[index] == 47 {
						//ctx.pathVariables.nameMap = node.handlerChain.pathVariableNameMap
						//pathVariableName := node.handlerChain.pathVariableIndexMap[pathVariablesIndex]

						//pathVariables.Put(pathVariableName, path[processed:index])

						node = node.childMapping
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
				handlerChain = node.wildCardMapping.handlerChain
			}
			break
		}

		if path[index] != node.path[index-processed] {
			if existLastWildcard {
				//ctx.pathVariables.nameMap = node.handlerChain.pathVariableNameMap
				//pathVariableName := node.handlerChain.pathVariableIndexMap[pathVariablesIndex]

				//pathVariables.Put(pathVariableName, path[lastWildcard:])
				handlerChain = lastWildcardMapping.handlerChain
			}
			break
		}

		index++

	}

	//ctx.pathVariables.nameMap = node.handlerChain.pathVariableNameMap
	return handlerChain
}
