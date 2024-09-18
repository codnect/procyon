package property

import (
	"strconv"
	"strings"
)

// MapSource struct represents a source of properties that are stored in a map.
type MapSource struct {
	name   string
	source map[string]any
}

// NewMapSource function creates a new MapSource with the given name and key-value pair map.
func NewMapSource(name string, source map[string]any) *MapSource {
	if strings.TrimSpace(name) == "" {
		panic("cannot create map source with empty or blank name")
	}

	if source == nil {
		panic("nil source")
	}

	return &MapSource{
		name:   name,
		source: flatMap(source),
	}
}

// Name method returns the name of the source.
func (m *MapSource) Name() string {
	return m.name
}

// Source method returns the source.
func (m *MapSource) Source() any {
	return m.source
}

// ContainsProperty checks if the given property name exists in the source.
func (m *MapSource) ContainsProperty(name string) bool {
	if _, exists := m.source[name]; exists {
		return true
	}

	return false
}

// Property returns the value of the given property name from the source.
// If the property does not exist, it returns false.
func (m *MapSource) Property(name string) (any, bool) {
	if value, exists := m.source[name]; exists {
		return value, true
	}

	return nil, false
}

// PropertyOrDefault returns the value of the given property name from the source.
// If the property does not exist, it returns the default value.
func (m *MapSource) PropertyOrDefault(name string, defaultValue any) any {
	value, exists := m.Property(name)
	if !exists {
		return defaultValue
	}

	return value
}

// PropertyNames returns the property names in the source.
func (m *MapSource) PropertyNames() []string {
	names := make([]string, 0)

	for name, _ := range m.source {
		names = append(names, name)
	}

	return names
}

// flatMap function flattens a map that contains nested maps or slices.
// It returns a new map where each key is a path to a nested property.
func flatMap(m map[string]interface{}) map[string]interface{} {
	flattenMap := map[string]interface{}{}

	for key, value := range m {
		switch child := value.(type) {
		case map[string]interface{}:
			nm := flatMap(child)

			for nk, nv := range nm {
				flattenMap[key+"."+nk] = nv
			}
		case []interface{}:
			for i := 0; i < len(child); i++ {
				flattenMap[key+"."+strconv.Itoa(i)] = child[i]
			}
		default:
			flattenMap[key] = value
		}
	}

	return flattenMap
}
