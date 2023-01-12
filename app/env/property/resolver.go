package property

import "fmt"

type Resolver interface {
	ContainsProperty(name string) bool
	Property(name string) (string, bool)
	PropertyOrDefault(name string, defaultValue string) string
	ResolvePlaceholders(text string) string
	ResolveRequiredPlaceholders(text string) (string, error)
}

type SourcesResolver struct {
	sources *Sources
}

func NewSourcesResolver(sources *Sources) *SourcesResolver {
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

func (r *SourcesResolver) ResolvePlaceholders(s string) string {
	result, _ := r.resolveRequiredPlaceHolders(s, true)
	return result
}

func (r *SourcesResolver) ResolveRequiredPlaceholders(s string) (string, error) {
	return r.resolveRequiredPlaceHolders(s, false)
}

func (r *SourcesResolver) resolveRequiredPlaceHolders(s string, continueOnError bool) (string, error) {
	var buf []byte

	i := 0
	for j := 0; j < len(s); j++ {
		if s[j] == '$' && j+1 < len(s) {
			if buf == nil {
				buf = make([]byte, 0, 2*len(s))
			}

			buf = append(buf, s[i:j]...)
			name, w := r.getPlaceholderName(s[j+1:])

			if name == "" && w > 0 {
			} else if name == "" {
				buf = append(buf, s[j])
			} else {
				value, ok := r.Property(name)

				if !ok && !continueOnError {
					return "", fmt.Errorf("could not resolve placeholder '%s'", s[j:i+w+1])
				}

				if continueOnError {
					buf = append(buf, s[j:i+w+1]...)
				} else {
					buf = append(buf, value...)
				}
			}

			j += w
			i = j + 1
		}
	}

	if buf == nil {
		return s, nil
	}

	return string(buf) + s[i:], nil
}

func (r *SourcesResolver) getPlaceholderName(s string) (string, int) {
	switch {
	case s[0] == '{':
		if len(s) > 2 && isSpecialVar(s[1]) && s[2] == '}' {
			return s[1:2], 3
		}

		for i := 1; i < len(s); i++ {
			if s[i] == '}' {
				if i == 1 {
					return "", 2
				}
				return s[1:i], i + 1
			}
		}
		return "", 1
	case isSpecialVar(s[0]):
		return s[0:1], 1
	}

	var i int
	for i = 0; i < len(s) && isAlphaNum(s[i]); i++ {
	}

	return s[:i], i
}

func isSpecialVar(c uint8) bool {
	switch c {
	case '*', '#', '$', '@', '!', '?', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	}
	return false
}

func isAlphaNum(c uint8) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}
