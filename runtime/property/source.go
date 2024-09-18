package property

import "sync"

// Source interface provides methods for handling property sources.
type Source interface {
	// Name returns the name of the source.
	Name() string
	// Source returns the source.
	Source() any
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

// Sources struct is a collection of property sources.
type Sources struct {
	sources []Source
	mu      sync.RWMutex
}

// NewSources function creates a new Sources.
func NewSources() *Sources {
	return &Sources{
		sources: make([]Source, 0),
		mu:      sync.RWMutex{},
	}
}

// Contains checks if the given source name exists in the sources.
func (s *Sources) Contains(name string) bool {
	defer s.mu.Unlock()
	s.mu.Lock()

	for _, source := range s.sources {
		if source != nil && source.Name() == name {
			return true
		}
	}

	return false
}

// Find returns the source with the given name.
func (s *Sources) Find(name string) (Source, bool) {
	defer s.mu.Unlock()
	s.mu.Lock()

	for _, source := range s.sources {
		if source != nil && source.Name() == name {
			return source, true
		}
	}

	return nil, false
}

// AddFirst adds the source to the beginning of the sources.
func (s *Sources) AddFirst(source Source) {
	defer s.mu.Unlock()
	s.mu.Lock()

	s.removeIfPresent(source)
	if len(s.sources) == 0 {
		s.sources = append(s.sources, source)
		return
	}

	s.sources = append(s.sources[:1], s.sources[0:]...)
	s.sources[0] = source
}

// AddLast adds a source to the end of the sources.
func (s *Sources) AddLast(source Source) {
	defer s.mu.Unlock()
	s.mu.Lock()

	s.removeIfPresent(source)
	s.sources = append(s.sources, source)
}

// AddAtIndex adds the source to the sources at the given index.
func (s *Sources) AddAtIndex(index int, source Source) {
	defer s.mu.Unlock()
	s.mu.Lock()

	s.removeIfPresent(source)

	if len(s.sources) == index {
		s.mu.Unlock()
		s.AddLast(source)
		return
	}

	s.sources = append(s.sources[:index+1], s.sources[index:]...)
	s.sources[index] = source
}

// Remove removes the source with the given name from the sources.
func (s *Sources) Remove(name string) Source {
	source, index := s.findPropertySourceByName(name)

	if index != -1 {
		s.sources = append(s.sources[:index], s.sources[index+1:]...)
	} else {
		return nil
	}

	return source
}

// Replace replaces a source with the given name in the sources with a new source.
func (s *Sources) Replace(name string, source Source) {
	_, index := s.findPropertySourceByName(name)

	if index != -1 {
		s.sources[index] = source
	}
}

// Count returns the number of sources.
func (s *Sources) Count() int {
	return len(s.sources)
}

// PrecedenceOf returns the index of a source in the sources
func (s *Sources) PrecedenceOf(source Source) int {
	if source == nil {
		return -1
	}

	_, index := s.findPropertySourceByName(source.Name())
	return index
}

// ToSlice returns the sources as a slice.
func (s *Sources) ToSlice() []Source {
	sources := make([]Source, len(s.sources))
	copy(sources, s.sources)
	return sources
}

// removeIfPresent removes a source from the sources if it exists.
func (s *Sources) removeIfPresent(source Source) {
	if source == nil {
		return
	}

	_, index := s.findPropertySourceByName(source.Name())

	if index != -1 {
		s.sources = append(s.sources[:index], s.sources[index+1:]...)
	}
}

// findPropertySourceByName finds a source by name in the sources.
func (s *Sources) findPropertySourceByName(name string) (Source, int) {
	for index, source := range s.sources {

		if source.Name() == name {
			return source, index
		}

	}

	return nil, -1
}
