package env

import "errors"

type PropertySource interface {
	GetName() string
	GetSource() interface{}
	GetProperty(name string) interface{}
	ContainsProperty(name string) bool
}

type AbstractPropertySource struct {
	PropertySource
	name   string
	source interface{}
}

func NewAbstractPropertySourceWithSource(name string, source interface{}) AbstractPropertySource {
	propertySource := AbstractPropertySource{
		name:   name,
		source: source,
	}
	return propertySource
}

func (source AbstractPropertySource) GetName() string {
	return source.name
}

func (source AbstractPropertySource) GetSource() interface{} {
	return source.source
}

type EnumerablePropertySource interface {
	GetPropertyNames() []string
}

type AbstractEnumerablePropertySource struct {
	EnumerablePropertySource
	AbstractPropertySource
}

func NewAbstractEnumerablePropertySourceWithSource(name string, source interface{}) AbstractEnumerablePropertySource {
	propertySource := AbstractEnumerablePropertySource{
		AbstractPropertySource: NewAbstractPropertySourceWithSource(name, source),
	}
	propertySource.PropertySource = propertySource
	return propertySource
}

func (source AbstractEnumerablePropertySource) ContainsProperty(name string) bool {
	for _, propertyName := range source.GetPropertyNames() {
		if propertyName == name {
			return true
		}
	}
	return false
}

type PropertySources struct {
	sources []PropertySource
}

func NewPropertySources() PropertySources {
	return PropertySources{
		sources: make([]PropertySource, 0),
	}
}

func (o PropertySources) Get(name string) (PropertySource, error) {
	for _, source := range o.sources {
		if source.GetName() == name {
			return source, nil
		}
	}
	return nil, errors.New("Property not found : " + name)
}

func (o PropertySources) Add(propertySource PropertySource) {
	o.RemoveIfPresent(propertySource)
	o.sources = append(o.sources, propertySource)
}

func (o PropertySources) Remove(name string) PropertySource {
	source, index := o.findPropertySourceByName(name)
	if index != -1 {
		o.sources = append(o.sources[:index], o.sources[index+1:]...)
	}
	return source
}

func (o PropertySources) Replace(name string, propertySource PropertySource) {
	_, index := o.findPropertySourceByName(name)
	if index != -1 {
		o.sources[index] = propertySource
	}
}

func (o PropertySources) RemoveIfPresent(propertySource PropertySource) {
	if propertySource == nil {
		return
	}
	_, index := o.findPropertySourceByName(propertySource.GetName())
	if index != -1 {
		o.sources = append(o.sources[:index], o.sources[index+1:]...)
	}
}

func (o PropertySources) findPropertySourceByName(name string) (PropertySource, int) {
	for index, source := range o.sources {
		if source.GetName() == name {
			return source, index
		}
	}
	return nil, -1
}

func (o PropertySources) GetSize() int {
	return len(o.sources)
}
