package property

type Resolver interface {
	ContainsProperty(name string) bool
	Property(name string) (string, bool)
	PropertyOrDefault(name string, defaultValue string) string
	ResolvePlaceholders(text string) (string, error)
}
