package runtime

import (
	"os"
	"strings"
)

const (
	// NonOptionArgs represents the non-option arguments.
	NonOptionArgs = "nonOptionArgs"
)

// ArgumentsSource struct represents a source of arguments.
type ArgumentsSource struct {
	args *Arguments
}

// NewArgumentsSource function creates a new ArgumentsSource with the given arguments.
func NewArgumentsSource(args *Arguments) *ArgumentsSource {
	return &ArgumentsSource{
		args: args,
	}
}

// Name method returns the name of the source.
func (s *ArgumentsSource) Name() string {
	return "commandLineArgs"
}

// Source method returns the source of the arguments.
func (s *ArgumentsSource) Source() any {
	return s.args
}

// ContainsProperty method checks whether the argument with the given name exists.
func (s *ArgumentsSource) ContainsProperty(name string) bool {
	if NonOptionArgs == name {
		return len(s.args.NonOptionArgs()) != 0
	}

	return s.args.ContainsOption(name)
}

// Property method returns the value of the argument with the given name.
func (s *ArgumentsSource) Property(name string) (any, bool) {
	if NonOptionArgs == name {
		nonOptValues := s.args.NonOptionArgs()

		if nonOptValues != nil {
			return strings.Join(nonOptValues, ","), true
		}

		return nil, false
	}

	optValues := s.args.OptionValues(name)

	if optValues != nil {
		return strings.Join(optValues, ","), true
	}

	return nil, false
}

// PropertyOrDefault returns the value of the given argument name from the source.
// If the argument does not exist, it returns the default value.
func (s *ArgumentsSource) PropertyOrDefault(name string, defaultValue any) any {
	val, ok := s.Property(name)
	if !ok {
		return defaultValue
	}

	return val
}

// PropertyNames method returns the names of the arguments.
func (s *ArgumentsSource) PropertyNames() []string {
	return s.args.OptionNames()
}

// OptionValues method returns the values of the option argument with the given name.
func (s *ArgumentsSource) OptionValues(name string) []string {
	return s.args.OptionValues(name)
}

// NonOptionArgs method returns the non-option arguments.
func (s *ArgumentsSource) NonOptionArgs() []string {
	return s.args.NonOptionArgs()
}

// EnvironmentSource struct represents a source of environment properties.
type EnvironmentSource struct {
	variables map[string]string
}

// NewEnvironmentSource function creates a new EnvironmentSource.
func NewEnvironmentSource() *EnvironmentSource {
	source := &EnvironmentSource{
		variables: make(map[string]string, 0),
	}

	variables := os.Environ()

	for _, variable := range variables {
		index := strings.Index(variable, "=")

		if index != -1 {
			source.variables[variable[:index]] = variable[index+1:]
		}
	}

	return source
}

// Name method returns the name of the source.
func (s *EnvironmentSource) Name() string {
	return "systemEnvironment"
}

// Source method returns the source of the environment properties.
func (s *EnvironmentSource) Source() any {
	copyOfVariables := make(map[string]string)
	for key, value := range s.variables {
		copyOfVariables[key] = value
	}

	return copyOfVariables
}

// Property method returns the value of the environment property with the given name.
func (s *EnvironmentSource) Property(name string) (any, bool) {
	propertyName, exists := s.checkPropertyName(strings.ToLower(name))

	if exists {
		if value, ok := s.variables[propertyName]; ok {
			return value, true
		}
	}

	propertyName, exists = s.checkPropertyName(strings.ToUpper(name))

	if exists {
		if value, ok := s.variables[propertyName]; ok {
			return value, true
		}
	}

	return nil, false
}

// PropertyOrDefault returns the value of the given environment property name from the source.
// If the environment property does not exist, it returns the default value.
func (s *EnvironmentSource) PropertyOrDefault(name string, defaultValue any) any {
	value, ok := s.Property(name)

	if !ok {
		return defaultValue
	}

	return value
}

// ContainsProperty method checks whether the environment property with the given name exists.
func (s *EnvironmentSource) ContainsProperty(name string) bool {
	_, exists := s.checkPropertyName(strings.ToUpper(name))
	if exists {
		return true
	}

	_, exists = s.checkPropertyName(strings.ToLower(name))
	if exists {
		return true
	}

	return false
}

// PropertyNames method returns the names of the environment properties.
func (s *EnvironmentSource) PropertyNames() []string {
	keys := make([]string, 0, len(s.variables))

	for key, _ := range s.variables {
		keys = append(keys, key)
	}

	return keys
}

// checkPropertyName method checks the given property name in the environment variables.
func (s *EnvironmentSource) checkPropertyName(name string) (string, bool) {
	if s.contains(name) {
		return name, true
	}

	noHyphenPropertyName := strings.ReplaceAll(name, "-", "_")

	if name != noHyphenPropertyName && s.contains(noHyphenPropertyName) {
		return noHyphenPropertyName, true
	}

	noDotPropertyName := strings.ReplaceAll(name, ".", "_")

	if name != noDotPropertyName && s.contains(noDotPropertyName) {
		return noDotPropertyName, true
	}

	noHyphenAndNoDotName := strings.ReplaceAll(noDotPropertyName, "-", "_")

	if noDotPropertyName != noHyphenAndNoDotName && s.contains(noHyphenAndNoDotName) {
		return noHyphenAndNoDotName, true
	}

	return "", false
}

// contains method checks whether the environment property with the given name exists.
func (s *EnvironmentSource) contains(name string) bool {
	if _, ok := s.variables[name]; ok {
		return true
	}

	return false
}
