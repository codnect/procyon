// Copyright 2025 Codnect
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

package property

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// MapSource struct represents a source of properties that are stored in a map.
type MapSource struct {
	name   string
	values map[string]any
}

// NewMapSource function creates a new MapSource with the given name and key-value pair map.
func NewMapSource(name string, values map[string]any) *MapSource {
	if strings.TrimSpace(name) == "" {
		panic("empty or blank name")
	}

	if values == nil {
		panic("nil map")
	}

	result := make(map[string]any)
	flatMap(result, "", values)

	return &MapSource{
		name:   name,
		values: result,
	}
}

// Name method returns the name of the source.
func (m *MapSource) Name() string {
	return m.name
}

// Underlying returns the underlying source object.
func (m *MapSource) Underlying() any {
	return m.values
}

// ContainsProperty checks if the given property name exists in the source.
func (m *MapSource) ContainsProperty(name string) bool {
	if _, exists := m.values[name]; exists {
		return true
	}

	return false
}

// Property returns the value of the given property name from the source.
// If the property does not exist, it returns false.
func (m *MapSource) Property(name string) (any, bool) {
	if value, exists := m.values[name]; exists {
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

	for name, _ := range m.values {
		names = append(names, name)
	}

	return names
}

func flatMap(dst map[string]any, prefix string, propVal any) {
	switch t := propVal.(type) {
	case map[string]any:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			child := t[k]
			key := join(prefix, k)
			flatMap(dst, key, child)
		}
	case []any:
		for i, child := range t {
			key := join(prefix, strconv.Itoa(i))
			flatMap(dst, key, child)
		}
	case nil:
		if prefix != "" {
			dst[prefix] = ""
		}
	default:
		if prefix != "" {
			dst[prefix] = fmt.Sprint(t)
		}
	}
}

func join(prefix, key string) string {
	if prefix == "" {
		return key
	}

	return fmt.Sprintf("%s.%s", prefix, key)
}
