package env

import (
	"fmt"
	"github.com/procyon-projects/procyon/app/env/property"
	"strings"
	"sync"
)

type Variables map[string]string

type Environment interface {
	ActiveProfiles() []string
	DefaultProfiles() []string
	IsProfileActive(profile string) bool

	SetActiveProfiles(profiles ...string) error
	AddActiveProfile(profile ...string) error
	SetDefaultProfiles(profiles ...string) error
	Merge(parent Environment)

	Variables() Variables
	PropertySources() *property.Sources
	PropertyResolver() property.Resolver
}

type environment struct {
	activeProfiles  map[string]struct{}
	defaultProfiles map[string]struct{}

	sources             *property.Sources
	resolver            property.Resolver
	activeProfilesOnce  sync.Once
	defaultProfilesOnce sync.Once
	resolverOnce        sync.Once
	mu                  sync.RWMutex
}

func New() Environment {
	return WithSources(property.NewPropertySources())
}

func WithSources(sources *property.Sources) Environment {
	if sources == nil {
		panic(fmt.Errorf("env: sources cannot be nil"))
	}

	return &environment{
		activeProfiles: map[string]struct{}{},
		defaultProfiles: map[string]struct{}{
			"default": {},
		},
		sources:             sources,
		activeProfilesOnce:  sync.Once{},
		defaultProfilesOnce: sync.Once{},
		mu:                  sync.RWMutex{},
	}
}

func (e *environment) validateProfile(profile string) error {
	if strings.TrimSpace(profile) == "" {
		return fmt.Errorf("env: `%s` is a invalid profile", profile)
	}

	return nil
}

func (e *environment) doGetActiveProfiles() {
	e.activeProfilesOnce.Do(func() {
		propertyValue, ok := e.PropertyResolver().Property("procyon.profiles.active")

		if ok {
			activeProfiles := strings.Split(strings.TrimSpace(propertyValue), ",")
			err := e.SetActiveProfiles(activeProfiles...)

			if err != nil {
				panic(err)
			}
		}
	})
}

func (e *environment) ActiveProfiles() []string {
	e.doGetActiveProfiles()

	defer e.mu.Unlock()
	e.mu.Lock()

	activeProfiles := make([]string, 0)
	for profile := range e.activeProfiles {
		activeProfiles = append(activeProfiles, profile)
	}

	return activeProfiles
}

func (e *environment) doGetDefaultProfiles() {
	e.defaultProfilesOnce.Do(func() {
		propertyValue, ok := e.PropertyResolver().Property("procyon.profiles.default")

		if ok {
			defaultProfiles := strings.Split(strings.TrimSpace(propertyValue), ",")
			err := e.SetDefaultProfiles(defaultProfiles...)

			if err != nil {
				panic(err)
			}
		}
	})
}

func (e *environment) DefaultProfiles() []string {
	e.doGetDefaultProfiles()

	defer e.mu.Unlock()
	e.mu.Lock()

	profiles := make([]string, 0)
	for profile := range e.defaultProfiles {
		profiles = append(profiles, profile)
	}

	return profiles
}

func (e *environment) IsProfileActive(profile string) bool {
	defer e.mu.Unlock()
	e.mu.Lock()

	if _, ok := e.activeProfiles[profile]; ok {
		return true
	}

	return false
}

func (e *environment) clearActiveProfiles() {
	defer e.mu.Unlock()
	e.mu.Lock()

	for profile := range e.activeProfiles {
		delete(e.activeProfiles, profile)
	}
}

func (e *environment) SetActiveProfiles(profiles ...string) error {
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
func (e *environment) AddActiveProfile(profiles ...string) error {
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

func (e *environment) clearDefaultProfiles() {
	defer e.mu.Unlock()
	e.mu.Lock()

	for profile := range e.defaultProfiles {
		delete(e.defaultProfiles, profile)
	}
}

func (e *environment) SetDefaultProfiles(profiles ...string) error {
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

func (e *environment) Merge(parent Environment) {
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

func (e *environment) Variables() Variables {
	return nil
}

func (e *environment) PropertySources() *property.Sources {
	return e.sources
}

func (e *environment) PropertyResolver() property.Resolver {
	defer e.mu.Unlock()
	e.mu.Lock()

	if e.resolver == nil {
		e.resolver = property.NewSourcesResolver(e.sources)
	}

	return e.resolver
}
