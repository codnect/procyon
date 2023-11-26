package router

import "codnect.io/procyon/web/http"

type mappingNodeType byte

const (
	pathSegmentMapping  mappingNodeType = 0
	pathVariableMapping mappingNodeType = 1
	pathWildcardMapping mappingNodeType = 2
)

type mappingNode struct {
	path                string
	length              uint
	mappingType         mappingNodeType
	handlerChain        *http.HandlerChain
	pathVariableName    string
	childMappings       []*mappingNode
	indices             []byte
	wildCardMapping     *mappingNode
	pathVariableMapping *mappingNode
	childMapping        *mappingNode
	hasWildcard         bool
	hasPathVariableNode bool
	childIndex          byte
	childStartIndex     byte
	childEndIndex       byte
}

func (m *mappingNode) addChildNode(childMapping *mappingNode) {
	character := childMapping.path[0]

	if len(m.childMappings) == 0 {
		m.childMappings = append(m.childMappings, nil)
		m.childStartIndex = character
		m.indices = []byte{0}
	} else {
		var difference byte
		if character < m.childStartIndex {
			difference = m.childStartIndex - character
		} else if character >= m.childEndIndex {
			difference = character - m.childEndIndex + 1
		}

		if character < m.childStartIndex {
			newIndices := make([]byte, difference+byte(len(m.indices)))
			copy(newIndices[difference:], m.indices)
			m.childStartIndex = character
			m.indices = newIndices
		} else if character >= m.childEndIndex {
			newIndices := make([]byte, difference+byte(len(m.indices)))
			copy(newIndices, m.indices)
			m.indices = newIndices
		}

	}

	m.childEndIndex = m.childStartIndex + byte(len(m.indices))
	index := m.indices[character-m.childStartIndex]

	if index == 0 {
		m.indices[character-m.childStartIndex] = byte(len(m.childMappings))
		m.childMappings = append(m.childMappings, childMapping)
	} else {
		m.childMappings[index] = childMapping
	}

	tempIndex := 47 - int(m.childStartIndex)
	if tempIndex >= 0 && len(m.indices) > tempIndex {
		m.childIndex = m.indices[tempIndex]
		m.childMapping = m.childMappings[m.childIndex]
	}
}

func (m *mappingNode) handlePathSegment(path string, chain *http.HandlerChain) {

search:
	for {

		if len(path) == 0 {
			m.handlerChain = chain
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
			if len(m.path) == 0 {
				m.path = path
				m.length = uint(len(path))
				m.handlerChain = chain
				break
			}

			child := &mappingNode{
				path:         path,
				length:       uint(len(path)),
				handlerChain: chain,
			}

			m.addChildNode(child)
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
			child := &mappingNode{
				path:   "*",
				length: 1,
			}

			if path[0] == ':' {
				child.mappingType = pathVariableMapping

				if len(pathVariableName) == 0 {
					panic("Path variable cannot be empty " + string(path))
				}

				//chain.pathVariableNameMap[pathVariableName] = len(chain.pathVariableNameMap)
				//chain.pathVariableIndexMap[len(chain.pathVariableIndexMap)] = pathVariableName

				m.pathVariableMapping = child
				m.hasPathVariableNode = true
				m = child
				path = path[tempIndex:]
				continue search
			}

			child.mappingType = pathWildcardMapping
			child.handlerChain = chain
			m.wildCardMapping = child
			m.hasWildcard = true
			break
		}

		if len(m.path) == 0 {
			m.path = path[:pathVariableIndex]
			m.length = uint(len(path[:pathVariableIndex]))
			path = path[pathVariableIndex:]
			continue search
		}

		child := &mappingNode{
			path:   path[:pathVariableIndex],
			length: uint(len(path[:pathVariableIndex])),
		}

		if child.path[0] == '/' {
			child.handlerChain = m.handlerChain
		}

		m.addChildNode(child)
		m = child
		path = path[pathVariableIndex:]
	}

}
