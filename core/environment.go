package core

import "os"

type Environment interface {
	PropertyResolver
}

type ConfigurableEnvironment interface {
	Environment
	GetPropertySources() PropertySources
	GetSystemEnvironment() []string
}

type StandardEnvironment struct {
	propertySources PropertySources
}

func NewStandardEnvironment() StandardEnvironment {
	return StandardEnvironment{
		propertySources: NewPropertySources(),
	}
}

func (env StandardEnvironment) GetPropertySources() PropertySources {
	return env.propertySources
}

func (env StandardEnvironment) GetSystemEnvironment() []string {
	return os.Environ()
}

func (env StandardEnvironment) ContainsProperty(name string) bool {
	return false
}

func (env StandardEnvironment) GetProperty(name string, defaultValue string) string {
	return ""
}
