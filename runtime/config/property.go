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

package config

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
)

// PropertySource interface represents a source of properties.
type PropertySource interface {
	// Name returns the name of the source.
	Name() string
	// Origin returns the origin of the property source.
	Origin() string
	// Value returns the value of the given property name from the source.
	// If the property does not exist, it returns false.
	Value(propertyName string) (any, bool)
	// ValueOrDefault returns the value of the given property name from the source.
	// If the property does not exist, it returns the default value.
	ValueOrDefault(propertyName string, defaultValue any) any
	// PropertyNames returns the property names in the source.
	PropertyNames() []string
}

// PropertySources struct is a collection of property sources.
type PropertySources struct {
	items []PropertySource
	mu    sync.RWMutex
}

// NewPropertySources function creates a new PropertySources.
func NewPropertySources(propertySources ...PropertySource) *PropertySources {
	items := make([]PropertySource, 0)
	items = append(items, propertySources...)
	return &PropertySources{
		items: items,
		mu:    sync.RWMutex{},
	}
}

// Has checks if the given source name exists in the sources.
func (s *PropertySources) Has(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, propertySource := range s.items {
		if propertySource != nil && propertySource.Name() == name {
			return true
		}
	}

	return false
}

// Get returns the source with the given name.
func (s *PropertySources) Get(name string) (PropertySource, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, propertySource := range s.items {
		if propertySource != nil && propertySource.Name() == name {
			return propertySource, true
		}
	}

	return nil, false
}

// PushFront adds the source to the beginning of the sources.
func (s *PropertySources) PushFront(propertySource PropertySource) {
	if propertySource == nil {
		panic("nil property source")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.removeIfPresent(propertySource)
	if len(s.items) == 0 {
		s.items = append(s.items, propertySource)
		return
	}

	s.items = append(s.items[:1], s.items[0:]...)
	s.items[0] = propertySource
}

// PushBack adds a source to the end of the sources.
func (s *PropertySources) PushBack(propertySource PropertySource) {
	if propertySource == nil {
		panic("nil property source")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.removeIfPresent(propertySource)
	s.items = append(s.items, propertySource)
}

// Insert adds the source to the sources at the given index.
func (s *PropertySources) Insert(index int, propertySource PropertySource) {
	if index < 0 {
		panic("negative index")
	}

	if propertySource == nil {
		panic("nil property source")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.removeIfPresent(propertySource)

	if len(s.items) == index {
		s.items = append(s.items, propertySource)
		return
	}

	s.items = append(s.items[:index+1], s.items[index:]...)
	s.items[index] = propertySource
}

// Remove removes the source with the given name from the sources.
func (s *PropertySources) Remove(name string) PropertySource {
	s.mu.Lock()
	defer s.mu.Unlock()

	source, index := s.findByName(name)

	if index != -1 {
		s.items = append(s.items[:index], s.items[index+1:]...)
		return source
	}

	return nil
}

// Replace replaces a source with the given name in the sources with a new source.
func (s *PropertySources) Replace(name string, propertySource PropertySource) {
	if propertySource == nil {
		panic("nil property source")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	_, index := s.findByName(name)

	if index != -1 {
		s.items[index] = propertySource
	}
}

// Len returns the number of sources.
func (s *PropertySources) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.items)
}

// IndexOf returns the index of a source in the sources
func (s *PropertySources) IndexOf(propertySource PropertySource) int {
	if propertySource == nil {
		return -1
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	_, index := s.findByName(propertySource.Name())
	return index
}

// Slice returns the sources as a slice.
func (s *PropertySources) Slice() []PropertySource {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sources := make([]PropertySource, len(s.items))
	copy(sources, s.items)
	return sources
}

// removeIfPresent removes a source from the sources if it exists.
func (s *PropertySources) removeIfPresent(propertySource PropertySource) {
	_, index := s.findByName(propertySource.Name())

	if index != -1 {
		s.items = append(s.items[:index], s.items[index+1:]...)
	}
}

// findPropertySourceByName finds a source by name in the sources.
func (s *PropertySources) findByName(name string) (PropertySource, int) {
	for index, propertySource := range s.items {

		if propertySource.Name() == name {
			return propertySource, index
		}

	}

	return nil, -1
}

// MapPropertySource struct represents a source of properties that are stored in a map.
type MapPropertySource struct {
	name   string
	values map[string]any
}

// NewMapPropertySource function creates a new MapPropertySource with the given name and key-value pair map.
func NewMapPropertySource(name string, values map[string]any) *MapPropertySource {
	if strings.TrimSpace(name) == "" {
		panic("empty or blank name")
	}

	if values == nil {
		panic("nil map")
	}

	result := make(map[string]any)
	flatMap(result, "", values)

	return &MapPropertySource{
		name:   name,
		values: result,
	}
}

// Name method returns the name of the source.
func (m *MapPropertySource) Name() string {
	return m.name
}

// Origin returns the origin of the property source.
func (m *MapPropertySource) Origin() string {
	return m.name
}

// Value returns the value of the given property key from the source.
// If the property does not exist, it returns false.
func (m *MapPropertySource) Value(key string) (any, bool) {
	if value, exists := m.values[key]; exists {
		return value, true
	}

	return nil, false
}

// ValueOrDefault returns the value of the given property key from the source.
// If the property does not exist, it returns the default value.
func (m *MapPropertySource) ValueOrDefault(key string, defaultValue any) any {
	value, exists := m.Value(key)
	if !exists {
		return defaultValue
	}

	return value
}

// PropertyNames returns the property keys in the source.
func (m *MapPropertySource) PropertyNames() []string {
	names := make([]string, 0)

	for name, _ := range m.values {
		names = append(names, name)
	}

	return names
}

// flatMap flattens a nested map into a flat map with dot-separated keys.
func flatMap(dst map[string]any, prefix string, propVal any) {
	if propVal == nil {
		if prefix != "" {
			dst[prefix] = nil
		}

		return
	}

	v := reflect.ValueOf(propVal)

	switch v.Kind() {
	case reflect.Map:
		if v.IsNil() {
			if prefix != "" {
				dst[prefix] = nil
			}
			return
		}

		keys := v.MapKeys()
		strKeys := make([]string, 0, len(keys))
		keyIndex := make(map[string]reflect.Value, len(keys))

		for _, k := range keys {
			ks := ""
			if k.IsValid() {
				if k.Kind() == reflect.String {
					ks = k.String()
				} else {
					ks = fmt.Sprint(k.Interface())
				}
			}
			strKeys = append(strKeys, ks)
			keyIndex[ks] = v.MapIndex(k)
		}

		if len(strKeys) == 0 {
			if prefix != "" {
				dst[prefix] = propVal
			}

			return
		}

		sort.Strings(strKeys)
		for _, ks := range strKeys {
			child := keyIndex[ks]
			if child.IsValid() {
				flatMap(dst, join(prefix, ks), child.Interface())
			} else {
				flatMap(dst, join(prefix, ks), nil)
			}
		}

	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			child := v.Index(i).Interface()
			key := fmt.Sprintf("%s[%d]", prefix, i)
			flatMap(dst, key, child)
		}

	default:
		if prefix != "" {
			dst[prefix] = propVal
		}
	}
}

// join joins the prefix and key with a dot.
// If the prefix is empty, it returns the key.
func join(prefix, key string) string {
	if prefix == "" {
		return key
	}

	return fmt.Sprintf("%s.%s", prefix, key)
}
