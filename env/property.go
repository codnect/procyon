package env

type PropertySource interface {
	GetProperty(name string) interface{}
	ContainsProperty(name string) bool
}

type BasePropertySource struct {
	PropertySource
	Name   string
	Source interface{}
}

func NewPropertySource(name string) *BasePropertySource {
	return &BasePropertySource{
		Name: name,
	}
}

func NewPropertySourceWithSource(name string, source interface{}) *BasePropertySource {
	return &BasePropertySource{
		Name:   name,
		Source: source,
	}
}

func (source *BasePropertySource) ContainsProperty(name string) bool {
	return source.GetProperty(name) != nil
}

type EnumerablePropertySource interface {
	GetPropertyNames() []string
}

type BaseEnumerablePropertySource struct {
	EnumerablePropertySource
	*BasePropertySource
}

func NewEnumerablePropertySource(name string) *BaseEnumerablePropertySource {
	return &BaseEnumerablePropertySource{
		BasePropertySource: NewPropertySource(name),
	}
}

func NewEnumerablePropertySourceWithSource(name string, source interface{}) *BaseEnumerablePropertySource {
	return &BaseEnumerablePropertySource{
		BasePropertySource: NewPropertySourceWithSource(name, source),
	}
}

func (source *BaseEnumerablePropertySource) ContainsProperty(name string) bool {
	for _, propertyName := range source.GetPropertyNames() {
		if propertyName == name {
			return true
		}
	}
	return false
}
