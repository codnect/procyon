package property

import "strings"

type MapPropertySource struct {
	name   string
	source map[string]interface{}
}

func NewMapPropertySource(name string, source map[string]interface{}) *MapPropertySource {
	if strings.TrimSpace(name) == "" {
		panic("env: name cannot be empty or blank")
	}

	if source == nil {
		panic("env: source cannot be nil")
	}

	return &MapPropertySource{
		name:   name,
		source: flatMap(source),
	}
}

func (m *MapPropertySource) Name() string {
	return m.name
}

func (m *MapPropertySource) Source() any {
	return m.source
}

func (m *MapPropertySource) ContainsProperty(name string) bool {
	if _, exists := m.source[name]; exists {
		return true
	}

	return false
}

func (m *MapPropertySource) Property(name string) (any, bool) {
	if value, exists := m.source[name]; exists {
		return value, true
	}

	return nil, false
}

func (m *MapPropertySource) PropertyOrDefault(name string, defaultValue any) any {
	value, exists := m.Property(name)
	if !exists {
		return defaultValue
	}

	return value
}

func (m *MapPropertySource) PropertyNames() []string {
	names := make([]string, 0)

	for name, _ := range m.source {
		names = append(names, name)
	}

	return names
}
