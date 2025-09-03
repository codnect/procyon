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

import "sync"

// Source interface provides methods for handling property sources.
type Source interface {
	// Name returns the name of the source.
	Name() string
	// Underlying returns the underlying source object.
	Underlying() any
	// ContainsProperty checks if the given property name exists in the source.
	ContainsProperty(name string) bool
	// Property returns the value of the given property name from the source.
	// If the property does not exist, it returns false.
	Property(name string) (any, bool)
	// PropertyOrDefault returns the value of the given property name from the source.
	// If the property does not exist, it returns the default value.
	PropertyOrDefault(name string, defaultValue any) any
	// PropertyNames returns the property names in the source.
	PropertyNames() []string
}

// SourceList struct is a collection of property sources.
type SourceList struct {
	sources []Source
	mu      sync.RWMutex
}

// SourcesAsList function creates a new SourceList.
func SourcesAsList(sources ...Source) *SourceList {
	sourceSlice := make([]Source, 0)
	sourceSlice = append(sourceSlice, sources...)
	return &SourceList{
		sources: sourceSlice,
		mu:      sync.RWMutex{},
	}
}

// Contains checks if the given source name exists in the sources.
func (s *SourceList) Contains(name string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, source := range s.sources {
		if source != nil && source.Name() == name {
			return true
		}
	}

	return false
}

// Find returns the source with the given name.
func (s *SourceList) Find(name string) (Source, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, source := range s.sources {
		if source != nil && source.Name() == name {
			return source, true
		}
	}

	return nil, false
}

// AddFirst adds the source to the beginning of the sources.
func (s *SourceList) AddFirst(source Source) {
	if source == nil {
		panic("nil source")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.removeIfPresent(source)
	if len(s.sources) == 0 {
		s.sources = append(s.sources, source)
		return
	}

	s.sources = append(s.sources[:1], s.sources[0:]...)
	s.sources[0] = source
}

// AddLast adds a source to the end of the sources.
func (s *SourceList) AddLast(source Source) {
	if source == nil {
		panic("nil source")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.removeIfPresent(source)
	s.sources = append(s.sources, source)
}

// AddAtIndex adds the source to the sources at the given index.
func (s *SourceList) AddAtIndex(index int, source Source) {
	if index < 0 {
		panic("negative index")
	}

	if source == nil {
		panic("nil source")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.removeIfPresent(source)

	if len(s.sources) == index {
		s.sources = append(s.sources, source)
		return
	}

	s.sources = append(s.sources[:index+1], s.sources[index:]...)
	s.sources[index] = source
}

// Remove removes the source with the given name from the sources.
func (s *SourceList) Remove(name string) Source {
	s.mu.Lock()
	defer s.mu.Unlock()

	source, index := s.findPropertySourceByName(name)

	if index != -1 {
		s.sources = append(s.sources[:index], s.sources[index+1:]...)
		return source
	}

	return nil
}

// Replace replaces a source with the given name in the sources with a new source.
func (s *SourceList) Replace(name string, source Source) {
	if source == nil {
		panic("nil source")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	_, index := s.findPropertySourceByName(name)

	if index != -1 {
		s.sources[index] = source
	}
}

// Count returns the number of sources.
func (s *SourceList) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.sources)
}

// PrecedenceOf returns the index of a source in the sources
func (s *SourceList) PrecedenceOf(source Source) int {
	if source == nil {
		return -1
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	_, index := s.findPropertySourceByName(source.Name())
	return index
}

// Slice returns the sources as a slice.
func (s *SourceList) Slice() []Source {
	s.mu.Lock()
	defer s.mu.Unlock()

	sources := make([]Source, len(s.sources))
	copy(sources, s.sources)
	return sources
}

// removeIfPresent removes a source from the sources if it exists.
func (s *SourceList) removeIfPresent(source Source) {
	_, index := s.findPropertySourceByName(source.Name())

	if index != -1 {
		s.sources = append(s.sources[:index], s.sources[index+1:]...)
	}
}

// findPropertySourceByName finds a source by name in the sources.
func (s *SourceList) findPropertySourceByName(name string) (Source, int) {
	for index, source := range s.sources {

		if source.Name() == name {
			return source, index
		}

	}

	return nil, -1
}
