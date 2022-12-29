package property

type Resolver interface {
	ContainsProperty(name string) bool
	Property(name string) (string, bool)
	PropertyOrDefault(name string, defaultValue string) string
	ResolvePlaceholders(text string) (string, error)
}

type SourcesResolver struct {
	sources Sources
}

func NewSourcesResolver(sources Sources) *SourcesResolver {
	if sources == nil {
		panic("property: sources cannot be nil")
	}

	return &SourcesResolver{
		sources: sources,
	}
}

func (r *SourcesResolver) ContainsProperty(name string) bool {
	return r.sources.Contains(name)
}

func (r *SourcesResolver) Property(name string) (string, bool) {
	for _, source := range r.sources.ToSlice() {
		if value, ok := source.Property(name); ok {
			return value.(string), true
		}
	}

	return "", false
}

func (r *SourcesResolver) PropertyOrDefault(name string, defaultValue string) string {
	for _, source := range r.sources.ToSlice() {
		if value, ok := source.Property(name); ok {
			return value.(string)
		}
	}

	return defaultValue
}

func (r *SourcesResolver) ResolvePlaceholders(text string) (string, error) {
	return "", nil
}
