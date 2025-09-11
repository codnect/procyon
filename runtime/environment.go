// Copyright 2025 Codnect
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package runtime

import (
	"os"
	"strings"

	"codnect.io/procyon/runtime/property"
)

// Environment interface represents the application environment.
// It provides methods for accessing active and default profiles, checking if a profile is active,
// setting and adding active profiles, setting default profiles, merging environments,
// and accessing the property sources and property resolver.
type Environment interface {
	// ActiveProfiles returns the active profiles.
	ActiveProfiles() []string
	// DefaultProfiles returns the default profiles.
	DefaultProfiles() []string
	// IsProfileActive checks if the given profile is active.
	IsProfileActive(profile string) bool
	// SetActiveProfiles sets the active profiles.
	// This will replace any existing active profiles.
	SetActiveProfiles(profiles ...string) error
	// AddActiveProfiles adds the given profiles to the active profiles.
	// This will not replace any existing active profiles.
	AddActiveProfiles(profiles ...string) error
	// SetDefaultProfiles sets the default profiles.
	// This will replace any existing default profiles.
	SetDefaultProfiles(profiles ...string) error
	// Merge merges the given environment into this environment.
	// The active and default profiles of the given environment will be added to this environment.
	// The property sources of the given environment will be added to this environment.
	Merge(other Environment) error
	// PropertySources returns the property sources.
	PropertySources() *property.SourceList
	// PropertyResolver returns the property resolver.
	PropertyResolver() property.Resolver
}

// EnvironmentCapable is an interface that indicates the ability to
// provide access to an Environment instance.
type EnvironmentCapable interface {
	// Environment returns the associated Environment.
	Environment() Environment
}

// EnvPropertySource struct represents a source of environment properties.
type EnvPropertySource struct {
	variables map[string]string
}

// NewEnvPropertySource function creates a new EnvPropertySource.
func NewEnvPropertySource() *EnvPropertySource {
	source := &EnvPropertySource{
		variables: make(map[string]string),
	}

	variables := os.Environ()

	for _, variable := range variables {
		index := strings.Index(variable, "=")
		source.variables[variable[:index]] = variable[index+1:]
	}

	return source
}

// Name method returns the name of the source.
func (s *EnvPropertySource) Name() string {
	return "systemEnvironment"
}

// Underlying returns the underlying source object.
func (s *EnvPropertySource) Underlying() any {
	copyOfVariables := make(map[string]string)
	for key, value := range s.variables {
		copyOfVariables[key] = value
	}

	return copyOfVariables
}

// ContainsProperty method checks whether the environment property with the given name exists.
func (s *EnvPropertySource) ContainsProperty(name string) bool {
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

// Property method returns the value of the environment property with the given name.
func (s *EnvPropertySource) Property(name string) (any, bool) {
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
func (s *EnvPropertySource) PropertyOrDefault(name string, defaultValue any) any {
	value, ok := s.Property(name)

	if !ok {
		return defaultValue
	}

	return value
}

// PropertyNames method returns the names of the environment properties.
func (s *EnvPropertySource) PropertyNames() []string {
	keys := make([]string, 0, len(s.variables))

	for key := range s.variables {
		keys = append(keys, key)
	}

	return keys
}

// checkPropertyName method checks the given property name in the environment variables.
func (s *EnvPropertySource) checkPropertyName(name string) (string, bool) {
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
func (s *EnvPropertySource) contains(name string) bool {
	if _, ok := s.variables[name]; ok {
		return true
	}

	return false
}
