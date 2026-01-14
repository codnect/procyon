package http

const (
	// PathValuesContextKey is the key for path values in the context
	PathValuesContextKey = "PathValues"
)

// PathValues represents the values of the path parameters.
type PathValues map[string]string

// Put adds a new path parameter with the provided name and value.
func (p PathValues) Put(name string, value string) {
	p[name] = value
}

// Value returns the value of the path parameter with the provided name.
func (p PathValues) Value(name string) (string, bool) {
	if val, ok := p[name]; ok {
		return val, true
	}

	return "", false
}

// Clear removes all path values.
func (p PathValues) Clear() {
	clear(p)
}
