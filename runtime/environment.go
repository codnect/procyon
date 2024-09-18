package runtime

import (
	"codnect.io/procyon/runtime/property"
	"fmt"
	"strings"
	"sync"
)

// Environment interface represents the application environment.
// It provides methods for accessing active and default profiles, checking if a profile is active,
// setting and adding active profiles, setting default profiles, merging environments,
// and accessing the property sources and property resolver.
type Environment interface {
	ActiveProfiles() []string
	DefaultProfiles() []string
	IsProfileActive(profile string) bool

	SetActiveProfiles(profiles ...string) error
	AddActiveProfile(profile ...string) error
	SetDefaultProfiles(profiles ...string) error
	Merge(parent Environment)

	PropertySources() *property.Sources
	PropertyResolver() property.Resolver
}

// DefaultEnvironment struct represents the default implementation of the Environment interface.
type DefaultEnvironment struct {
	activeProfiles  map[string]struct{}
	defaultProfiles map[string]struct{}

	sources             *property.Sources
	resolver            property.Resolver
	activeProfilesOnce  sync.Once
	defaultProfilesOnce sync.Once
	resolverOnce        sync.Once
	mu                  sync.RWMutex
}

// NewDefaultEnvironment function creates a new DefaultEnvironment.
func NewDefaultEnvironment() *DefaultEnvironment {
	return &DefaultEnvironment{
		activeProfiles: map[string]struct{}{},
		defaultProfiles: map[string]struct{}{
			"default": {},
		},
		sources:             property.NewSources(),
		activeProfilesOnce:  sync.Once{},
		defaultProfilesOnce: sync.Once{},
		mu:                  sync.RWMutex{},
	}
}

// validateProfile function validates a profile name
func (e *DefaultEnvironment) validateProfile(profile string) error {
	if strings.TrimSpace(profile) == "" {
		return fmt.Errorf("'%s' is a invalid profile", profile)
	}

	return nil
}

// doGetActiveProfiles function retrieves the active profiles from the property resolver.
func (e *DefaultEnvironment) doGetActiveProfiles() {
	e.activeProfilesOnce.Do(func() {
		propertyValue, ok := e.PropertyResolver().Property("procyon.profiles.active")

		if ok {
			activeProfiles := strings.Split(strings.TrimSpace(propertyValue.(string)), ",")
			err := e.SetActiveProfiles(activeProfiles...)

			if err != nil {
				panic(err)
			}
		}
	})
}

// ActiveProfiles method returns the active profiles.
func (e *DefaultEnvironment) ActiveProfiles() []string {
	e.doGetActiveProfiles()

	defer e.mu.Unlock()
	e.mu.Lock()

	activeProfiles := make([]string, 0)
	for profile := range e.activeProfiles {
		activeProfiles = append(activeProfiles, profile)
	}

	return activeProfiles
}

// doGetDefaultProfiles function retrieves the default profiles from the property resolver.
func (e *DefaultEnvironment) doGetDefaultProfiles() {
	e.defaultProfilesOnce.Do(func() {
		propertyValue, ok := e.PropertyResolver().Property("procyon.profiles.default")

		if ok {
			defaultProfiles := strings.Split(strings.TrimSpace(propertyValue.(string)), ",")
			err := e.SetDefaultProfiles(defaultProfiles...)

			if err != nil {
				panic(err)
			}
		}
	})
}

// DefaultProfiles method returns the default profiles.
func (e *DefaultEnvironment) DefaultProfiles() []string {
	e.doGetDefaultProfiles()

	defer e.mu.Unlock()
	e.mu.Lock()

	profiles := make([]string, 0)
	for profile := range e.defaultProfiles {
		profiles = append(profiles, profile)
	}

	return profiles
}

// IsProfileActive method checks if a profile is active.
func (e *DefaultEnvironment) IsProfileActive(profile string) bool {
	defer e.mu.Unlock()
	e.mu.Lock()

	if _, ok := e.activeProfiles[profile]; ok {
		return true
	}

	return false
}

// clearActiveProfiles method clears the active profiles.
func (e *DefaultEnvironment) clearActiveProfiles() {
	defer e.mu.Unlock()
	e.mu.Lock()

	for profile := range e.activeProfiles {
		delete(e.activeProfiles, profile)
	}
}

// SetActiveProfiles method sets the active profiles.
func (e *DefaultEnvironment) SetActiveProfiles(profiles ...string) error {
	e.clearActiveProfiles()

	defer e.mu.Unlock()
	e.mu.Lock()

	for _, profile := range profiles {
		err := e.validateProfile(profile)
		if err != nil {
			return err
		}

		e.activeProfiles[profile] = struct{}{}
	}

	return nil
}

// AddActiveProfile method adds active profiles.
func (e *DefaultEnvironment) AddActiveProfile(profiles ...string) error {
	for _, profile := range profiles {
		err := e.validateProfile(profile)

		if err != nil {
			return err
		}
	}

	e.doGetActiveProfiles()

	defer e.mu.Unlock()
	e.mu.Lock()

	for _, profile := range profiles {
		e.activeProfiles[profile] = struct{}{}
	}

	return nil
}

// clearDefaultProfiles method clears the default profiles.
func (e *DefaultEnvironment) clearDefaultProfiles() {
	defer e.mu.Unlock()
	e.mu.Lock()

	for profile := range e.defaultProfiles {
		delete(e.defaultProfiles, profile)
	}
}

// SetDefaultProfiles method sets the default profiles.
func (e *DefaultEnvironment) SetDefaultProfiles(profiles ...string) error {
	e.clearDefaultProfiles()

	defer e.mu.Unlock()
	e.mu.Lock()

	for _, profile := range profiles {
		err := e.validateProfile(profile)
		if err != nil {
			return err
		}

		e.defaultProfiles[profile] = struct{}{}
	}

	return nil
}

// Merge method merges the current environment with a parent environment.
func (e *DefaultEnvironment) Merge(parent Environment) {
	parentPropertySources := parent.PropertySources().ToSlice()

	if len(parentPropertySources) != 0 {
		for _, propertySource := range parentPropertySources {
			if !e.sources.Contains(propertySource.Name()) {
				e.sources.AddLast(propertySource)
			}
		}
	}

	defer e.mu.Unlock()
	e.mu.Lock()

	parentActiveProfiles := parent.ActiveProfiles()
	if len(parentActiveProfiles) != 0 {
		for _, profile := range parentActiveProfiles {
			e.activeProfiles[profile] = struct{}{}
		}
	}

	parentDefaultProfiles := parent.DefaultProfiles()
	if len(parentDefaultProfiles) != 0 {
		for _, profile := range parentDefaultProfiles {
			e.defaultProfiles[profile] = struct{}{}
		}
	}
}

// PropertySources method returns the property sources.
func (e *DefaultEnvironment) PropertySources() *property.Sources {
	return e.sources
}

// PropertyResolver method returns the property resolver.
func (e *DefaultEnvironment) PropertyResolver() property.Resolver {
	defer e.mu.Unlock()
	e.mu.Lock()

	if e.resolver == nil {
		e.resolver = property.NewSourcesResolver(e.sources)
	}

	return e.resolver
}
