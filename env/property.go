package env

type PropertySource interface {
	GetProperty(name string) interface{}
	ContainsProperty(name string) bool
}

type AbstractPropertySource struct {
	PropertySource
	Name   string
	Source interface{}
}

func NewAbstractPropertySourceWithSource(name string, source interface{}) AbstractPropertySource {
	propertySource := AbstractPropertySource{
		Name:   name,
		Source: source,
	}
	return propertySource
}

func (source AbstractPropertySource) ContainsProperty(name string) bool {
	return source.GetProperty(name) != nil
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
