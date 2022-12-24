package env

import (
	"github.com/procyon-projects/procyon/app/env/property"
)

type Variables map[string]string

type Environment interface {
	ActiveProfiles() []string
	DefaultProfiles() []string
	AcceptProfiles(profiles ...string)
	IsProfileActive(profile string)

	SetActiveProfiles(profiles ...string)
	AddActiveProfile(profile string)
	SetDefaultProfiles(profiles ...string)
	Merge(other Environment)

	Variables() Variables
	PropertySources() property.Sources
	PropertyResolver() property.Resolver
}

type environment struct {
	activeProfiles  map[string]struct{}
	defaultProfiles map[string]struct{}

	sources property.Sources
}

func New() Environment {
	return &environment{}
}

func WithSources(sources property.Sources) Environment {
	return &environment{}
}

func (e *environment) ActiveProfiles() []string {
	return nil
}

func (e *environment) DefaultProfiles() []string {
	return nil
}

func (e *environment) AcceptProfiles(profiles ...string) {

}
func (e *environment) IsProfileActive(profile string) {

}

func (e *environment) SetActiveProfiles(profiles ...string) {

}
func (e *environment) AddActiveProfile(profile string) {

}
func (e *environment) SetDefaultProfiles(profiles ...string) {

}
func (e *environment) Merge(other Environment) {

}

func (e *environment) Variables() Variables {
	return nil
}

func (e *environment) PropertySources() property.Sources {
	return nil
}

func (e *environment) PropertyResolver() property.Resolver {
	return nil
}
