package http

const (
	PathVariablesAttribute = "PathVariables"
)

type PathVariable struct {
	Name  string
	Value string
}

type PathVariables struct {
	variables    []PathVariable
	nameMap      map[string]int
	currentIndex int
}

func (p *PathVariables) clear() {
	p.currentIndex = 0
}

func (p *PathVariables) Put(name string, value string) {
	if p.currentIndex >= len(p.nameMap) {
		p.variables = append(p.variables, PathVariable{
			Name:  name,
			Value: value,
		})
	} else {
		p.variables[p.currentIndex].Name = name
		p.variables[p.currentIndex].Value = value
	}

	p.currentIndex++
}

func (p *PathVariables) Value(name string) (string, bool) {
	if index, ok := p.nameMap[name]; ok {
		if index >= p.currentIndex {
			return "", false
		}

		return p.variables[index].Value, true
	}

	return "", false
}
