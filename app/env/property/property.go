package property

import "sync"

type Source interface {
	Name() string
	Source() any
	ContainsProperty(name string) bool
	Property(name string) (any, bool)
	PropertyOrDefault(name string, defaultValue any) any
	PropertyNames() []string
}

type Sources interface {
	Contains(name string) bool
	Find(name string) (Source, bool)
	Size() int
	AddFirst(source Source)
	AddLast(source Source)
	AddAtIndex(index int, source Source)
	Remove(name string) Source
	Replace(name string, source Source)
	PrecendenceOf(source Source) int
	ToSlice() []Source
}

type propertySources struct {
	sources []Source
	mu      sync.RWMutex
}

func NewPropertySources() Sources {
	return &propertySources{
		sources: make([]Source, 0),
		mu:      sync.RWMutex{},
	}
}

func (p *propertySources) Contains(name string) bool {
	defer p.mu.Unlock()
	p.mu.Lock()

	for _, source := range p.sources {
		if source != nil && source.Name() == name {
			return true
		}
	}

	return false
}

func (p *propertySources) Find(name string) (Source, bool) {
	defer p.mu.Unlock()
	p.mu.Lock()

	for _, source := range p.sources {
		if source != nil && source.Name() == name {
			return source, true
		}
	}

	return nil, false
}

func (p *propertySources) AddFirst(source Source) {
	defer p.mu.Unlock()
	p.mu.Lock()

	p.removeIfPresent(source)
	if len(p.sources) == 0 {
		p.sources = append(p.sources, source)
		return
	}

	p.sources = append(p.sources[:1], p.sources[0:]...)
	p.sources[0] = source
}

func (p *propertySources) AddLast(source Source) {
	defer p.mu.Unlock()
	p.mu.Lock()

	p.removeIfPresent(source)
	p.sources = append(p.sources, source)
}

func (p *propertySources) AddAtIndex(index int, source Source) {
	defer p.mu.Unlock()
	p.mu.Lock()

	p.removeIfPresent(source)

	if len(p.sources) == index {
		p.mu.Unlock()
		p.AddLast(source)
		return
	}

	p.sources = append(p.sources[:index+1], p.sources[index:]...)
	p.sources[index] = source
}

func (p *propertySources) Remove(name string) Source {
	source, index := p.findPropertySourceByName(name)

	if index != -1 {
		p.sources = append(p.sources[:index], p.sources[index+1:]...)
	} else {
		return nil
	}

	return source
}

func (p *propertySources) Replace(name string, source Source) {
	_, index := p.findPropertySourceByName(name)

	if index != -1 {
		p.sources[index] = source
	}
}

func (p *propertySources) Size() int {
	return len(p.sources)
}

func (p *propertySources) PrecendenceOf(source Source) int {
	return 0
}

func (p *propertySources) ToSlice() []Source {
	sources := make([]Source, len(p.sources))
	copy(sources, p.sources)
	return sources
}

func (p *propertySources) removeIfPresent(source Source) {
	if source == nil {
		return
	}

	_, index := p.findPropertySourceByName(source.Name())

	if index != -1 {
		p.sources = append(p.sources[:index], p.sources[index+1:]...)
	}
}

func (p *propertySources) findPropertySourceByName(name string) (Source, int) {
	for index, source := range p.sources {

		if source.Name() == name {
			return source, index
		}

	}

	return nil, -1
}
